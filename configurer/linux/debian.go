package linux

import (
	"github.com/deepsquare-io/cfctl/configurer"
	"github.com/k0sproject/rig"
	"github.com/k0sproject/rig/os/linux"
	"github.com/k0sproject/rig/os/registry"
)

// Debian provides OS support for Debian systems
type Debian struct {
	linux.Ubuntu
	configurer.Linux
}

func init() {
	registry.RegisterOSModule(
		func(os rig.OSVersion) bool {
			return os.ID == "debian"
		},
		func() interface{} {
			return &Debian{}
		},
	)
}
