package ruby

import (
	"github.com/specgen-io/specgen/spec/v2"
	"github.com/specgen-io/specgen/specgen/v2/generator"
)

var Models = generator.Generator{
	"models-ruby",
	"Ruby Models",
	"Generate Ruby models source code",
	[]generator.GeneratorArg{
		{Arg: generator.ArgSpecFile, Required: true},
		{Arg: generator.ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params generator.GeneratorArgsValues) *generator.Sources {
		return GenerateModels(specification, params[generator.ArgGeneratePath])
	},
}

var Client = generator.Generator{
	"client-ruby",
	"Ruby Client",
	"Generate Ruby client source code",
	[]generator.GeneratorArg{
		{Arg: generator.ArgSpecFile, Required: true},
		{Arg: generator.ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params generator.GeneratorArgsValues) *generator.Sources {
		return GenerateClient(specification, params[generator.ArgGeneratePath])
	},
}
