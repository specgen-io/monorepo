package gen

import (
	"github.com/specgen-io/specgen/v2/gengo"
	"github.com/specgen-io/specgen/v2/genjava"
	"github.com/specgen-io/specgen/v2/genkotlin"
	"github.com/specgen-io/specgen/v2/genruby"
	"github.com/specgen-io/specgen/v2/genscala"
	"github.com/specgen-io/specgen/v2/gents"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

var ModelsGo = Generator{
	"models-go",
	"Generate Go models source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgModuleName, Required: true},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[string]string) *sources.Sources {
		return gengo.GenerateModels(specification, params[ModuleName], params[GeneratePath])
	},
}

var ClientGo = Generator{
	"client-go",
	"Generate Go client source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgModuleName, Required: true},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[string]string) *sources.Sources {
		return gengo.GenerateGoClient(specification, params[ModuleName], params[GeneratePath])
	},
}

var ServiceGo = Generator{
	"service-go",
	"Generate Go service source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgModuleName, Required: true},
		{Arg: ArgSwaggerPath, Required: false},
		{Arg: ArgGeneratePath, Required: true},
		{Arg: ArgServicesPath, Required: false},
	},
	func(specification *spec.Spec, params map[string]string) *sources.Sources {
		return gengo.GenerateService(specification, params[ModuleName], params[SwaggerPath], params[GeneratePath], params[ServicesPath])
	},
}

var ModelsJava = Generator{
	"models-java",
	"Generate Java models source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgPackageName, Required: false, Default: ""},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[string]string) *sources.Sources {
		return genjava.GenerateModels(specification, params[PackageName], params[GeneratePath])
	},
}

var ClientJava = Generator{
	"client-java",
	"Generate Java client source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgPackageName, Required: true},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[string]string) *sources.Sources {
		return genjava.GenerateClient(specification, params[PackageName], params[GeneratePath])
	},
}

var ServiceJava = Generator{
	"service-java",
	"Generate Java service source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgPackageName, Required: false, Default: ""},
		{Arg: ArgSwaggerPath, Required: false, Default: ""},
		{Arg: ArgGeneratePath, Required: true},
		{Arg: ArgServicesPath, Required: false, Default: ""},
	},
	func(specification *spec.Spec, params map[string]string) *sources.Sources {
		return genjava.GenerateService(specification, params[PackageName], params[SwaggerPath], params[GeneratePath], params[ServicesPath])
	},
}

var ModelsKotlin = Generator{
	"models-kotlin",
	"Generate Kotlin models source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgJsonlib, Required: false, Default: "jackson"},
		{Arg: ArgPackageName, Required: false, Default: ""},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[string]string) *sources.Sources {
		return genkotlin.GenerateModels(specification, params[PackageName], params[GeneratePath])
	},
}

var ClientKotlin = Generator{
	"client-kotlin",
	"Generate Kotlin client source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgPackageName, Required: true},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[string]string) *sources.Sources {
		return genkotlin.GenerateClient(specification, params[PackageName], params[GeneratePath])
	},
}

var ModelsRuby = Generator{
	"models-ruby",
	"Generate Ruby models source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[string]string) *sources.Sources {
		return genruby.GenerateModels(specification, params[GeneratePath])
	},
}

var ClientRuby = Generator{
	"client-ruby",
	"Generate Ruby client source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[string]string) *sources.Sources {
		return genruby.GenerateClient(specification, params[GeneratePath])
	},
}

var ModelsScala = Generator{
	"models-scala",
	"Generate Scala models source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[string]string) *sources.Sources {
		return genscala.GenerateCirceModels(specification, "", params[GeneratePath])
	},
}

var ClientScala = Generator{
	"client-scala",
	"Generate Scala client source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[string]string) *sources.Sources {
		return genscala.GenerateSttpClient(specification, "", params[GeneratePath])
	},
}

var ServiceScala = Generator{
	"service-scala",
	"Generate Scala service source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgSwaggerPath, Required: false, Default: ""},
		{Arg: ArgGeneratePath, Required: true},
		{Arg: ArgServicesPath, Required: false, Default: ""},
	},
	func(specification *spec.Spec, params map[string]string) *sources.Sources {
		return genscala.GeneratePlayService(specification, params[SwaggerPath], params[GeneratePath], params[ServicesPath])
	},
}

var ModelsTs = Generator{
	"models-ts",
	"Generate TypeScript models source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgValidation, Required: true},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[string]string) *sources.Sources {
		return gents.GenerateModels(specification, params[Validation], params[GeneratePath])
	},
}

var ClientTs = Generator{
	"client-ts",
	"Generate TypeScript client source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgClient, Required: true},
		{Arg: ArgValidation, Required: true},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[string]string) *sources.Sources {
		return gents.GenerateClient(specification, params[GeneratePath], params[Client], params[Validation])
	},
}

var ServiceTs = Generator{
	"service-ts",
	"Generate TypeScript client source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgServer, Required: true},
		{Arg: ArgValidation, Required: true},
		{Arg: ArgGeneratePath, Required: true},
		{Arg: ArgServicesPath, Required: false},
		{Arg: ArgSwaggerPath, Required: false},
	},
	func(specification *spec.Spec, params map[string]string) *sources.Sources {
		return gents.GenerateService(specification, params[SwaggerPath], params[GeneratePath], params[ServicesPath], params[Server], params[Validation])
	},
}

var Generators = []Generator{
	ModelsGo,
	ClientGo,
	ServiceGo,
	ModelsJava,
	ClientJava,
	ServiceJava,
	ModelsKotlin,
	ClientKotlin,
	ModelsRuby,
	ClientRuby,
	ModelsScala,
	ClientScala,
	ServiceScala,
	ModelsTs,
	ClientTs,
	ServiceTs,
}
