package genscala

import (
	"github.com/ModaOperandi/spec"
	"github.com/vsapronov/gopoetry/scala"
	"path/filepath"
	"specgen/gen"
	"strings"
)

func GenerateSttpClient(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}

	clientPackage := clientPackageName(specification.ServiceName)

	jsonFile := GenerateJsonObject(clientPackage, generatePath)
	operationResultFile := GenerateOperationResult(clientPackage, generatePath)
	stringParamsFile := GenerateStringParams(clientPackage, generatePath)
	backendFile := GenerateSttpBackend(clientPackage, generatePath)

	modelsFile := GenerateCirceModels(specification, clientPackage, generatePath)
	interfacesFile := generateClientApisInterfaces(specification, clientPackage, generatePath)
	implsFile := generateClientApiImplementations(specification, clientPackage, generatePath)

	sourceManaged := []gen.TextFile{
		*jsonFile,
		*operationResultFile,
		*stringParamsFile,
		*backendFile,
		*modelsFile,
		*interfacesFile,
		*implsFile,
	}

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
		Import("json._")

	if len(specification.Apis) > 1 {
		unit.AddDeclarations(generateClientSuperClass(specification))
	}

	for _, api := range specification.Apis {
		apiTrait := generateClientApiClass(api)
		unit.AddDeclarations(apiTrait)
	}

	return &gen.TextFile{
		Path:    filepath.Join(outPath, "Client.scala"),
		Content: unit.Code(),
	}
}

func generateClientApisInterfaces(specification *spec.Spec, packageName string, outPath string) *gen.TextFile {
	unit := scala.Unit(packageName)

	unit.
		Import("scala.concurrent._").
		Import("json._")

	for _, api := range specification.Apis {
		apiTrait := generateClientApiTrait(api)
		unit.AddDeclarations(apiTrait)
	}

	for _, api := range specification.Apis {
		apiObject := generateApiInterfaceResponse(api, clientTraitName(api.Name))
		unit.AddDeclarations(apiObject)
	}

	return &gen.TextFile{
		Path:    filepath.Join(outPath, "Interfaces.scala"),
		Content: unit.Code(),
	}
}

func addParams(method *scala.MethodDeclaration, params []spec.NamedParam, defaulted bool) {
	for _, param := range params {
		if !defaulted && param.Default == nil {
			method.Param(param.Name.CamelCase(), ScalaType(&param.Type.Definition))
		}
		if defaulted && param.Default != nil {
			defaultValue := DefaultValue(&param.Type.Definition, *param.Default)
			method.Param(param.Name.CamelCase(), ScalaType(&param.Type.Definition)).Init(scala.Code(defaultValue))
		}
	}
}

func generateClientOperationSignature(operation spec.NamedOperation) *scala.MethodDeclaration {
	returnType := "Future[" + responseType(operation) + "]"
	method := scala.Method(operation.Name.CamelCase()).Returns(returnType)
	addParams(method, operation.HeaderParams, false)
	if operation.Body != nil {
		method.Param("body", ScalaType(&operation.Body.Type.Definition))
	}
	for _, param := range operation.Endpoint.UrlParams {
		method.Param(param.Name.CamelCase(), ScalaType(&param.Type.Definition))
	}
	addParams(method, operation.QueryParams, false)
	addParams(method, operation.HeaderParams, true)
	addParams(method, operation.QueryParams, true)
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

func addParamsWriting(code *scala.StatementsDeclaration, params []spec.NamedParam, paramsName string) {
	if params != nil && len(params) > 0 {
		code.AddLn("val " + paramsName + " = new StringParamsWriter()")
		for _, p := range params {
			code.AddLn(paramsName + `.write("` + p.Name.Source + `", ` + p.Name.CamelCase() + `)`)
		}
	}
}

func generateClientOperationImplementation(method *scala.MethodDeclaration, operation spec.NamedOperation) {
	httpMethod := strings.ToLower(operation.Endpoint.Method)
	url := operation.Endpoint.Url
	for _, param := range operation.Endpoint.UrlParams {
		url = strings.Replace(url, spec.UrlParamStr(param.Name.Source), "$"+param.Name.CamelCase(), -1)
	}

	method_ := method.Define().Block(true)

	addParamsWriting(method_, operation.QueryParams, "query")

	if operation.QueryParams != nil && len(operation.QueryParams) > 0 {
		method_.AddLn(`val url = Uri.parse(baseUrl+s"` + url + `").get.params(query.params)`)
	} else {
		method_.AddLn(`val url = Uri.parse(baseUrl+s"` + url + `").get`)
	}

	addParamsWriting(method_, operation.HeaderParams, "headers")

	if operation.Body != nil {
		method_.AddLn("val bodyJson = Jsoner.write(body)")
	}

	method_.AddLn("val response: Future[Response[String]] =")
	httpCall := method_.Block(false).AddLn("sttp").Block(false)
	httpCall.AddLn(`.` + httpMethod + `(url)`)
	if operation.HeaderParams != nil && len(operation.HeaderParams) > 0 {
		httpCall.AddLn(`.headers(headers.params)`)
	}
	if operation.Body != nil {
		httpCall.AddLn(`.header("Content-Type", "application/json")`)
		httpCall.AddLn(`.body(bodyJson)`)
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
		AddLn(`case Left(errorData) => throw new RuntimeException(s"Request failed, status code: ${response.code}, body: ${new String(errorData)}")`)
}

func generateClientApiClass(api spec.Api) *scala.ClassDeclaration {
	apiClassName := clientClassName(api.Name)
	apiTraitName := clientTraitName(api.Name)
	apiClass := scala.Class(apiClassName)
	apiClass.Extends(apiTraitName)
	apiClassCtor := apiClass.Contructor()
	apiClassCtor.Param("baseUrl", "String")
	apiClassCtor.ImplicitParam("backend", "SttpBackend[Future, Source[ByteString, Any]]")
	apiClass_ := apiClass.Define(true)
	apiClass_.AddCode(scala.Import(apiTraitName + "._"))
	apiClass_.AddCode(scala.Import("ExecutionContext.Implicits.global"))
	for _, operation := range api.Operations {
		method := generateClientOperationSignature(operation)
		generateClientOperationImplementation(method, operation)
		apiClass_.AddCode(method)
	}
	return apiClass
}

func generateClientSuperClass(specification *spec.Spec) *scala.ClassDeclaration {
	clientClass := scala.Class(specification.ServiceName.PascalCase() + "Client")
	clientClassCtor := clientClass.Contructor()
	clientClassCtor.Param("baseUrl", "String")
	clientClassCtor.ImplicitParam("backend", "SttpBackend[Future, Source[ByteString, Any]]")
	clientClass_ := clientClass.Define(true)
	for _, api := range specification.Apis {
		clientClass_.Val(api.Name.CamelCase(), clientTraitName(api.Name)).Init(scala.Code("new " + clientClassName(api.Name) + "(baseUrl)"))
		clientClass_.AddLn("")
	}
	return clientClass
}
