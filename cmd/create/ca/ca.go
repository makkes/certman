package ca

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"go.e13.dev/certman/pkg/ca"
	"go.e13.dev/certman/pkg/flags"
)

func NewCommand() *cobra.Command {
	certOutFlag := flags.Flag{
		Name:      "cert-out",
		Validator: flags.DisallowEmpty{},
	}
	keyOutFlag := flags.Flag{
		Name:      "key-out",
		Validator: flags.DisallowEmpty{},
	}
	cfgFlag := flags.Flag{
		Name:      "config",
		Validator: flags.DisallowEmpty{},
	}
	var cfg ca.CAConfig

	cmd := cobra.Command{
		Use:   "ca",
		Short: "Create a new certificate authoritiy",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := flags.Validate(certOutFlag, keyOutFlag, cfgFlag); err != nil {
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
			ca, err := ca.CreateCA(&cfg)
			if err != nil {
				return fmt.Errorf("failed to create CA: %w", err)
			}

			err = os.WriteFile(certOutFlag.Val, ca.Certificate, 0o644) //nolint:gosec // doesn't contain sensitive data
			if err != nil {
				return fmt.Errorf("failed to write the certificate file: %w", err)
			}
			if err := os.WriteFile(keyOutFlag.Val, ca.PrivateKey, 0o600); err != nil {
				return fmt.Errorf("failed to write the private key file: %w", err)
			}

			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&certOutFlag.Val, certOutFlag.Name, "", "filename to use for storing the certificate")
	cmd.PersistentFlags().StringVar(&keyOutFlag.Val, keyOutFlag.Name, "", "filename to use for storing the private key")
	cmd.PersistentFlags().StringVar(&cfgFlag.Val, cfgFlag.Name, "", "filename of the input configuration for the CA")

	return &cmd
}
