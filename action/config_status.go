package action

import (
	"fmt"
	"io"

	"github.com/deepsquare-io/cfctl/analytics"
	"github.com/deepsquare-io/cfctl/pkg/apis/cfctl.clusterfactory.io/v1beta1"

	"github.com/k0sproject/rig/exec"
)

type ConfigStatus struct {
	Config      *v1beta1.Cluster
	Concurrency int
	Format      string
	Writer      io.Writer
}

func (c ConfigStatus) Run() error {
	analytics.Client.Publish("config-status-start", map[string]interface{}{})

	h := c.Config.Spec.K0sLeader()

	if err := h.Connect(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer h.Disconnect()

	if err := h.ResolveConfigurer(); err != nil {
		return err
	}
	format := c.Format
	if format != "" {
		format = "-o " + format
	}

	output, err := h.ExecOutput(
		h.Configurer.KubectlCmdf(
			h,
			h.K0sDataDir(),
			"-n kube-system get event --field-selector involvedObject.name=k0s %s",
			format,
		),
		exec.Sudo(h),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", h, err)
	}
	fmt.Fprintln(c.Writer, output)

	return nil
}
