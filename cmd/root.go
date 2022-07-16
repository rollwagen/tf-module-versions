package cmd

import (
	"github.com/spf13/cobra"
)

var (
	FlagQuiet bool

	rootCmd = &cobra.Command{
		Use:   "tf-module-versions",
		Short: "A tool that check and compares used vs. available terraform module versions in git repositories",
		Long: `A tool that validates and compares used vs. available terraform module version
in git repositories, specific modules hosted in Gitlab repositories`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&FlagQuiet, "quiet", "q", false, "Suppress debug output")
}
