package phase

import (
	"github.com/SquareFactory/cfctl/pkg/apis/cfctl.clusterfactory.io/v1beta1"
	"github.com/SquareFactory/cfctl/pkg/apis/cfctl.clusterfactory.io/v1beta1/cluster"
	"github.com/k0sproject/rig/exec"
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

// Run the phase
func (p *SymlinkKubelet) Run() error {
	for _, h := range p.hosts {
		if err := ensureDir(h, "/var/lib/k0s/kubelet", "0755", "0"); err != nil {
			return err
		}
		if err := h.Exec("if [ -L /var/lib/kubelet ]; then echo symlink already exists; else rm -rf /var/lib/kubelet && ln -s /var/lib/k0s/kubelet /var/lib/kubelet; fi", exec.Sudo(h)); err != nil {
			return err
		}
	}
	return nil
}
