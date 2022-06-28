package generator

const ModuleName = "module-name"
const ModuleNameTitle = "Module name"
const ModuleNameDescription = "module name"

const SpecFile = "spec-file"
const SpecFileTitle = "Spec file path"
const SpecFileDescription = "path to specification file"

const GeneratePath = "generate-path"
const GeneratePathTitle = "Generate path"
const GeneratePathDescription = "path to generate source code into"

const SwaggerPath = "swagger-path"
const SwaggerPathTitle = "Swagger path"
const SwaggerPathDescription = "path of generated OpenAPI (Swagger) specification file"

const ServicesPath = "services-path"
const ServicesPathTitle = "Services path"
const ServicesPathDescription = "path to scaffold services code"

const PackageName = "package-name"
const PackageNameTitle = "Package name"
const PackageNameDescription = "package name"

const Jsonlib = "jsonlib"
const JsonlibTitle = "JSON library"
const JsonlibDescription = "json serialization/deserialization library"

const Validation = "validation"
const ValidationTitle = "Type validation library"
const ValidationDescription = "type validation library"

const Client = "client"
const ClientTitle = "Client library"
const ClientDescription = "client HTTP library"

const Server = "server"
const ServerTitle = "Server library"
const ServerDescription = "server HTTP library/framework"

const OutFile = "out-file"
const OutFileTitle = "Output file path"
const OutFileDescription = "path to output file"

var ArgSpecFile = Arg{SpecFile, SpecFileTitle, SpecFileDescription}
var ArgOutFile = Arg{OutFile, OutFileTitle, OutFileDescription}
var ArgModuleName = Arg{ModuleName, ModuleNameTitle, ModuleNameDescription}
var ArgGeneratePath = Arg{GeneratePath, GeneratePathTitle, GeneratePathDescription}
var ArgSwaggerPath = Arg{SwaggerPath, SwaggerPathTitle, SwaggerPathDescription}
var ArgServicesPath = Arg{ServicesPath, ServicesPathTitle, ServicesPathDescription}
var ArgPackageName = Arg{PackageName, PackageNameTitle, PackageNameDescription}
var ArgJsonlib = Arg{Jsonlib, JsonlibTitle, JsonlibDescription}
var ArgValidation = Arg{Validation, ValidationTitle, ValidationDescription}
var ArgClient = Arg{Client, ClientTitle, ClientDescription}
var ArgServer = Arg{Server, ServerTitle, ServerDescription}
