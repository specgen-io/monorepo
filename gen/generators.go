package gen

import (
	"github.com/specgen-io/specgen/v2/gengo"
	"github.com/specgen-io/specgen/v2/genjava"
	"github.com/specgen-io/specgen/v2/genkotlin"
	"github.com/specgen-io/specgen/v2/genopenapi"
	"github.com/specgen-io/specgen/v2/genruby"
	"github.com/specgen-io/specgen/v2/genscala"
	"github.com/specgen-io/specgen/v2/gents"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

var Openapi = Generator{
	"openapi",
	"OpenAPI v3",
	"Generate OpenAPI specification",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgOutFile, Required: true},
	},
	func(specification *spec.Spec, params map[Arg]string) *sources.Sources {
		openapiFile := genopenapi.GenerateOpenapi(specification, params[ArgOutFile])
		sources := sources.NewSources()
		sources.AddGenerated(openapiFile)
		return sources
	},
}

var ModelsGo = Generator{
	"models-go",
	"Go Models",
	"Generate Go models source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgModuleName, Required: true},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[Arg]string) *sources.Sources {
		return gengo.GenerateModels(specification, params[ArgModuleName], params[ArgGeneratePath])
	},
}

var ClientGo = Generator{
	"client-go",
	"Go Client",
	"Generate Go client source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgModuleName, Required: true},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[Arg]string) *sources.Sources {
		return gengo.GenerateClient(specification, params[ArgModuleName], params[ArgGeneratePath])
	},
}

var ServiceGo = Generator{
	"service-go",
	"Go Service",
	"Generate Go service source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgModuleName, Required: true},
		{Arg: ArgSwaggerPath, Required: false},
		{Arg: ArgGeneratePath, Required: true},
		{Arg: ArgServicesPath, Required: false},
	},
	func(specification *spec.Spec, params map[Arg]string) *sources.Sources {
		return gengo.GenerateService(specification, params[ArgModuleName], params[ArgSwaggerPath], params[ArgGeneratePath], params[ArgServicesPath])
	},
}

var JsonlibJavaValues = []string{"jackson", "moshi"}

var ModelsJava = Generator{
	"models-java",
	"Java Models",
	"Generate Java models source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgJsonlib, Required: true, Values: JsonlibJavaValues},
		{Arg: ArgPackageName, Required: false},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[Arg]string) *sources.Sources {
		return genjava.GenerateModels(specification, params[ArgJsonlib], params[ArgPackageName], params[ArgGeneratePath])
	},
}

var ClientJava = Generator{
	"client-java",
	"Java Client",
	"Generate Java client source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgJsonlib, Required: true, Values: JsonlibJavaValues},
		{Arg: ArgPackageName, Required: false},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[Arg]string) *sources.Sources {
		return genjava.GenerateClient(specification, params[ArgJsonlib], params[ArgPackageName], params[ArgGeneratePath])
	},
}

var ServiceJava = Generator{
	"service-java",
	"Java Service",
	"Generate Java service source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgJsonlib, Required: true, Values: JsonlibJavaValues},
		{Arg: ArgPackageName, Required: false},
		{Arg: ArgSwaggerPath, Required: false},
		{Arg: ArgGeneratePath, Required: true},
		{Arg: ArgServicesPath, Required: false},
	},
	func(specification *spec.Spec, params map[Arg]string) *sources.Sources {
		return genjava.GenerateService(specification, params[ArgJsonlib], params[ArgPackageName], params[ArgSwaggerPath], params[ArgGeneratePath], params[ArgServicesPath])
	},
}

var JsonlibKotlinValues = []string{"jackson"}

var ModelsKotlin = Generator{
	"models-kotlin",
	"Kotlin Models",
	"Generate Kotlin models source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgJsonlib, Required: false, Values: JsonlibKotlinValues},
		{Arg: ArgPackageName, Required: false},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[Arg]string) *sources.Sources {
		return genkotlin.GenerateModels(specification, params[ArgJsonlib], params[ArgPackageName], params[ArgGeneratePath])
	},
}

var ClientKotlin = Generator{
	"client-kotlin",
	"Kotlin Client",
	"Generate Kotlin client source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgJsonlib, Required: false, Values: JsonlibKotlinValues},
		{Arg: ArgPackageName, Required: false},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[Arg]string) *sources.Sources {
		return genkotlin.GenerateClient(specification, params[ArgJsonlib], params[ArgPackageName], params[ArgGeneratePath])
	},
}

var ModelsRuby = Generator{
	"models-ruby",
	"Ruby Models",
	"Generate Ruby models source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[Arg]string) *sources.Sources {
		return genruby.GenerateModels(specification, params[ArgGeneratePath])
	},
}

var ClientRuby = Generator{
	"client-ruby",
	"Ruby Client",
	"Generate Ruby client source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[Arg]string) *sources.Sources {
		return genruby.GenerateClient(specification, params[ArgGeneratePath])
	},
}

var ModelsScala = Generator{
	"models-scala",
	"Scala Models",
	"Generate Scala models source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[Arg]string) *sources.Sources {
		return genscala.GenerateCirceModels(specification, "", params[ArgGeneratePath])
	},
}

var ClientScala = Generator{
	"client-scala",
	"Scala Client",
	"Generate Scala client source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[Arg]string) *sources.Sources {
		return genscala.GenerateSttpClient(specification, "", params[ArgGeneratePath])
	},
}

var ServiceScala = Generator{
	"service-scala",
	"Scala Service",
	"Generate Scala service source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgSwaggerPath, Required: false},
		{Arg: ArgGeneratePath, Required: true},
		{Arg: ArgServicesPath, Required: false},
	},
	func(specification *spec.Spec, params map[Arg]string) *sources.Sources {
		return genscala.GeneratePlayService(specification, params[ArgSwaggerPath], params[ArgGeneratePath], params[ArgServicesPath])
	},
}

var ValidationTsValues = []string{"superstruct", "io-ts"}
var ClientTsValues = []string{"axios", "node-fetch", "browser-fetch"}
var ServerTsValues = []string{"express", "koa"}

var ModelsTs = Generator{
	"models-ts",
	"TypeScript Models",
	"Generate TypeScript models source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgValidation, Required: true, Values: ValidationTsValues},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[Arg]string) *sources.Sources {
		return gents.GenerateModels(specification, params[ArgValidation], params[ArgGeneratePath])
	},
}

var ClientTs = Generator{
	"client-ts",
	"TypeScript Client",
	"Generate TypeScript client source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgClient, Required: true, Values: ClientTsValues},
		{Arg: ArgValidation, Required: true, Values: ValidationTsValues},
		{Arg: ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params map[Arg]string) *sources.Sources {
		return gents.GenerateClient(specification, params[ArgGeneratePath], params[ArgClient], params[ArgValidation])
	},
}

var ServiceTs = Generator{
	"service-ts",
	"TypeScript Service",
	"Generate TypeScript client source code",
	[]GeneratorArg{
		{Arg: ArgSpecFile, Required: true},
		{Arg: ArgServer, Required: true, Values: ServerTsValues},
		{Arg: ArgValidation, Required: true, Values: ValidationTsValues},
		{Arg: ArgGeneratePath, Required: true},
		{Arg: ArgServicesPath, Required: false},
		{Arg: ArgSwaggerPath, Required: false},
	},
	func(specification *spec.Spec, params map[Arg]string) *sources.Sources {
		return gents.GenerateService(specification, params[ArgSwaggerPath], params[ArgGeneratePath], params[ArgServicesPath], params[ArgServer], params[ArgValidation])
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
	Openapi,
}
