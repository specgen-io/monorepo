package cmd

import (
	"fmt"
	"os"
	"specgen/version"

	"github.com/spf13/cobra"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:     "specgen",
	Version: version.Current,
	Short:   "Code generation based of service specification",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}
