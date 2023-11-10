package phase

import (
	"strings"

	"github.com/deepsquare-io/cfctl/pkg/apis/cfctl.clusterfactory.io/v1beta1/cluster"

	// anonymous import is needed to load the os configurers
	_ "github.com/deepsquare-io/cfctl/configurer"
	// anonymous import is needed to load the os configurers
	_ "github.com/deepsquare-io/cfctl/configurer/linux"
	// anonymous import is needed to load the os configurers
	_ "github.com/deepsquare-io/cfctl/configurer/linux/enterpriselinux"

	log "github.com/sirupsen/logrus"
)

// DetectOS performs remote OS detection
type DetectOS struct {
	GenericPhase
}

// Title for the phase
func (p *DetectOS) Title() string {
	return "Detect host operating systems"
}

// Run the phase
func (p *DetectOS) Run() error {
	return p.parallelDo(p.Config.Spec.Hosts, func(h *cluster.Host) error {
		if h.OSIDOverride != "" {
			log.Infof("%s: OS ID has been manually set to %s", h, h.OSIDOverride)
		}
		if err := h.ResolveConfigurer(); err != nil {
			p.SetProp("missing-support", h.OSVersion.String())
			if h.OSVersion.IDLike != "" {
				log.Debugf(
					"%s: trying to find a fallback OS support module for %s using os-release ID_LIKE '%s'",
					h,
					h.OSVersion.String(),
					h.OSVersion.IDLike,
				)
				for _, id := range strings.Split(h.OSVersion.IDLike, " ") {
					h.OSVersion.ID = id
					if err := h.ResolveConfigurer(); err == nil {
						log.Warnf(
							"%s: using '%s' as OS support fallback for %s",
							h,
							id,
							h.OSVersion.String(),
						)
						return nil
					}
				}
			}
			return err
		}
		os := h.OSVersion.String()
		p.IncProp(os)
		log.Infof("%s: is running %s", h, os)

		return nil
	})
}
