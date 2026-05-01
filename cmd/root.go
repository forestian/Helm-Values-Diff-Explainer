package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hvediff",
	Short: "Helm Values Diff Explainer",
	Long:  "Compare two Helm values YAML files and explain changes with operational impact hints.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(diffCmd)
}
