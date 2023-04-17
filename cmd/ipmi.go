package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var user, address string
var ipmiCommand = &cli.Command{
	Name:      "ipmi",
	ArgsUsage: "host action",
	Usage:     "Manage compute nodes using ipmi-api",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "user",
			Destination: &user,
			Usage:       "IPMI user provided as 'username:password'",
			EnvVars:     []string{"IPMIUSER"},
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

		if _, ipmiServer := os.LookupEnv("IPMIADDRESS"); !ipmiServer {
			return errors.New("ipmi server not defined, use --help")
		}
		requestPath := fmt.Sprintf("%s/host/%s/%s", os.Getenv("IPMIADDRESS"), host, action)

		if _, ipmiUser := os.LookupEnv("IPMIUSER"); !ipmiUser {
			return errors.New("ipmi user not defined, use --help")
		}
		userCreds := strings.Split(user, ":")
		credential := Credential{
			Username: userCreds[0],
			Password: userCreds[1],
		}

		postData, _ := json.Marshal(credential)

		resp, err := http.Post(requestPath, "application/json", bytes.NewBuffer(postData))
		if err != nil {
			return errors.New("bad request")
		}
		fmt.Printf("Status code : %d\n", resp.StatusCode)

		return nil
	},
}
