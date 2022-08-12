package linux

import (
	"github.com/SquareFactory/cfctl/configurer"
	"github.com/k0sproject/rig"
	"github.com/k0sproject/rig/os/registry"
)

// OpenSUSE provides OS support for OpenSUSE
type OpenSUSE struct {
	SLES
}

func init() {
	registry.RegisterOSModule(
		func(os rig.OSVersion) bool {
			return os.ID == "opensuse" || os.ID == "opensuse-microos"
		},
		func() interface{} {
			linuxType := &OpenSUSE{}
			linuxType.PathFuncs = interface{}(linuxType).(configurer.PathFuncs)
			return linuxType
		},
	)
}
