package enterpriselinux

import (
	"github.com/SquareFactory/cfctl/configurer"
	k0slinux "github.com/SquareFactory/cfctl/configurer/linux"
	"github.com/k0sproject/rig"
	"github.com/k0sproject/rig/os/registry"
	"strings"
)

// RHEL provides OS support for RedHat Enterprise Linux
type RHEL struct {
	k0slinux.EnterpriseLinux
}

func init() {
	registry.RegisterOSModule(
		func(os rig.OSVersion) bool {
			return os.ID == "rhel" && !strings.Contains(os.Name, "CoreOS")
		},
		func() interface{} {
			linuxType := &RHEL{}
			linuxType.PathFuncs = interface{}(linuxType).(configurer.PathFuncs)
			return linuxType
		},
	)
}
