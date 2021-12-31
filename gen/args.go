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
const JsonlibDescription = "json serialization/deserialization library"

const Validation = "validation"
const ValidationDescription = "type validation library"

const Client = "client"
const ClientDescription = "client HTTP library"

const Server = "server"
const ServerDescription = "server HTTP library/framework"

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
