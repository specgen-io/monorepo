package cmd

import (
	"generator"
	"specgen/v2/generators"
)

func init() {
	generator.AddCobraCommands(rootCmd, generators.All)
}
