package version

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"go.e13.dev/certman/pkg/config"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display the version of certman",
		RunE: func(cmd *cobra.Command, args []string) error {
			var version strings.Builder
			fmt.Fprintf(&version, "%s", config.GetVersion())
			fmt.Println(version.String())
			return nil
		},
	}
}
