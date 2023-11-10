package cmd

import (
	"github.com/SquareFactory/cfctl/action"
	"github.com/SquareFactory/cfctl/pkg/apis/cfctl.clusterfactory.io/v1beta1"

	"github.com/urfave/cli/v2"
)

var configStatusCommand = &cli.Command{
	Name:  "status",
	Usage: "Show k0s dynamic config reconciliation events",
	Flags: []cli.Flag{
		configFlag,
		debugFlag,
		traceFlag,
		redactFlag,
		analyticsFlag,
		upgradeCheckFlag,
		&cli.StringFlag{
			Name:    "output",
			Usage:   "kubectl output formatting",
			Aliases: []string{"o"},
		},
	},
	Before: actions(initLogging, startCheckUpgrade, initConfig, initAnalytics),
	After:  actions(reportCheckUpgrade, closeAnalytics),
	Action: func(ctx *cli.Context) error {
		configStatusAction := action.ConfigStatus{
			Config: ctx.Context.Value(ctxConfigKey{}).(*v1beta1.Cluster),
			Format: ctx.String("output"),
			Writer: ctx.App.Writer,
		}

		return configStatusAction.Run()
	},
}
