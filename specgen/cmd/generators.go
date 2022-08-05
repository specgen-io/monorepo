package cmd

import (
	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/v2/generators"
)

func init() {
	generator.AddCobraCommands(rootCmd, generators.All)
}
