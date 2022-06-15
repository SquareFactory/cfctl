package linux

import (
	"testing"

	"github.com/SquareFactory/cfctl/pkg/apis/cfctl.clusterfactory.io/v1beta1/cluster"
)

func TestAlpineConfigurerInterface(t *testing.T) {
	h := cluster.Host{}
	h.Configurer = Alpine{}
}
