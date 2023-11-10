package cmd

import (
	"fmt"

	"github.com/SquareFactory/cfctl/action"
	"github.com/SquareFactory/cfctl/phase"
	"github.com/urfave/cli/v2"
)

var backupCommand = &cli.Command{
	Name:  "backup",
	Usage: "Take backup of existing clusters state",
	Flags: []cli.Flag{
		configFlag,
		dryRunFlag,
		concurrencyFlag,
		debugFlag,
		traceFlag,
		redactFlag,
		retryIntervalFlag,
		retryTimeoutFlag,
		analyticsFlag,
		upgradeCheckFlag,
	},
	Before: actions(
		initLogging,
		startCheckUpgrade,
		initConfig,
		initManager,
		displayLogo,
		initAnalytics,
		displayCopyright,
	),
	After: actions(reportCheckUpgrade, closeAnalytics),
	Action: func(ctx *cli.Context) error {
		backupAction := action.Backup{
			Manager: ctx.Context.Value(ctxManagerKey{}).(*phase.Manager),
		}

		if err := backupAction.Run(); err != nil {
			return fmt.Errorf(
				"backup failed - log file saved to %s: %w",
				ctx.Context.Value(ctxLogFileKey{}).(string),
				err,
			)
		}

		return nil
	},
}
