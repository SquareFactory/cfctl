package phase

import (
	"context"
	"fmt"
	gos "os"
	"sync"
	"time"

	"github.com/deepsquare-io/cfctl/analytics"
	"github.com/deepsquare-io/cfctl/pkg/apis/cfctl.clusterfactory.io/v1beta1"
	"github.com/deepsquare-io/cfctl/pkg/apis/cfctl.clusterfactory.io/v1beta1/cluster"
	"github.com/deepsquare-io/cfctl/pkg/retry"
	"github.com/k0sproject/rig/exec"
	log "github.com/sirupsen/logrus"
)

// Lock acquires an exclusive cfctl lock on hosts
type Lock struct {
	GenericPhase
	cfs        []func()
	instanceID string
	m          sync.Mutex
	wg         sync.WaitGroup
}

// Prepare the phase
func (p *Lock) Prepare(c *v1beta1.Cluster) error {
	p.Config = c
	mid, _ := analytics.MachineID()
	p.instanceID = fmt.Sprintf("%s-%d", mid, gos.Getpid())
	return nil
}

// Title for the phase
func (p *Lock) Title() string {
	return "Acquire exclusive host lock"
}

// Cancel releases the lock
func (p *Lock) Cancel() {
	p.m.Lock()
	defer p.m.Unlock()
	for _, f := range p.cfs {
		f()
	}
	p.wg.Wait()
}

// CleanUp calls Cancel to release the lock
func (p *Lock) CleanUp() {
	p.Cancel()
}

// Run the phase
func (p *Lock) Run() error {
	if err := p.parallelDo(p.Config.Spec.Hosts, p.startLock); err != nil {
		return err
	}
	return p.Config.Spec.Hosts.ParallelEach(p.startTicker)
}

func (p *Lock) startTicker(h *cluster.Host) error {
	p.wg.Add(1)
	lfp := h.Configurer.CfctlLockFilePath(h)
	ticker := time.NewTicker(10 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	p.m.Lock()
	p.cfs = append(p.cfs, cancel)
	p.m.Unlock()

	go func() {
		log.Debugf("%s: started periodic update of lock file %s timestamp", h, lfp)
		for {
			select {
			case <-ticker.C:
				if err := h.Configurer.Touch(h, lfp, time.Now(), exec.Sudo(h)); err != nil {
					log.Warnf("%s: failed to touch lock file: %s", h, err)
				}
			case <-ctx.Done():
				log.Debugf("%s: stopped lock cycle, removing file", h)
				if err := h.Configurer.DeleteFile(h, lfp); err != nil {
					log.Warnf("%s: failed to remove host lock file: %s", h, err)
				}
				p.wg.Done()
				return
			}
		}
	}()

	return nil
}

func (p *Lock) startLock(h *cluster.Host) error {
	return retry.Times(context.TODO(), 10, func(_ context.Context) error {
		return p.tryLock(h)
	})
}

func (p *Lock) tryLock(h *cluster.Host) error {
	lfp := h.Configurer.CfctlLockFilePath(h)

	if err := h.Configurer.UpsertFile(h, lfp, p.instanceID); err != nil {
		stat, err := h.Configurer.Stat(h, lfp, exec.Sudo(h))
		if err != nil {
			return fmt.Errorf("lock file disappeared: %w", err)
		}
		content, err := h.Configurer.ReadFile(h, lfp)
		if err != nil {
			return fmt.Errorf("failed to read lock file:  %w", err)
		}
		if content != p.instanceID {
			if time.Since(stat.ModTime()) < 30*time.Second {
				return fmt.Errorf("another instance of cfctl is currently operating on the host")
			}
			_ = h.Configurer.DeleteFile(h, lfp)
			return fmt.Errorf("removed existing expired lock file")
		}
	}

	return nil
}
