package cmd

import (
	"generator"
	"github.com/spf13/cobra"
	"specgen/generators"
)

func init() {
	generator.AddCobraCommands(cmdCodegen, generators.All)
	rootCmd.AddCommand(cmdCodegen)
}

var cmdCodegen = &cobra.Command{
	Use:   "codegen",
	Short: "Generate code from spec using specgen code generators",
}
