package gen

const ModuleName = "module-name"
const ModuleNameDescription = "module name"

const SpecFile = "spec-file"
const SpecFileDescription = "path to specification file"

const GeneratePath = "generate-path"
const GeneratePathDescription = "path to generate source code into"

const SwaggerPath = "swagger-path"
const SwaggerPathDescription = "path of generated OpenAPI (Swagger) specification file"

const ServicesPath = "services-path"
const ServicesPathDescription = "path to scaffold services code"

const PackageName = "package-name"
const PackageNameDescription = "package name"

const Jsonlib = "jsonlib"
const JsonlibDescription = "JSON library: jackson, moshi"

const Validation = "validation"
const ValidationDescription = "validation TypeScript library: superstruct, io-ts"

const Client = "client"
const ClientDescription = "client TypeScript library: axios, node-fetch, browser-fetch"

const Server = "server"
const ServerDescription = "server TypeScript library: express, koa"

type Arg struct {
	Name        string
	Description string
}

var ArgSpecFile = Arg{SpecFile, SpecFileDescription}
var ArgModuleName = Arg{ModuleName, ModuleNameDescription}
var ArgGeneratePath = Arg{GeneratePath, GeneratePathDescription}
var ArgSwaggerPath = Arg{SwaggerPath, SwaggerPathDescription}
var ArgServicesPath = Arg{ServicesPath, ServicesPathDescription}
var ArgPackageName = Arg{PackageName, PackageNameDescription}
var ArgJsonlib = Arg{Jsonlib, JsonlibDescription}
var ArgValidation = Arg{Validation, ValidationDescription}
var ArgClient = Arg{Client, ClientDescription}
var ArgServer = Arg{Server, ServerDescription}
