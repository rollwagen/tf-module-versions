package cmd

import (
	"github.com/spf13/cobra"
)

var (
	FlagQuiet   bool
	FlagVerbose bool

	rootCmd = &cobra.Command{
		Use:   "tfm",
		Short: "A tool that check and compares used vs. available terraform module versions in git repositories",
		Long: `A tool that validates and compares used vs. available terraform module version
in git repositories, specific modules hosted in Gitlab repositories`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&FlagQuiet, "quiet", "q", false, "DEPRECATED: log output suppressed by default now; use --verbose to enable")
	rootCmd.PersistentFlags().BoolVarP(&FlagVerbose, "verbose", "v", false, "Suppress debug output")
}
