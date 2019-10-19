package genscala

import (
	"github.com/ModaOperandi/spec"
	"github.com/vsapronov/gopoetry/scala"
	"path/filepath"
	"specgen/gen"
	"strings"
)

func GenerateSttpClient(serviceFile string, sourceManagedPath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}

	clientPackage := clientPackageName(specification.ServiceName)
	modelsFile := GenerateModels(specification, clientPackage, sourceManagedPath)
	interfacesFile := generateClientApisInterfaces(specification.Apis, clientPackage, sourceManagedPath)
	implsFile := generateClientApiImplementations(specification, clientPackage, sourceManagedPath)

	sourceManaged := []gen.TextFile{*modelsFile, *interfacesFile, *implsFile}

	err = gen.WriteFiles(sourceManaged, true)
	if err != nil {
		return err
	}

	return nil
}

func clientPackageName(name spec.Name) string {
	return name.FlatCase() + ".client"
}

func generateClientApiImplementations(specification *spec.Spec, packageName string, outPath string) *gen.TextFile {
	unit := scala.Unit(packageName)

	unit.
		Import("scala.concurrent._").
		Import("akka.stream.scaladsl.Source").
		Import("akka.util.ByteString").
		Import("com.softwaremill.sttp._").
		Import("com.softwaremill.sttp.akkahttp._").
		Import("spec.jackson._").
		Import("spec.http._")

	unit.AddDeclarations(generateClientClass(specification))
	for _, api := range specification.Apis {
		apiTrait := generateClientApiClass(api)
		unit.AddDeclarations(apiTrait)
	}

	return &gen.TextFile{
		Path:    filepath.Join(outPath, "Client.scala"),
		Content: unit.Code(),
	}
}

func generateClientApisInterfaces(apis spec.Apis, packageName string, outPath string) *gen.TextFile {
	unit := scala.Unit(packageName)

	unit.
		Import("scala.concurrent._").
		Import("spec.jackson._").
		Import("spec.http._")

	for _, api := range apis {
		apiTrait := generateClientApiTrait(api)
		unit.AddDeclarations(apiTrait)
	}

	for _, api := range apis {
		apiObject := generateApiInterfaceResponse(api, clientTraitName(api.Name))
		unit.AddDeclarations(apiObject)
	}

	return &gen.TextFile{
		Path:    filepath.Join(outPath, "Interfaces.scala"),
		Content: unit.Code(),
	}
}

func generateClientOperationSignature(operation spec.NamedOperation) *scala.MethodDeclaration {
	returnType := "Future[" + responseType(operation) + "]"
	method := scala.Method(operation.Name.CamelCase()).Returns(returnType)
	for _, param := range operation.HeaderParams {
		method.Param(param.Name.CamelCase(), ScalaType(&param.Type))
	}
	if operation.Body != nil {
		method.Param("body", ScalaType(&operation.Body.Type))
	}
	for _, param := range operation.UrlParams {
		method.Param(param.Name.CamelCase(), ScalaType(&param.Type))
	}
	for _, param := range operation.QueryParams {
		method.Param(param.Name.CamelCase(), ScalaType(&param.Type))
	}
	return method
}

func generateClientApiTrait(api spec.Api) *scala.TraitDeclaration {
	apiTraitName := clientTraitName(api.Name)
	apiTrait := scala.Trait(apiTraitName)
	apiTrait_ := apiTrait.Define(true)
	apiTrait_.AddCode(scala.Import(apiTraitName + "._"))
	for _, operation := range api.Operations {
		apiTrait_.AddCode(generateClientOperationSignature(operation))
	}
	return apiTrait
}

func clientTraitName(apiName spec.Name) string {
	return "I" + apiName.PascalCase() + "Client"
}

func clientClassName(apiName spec.Name) string {
	return apiName.PascalCase() + "Client"
}

