package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tf-modver",
	Short: "A tool that check and compares used vs. available terraform module versions in git repositories",
	Long: `A tool that validates and compares used vs. available terraform module version
in git repositories, specific modules hosted in Gitlab repositories`,
}

func Execute() error {
	return rootCmd.Execute()
}
