package phase

import (
	"github.com/k0sproject/k0sctl/pkg/apis/k0sctl.k0sproject.io/v1beta1"
	"github.com/k0sproject/k0sctl/pkg/apis/k0sctl.k0sproject.io/v1beta1/cluster"
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
		if err := h.Exec("sh -c \"if [ -L /var/lib/kubelet ]; then echo symlink already exists; else rm -rf /var/lib/kubelet && ln -s /var/lib/k0s/kubelet /var/lib/kubelet; fi\""); err != nil {
			return err
		}
	}
	return nil
}
