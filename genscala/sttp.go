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
	modelsFile := GenerateCirceModels(specification, clientPackage, generatePath)
	interfacesFile := generateClientApisInterfaces(specification, clientPackage, generatePath)
	implsFile := generateClientApiImplementations(specification, clientPackage, generatePath)

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
		Import("spec.circe.json._").
		Import("spec.http._")

	if len(specification.Apis) > 1 {
		unit.AddDeclarations(generateClientSuperClass(specification))
	}

	modelsMap := buildModelsMap(specification.Models)
	for _, api := range specification.Apis {
		apiTrait := generateClientApiClass(api, modelsMap)
		unit.AddDeclarations(apiTrait)
	}

	return &gen.TextFile{
		Path:    filepath.Join(outPath, "Client.scala"),
		Content: unit.Code(),
	}
}

func generateClientApisInterfaces(specification *spec.Spec, packageName string, outPath string) *gen.TextFile {
	modelsMap := buildModelsMap(specification.Models)

	unit := scala.Unit(packageName)

	unit.
		Import("scala.concurrent._").
		Import("spec.http._")

	for _, api := range specification.Apis {
		apiTrait := generateClientApiTrait(modelsMap, api)
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

func addParams(modelsMap ModelsMap, method *scala.MethodDeclaration, params []spec.NamedParam, defaulted bool) {
	for _, param := range params {
		if !defaulted && param.Default == nil {
			method.Param(param.Name.CamelCase(), ScalaType(&param.Type))
		}
		if defaulted && param.Default != nil {
			method.Param(param.Name.CamelCase(), ScalaType(&param.Type)).Init(scala.Code(DefaultValue(&param.Type, *param.Default, modelsMap)))
		}
	}
}

func generateClientOperationSignature(modelsMap ModelsMap, operation spec.NamedOperation) *scala.MethodDeclaration {
	returnType := "Future[" + responseType(operation) + "]"
	method := scala.Method(operation.Name.CamelCase()).Returns(returnType)
	addParams(modelsMap, method, operation.HeaderParams, false)
	if operation.Body != nil {
		method.Param("body", ScalaType(&operation.Body.Type))
	}
	for _, param := range operation.Endpoint.UrlParams {
		method.Param(param.Name.CamelCase(), ScalaType(&param.Type))
	}
	addParams(modelsMap, method, operation.QueryParams, false)
	addParams(modelsMap, method, operation.HeaderParams, true)
	addParams(modelsMap, method, operation.QueryParams, true)
	return method
}

func generateClientApiTrait(modelsMap ModelsMap, api spec.Api) *scala.TraitDeclaration {
	apiTraitName := clientTraitName(api.Name)
	apiTrait := scala.Trait(apiTraitName)
	apiTrait_ := apiTrait.Define(true)
	apiTrait_.AddCode(scala.Import(apiTraitName + "._"))
	for _, operation := range api.Operations {
		apiTrait_.AddCode(generateClientOperationSignature(modelsMap, operation))
	}
	return apiTrait
}

func clientTraitName(apiName spec.Name) string {
	return "I" + apiName.PascalCase() + "Client"
}

func clientClassName(apiName spec.Name) string {
	return apiName.PascalCase() + "Client"
}

func addParamsWriting(modelsMap ModelsMap, code *scala.StatementsDeclaration, params []spec.NamedParam, paramsName string) {
	if params != nil && len(params) > 0 {
		code.AddLn("val " + paramsName + " = new StringParamsWriter()")
		for _, p := range params {
			paramBaseType := p.Type.BaseType()
			if model, ok := modelsMap[paramBaseType.PlainType]; ok {
				if model.IsEnum() {
					code.AddLn(paramsName + `.write("` + p.Name.Source + `", ` + p.Name.CamelCase() + `.value)`)
				}
			} else {
				code.AddLn(paramsName + `.write("` + p.Name.Source + `", ` + p.Name.CamelCase() + `)`)
			}
		}
	}
}

func generateClientOperationImplementation(modelsMap ModelsMap, operation spec.NamedOperation, method *scala.MethodDeclaration) {
	httpMethod := strings.ToLower(operation.Endpoint.Method)
	url := operation.Endpoint.Url
	for _, param := range operation.Endpoint.UrlParams {
		url = strings.Replace(url, spec.UrlParamStr(param.Name.Source), "$"+param.Name.CamelCase(), -1)
	}

	method_ := method.Define().Block(true)

	method_.AddLn("implicit val jsonConfig = Json.config")

	addParamsWriting(modelsMap, method_, operation.QueryParams, "query")

	if operation.QueryParams != nil && len(operation.QueryParams) > 0 {
		method_.AddLn(`val url = Uri.parse(baseUrl+s"` + url + `").get.params(query.params)`)
	} else {
		method_.AddLn(`val url = Uri.parse(baseUrl+s"` + url + `").get`)
	}

	addParamsWriting(modelsMap, method_, operation.HeaderParams, "headers")

	method_.AddLn("val response: Future[Response[String]] =")
	httpCall := method_.Block(false).AddLn("sttp").Block(false)
	httpCall.AddLn(`.` + httpMethod + `(url)`)
	if operation.HeaderParams != nil && len(operation.HeaderParams) > 0 {
		httpCall.AddLn(`.headers(headers.params)`)
	}
	if operation.Body != nil {
		httpCall.AddLn(`.body(Json.write(body))`)
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

func generateClientApiClass(api spec.Api, modelsMap ModelsMap) *scala.ClassDeclaration {
	apiClassName := clientClassName(api.Name)
	apiTraitName := clientTraitName(api.Name)
	apiClass := scala.Class(apiClassName).Extends(apiTraitName)
	apiClassCtor := apiClass.Contructor()
	apiClassCtor.Param("baseUrl", "String")
	apiClassCtor.ImplicitParam("backend", "SttpBackend[Future, Source[ByteString, Any]]")
	apiClass_ := apiClass.Define(true)
	apiClass_.AddCode(scala.Import(apiTraitName + "._"))
	apiClass_.AddCode(scala.Import("ExecutionContext.Implicits.global"))
	for _, operation := range api.Operations {
		method := generateClientOperationSignature(modelsMap, operation)
		generateClientOperationImplementation(modelsMap, operation, method)
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
