package enterpriselinux

import (
	"strings"

	"github.com/SquareFactory/cfctl/configurer"
	k0slinux "github.com/SquareFactory/cfctl/configurer/linux"

	"github.com/k0sproject/rig"
	"github.com/k0sproject/rig/os/registry"
)

// Fedora provides OS support for Fedora
type Fedora struct {
	k0slinux.EnterpriseLinux
	configurer.Linux
}

func init() {
	registry.RegisterOSModule(
		func(os rig.OSVersion) bool {
			return os.ID == "fedora" && !strings.Contains(os.Name, "CoreOS")
		},
		func() interface{} {
			return &Fedora{}
		},
	)
}
