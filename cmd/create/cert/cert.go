package cert

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"go.e13.dev/certman/pkg/ca"
	"go.e13.dev/certman/pkg/flags"
)

func NewCommand() *cobra.Command {
	csrFlag := flags.Flag{
		Name:      "csr",
		Validator: flags.DisallowEmpty{},
	}
	caKeyFlag := flags.Flag{
		Name:      "ca-key",
		Validator: flags.DisallowEmpty{},
	}
	caCertFlag := flags.Flag{
		Name:      "ca-cert",
		Validator: flags.DisallowEmpty{},
	}
	outFlag := flags.Flag{
		Name:      "out",
		Validator: flags.DisallowEmpty{},
	}

	cmd := cobra.Command{
		Use:   "cert",
		Short: "Create a certificate from a CSR signed by a CA",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return flags.Validate(csrFlag, caKeyFlag, caCertFlag, outFlag)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			csr, err := os.ReadFile(csrFlag.Val)
			if err != nil {
				return fmt.Errorf("failed to read CSR file '%s': %w", csrFlag.Val, err)
			}
			caKey, err := os.ReadFile(caKeyFlag.Val)
			if err != nil {
				return fmt.Errorf("failed to read CA key file '%s': %w", caKeyFlag.Val, err)
			}
			caCert, err := os.ReadFile(caCertFlag.Val)
			if err != nil {
				return fmt.Errorf("failed to read CA key file '%s': %w", caKeyFlag.Val, err)
			}

			cert, err := ca.CreateCertificate(csr, caCert, caKey)
			if err != nil {
				return fmt.Errorf("failed to create certificate: %w", err)
			}

			err = os.WriteFile(outFlag.Val, cert, 0o644) //nolint:gosec // doesn't contain sensitive data
			if err != nil {
				return fmt.Errorf("failed to write the certificate file: %w", err)
			}

			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&csrFlag.Val, csrFlag.Name, "", "filename of the CSR")
	cmd.PersistentFlags().StringVar(&caKeyFlag.Val, caKeyFlag.Name, "", "filename of the CA's private key")
	cmd.PersistentFlags().StringVar(&caCertFlag.Val, caCertFlag.Name, "", "filename of the CA's certificate")
	cmd.PersistentFlags().StringVar(&outFlag.Val, outFlag.Name, "", "filename to use for storing the certificate")

	return &cmd
}
