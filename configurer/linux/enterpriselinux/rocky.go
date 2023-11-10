package enterpriselinux

import (
	"github.com/deepsquare-io/cfctl/configurer"
	k0slinux "github.com/deepsquare-io/cfctl/configurer/linux"
	"github.com/k0sproject/rig"
	"github.com/k0sproject/rig/os/registry"
)

// RockyLinux provides OS support for RockyLinux
type RockyLinux struct {
	k0slinux.EnterpriseLinux
	configurer.Linux
}

func init() {
	registry.RegisterOSModule(
		func(os rig.OSVersion) bool {
			return os.ID == "rocky"
		},
		func() interface{} {
			return &RockyLinux{}
		},
	)
}
