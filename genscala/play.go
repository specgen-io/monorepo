package genscala

import (
	"fmt"
	spec "github.com/specgen-io/spec.v2"
	"github.com/vsapronov/gopoetry/scala"
	"path/filepath"
	"specgen/gen"
	"specgen/genopenapi"
	"specgen/static"
	"strings"
)

func GeneratePlayService(serviceFile string, swaggerPath string, generatePath string, servicesPath string) (err error) {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return
	}

	modelsPackage := modelsPackage(specification)
	controllersPackage := controllersPackage(specification)
	servicesPackage := servicesPackage(specification)

	scalaCirceFiles, err := static.RenderTemplate("scala-circe", generatePath, static.ScalaStaticCode{PackageName: "spec.models"})
	if err != nil {
		return
	}
	scalaPlayStaticFiles, err := static.RenderTemplate("scala-play", generatePath, static.ScalaStaticCode{PackageName: "spec.controllers"})
	if err != nil {
		return
	}
	scalaHttpStaticFiles, err := static.RenderTemplate("scala-http", generatePath, static.ScalaStaticCode{PackageName: "spec.controllers"})
	if err != nil {
		return
	}

	modelsFiles := GenerateCirceModels(specification, modelsPackage, generatePath)

	sourceManaged := modelsFiles
	sourceManaged = append(sourceManaged, scalaPlayStaticFiles...)
	sourceManaged = append(sourceManaged, scalaHttpStaticFiles...)
	sourceManaged = append(sourceManaged, scalaCirceFiles...)

	source := []gen.TextFile{}

	for _, versionedApis := range specification.Http.Versions {
		apisSourceManaged := generateApis(versionedApis, servicesPackage, controllersPackage, generatePath)
		sourceManaged = append(sourceManaged, apisSourceManaged...)
		apiControllerFile := generateApiControllers(versionedApis, controllersPackage, generatePath)
		sourceManaged = append(sourceManaged, *apiControllerFile)
		apiRoutersFile := generateRouter(&versionedApis, "app", generatePath)
		sourceManaged = append(sourceManaged, *apiRoutersFile)
		servicesSource := generateApisServices(versionedApis, servicesPackage, servicesPath)
		source = append(source, servicesSource...)
	}

	routesFile := generateMainRouter(specification.Http.Versions, "app", generatePath)
	sourceManaged = append(sourceManaged, *routesFile)

	genopenapi.GenerateSpecification(serviceFile, filepath.Join(swaggerPath, "swagger.yaml"))

	err = gen.WriteFiles(source, false)
	if err != nil {
		return
	}

	err = gen.WriteFiles(sourceManaged, true)
	if err != nil {
		return
	}

	return
}

func generateApis(versionedApis spec.VersionedApis, servicesPackage string, controllersPackage string, generatePath string) []gen.TextFile {
	sourceManaged := []gen.TextFile{}
	for _, api := range versionedApis.Apis {
		apiTraitFile := generateApiInterface(versionedApis.Version, api, servicesPackage, generatePath)
		sourceManaged = append(sourceManaged, *apiTraitFile)
	}
	return sourceManaged
}

