package main

import (
	"fmt"
	"generator"
	"generator/console"
	"github.com/spf13/cobra"
	"os"
	"ruby/generators"
	"ruby/version"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "specgen",
		Version: version.Current,
		Short:   "Code generation based on specification",
	}
	generator.AddCobraCommands(rootCmd, generators.All)
	cobra.OnInitialize()
	console.PrintLnF("Running specgen ruby, version: %s", version.Current)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
