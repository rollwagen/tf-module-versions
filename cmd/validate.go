package cmd

import (
	"github.com/rollwagen/tf-module-versions/pkg/validater"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "validate",
	Short: "Print module version validation on stdout as logs",
	Run: func(cmd *cobra.Command, args []string) {
		validater.Validate(".", FlagQuiet)
	},
}
