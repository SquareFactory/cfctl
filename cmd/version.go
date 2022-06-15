package cmd

import (
	"fmt"
	"os"

	"github.com/SquareFactory/cfctl/analytics"
	"github.com/SquareFactory/cfctl/integration/github"
	"github.com/SquareFactory/cfctl/version"
	"github.com/urfave/cli/v2"
)

var versionCommand = &cli.Command{
	Name:  "version",
	Usage: "Output cfctl version",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:   "machine-id",
			Hidden: true,
		},
		&cli.BoolFlag{
			Name:  "k0s",
			Usage: "Retrieve the latest k0s version number",
		},
		&cli.BoolFlag{
			Name:  "pre",
			Usage: "When used in conjunction with --k0s, a pre release is accepted as the latest version",
		},
	},
	Before: func(ctx *cli.Context) error {
		if ctx.Bool("k0s") {
			v, err := github.LatestK0sVersion(ctx.Bool("pre"))
			if err != nil {
				return err
			}
			fmt.Println(v)
			os.Exit(0)
		}

		if ctx.Bool("machine-id") {
			id, err := analytics.MachineID()
			if err != nil {
				println(err.Error())
				os.Exit(1)
			}
			fmt.Println(id)
			os.Exit(0)
		}

		return nil
	},
	Action: func(ctx *cli.Context) error {
		fmt.Printf("version: %s\n", version.Version)
		fmt.Printf("commit: %s\n", version.GitCommit)
		return nil
	},
}
