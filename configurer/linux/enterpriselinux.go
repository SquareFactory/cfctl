package linux

import (
	"github.com/SquareFactory/cfctl/configurer"
	"github.com/k0sproject/rig/os/linux"
)

// EnterpriseLinux is a base package for several RHEL-like enterprise linux distributions
type EnterpriseLinux struct {
	linux.EnterpriseLinux
	configurer.Linux
}
