package cmd

import (
	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/golang/v2"
	"github.com/specgen-io/specgen/java/v2"
	"github.com/specgen-io/specgen/kotlin/v2"
	"github.com/specgen-io/specgen/openapi/v2"
	"github.com/specgen-io/specgen/ruby/v2"
	"github.com/specgen-io/specgen/scala/v2"
	"github.com/specgen-io/specgen/typescript/v2"
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
