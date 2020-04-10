package genscala

import (
	"github.com/ModaOperandi/spec"
	"github.com/vsapronov/gopoetry/scala"
	"path/filepath"
	"specgen/gen"
	"specgen/genopenapi"
	"strings"
)

func GeneratePlayService(serviceFile string, swaggerPath string, generatePath string, servicesPath string, routesPath string) (err error) {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return
	}

	modelsPackage := modelsPackage(specification)
	controllersPackage := controllersPackage(specification)
	jsonFile := GenerateJsonObject(modelsPackage, generatePath)
	servicesPackage := servicesPackage(specification)

	operationResultFile := GenerateOperationResult(servicesPackage, generatePath)
	resultHelpersFile := GeneratePlayResultHelpers(controllersPackage, generatePath)
	stringParamsFile := GenerateStringParams(controllersPackage, generatePath)

	modelsFile := GenerateCirceModels(specification, modelsPackage, generatePath)

	source := []gen.TextFile{}
	sourceManaged := []gen.TextFile{
		*jsonFile,
		*operationResultFile,
		*resultHelpersFile,
		*stringParamsFile,
		*modelsFile,
	}

	apis := specification.Apis

	for _, api := range apis {
		apiTraitFile := generateApiInterface(api, servicesPackage, generatePath)
		apiControllerFile := generateApiController(api, controllersPackage, generatePath)
		sourceManaged = append(sourceManaged, *apiTraitFile, *apiControllerFile)
		apiClassFile := generateApiClass(api, servicesPackage, servicesPath)
		source = append(source, *apiClassFile)
	}

	openApiRoutes := generateCommonRoutes()
	controllersRoutes := generateControllersRoutes(specification)

	routesFile :=
		&gen.TextFile{
			Path:    filepath.Join(routesPath, "routes"),
			Content: openApiRoutes + controllersRoutes,
		}

	resource := []gen.TextFile{*routesFile}

	swaggerFile := filepath.Join(swaggerPath, "swagger.yaml")
	genopenapi.GenerateSpecification(serviceFile, swaggerFile)

	err = gen.WriteFiles(source, false)
	if err != nil {
		return
	}

	err = gen.WriteFiles(sourceManaged, true)
	if err != nil {
		return
	}

	err = gen.WriteFiles(resource, true)
	if err != nil {
		return
	}

	return
}

func controllerType(apiName spec.Name) string {
	return apiName.PascalCase() + "Controller"
}

func apiTraitType(apiName spec.Name) string {
	return "I" + apiName.PascalCase() + "Service"
}

func apiClassType(apiName spec.Name) string {
	return apiName.PascalCase() + "Service"
}

func controllersPackage(specification *spec.Spec) string {
	return "controllers"
}

func servicesPackage(specification *spec.Spec) string {
	return "services"
}

func modelsPackage(specification *spec.Spec) string {
	return "models"
}

func operationSignature(operation spec.NamedOperation) *scala.MethodDeclaration {
	returnType := "Future[" + responseType(operation) + "]"
	method := scala.Method(operation.Name.CamelCase()).Returns(returnType)
	for _, param := range operation.HeaderParams {
		method.Param(param.Name.CamelCase(), ScalaType(&param.Type.Definition))
	}
	if operation.Body != nil {
		method.Param("body", ScalaType(&operation.Body.Type.Definition))
	}
	for _, param := range operation.Endpoint.UrlParams {
		method.Param(param.Name.CamelCase(), ScalaType(&param.Type.Definition))
	}
	for _, param := range operation.QueryParams {
		method.Param(param.Name.CamelCase(), ScalaType(&param.Type.Definition))
	}
	return method
}

func generateApiInterface(api spec.Api, packageName string, outPath string) *gen.TextFile {
	unit := scala.Unit(packageName)
	unit.
		Import("com.google.inject.ImplementedBy").
		Import("scala.concurrent.Future").
		Import("models._").
		Import("json._")

	apiTraitName := apiTraitType(api.Name)

	apiTrait := generateApiInterfaceTrait(api, apiTraitName)
	unit.AddDeclarations(apiTrait)

	apiObject := generateApiInterfaceResponse(api, apiTraitName)
	unit.AddDeclarations(apiObject)

	return &gen.TextFile{
		Path:    filepath.Join(outPath, apiTraitName+".scala"),
		Content: unit.Code(),
	}
}

