package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/urfave/cli/v2"
)

var secretFileRegex = regexp.MustCompile(`^(.*)-secret.(yml|yaml).local$`)
var kubesealCommand = &cli.Command{
	Name:  "kubeseal",
	Usage: "Kubeseal every '-secret.yaml.local' files recursively",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "controller-namespace",
			Value:   "sealed-secrets",
			Usage:   "The namespace where the sealed secrets controller resides.",
			EnvVars: []string{"SEALED_SECRETS_CONTROLLER_NAMESPACE"},
		},
		&cli.StringFlag{
			Name:  "controller-name",
			Value: "sealed-secrets",
			Usage: "The name of the sealed secrets controller.",
		},
	},
	Action: func(ctx *cli.Context) error {
		err := filepath.Walk(".",
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				if res := secretFileRegex.FindStringSubmatch(path); len(res) == 3 {
					sealedFilePath := fmt.Sprintf("%s-sealed-secret.%s", res[1], res[2])
					if _, err := os.Stat(sealedFilePath); err == nil {
						return nil
					}
					fmt.Printf("Processing %s\n", path)
					_, err := exec.Command(
						"kubeseal",
						"--controller-namespace",
						ctx.String("controller-namespace"),
						"--controller-name",
						ctx.String("controller-name"),
						"--format",
						"yaml",
						"--secret-file",
						path,
						"--sealed-secret-file",
						sealedFilePath,
					).Output()
					if err != nil {
						fmt.Printf("%v\n", err)
					}
				}
				return nil
			})
		if err != nil {
			log.Println(err)
		}
		return nil
	},
}
