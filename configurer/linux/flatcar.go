package linux

import (
	"errors"

	"github.com/k0sproject/rig"
	"github.com/k0sproject/rig/os"
	"github.com/k0sproject/rig/os/registry"
)

type Flatcar struct {
	BaseLinux
	os.Linux
}

func init() {
	registry.RegisterOSModule(
		func(os rig.OSVersion) bool {
			return os.ID == "flatcar"
		},
		func() interface{} {
			fc := &Flatcar{}
			fc.SetPath("K0sBinaryPath", "/opt/bin/k0s")
			return fc
		},
	)
}

func (l *Flatcar) InstallPackage(h os.Host, pkg ...string) error {
	return errors.New("FlatcarContainerLinux does not support installing packages manually")
}