func generateApiInterfaceTrait(api spec.Api, apiTraitName string) *scala.TraitDeclaration {
	apiTrait := scala.Trait(apiTraitName).Attribute("ImplementedBy(classOf[" + apiClassType(api.Name) + "])")
	apiTrait_ := apiTrait.Define(true)
	apiTrait_.AddCode(scala.Import(apiTraitName + "._"))
	for _, operation := range api.Operations {
		method := operationSignature(operation)
		apiTrait_.AddCode(method)
	}
	return apiTrait
}

func generateApiClass(api spec.Api, packageName string, outPath string) *gen.TextFile {
	unit := scala.Unit(packageName)
	unit.
		Import("javax.inject._").
		Import("scala.concurrent._").
		Import("models._")

	apiClassName := apiClassType(api.Name)
	apiTraitName := apiTraitType(api.Name)
	class := scala.Class(apiClassName).Attribute("Singleton")
	class.Extends(apiTraitName)
	ctor := class.Contructor()
	ctor.ImplicitParam("ec", "ExecutionContext")
	ctor.Attribute("Inject()")
	class_ := class.Define(true)
	class_.AddCode(scala.Import(apiTraitName + "._"))

	for _, operation := range api.Operations {
		method := operationSignature(operation).Override()
		method.Define().AddLn("Future { ??? }")
		class_.AddCode(method)
	}

	unit.AddDeclarations(class)

	return &gen.TextFile{
		Path:    filepath.Join(outPath, apiClassName+".scala"),
		Content: unit.Code(),
	}
}

func addParamsParsing(code *scala.StatementsDeclaration, params []spec.NamedParam, paramsName string, readingFun string) {
	if params != nil && len(params) > 0 {
		code.AddLn(`val ` + paramsName + ` = new StringParamsReader(` + readingFun + `)`)
		for _, param := range params {
			paramBaseType := param.Type.Definition.BaseType()
			method := "read"
			if paramBaseType.Info.Model != nil && paramBaseType.Info.Model.IsEnum() {
				method = "readEnum"
			}
			code.Add(`val ` + param.Name.CamelCase() + ` = ` + paramsName + `.` + method + `[` + ScalaType(paramBaseType) + `]("` + param.Name.Source + `")`)
			if !param.Type.Definition.IsNullable() {
				if param.Default != nil {
					code.Add(`.getOrElse(` + DefaultValue(&param.Type.Definition, *param.Default) + `)`)
				} else {
					code.Add(".get")
				}
			}
			code.AddLn("")
		}
	}

}

func generateApiController(api spec.Api, packageName string, outPath string) *gen.TextFile {
	unit := scala.Unit(packageName)

	unit.
		Import("javax.inject._").
		Import("scala.util._").
		Import("scala.concurrent._").
		Import("play.api.mvc._").
		Import("controllers.PlayResultHelpers._").
		Import("models._").
		Import("json._").
		Import("services._")

	class := scala.Class(controllerType(api.Name)).Attribute("Singleton")
	ctor := class.Contructor()
	ctor.Attribute("Inject()")
	ctor.Param("api", apiTraitType(api.Name))
	ctor.Param("cc", "ControllerComponents")
	ctor.ImplicitParam("ec", "ExecutionContext")

	class.Extends("AbstractController(cc)")
	class_ := class.Define(true)

	for _, operation := range api.Operations {
		method := class_.Def(operation.Name.CamelCase())
		for _, param := range operation.Endpoint.UrlParams {
			method.Param(param.Name.CamelCase(), ScalaType(&param.Type.Definition))
		}
		definition := method.Define()
		if operation.Body != nil {
			definition.Add("Action(parse.byteString).async ")
		} else {
			definition.Add("Action.async ")
		}
		lambda := definition.Block(true).AddLn("implicit request =>").Block(false)

		parseParams := getOperationParams(operation, false)

		if len(parseParams) > 0 {
			tryBlock := lambda.Add("val params = Try ").Block(true)
			addParamsParsing(tryBlock, operation.HeaderParams, "header", "request.headers.get")
			addParamsParsing(tryBlock, operation.QueryParams, "query", "request.getQueryString")
			if operation.Body != nil {
				tryBlock.AddLn("val body = Jsoner.read[" + ScalaType(&operation.Body.Type.Definition) + "](request.body.utf8String)")
			}
			tryBlock.AddLn("(" + strings.Join(parseParams, ", ") + ")")
		}

		if len(parseParams) > 0 {
			match := lambda.Add("params match ").Block(true)
			match.AddLn("case Failure(ex) => Future { BadRequest }")
			success := match.AddLn("case Success(params) => ").Block(false)
			success.AddLn("val " + "(" + strings.Join(parseParams, ", ") + ") = params")
			callApi(success, operation)
		} else {
			callApi(lambda, operation)
		}
	}

	unit.AddDeclarations(class)

	return &gen.TextFile{
		Path:    filepath.Join(outPath, controllerType(api.Name)+".scala"),
		Content: unit.Code(),
	}
}

