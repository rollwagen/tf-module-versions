package cmd

import (
	"github.com/rollwagen/tf-module-versions/pkg/validater"
	"github.com/spf13/cobra"
)

var outputFormat string

var versionCmd = &cobra.Command{
	Use:   "validate",
	Short: "Print module version validation on stdout as logs",
	Run: func(cmd *cobra.Command, args []string) {
		validater.Validate(".", outputFormat, FlagVerbose)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Select output format (table or json)")
}
