package main

import (
	"fmt"
	"generator"
	"generator/console"
	"github.com/spf13/cobra"
	"golang/generators"
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