func callApi(lambda *scala.StatementsDeclaration, operation spec.NamedOperation) {
	allParams := getOperationParams(operation, true)
	lambda.AddLn("val result = api." + operation.Name.CamelCase() + "(" + strings.Join(allParams, ", ") + ")")
	lambda.AddLn("result.map(_.toResult.toPlay).recover { case _: Exception => InternalServerError }")
}

func getOperationParams(operation spec.NamedOperation, includeUrl bool) []string {
	params := []string{}
	if operation.HeaderParams != nil {
		for _, param := range operation.HeaderParams {
			params = append(params, param.Name.CamelCase())
		}
	}
	if operation.Body != nil {
		params = append(params, "body")
	}
	if includeUrl {
		for _, param := range operation.Endpoint.UrlParams {
			params = append(params, param.Name.CamelCase())
		}
	}
	if operation.QueryParams != nil {
		for _, param := range operation.QueryParams {
			params = append(params, param.Name.CamelCase())
		}
	}
	return params
}

func tail(s string, length int) string {
	return s + strings.Repeat(" ", length-len(s))
}

func generateCommonRoutes() string {
	return `#  documentation
GET      /docs                controllers.Default.redirect(to="/docs/index.html?url=swagger.yaml")
GET      /docs/swagger.yaml   controllers.Assets.at(path="/public", file="swagger.yaml")
GET      /docs/*file          controllers.Assets.at(path="/public/lib/swagger-ui", file)

#  static files
GET      /*file   controllers.Assets.at(path="/public", file)
`
}

func generateControllersRoutes(specification *spec.Spec) string {
	builder := strings.Builder{}

	routeLength := 0
	for _, api := range specification.Apis {
		for _, operation := range api.Operations {
			thisRouteLength := len(routeUrl(operation))
			if thisRouteLength > routeLength {
				routeLength = thisRouteLength
			}
		}
	}

	for _, api := range specification.Apis {
		builder.WriteString("\n")
		builder.WriteString("#  " + api.Name.Source + "\n")
		for _, operation := range api.Operations {
			controllerEndpoint := controllersPackage(specification) + "." + controllerType(api.Name) + "." + operation.Name.CamelCase()
			params := []string{}
			for _, param := range operation.Endpoint.UrlParams {
				params = append(params, param.Name.CamelCase()+": "+ScalaType(&param.Type.Definition))
			}
			route := tail(operation.Endpoint.Method, 8) + " " + tail(routeUrl(operation), routeLength) + "   " + controllerEndpoint + "(" + strings.Join(params, ", ") + ")\n"
			builder.WriteString(route)
		}
	}

	return builder.String()
}

func routeUrl(operation spec.NamedOperation) string {
	routeUrl := operation.Endpoint.Url
	for _, param := range operation.Endpoint.UrlParams {
		routeUrl = strings.Replace(routeUrl, spec.UrlParamStr(param.Name.Source), ":"+param.Name.CamelCase(), 1)
	}
	return routeUrl
}
