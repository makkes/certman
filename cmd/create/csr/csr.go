package csr

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

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

	cmd := cobra.Command{
		Use:   "csr",
		Short: "Create a new certificate signing request (CSR)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return flags.Validate(csrOutFlag, privkeyOutFlag)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			csr, err := cert.CreateCSR()
			if err != nil {
				return fmt.Errorf("failed to create CA: %w", err)
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

	return &cmd
}
