package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/SquareFactory/cfctl/utils/generators"
	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var credential = Credential{}
var address string
var ipmiCommand = &cli.Command{
	Name:        "ipmi",
	ArgsUsage:   "hostnames action",
	Usage:       "Manage compute nodes using ipmi-api",
	Description: "Send action to IPMI API. Available actions: on, off, cycle, status, soft, reset.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "user",
			Destination: &credential.Username,
			Required:    true,
			Usage:       "IPMI user provided",
			EnvVars:     []string{"IPMIUSER"},
		},
		&cli.StringFlag{
			Name:        "password",
			Destination: &credential.Password,
			Required:    true,
			Usage:       "IPMI password",
			EnvVars:     []string{"IPMIPASS"},
		},
		&cli.StringFlag{
			Name:        "address",
			Destination: &address,
			Usage:       "API address",
			Value:       "https://ipmi.internal",
			EnvVars:     []string{"IPMIADDRESS"},
		},
	},
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() != 2 {
			return errors.New("not enough arguments, use --help")
		}

		arg := ctx.Args().Get(0)
		hostnamesRanges := generators.SplitCommaOutsideOfBrackets(arg)

		var hostnames []string
		for _, hostnamesRange := range hostnamesRanges {
			h := generators.ExpandBrackets(hostnamesRange)
			hostnames = append(hostnames, h...)
		}

		action := ctx.Args().Get(1)

		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}

		postData, err := json.Marshal(credential)
		if err != nil {
			return err
		}

		for _, host := range hostnames {

			requestPath := fmt.Sprintf("%s/host/%s/%s", address, host, action)

			resp, err := client.Post(requestPath, "application/json", bytes.NewBuffer(postData))
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				log.WithError(err).Warn("ipmi response body couldn't be read")
			}

			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				log.WithFields(log.Fields{
					"status": resp.StatusCode,
					"body":   string(b),
				}).Error("ipmi API returned non-OK status code")
				return errors.New("ipmi API returned non-OK status code")
			}
			log.Info(string(b))

		}

		return nil
	},
}
