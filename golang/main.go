package main

import (
	"fmt"
	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/generator/v2/console"
	"github.com/specgen-io/specgen/golang/v2/generators"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "specgen",
		Short: "Code generation based on specification",
	}
	generator.AddCobraCommands(rootCmd, generators.All)
	cobra.OnInitialize()
	console.PrintLn("Running specgen")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
