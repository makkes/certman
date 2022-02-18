package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"go.e13.dev/certman/cmd/create"
	"go.e13.dev/certman/cmd/version"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:          "certman",
		SilenceUsage: true,
	}

	rootCmd.AddCommand(create.NewCommand())
	rootCmd.AddCommand(version.NewCommand())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
