package kotlin

import (
	"github.com/specgen-io/specgen/v2/gen/kotlin/client"
	"github.com/specgen-io/specgen/v2/gen/kotlin/models"
	"github.com/specgen-io/specgen/v2/gen/kotlin/service"
	"github.com/specgen-io/specgen/v2/generator"
	"github.com/specgen-io/specgen/v2/spec"
)

var JsonlibKotlinValues = []string{"jackson", "moshi"}

var Models = generator.Generator{
	"models-kotlin",
	"Kotlin Models",
	"Generate Kotlin models source code",
	[]generator.GeneratorArg{
		{Arg: generator.ArgSpecFile, Required: true},
		{Arg: generator.ArgJsonlib, Required: false, Values: JsonlibKotlinValues},
		{Arg: generator.ArgPackageName, Required: false},
		{Arg: generator.ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params generator.GeneratorArgsValues) *generator.Sources {
		return models.Generate(specification, params[generator.ArgJsonlib], params[generator.ArgPackageName], params[generator.ArgGeneratePath])
	},
}

var ClientKotlinValues = []string{"okhttp", "micronaut-declarative", "micronaut-low-level"}

var Client = generator.Generator{
	"client-kotlin",
	"Kotlin Client",
	"Generate Kotlin client source code",
	[]generator.GeneratorArg{
		{Arg: generator.ArgSpecFile, Required: true},
		{Arg: generator.ArgJsonlib, Required: false, Values: JsonlibKotlinValues},
		{Arg: generator.ArgClient, Required: true, Values: ClientKotlinValues},
		{Arg: generator.ArgPackageName, Required: false},
		{Arg: generator.ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params generator.GeneratorArgsValues) *generator.Sources {
		return client.Generate(specification, params[generator.ArgJsonlib], params[generator.ArgClient], params[generator.ArgPackageName], params[generator.ArgGeneratePath])
	},
}

var ServerKotlinValues = []string{"micronaut"}

var Service = generator.Generator{
	"service-kotlin",
	"Kotlin Service",
	"Generate Kotlin service source code",
	[]generator.GeneratorArg{
		{Arg: generator.ArgSpecFile, Required: true},
		{Arg: generator.ArgJsonlib, Required: true, Values: JsonlibKotlinValues},
		{Arg: generator.ArgServer, Required: true, Values: ServerKotlinValues},
		{Arg: generator.ArgPackageName, Required: false},
		{Arg: generator.ArgSwaggerPath, Required: false},
		{Arg: generator.ArgGeneratePath, Required: true},
		{Arg: generator.ArgServicesPath, Required: false},
	},
	func(specification *spec.Spec, params generator.GeneratorArgsValues) *generator.Sources {
		return service.Generate(specification, params[generator.ArgJsonlib], params[generator.ArgServer], params[generator.ArgPackageName], params[generator.ArgSwaggerPath], params[generator.ArgGeneratePath], params[generator.ArgServicesPath])
	},
}
