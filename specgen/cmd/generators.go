package cmd

import (
	"github.com/specgen-io/specgen/generator/v2"
	golang "github.com/specgen-io/specgen/golang/v2/generators"
	java "github.com/specgen-io/specgen/java/v2/generators"
	kotlin "github.com/specgen-io/specgen/kotlin/v2/generators"
	"github.com/specgen-io/specgen/openapi/v2"
	ruby "github.com/specgen-io/specgen/ruby/v2/generators"
	scala "github.com/specgen-io/specgen/scala/v2/generators"
	typescript "github.com/specgen-io/specgen/typescript/v2/generators"
)

func init() {
	var generators = []generator.Generator{
		golang.Models,
		golang.Client,
		golang.Service,
		java.Models,
		java.Client,
		java.Service,
		kotlin.Models,
		kotlin.Client,
		kotlin.Service,
		ruby.Models,
		ruby.Client,
		scala.Models,
		scala.Client,
		scala.Service,
		typescript.Models,
		typescript.Client,
		typescript.Service,
		openapi.Openapi,
	}

	generator.AddCobraCommands(rootCmd, generators)
}
