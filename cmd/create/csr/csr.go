package csr

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"go.e13.dev/certman/pkg/cert"
	"go.e13.dev/certman/pkg/flags"
)

func NewCommand() *cobra.Command {
	csrOutFlag := flags.Flag{
		Name:      "out",
		Validator: flags.DisallowEmpty{},
	}
	privkeyOutFlag := flags.Flag{
		Name:      "privkey-out",
		Validator: flags.DisallowEmpty{},
	}
	cfgFlag := flags.Flag{
		Name:      "config",
		Validator: flags.DisallowEmpty{},
	}
	var cfg cert.CSRConfig

	cmd := cobra.Command{
		Use:   "csr",
		Short: "Create a new certificate signing request (CSR)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := flags.Validate(csrOutFlag, privkeyOutFlag, cfgFlag); err != nil {
				return err
			}

			cfgData, err := os.ReadFile(cfgFlag.Val)
			if err != nil {
				return fmt.Errorf("failed to read configuration file '%s': %w", cfgFlag.Val, err)
			}
			if err := yaml.Unmarshal(cfgData, &cfg); err != nil {
				return fmt.Errorf("failed to parse configuration file '%s': %w", cfgFlag.Val, err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			csr, err := cert.CreateCSR(&cfg)
			if err != nil {
				return fmt.Errorf("failed to create CSR: %w", err)
			}

			err = os.WriteFile(csrOutFlag.Val, csr.Request, 0o644) //nolint:gosec // doesn't contain sensitive data
			if err != nil {
				return fmt.Errorf("failed to write the CSR file: %w", err)
			}
			if err := os.WriteFile(privkeyOutFlag.Val, csr.PrivateKey, 0o600); err != nil {
				return fmt.Errorf("failed to write the private key file: %w", err)
			}

			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&csrOutFlag.Val, csrOutFlag.Name, "", "filename to use for storing the CSR")
	cmd.PersistentFlags().StringVar(
		&privkeyOutFlag.Val, privkeyOutFlag.Name, "", "filename to use for storing the private key")
	cmd.PersistentFlags().StringVar(&cfgFlag.Val, cfgFlag.Name, "", "filename of the input configuration for the CSR")

	return &cmd
}
