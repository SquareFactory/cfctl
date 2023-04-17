package cmd

import (
	"fmt"
	"os/exec"

	"github.com/urfave/cli/v2"
)

var ipmiCommand = &cli.Command{
	Name:  "ipmi",
	Usage: "Manage compute nodes using ipmi-api",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "action",
			Usage:    "IPMI action to perform (power-on, power-off, cycle, soft, reset, status)",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "host",
			Usage:    "Hostname",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "user",
			Usage:    "IPMI user provided as 'username:password'",
			Required: true,
		},
	},
	Action: func(ctx *cli.Context) error {
		host := ctx.String("host")
		action := ctx.String("action")
		user := ctx.String("user")

		switch action {
		case "power-on":
			requestPath := fmt.Sprintf("10.10.2.162:8080/host/%s/on", host)
			out, err := exec.Command(
				"curl",
				"--request",
				"POST",
				requestPath,
				"--user",
				user,
			).CombinedOutput()
			if err != nil {
				fmt.Printf("%v: %s\n", err, out)
			}

		case "power-off":
			requestPath := fmt.Sprintf("10.10.2.162:8080/host/%s/off", host)
			out, err := exec.Command(
				"curl",
				"--request",
				"POST",
				requestPath,
				"--user",
				user,
			).CombinedOutput()

			if err != nil {
				fmt.Printf("%v: %s\n", err, out)
			}

		case "cycle":
			requestPath := fmt.Sprintf("10.10.2.162:8080/host/%s/cycle", host)
			out, err := exec.Command(
				"curl",
				"--request",
				"POST",
				requestPath,
				"--user",
				user,
			).CombinedOutput()

			if err != nil {
				fmt.Printf("%v: %s\n", err, out)
			}

		case "soft":
			requestPath := fmt.Sprintf("10.10.2.162:8080/host/%s/soft", host)
			out, err := exec.Command(
				"curl",
				"--request",
				"POST",
				requestPath,
				"--user",
				user,
			).CombinedOutput()

			if err != nil {
				fmt.Printf("%v: %s\n", err, out)
			}

		case "reset":
			requestPath := fmt.Sprintf("10.10.2.162:8080/host/%s/reset", host)
			out, err := exec.Command(
				"curl",
				"--request",
				"POST",
				requestPath,
				"--user",
				user,
			).CombinedOutput()

			if err != nil {
				fmt.Printf("%v: %s\n", err, out)
			}

		case "status":
			requestPath := fmt.Sprintf("10.10.2.162:8080/host/%s/status", host)
			out, err := exec.Command(
				"curl",
				requestPath,
				"--user",
				user,
			).CombinedOutput()

			if err != nil {
				fmt.Printf("%v: %s\n", err, out)
			}

		default:
			return cli.Exit("Invalid action specified", 1)
		}

		return nil
	},
}
