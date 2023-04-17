package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

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
	Name:      "ipmi",
	ArgsUsage: "host action",
	Usage:     "Manage compute nodes using ipmi-api",
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

		host := ctx.Args().Get(0)
		action := ctx.Args().Get(1)

		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}

		requestPath := fmt.Sprintf("%s/host/%s/%s", address, host, action)

		postData, err := json.Marshal(credential)
		if err != nil {
			return err
		}

		resp, err := client.Post(requestPath, "application/json", bytes.NewBuffer(postData))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.WithError(err).Warn("ipmi response body couldn't be read")
		}

		switch {
		case resp.StatusCode == 404:
			return fmt.Errorf("action '%s' not found", action)
		case resp.StatusCode < 200 || resp.StatusCode >= 300:
			log.WithFields(log.Fields{
				"status": resp.StatusCode,
				"body":   string(b),
			}).Error("ipmi API returned non-OK status code")
			return errors.New("ipmi API returned non-OK status code")
		}
		log.Info(string(b))

		return nil
	},
}
