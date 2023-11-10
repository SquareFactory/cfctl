package cmd

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var secretFileRegex = regexp.MustCompile(`^(.*)-secret.(yml|yaml).local$`)
var kubesealCommand = &cli.Command{
	Name:  "kubeseal",
	Usage: "Kubeseal every '-secret.yaml.local' files recursively",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "cert",
			Value:   "kubeseal.crt",
			Usage:   "The name of the sealed secrets certificate, used for encryption.",
			EnvVars: []string{"SEALED_SECRETS_CERTIFICATE"},
		},
		&cli.StringFlag{
			Name:    "controller-namespace",
			Value:   "sealed-secrets",
			Usage:   "The namespace where the sealed secrets controller resides (not needed if certificate present).",
			EnvVars: []string{"SEALED_SECRETS_CONTROLLER_NAMESPACE"},
		},
		&cli.StringFlag{
			Name:    "controller-name",
			Value:   "sealed-secrets",
			Usage:   "The name of the sealed secrets controller (not needed if certificate present).",
			EnvVars: []string{"SEALED_SECRETS_CONTROLLER_NAME"},
		},
	},
	Action: func(ctx *cli.Context) error {
		cert := ctx.String("cert")
		if _, err := os.Stat(cert); errors.Is(err, os.ErrNotExist) {
			if err := fetchCertificate(ctx.Context, cert, ctx.String("controller-name"), ctx.String("controller-namespace")); err != nil {
				return err
			}
		}

		// Check certificate expiration
		certData, err := os.ReadFile(cert)
		if err != nil {
			return err
		}

		block, _ := pem.Decode(certData)
		if block == nil {
			logrus.Error("failed to decode certificate")
			return errors.New("failed to decode certificate")
		}

		certParsed, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return err
		}

		// Check if the certificate has expired
		if time.Now().After(certParsed.NotAfter) {
			logrus.Warn("Certificate has expired. Fetching new certificate.")
			if err := fetchCertificate(ctx.Context, cert, ctx.String("controller-name"), ctx.String("controller-namespace")); err != nil {
				return err
			}
		} else {
			logrus.WithField("expirationDate", certParsed.NotAfter).Info("Certificate has not expired")
		}

		logrus.Info("Processing...")
		if err := filepath.Walk(".",
			func(path string, info os.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}
				if res := secretFileRegex.FindStringSubmatch(path); len(res) == 3 {
					sealedFilePath := fmt.Sprintf("%s-sealed-secret.%s", res[1], res[2])
					if _, err := os.Stat(sealedFilePath); err == nil {
						return nil
					}
					out, err := exec.Command(
						"kubeseal",
						"--cert",
						cert,
						"--format",
						"yaml",
						"--secret-file",
						path,
						"--sealed-secret-file",
						sealedFilePath,
					).CombinedOutput()
					if err != nil {
						logrus.Errorf("%v: %s\n", err, out)
					}
					logrus.Printf("Sealed at %s\n", sealedFilePath)
				}
				return nil
			},
		); err != nil {
			logrus.Println(err)
		}
		return nil
	},
}

func fetchCertificate(
	ctx context.Context,
	cert string,
	ctrlName string,
	ctrlNamespace string,
) error {
	logrus.WithField("path", cert).Warn(
		"sealed secret certificate not found, trying to fetching it",
	)
	// Fetch certificate
	data, err := exec.CommandContext(
		ctx,
		"kubeseal",
		"--controller-namespace",
		ctrlNamespace,
		"--controller-name",
		ctrlName,
		"--fetch-cert",
	).CombinedOutput()
	if err != nil {
		fmt.Printf("%v: %s\n", err, cert)
		return err
	}
	return os.WriteFile(cert, data, os.ModePerm)
}
