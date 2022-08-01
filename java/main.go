package main

import (
	"fmt"
	"github.com/specgen-io/specgen/console/v2"
	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/java/v2/generators"
	"github.com/specgen-io/specgen/java/v2/version"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "specgen",
		Version: version.Current,
		Short:   "Code generation based on specification",
	}
	generator.AddCobraCommands(rootCmd, generators.All)
	cobra.OnInitialize()
	console.PrintLnF("Running specgen java, version: %s", version.Current)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
