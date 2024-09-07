package cmd

import (
	"fmt"

	"github.com/dtrejod/airgradient-exporter/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Run:   versionRun,
	Short: "Print the version",
}

func versionRun(_ *cobra.Command, _ []string) {
	fmt.Printf("version: %s\n", version.Version())
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
