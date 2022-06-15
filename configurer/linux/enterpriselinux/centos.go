package enterpriselinux

import (
	"github.com/SquareFactory/cfctl/configurer"
	k0slinux "github.com/SquareFactory/cfctl/configurer/linux"
	"github.com/k0sproject/rig"
	"github.com/k0sproject/rig/os/registry"
)

// CentOS provides OS support for CentOS
type CentOS struct {
	k0slinux.EnterpriseLinux
	configurer.Linux
}

func init() {
	registry.RegisterOSModule(
		func(os rig.OSVersion) bool {
			return os.ID == "centos"
		},
		func() interface{} {
			linuxType := &CentOS{}
			linuxType.PathFuncs = interface{}(linuxType).(configurer.PathFuncs)
			return linuxType
		},
	)
}
