package cmd

import (
	"fmt"
	"github.com/specgen-io/specgen/generator/v2/console"
	"github.com/specgen-io/specgen/v2/version"
	"github.com/spf13/cobra"
	"os"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:     "specgen",
	Version: version.Current,
	Short:   "Code generation based on specification",
}

func Execute() {
	console.PrintLnF("Running specgen version: %s", version.Current)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}
