package phase

import (
	"fmt"

	"github.com/deepsquare-io/cfctl/pkg/apis/cfctl.clusterfactory.io/v1beta1"
	"github.com/deepsquare-io/cfctl/pkg/apis/cfctl.clusterfactory.io/v1beta1/cluster"
	"github.com/k0sproject/rig/exec"
	log "github.com/sirupsen/logrus"
)

var _ phase = &SymlinkKubelet{}

type SymlinkKubelet struct {
	GenericPhase
	hosts []*cluster.Host
}

// Title returns the title for the phase
func (p *SymlinkKubelet) Title() string {
	return "Symlink kubelet"
}

// Prepare the phase
func (p *SymlinkKubelet) Prepare(config *v1beta1.Cluster) error {
	p.Config = config
	p.hosts = p.Config.Spec.Hosts
	return nil
}

func (p *SymlinkKubelet) ensureDir(h *cluster.Host, dir, perm, owner string) error {
	log.Debugf("%s: ensuring directory %s", h, dir)
	if h.Configurer.FileExist(h, dir) {
		return nil
	}

	err := p.Wet(
		h,
		fmt.Sprintf("create a directory for uploading: `mkdir -p \"%s\"`", dir),
		func() error {
			return h.Configurer.MkDir(h, dir, exec.Sudo(h))
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	if perm == "" {
		perm = "0755"
	}

	err = p.Wet(h, fmt.Sprintf("set permissions for directory %s to %s", dir, perm), func() error {
		return h.Configurer.Chmod(h, dir, perm, exec.Sudo(h))
	})
	if err != nil {
		return fmt.Errorf("failed to set permissions for directory %s: %w", dir, err)
	}

	if owner != "" {
		err = p.Wet(h, fmt.Sprintf("set owner for directory %s to %s", dir, owner), func() error {
			return h.Execf(`chown "%s" "%s"`, owner, dir, exec.Sudo(h))
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// Run the phase
func (p *SymlinkKubelet) Run() error {
	for _, h := range p.hosts {
		if err := p.ensureDir(h, "/var/lib/k0s/kubelet", "0755", "0"); err != nil {
			return err
		}
		if err := h.Exec("sh -c 'if [ -L /var/lib/kubelet ]; then echo symlink already exists; else rm -rf /var/lib/kubelet && ln -s /var/lib/k0s/kubelet /var/lib/kubelet; fi'", exec.Sudo(h)); err != nil {
			return err
		}
	}
	return nil
}