func generateApisServices(versionedApis spec.VersionedApis, servicesPackage string, servicesPath string) []gen.TextFile {
	source := []gen.TextFile{}
	for _, api := range versionedApis.Apis {
		apiClassFile := generateApiClass(versionedApis.Version, api, servicesPackage, servicesPath)
		source = append(source, *apiClassFile)
	}
	return source
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

func controllerMethodName(operation spec.NamedOperation) string {
	return operation.Name.CamelCase()
}

func operationSignature(operation spec.NamedOperation) *scala.MethodDeclaration {
	returnType := "Future[" + responseType(operation) + "]"
	method := Def(controllerMethodName(operation)).Returns(returnType)
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

func generateApiInterface(version spec.Name, api spec.Api, packageName string, outPath string) *gen.TextFile {
	unit := Unit(versionedPackage(version, packageName))

	modelsPackage := versionedPackage(version, "models")

	unit.
		Import("com.google.inject.ImplementedBy").
		Import("scala.concurrent.Future").
		Import(modelsPackage + "._")

	apiTraitName := apiTraitType(api.Name)

	apiTrait := generateApiInterfaceTrait(api, apiTraitName)
	unit.AddDeclarations(apiTrait)

	apiObject := generateApiInterfaceResponse(api, apiTraitName)
	unit.AddDeclarations(apiObject)

	return &gen.TextFile{
		Path:    filepath.Join(outPath, fmt.Sprintf("%s%s.scala", apiTraitName, version.PascalCase())),
		Content: unit.Code(),
	}
}

func generateApiInterfaceTrait(api spec.Api, apiTraitName string) *scala.TraitDeclaration {
	apiTrait := Trait(apiTraitName).Attribute("ImplementedBy(classOf[" + apiClassType(api.Name) + "])")
	apiTrait.Add(Import(apiTraitName + "._"))
	for _, operation := range api.Operations {
		apiTrait.Add(operationSignature(operation))
	}
	return apiTrait
}

func generateApiClass(version spec.Name, api spec.Api, packageName string, outPath string) *gen.TextFile {
	unit := Unit(versionedPackage(version, packageName))

	modelsPackage := versionedPackage(version, "models")

	unit.
		Import("javax.inject._").
		Import("scala.concurrent._").
		Import(modelsPackage + "._")

	apiClassName := apiClassType(api.Name)
	apiTraitName := apiTraitType(api.Name)
	class :=
		Class(apiClassName).Attribute("Singleton").Extends(apiTraitName).
			Constructor(Constructor().
				Attribute("Inject()").
				ImplicitParam("ec", "ExecutionContext"),
			).
			Add(Import(apiTraitName + "._"))

	for _, operation := range api.Operations {
		method := operationSignature(operation).Override().BodyInline(Code("Future { ??? }"), Eol())
		class.Add(method)
	}

	unit.AddDeclarations(class)

	return &gen.TextFile{
		Path:    filepath.Join(outPath, fmt.Sprintf("%s%s.scala", apiClassName, version.PascalCase())),
		Content: unit.Code(),
	}
}

func addParamsParsing(params []spec.NamedParam, paramsName string, readingFun string) *scala.StatementsDeclaration {
	code := Statements()
	if params != nil && len(params) > 0 {
		code.Add(Line(`val %s = new StringParamsReader(%s)`, paramsName, readingFun))
		for _, param := range params {
			paramBaseType := param.Type.Definition.BaseType()
			method := "read"
			if paramBaseType.Info.Model != nil && paramBaseType.Info.Model.IsEnum() {
				method = "readEnum"
			}
			code.Add(Code(`val %s = %s.%s[%s]("%s")`, param.Name.CamelCase(), paramsName, method, ScalaType(paramBaseType), param.Name.Source))
			if !param.Type.Definition.IsNullable() {
				if param.Default != nil {
					code.Add(Line(`.getOrElse(%s)`, DefaultValue(&param.Type.Definition, *param.Default)))
				} else {
					code.Add(Line(".get"))
				}
			}
		}
	}
	return code
}

func generateApiControllers(versionedApis spec.VersionedApis, packageName string, outPath string) *gen.TextFile {
	unit := Unit(versionedPackage(versionedApis.Version, packageName))

	modelsPackage := versionedPackage(versionedApis.Version, "models")
	servicePackage := versionedPackage(versionedApis.Version, "services")
	unit.
		Import("javax.inject._").
		Import("scala.util._").
		Import("scala.concurrent._").
		Import("play.api.mvc._").
		Import("spec.controllers.ParamsTypesBindings._").
		Import("spec.models.Jsoner").
		Import(servicePackage + "._").
		Import(modelsPackage + "._").
		Import(modelsPackage + ".json._")

	for _, api := range versionedApis.Apis {
		class :=
			Class(controllerType(api.Name)).Attribute("Singleton").Extends("AbstractController(cc)").
				Constructor(Constructor().
					Attribute("Inject()").
					Param("api", apiTraitType(api.Name)).
					Param("cc", "ControllerComponents").
					ImplicitParam("ec", "ExecutionContext"),
				)

		class.Add(Import(apiTraitType(api.Name) + "._"))

		for _, operation := range api.Operations {
			class.Add(generateControllerMethod(operation))
		}
		unit.AddDeclarations(class)
	}

	return &gen.TextFile{
		Path:    filepath.Join(outPath, fmt.Sprintf("%sControllers.scala", versionedApis.Version.PascalCase())),
		Content: unit.Code(),
	}
}

func generateControllerMethod(operation spec.NamedOperation) *scala.MethodDeclaration {
	parseParams := getParsedOperationParams(operation)
	allParams := getOperationCallParams(operation)

	method := Def(operation.Name.CamelCase())

	for _, param := range operation.Endpoint.UrlParams {
		method.Param(param.Name.CamelCase(), ScalaType(&param.Type.Definition))
	}
	for _, param := range operation.QueryParams {
		method.Param(param.Name.Source, ScalaType(&param.Type.Definition))
	}

	method.BodyInline(
		Statements(Dynamic(func(code *scala.WritableList) {
			if operation.Body != nil {
				code.Add(Code("Action(parse.byteString).async "))
			} else {
				code.Add(Code("Action.async "))
			}
		})...),
		Scope(
			Line("implicit request =>"),
			Block(Dynamic(func(code *scala.WritableList) {
				if len(parseParams) > 0 {
					code.Add(
						Code("val params = Try "),
						Scope(
							addParamsParsing(operation.HeaderParams, "header", "request.headers.get"),
							Statements(Dynamic(func(code *scala.WritableList) {
								if operation.Body != nil {
									code.Add(Line("val body = Jsoner.read[%s](request.body.utf8String)", ScalaType(&operation.Body.Type.Definition)))
								}
							})...),
							Line("(%s)", JoinParams(parseParams)),
						),
						Code("params match "),
						Scope(
							Line("case Failure(ex) => Future { BadRequest }"),
							Line("case Success(params) => "),
							Block(
								Line("val (%s) = params", JoinParams(parseParams)),
								Line("val result = api.%s(%s)", operation.Name.CamelCase(), JoinParams(allParams)),
								Code("val response = result.map "),
								Scope(
									Dynamic(func(code *scala.WritableList) { genResponseCases(code, operation)})...,
								),
								Line("response.recover { case _: Exception => InternalServerError }"),
							),
						),
					)
				} else {
					code.Add(
						Line("val result = api.%s(%s)", operation.Name.CamelCase(), JoinParams(allParams)),
						Code("val response = result.map "),
						Scope(
							Dynamic(func(code *scala.WritableList) { genResponseCases(code, operation)})...,
						),
						Line("response.recover { case _: Exception => InternalServerError }"),
					)
				}
			})...),
		),
	)
	return method
}

func genResponseCases(code *scala.WritableList, operation spec.NamedOperation) {
	for _, r := range operation.Responses {
		if !r.Type.Definition.IsEmpty() {
			code.Add(Line("case %s.%s(body) => new Status(%s)(Jsoner.write(body))", responseType(operation), r.Name.PascalCase(), spec.HttpStatusCode(r.Name)))
		} else {
			code.Add(Line("case %s.%s() => new Status(%s)", responseType(operation), r.Name.PascalCase(), spec.HttpStatusCode(r.Name)))
		}
	}
}

func getOperationCallParams(operation spec.NamedOperation) []string {
	params := []string{}
	if operation.HeaderParams != nil {
		for _, param := range operation.HeaderParams {
			params = append(params, param.Name.CamelCase())
		}
	}
	if operation.Body != nil {
		params = append(params, "body")
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, param.Name.CamelCase())
	}
	if operation.QueryParams != nil {
		for _, param := range operation.QueryParams {
			params = append(params, param.Name.Source)
		}
	}
	return params
}

func getParsedOperationParams(operation spec.NamedOperation) []string {
	params := []string{}
	if operation.HeaderParams != nil {
		for _, param := range operation.HeaderParams {
			params = append(params, param.Name.CamelCase())
		}
	}
	if operation.Body != nil {
		params = append(params, "body")
	}
	return params
}

func generateRouter(versionedApis *spec.VersionedApis, packageName string, outPath string) *gen.TextFile {
	packageName = versionedPackage(versionedApis.Version, packageName)
	controllersPackage := versionedPackage(versionedApis.Version, "controllers")
	modelsPackage := versionedPackage(versionedApis.Version, "models")

	unit :=
		Unit(packageName).
			Import("javax.inject._").
			Import("play.api.mvc._").
			Import("play.api.routing._").
			Import("play.core.routing._").
			Import("spec.controllers.ParamsTypesBindings._").
			Import("spec.controllers.PlayParamsTypesBindings._").
			Import(controllersPackage + "._").
			Import(modelsPackage + "._")

	for _, api := range versionedApis.Apis {
		unit.AddDeclarations(generateApiRouter(versionedApis, api))
	}

	return &gen.TextFile{
		Path:    filepath.Join(outPath, fmt.Sprintf("%sRouters.scala", versionedApis.Version.PascalCase())),
		Content: unit.Code(),
	}
}

func generateMainRouter(versions []spec.VersionedApis, packageName string, outPath string) *gen.TextFile {
	unit :=
		Unit(packageName).
			Import("javax.inject._").
			Import("play.api.routing._")

	unit.AddDeclarations(generateSpecRouterMainClass(versions))

	return &gen.TextFile{
		Path:    filepath.Join(outPath, "SpecRouter.scala"),
		Content: unit.Code(),
	}
}

func generateSpecRouterMainClass(versions []spec.VersionedApis) *scala.ClassDeclaration {
	class :=
		Class(`SpecRouter`).Extends(`SimpleRouter`).
			Constructor(
				Constructor().
					Attribute(`Inject()`).
					AddParams(Dynamic(func(code *scala.WritableList) {
						for _, versionedApis := range versions {
							for _, api := range versionedApis.Apis {
								apiParamName := api.Name.CamelCase() + versionedApis.Version.PascalCase()
								apiTypeName := versionedTypeName(versionedApis.Version, routerType(api.Name))
								code.Add(Param(apiParamName, apiTypeName))
							}
						}
					})...),
			).
			Add(
				Line(`def routes: Router.Routes =`),
				Block(
					Line(`Seq(`),
					Block(Dynamic(func(code *scala.WritableList) {
						for _, versionedApis := range versions {
							for _, api := range versionedApis.Apis {
								apiParamName := api.Name.CamelCase() + versionedApis.Version.PascalCase()
								code.Add(Line(`%s.routes,`, apiParamName))
							}
						}
					})...),
					Line(`).reduce { (r1, r2) => r1.orElse(r2) }`),
				),
			)
	return class
}

func routerType(apiName spec.Name) string {
	return fmt.Sprintf("%sRouter", apiName.PascalCase())
}

func routeName(operationName spec.Name) string {
	return fmt.Sprintf("route%s", operationName.PascalCase())
}

func generateApiRouter(versionedApis *spec.VersionedApis, api spec.Api) *scala.ClassDeclaration {
	class :=
		Class(routerType(api.Name)).Extends(`SimpleRouter`).
			Constructor(Constructor().
				Attribute(`Inject()`).
				Param(`Action`, `DefaultActionBuilder`).
				Param(`controller`, controllerType(api.Name)),
			)

	for _, operation := range api.Operations {
		class.Add(
			Line(`lazy val %s = Route("%s", PathPattern(List(`, routeName(operation.Name), operation.Endpoint.Method),
			Block(Dynamic(func(code *scala.WritableList) {
				reminder := versionedApis.GetUrl() + operation.Endpoint.Url
				for _, param := range operation.Endpoint.UrlParams {
					parts := strings.Split(reminder, spec.UrlParamStr(param.Name.Source))
					code.Add(Line(`StaticPart("%s"),`, parts[0]))
					code.Add(Line(`DynamicPart("%s", """[^/]+""", true),`, param.Name.Source))
					reminder = parts[1]
				}
				if reminder != `` {
					code.Add(Line(`StaticPart("%s"),`, reminder))
				}
			})...),
			Line(`)))`),
		)
	}

	class.Add(Code(`def routes: Router.Routes = `))

	cases := Scope()
	for _, operation := range api.Operations {
		arguments := JoinParams(getControllerParams(operation))
		cases.Add(
			Line(`case %s(params@_) =>`, routeName(operation.Name)),
			Block(Dynamic(func(code *scala.WritableList) {
				if len(arguments) > 0 {
					code.Add(
						Line(`val arguments =`),
						Block(
							Code(`for `),
							Scope(Dynamic(func(code *scala.WritableList) {
								for _, p := range operation.Endpoint.UrlParams {
									code.Add(
										Line(`%s <- params.fromPath[%s]("%s").value`, p.Name.CamelCase(), ScalaType(&p.Type.Definition), p.Name.Source),
									)
								}
								for _, p := range operation.QueryParams {
									defaultValue := `None`
									if p.Default != nil {
										defaultValue = fmt.Sprintf(`Some(%s)`, DefaultValue(&p.Type.Definition, *p.Default))
									}
									code.Add(Line(`%s <- params.fromQuery[%s]("%s", %s).value`, p.Name.CamelCase(), ScalaType(&p.Type.Definition), p.Name.Source, defaultValue))
								}
							})...),
							Line(`yield (%s)`, arguments),
						),
						Code(`arguments match`),
						Scope(
							Line(`case Left(_) => Action { Results.BadRequest }`),
							Line(`case Right((%s)) => controller.%s(%s)`, arguments, controllerMethodName(operation), arguments),
						),
					)
				} else {
					code.Add(
						Line(`controller.%s(%s)`, controllerMethodName(operation), arguments),
					)
				}
			})...),
		)
	}
	class.Add(cases)

	return class
}

func getControllerParams(operation spec.NamedOperation) []string {
	params := []string{}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, param.Name.CamelCase())
	}
	if operation.QueryParams != nil {
		for _, param := range operation.QueryParams {
			params = append(params, param.Name.CamelCase())
		}
	}
	return params
}

func versionedPackage(version spec.Name, packageName string) string {
	if version.Source != "" {
		return fmt.Sprintf("%s.%s", packageName, version.FlatCase())
	} else {
		return packageName
	}
}

func versionedTypeName(version spec.Name, typeName string) string {
	if version.Source != "" {
		return fmt.Sprintf("%s.%s", version.FlatCase(), typeName)
	} else {
		return typeName
	}
}
