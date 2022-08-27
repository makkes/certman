package ca

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"go.e13.dev/certman/pkg/ca"
	"go.e13.dev/certman/pkg/flags"
)

func NewCommand() *cobra.Command {
	outFlag := flags.Flag{
		Name:      "out",
		Shorthand: "o",
		Validator: flags.DisallowEmpty{},
	}

	cmd := cobra.Command{
		Use:   "ca",
		Short: "Bootstrap the configuration of a CA",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := flags.Validate(outFlag); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := ca.CAConfig{
				APIVersion:   ca.APIVersionV1,
				CommonName:   "My awesome CA",
				Organization: "Me",
				Country:      "US",
				Province:     "California",
				Locality:     "San Francisco",
				NotBefore:    time.Now(),
				NotAfter:     time.Now().AddDate(10, 0, 0),
			}

			outData, err := yaml.Marshal(cfg)
			if err != nil {
				return fmt.Errorf("failed to marshal bootstrap config: %w", err)
			}

			err = os.WriteFile(outFlag.Val, outData, 0o644) //nolint:gosec // doesn't contain sensitive data
			if err != nil {
				return fmt.Errorf("failed to write the bootstrap configuration file: %w", err)
			}

			fmt.Fprintf(os.Stderr, "bootstrap CA configuration written to %s\n", outFlag.Val)

			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&outFlag.Val, outFlag.Name, outFlag.Shorthand, "ca.yaml",
		"filename to use for storing the certificate.")

	return &cmd
}
