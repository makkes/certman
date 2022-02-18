package create

import (
	"github.com/spf13/cobra"

	"go.e13.dev/certman/cmd/create/ca"
	"go.e13.dev/certman/cmd/create/cert"
	"go.e13.dev/certman/cmd/create/csr"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a resource such as a CA or a CSR",
	}

	cmd.AddCommand(ca.NewCommand())
	cmd.AddCommand(csr.NewCommand())
	cmd.AddCommand(cert.NewCommand())

	return cmd
}
