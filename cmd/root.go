package cmd

import (
	"github.com/urfave/cli/v2"
)

// App is the main urfave/cli.App for cfctl
var App = &cli.App{
	Name:  "cfctl",
	Usage: "k0s cluster management tool",
	Flags: []cli.Flag{
		debugFlag,
		traceFlag,
		redactFlag,
	},
	Commands: []*cli.Command{
		versionCommand,
		applyCommand,
		kubeconfigCommand,
		initCommand,
		resetCommand,
		backupCommand,
		{
			Name:  "config",
			Usage: "Configuration related sub-commands",
			Subcommands: []*cli.Command{
				configEditCommand,
				configStatusCommand,
			},
		},
		kubesealCommand,
		completionCommand,
	},
	EnableBashCompletion: true,
}
