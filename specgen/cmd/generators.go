package cmd

import (
	"generator"
	"specgen/generators"
)

func init() {
	generator.AddCobraCommands(rootCmd, generators.All)
}
