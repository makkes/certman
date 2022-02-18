package cert

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"go.e13.dev/certman/pkg/ca"
	"go.e13.dev/certman/pkg/cert"
	"go.e13.dev/certman/pkg/flags"
)

func NewCommand() *cobra.Command {
	csrFlag := flags.Flag{
		Name: "csr",
	}
	csrCfgFlag := flags.Flag{
		Name: "csr-config",
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
	privkeyOutFlag := flags.Flag{
		Name: "privkey-out",
	}

	var csrConfig cert.CSRConfig

	cmd := cobra.Command{
		Use:   "cert",
		Short: "Create a certificate from a CSR signed by a CA",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := flags.Validate(
				flags.Exclusive(csrFlag, csrCfgFlag),
				flags.And(csrCfgFlag, privkeyOutFlag),
				caKeyFlag,
				caCertFlag,
				outFlag); err != nil {
				return err
			}
			if csrCfgFlag.Val != "" {
				cfgData, err := os.ReadFile(csrCfgFlag.Val)
				if err != nil {
					return fmt.Errorf("failed to read configuration file '%s': %w", csrCfgFlag.Val, err)
				}
				if err := yaml.Unmarshal(cfgData, &csrConfig); err != nil {
					return fmt.Errorf("failed to parse configuration file '%s': %w", csrCfgFlag.Val, err)
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var csrPEM []byte
			if csrCfgFlag.Val != "" {
				csr, err := cert.CreateCSR(&csrConfig)
				if err != nil {
					return fmt.Errorf("failed to create CSR: %w", err)
				}
				if err := os.WriteFile(privkeyOutFlag.Val, csr.PrivateKey, 0o600); err != nil {
					return fmt.Errorf("failed to write the private key file: %w", err)
				}
				csrPEM = csr.Request
			} else {
				var err error
				csrPEM, err = os.ReadFile(csrFlag.Val)
				if err != nil {
					return fmt.Errorf("failed to read CSR file '%s': %w", csrFlag.Val, err)
				}
			}
			caKey, err := os.ReadFile(caKeyFlag.Val)
			if err != nil {
				return fmt.Errorf("failed to read CA key file '%s': %w", caKeyFlag.Val, err)
			}
			caCert, err := os.ReadFile(caCertFlag.Val)
			if err != nil {
				return fmt.Errorf("failed to read CA key file '%s': %w", caKeyFlag.Val, err)
			}

			cert, err := ca.CreateCertificate(csrPEM, caCert, caKey)
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
	cmd.PersistentFlags().StringVar(
		&csrCfgFlag.Val, csrCfgFlag.Name, "", "filename of the input configuration for the CSR")
	cmd.PersistentFlags().StringVar(
		&privkeyOutFlag.Val, privkeyOutFlag.Name, "", "filename to use for storing the private key")
	cmd.PersistentFlags().StringVar(&caKeyFlag.Val, caKeyFlag.Name, "", "filename of the CA's private key")
	cmd.PersistentFlags().StringVar(&caCertFlag.Val, caCertFlag.Name, "", "filename of the CA's certificate")
	cmd.PersistentFlags().StringVar(&outFlag.Val, outFlag.Name, "", "filename to use for storing the certificate")

	return &cmd
}
