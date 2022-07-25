package java

import (
	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/spec/v2"
	"github.com/specgen-io/specgen/v2/gen/java/client"
	"github.com/specgen-io/specgen/v2/gen/java/models"
	"github.com/specgen-io/specgen/v2/gen/java/service"
)

var JsonlibJavaValues = []string{"jackson", "moshi"}

var Models = generator.Generator{
	"models-java",
	"Java Models",
	"Generate Java models source code",
	[]generator.GeneratorArg{
		{Arg: generator.ArgSpecFile, Required: true},
		{Arg: generator.ArgJsonlib, Required: true, Values: JsonlibJavaValues},
		{Arg: generator.ArgPackageName, Required: false},
		{Arg: generator.ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params generator.GeneratorArgsValues) *generator.Sources {
		return models.Generate(specification, params[generator.ArgJsonlib], params[generator.ArgPackageName], params[generator.ArgGeneratePath])
	},
}

var Client = generator.Generator{
	"client-java",
	"Java Client",
	"Generate Java client source code",
	[]generator.GeneratorArg{
		{Arg: generator.ArgSpecFile, Required: true},
		{Arg: generator.ArgJsonlib, Required: true, Values: JsonlibJavaValues},
		{Arg: generator.ArgPackageName, Required: false},
		{Arg: generator.ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params generator.GeneratorArgsValues) *generator.Sources {
		return client.Generate(specification, params[generator.ArgJsonlib], params[generator.ArgPackageName], params[generator.ArgGeneratePath])
	},
}

var ServerJavaValues = []string{"spring", "micronaut"}

var Service = generator.Generator{
	"service-java",
	"Java Service",
	"Generate Java service source code",
	[]generator.GeneratorArg{
		{Arg: generator.ArgSpecFile, Required: true},
		{Arg: generator.ArgJsonlib, Required: true, Values: JsonlibJavaValues},
		{Arg: generator.ArgServer, Required: true, Values: ServerJavaValues},
		{Arg: generator.ArgPackageName, Required: false},
		{Arg: generator.ArgSwaggerPath, Required: false},
		{Arg: generator.ArgGeneratePath, Required: true},
		{Arg: generator.ArgServicesPath, Required: false},
	},
	func(specification *spec.Spec, params generator.GeneratorArgsValues) *generator.Sources {
		return service.Generate(specification, params[generator.ArgJsonlib], params[generator.ArgServer], params[generator.ArgPackageName], params[generator.ArgSwaggerPath], params[generator.ArgGeneratePath], params[generator.ArgServicesPath])
	},
}