func generateClientOperationImplementation(operation spec.NamedOperation, method *scala.MethodDeclaration) {
	httpMethod := strings.ToLower(operation.Method)
	url := operation.Url
	for _, param := range operation.UrlParams {
		url = strings.Replace(url, spec.UrlParamStr(param.Name.Source), "$"+param.Name.CamelCase(), -1)
	}

	method_ := method.Define().Block(true)

	if len(operation.QueryParams) > 0 {
		method_.AddLn("val query = new StringParamsWriter()")
		for _, p := range operation.QueryParams {
			method_.AddLn(`query.write("` + p.Name.Source + `", ` + p.Name.CamelCase() + `)`)
		}
		method_.AddLn(`val url = Uri.parse(baseUrl+s"` + url + `").get.params(query.params)`)
	} else {
		method_.AddLn(`val url = Uri.parse(baseUrl+s"` + url + `").get`)
	}

	if len(operation.HeaderParams) > 0 {
		method_.AddLn("val headers = new StringParamsWriter()")
		for _, p := range operation.QueryParams {
			method_.AddLn(`headers.write("` + p.Name.Source + `", ` + p.Name.CamelCase() + `)`)
		}
	}

	method_.AddLn("val response: Future[Response[String]] =")
	httpCall := method_.Block(false).AddLn("sttp").Block(false)
	httpCall.AddLn(`.` + httpMethod + `(url)`)
	if len(operation.HeaderParams) > 0 {
		httpCall.AddLn(`.headers(headers.params)`)
	}
	if operation.Body != nil {
		httpCall.AddLn(`.body(json.writeValueAsString(body))`)
	}
	httpCall.AddLn(`.parseResponseIf { status => status < 500 }`) // TODO: Allowed statuses from spec
	httpCall.AddLn(`.send()`)

	responseName := responseType(operation)
	method_.
		Add(`response.map `).Block(true).
		AddLn(`response: Response[String] =>`).Block(false).
		Add(`response.body match `).Block(true).
		AddLn(`case Right("") => ` + responseName + `.fromResult(OperationResult(response.code, None))`).
		AddLn(`case Right(bodyStr) => ` + responseName + `.fromResult(OperationResult(response.code, Some(bodyStr)))`).
		AddLn(`case Left(errorData) => throw new RuntimeException("Request failed")`)
}

func generateClientApiClass(api spec.Api) *scala.ClassDeclaration {
	apiClassName := clientClassName(api.Name)
	apiTraitName := clientTraitName(api.Name)
	apiClass := scala.Class(apiClassName).Extends(apiTraitName)
	apiClassCtor := apiClass.Contructor()
	apiClassCtor.Param("baseUrl", "String")
	apiClassCtor.ImplicitParam("json", "JsonMapper")
	apiClassCtor.ImplicitParam("backend", "SttpBackend[Future, Source[ByteString, Any]]")
	apiClass_ := apiClass.Define(true)
	apiClass_.AddCode(scala.Import(apiTraitName + "._"))
	apiClass_.AddCode(scala.Import("ExecutionContext.Implicits.global"))
	for _, operation := range api.Operations {
		method := generateClientOperationSignature(operation)
		generateClientOperationImplementation(operation, method)
		apiClass_.AddCode(method)
	}
	return apiClass
}

func generateClientClass(specification *spec.Spec) *scala.ClassDeclaration {
	clientClass := scala.Class(specification.ServiceName.PascalCase() + "Client")
	clientClassCtor := clientClass.Contructor()
	clientClassCtor.Param("baseUrl", "String")
	clientClassCtor.ImplicitParam("json", "JsonMapper")
	clientClassCtor.ImplicitParam("backend", "SttpBackend[Future, Source[ByteString, Any]]")
	clientClass_ := clientClass.Define(true)
	for _, api := range specification.Apis {
		clientClass_.Val(api.Name.CamelCase(), clientTraitName(api.Name)).Val().Init(scala.Code("new " + clientClassName(api.Name) + "(baseUrl)"))
		clientClass_.AddLn("")
	}
	return clientClass
}
