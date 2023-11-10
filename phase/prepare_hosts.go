package phase

import (
	"fmt"

	"github.com/SquareFactory/cfctl/pkg/apis/cfctl.clusterfactory.io/v1beta1/cluster"
	"github.com/k0sproject/rig/os"
	"github.com/k0sproject/version"
	log "github.com/sirupsen/logrus"
)

var iptablesEmbeddedSince = version.MustConstraint(">= v1.22.1+k0s.0")

// PrepareHosts installs required packages and so on on the hosts.
type PrepareHosts struct {
	GenericPhase
}

// Title for the phase
func (p *PrepareHosts) Title() string {
	return "Prepare hosts"
}

// Run the phase
func (p *PrepareHosts) Run() error {
	return p.parallelDo(p.Config.Spec.Hosts, p.prepareHost)
}

type prepare interface {
	Prepare(os.Host) error
}

func (p *PrepareHosts) prepareHost(h *cluster.Host) error {
	if c, ok := h.Configurer.(prepare); ok {
		if err := c.Prepare(h); err != nil {
			return err
		}
	}

	if len(h.Environment) > 0 {
		log.Infof("%s: updating environment", h)
		if err := h.Configurer.UpdateEnvironment(h, h.Environment); err != nil {
			return err
		}
	}

	var pkgs []string

	if h.NeedCurl() {
		pkgs = append(pkgs, "curl")
	}

	// iptables is only required for very old versions of k0s
	if p.Config.Spec.K0s.Version != nil && !iptablesEmbeddedSince.Check(p.Config.Spec.K0s.Version) && h.NeedIPTables() { //nolint:staticcheck
		pkgs = append(pkgs, "iptables")
	}

	if h.NeedInetUtils() {
		pkgs = append(pkgs, "inetutils")
	}

	for _, pkg := range pkgs {
		err := p.Wet(h, fmt.Sprintf("install package %s", pkg), func() error {
			log.Infof("%s: installing package %s", h, pkg)
			return h.Configurer.InstallPackage(h, pkg)
		})
		if err != nil {
			return err
		}
	}

	if h.Configurer.IsContainer(h) {
		log.Infof("%s: is a container, applying a fix", h)
		if err := h.Configurer.FixContainer(h); err != nil {
			return err
		}
	}

	return nil
}
