package action

import (
	"github.com/SquareFactory/cfctl/phase"
	"github.com/SquareFactory/cfctl/pkg/apis/cfctl.clusterfactory.io/v1beta1/cluster"
)

type Kubeconfig struct {
	// Manager is the phase manager
	Manager              *phase.Manager
	KubeconfigAPIAddress string

	Kubeconfig string
}

func (k *Kubeconfig) Run() error {
	// Change so that the internal config has only single controller host as we
	// do not need to connect to all nodes
	k.Manager.Config.Spec.Hosts = cluster.Hosts{k.Manager.Config.Spec.K0sLeader()}

	k.Manager.AddPhase(
		&phase.Connect{},
		&phase.DetectOS{},
		&phase.GetKubeconfig{APIAddress: k.KubeconfigAPIAddress},
		&phase.Disconnect{},
	)

	return k.Manager.Run()
}
