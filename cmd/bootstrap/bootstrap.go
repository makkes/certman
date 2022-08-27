package bootstrap

import (
	"github.com/spf13/cobra"

	"go.e13.dev/certman/cmd/bootstrap/ca"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstrap a resource, such as a CA",
	}

	cmd.AddCommand(ca.NewCommand())

	return cmd
}
