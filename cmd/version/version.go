package version

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"go.e13.dev/certman/pkg/config"
	"go.e13.dev/certman/pkg/flags"
)

const (
	outJSON  string = "json"
	outPlain string = "plain"
)

func NewCommand() *cobra.Command {
	outFlag := flags.Flag{
		Name:      "out",
		Shorthand: "o",
		Validator: flags.AllowSet("plain", "json"),
	}
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display the version of certman",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := flags.Validate(outFlag); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var out strings.Builder
			v := config.GetVersion()
			switch outFlag.Val {
			case outJSON:
				fmt.Fprint(&out, v.String())
			case outPlain:
				fmt.Fprintf(&out, "%s.%s.%s", v.Major, v.Minor, v.Patch)
			default:
				// this should be caught by validation and never happen
			}
			fmt.Println(out.String())
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&outFlag.Val, outFlag.Name, outFlag.Shorthand, "plain",
		"output format. One of [plain, json]")

	return cmd
}
