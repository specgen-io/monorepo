package generators

import (
	"generator"
	"kotlin/client"
	"kotlin/models"
	"kotlin/service"
	"spec"
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

var ClientKotlinValues = []string{"okhttp", "micronaut", "micronaut-declarative"}

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

var ServerKotlinValues = []string{"spring", "micronaut"}

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

var All = []generator.Generator{
	Models,
	Client,
	Service,
}
