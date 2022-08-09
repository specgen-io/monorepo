package openapi

import (
	"generator"
	"spec"
)

var Openapi = generator.Generator{
	"openapi",
	"OpenAPI v3",
	"Generate OpenAPI specification",
	[]generator.GeneratorArg{
		{Arg: generator.ArgSpecFile, Required: true},
		{Arg: generator.ArgOutFile, Required: true},
	},
	func(specification *spec.Spec, params generator.GeneratorArgsValues) *generator.Sources {
		openapiFile := GenerateOpenapi(specification, params[generator.ArgOutFile])
		sources := generator.NewSources()
		sources.AddGenerated(openapiFile)
		return sources
	},
}
