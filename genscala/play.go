package genscala

import (
    "github.com/ModaOperandi/spec"
    "github.com/vsapronov/gopoetry/scala"
    "path/filepath"
    "specgen/gen"
    "specgen/genopenapi"
    "specgen/static"
    "strings"
)

func GeneratePlayService(serviceFile string, swaggerPath string, generatePath string, servicesPath string, routesPath string) (err error) {
    specification, err := spec.ReadSpec(serviceFile)
    if err != nil {
        return
    }

    modelsPackage := modelsPackage(specification)
    controllersPackage := controllersPackage(specification)
    servicesPackage := servicesPackage(specification)

    scalaCirceFiles, err := static.RenderTemplate("scala-circe", generatePath, static.ScalaStaticCode{ PackageName: modelsPackage })
    if err != nil {
        return
    }
    scalaResponseFiles, err := static.RenderTemplate("scala-response", generatePath, static.ScalaStaticCode{ PackageName: servicesPackage })
    if err != nil {
        return
    }
    scalaPlayStaticFiles, err := static.RenderTemplate("scala-play", generatePath, static.ScalaStaticCode{ PackageName: controllersPackage })
    if err != nil {
        return
    }
    scalaHttpStaticFiles, err := static.RenderTemplate("scala-http", generatePath, static.ScalaStaticCode{ PackageName: controllersPackage })
    if err != nil {
        return
    }

    modelsFile := GenerateCirceModels(specification, modelsPackage, generatePath)

    sourceManaged := []gen.TextFile{ *modelsFile }
    sourceManaged = append(sourceManaged, scalaPlayStaticFiles...)
    sourceManaged = append(sourceManaged, scalaHttpStaticFiles...)
    sourceManaged = append(sourceManaged, scalaResponseFiles...)
    sourceManaged = append(sourceManaged, scalaCirceFiles...)

    apis := specification.Apis

    for _, api := range apis {
        apiTraitFile := generateApiInterface(api, servicesPackage, generatePath)
        apiControllerFile := generateApiController(api, controllersPackage, generatePath)
        sourceManaged = append(sourceManaged, *apiTraitFile, *apiControllerFile)
    }

    source := []gen.TextFile{}
    for _, api := range apis {
        apiClassFile := generateApiClass(api, servicesPackage, servicesPath)
        source = append(source, *apiClassFile)
    }

    routesFile := generateRoutesFile(specification, filepath.Join(routesPath, "routes"))
    resource := []gen.TextFile{*routesFile}

    genopenapi.GenerateSpecification(serviceFile, filepath.Join(swaggerPath, "swagger.yaml"))

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

func generateRoutesFile(specification *spec.Spec, routesPath string) *gen.TextFile {
    openApiRoutes := generateCommonRoutes()
    controllersRoutes := generateControllersRoutes(specification)
    routes := controllersRoutes + openApiRoutes
    routesFile := &gen.TextFile{ Path:    routesPath, Content: routes }
    return routesFile
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
    method := Def(operation.Name.CamelCase()).Returns(returnType)
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
    unit := Unit(packageName)
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
    apiTrait := Trait(apiTraitName).Attribute("ImplementedBy(classOf[" + apiClassType(api.Name) + "])")
    apiTrait.Add(Import(apiTraitName + "._"))
    for _, operation := range api.Operations {
        apiTrait.Add(operationSignature(operation))
    }
    return apiTrait
}

func generateApiClass(api spec.Api, packageName string, outPath string) *gen.TextFile {
    unit := Unit(packageName)
    unit.
        Import("javax.inject._").
        Import("scala.concurrent._").
        Import("models._")

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
        method := operationSignature(operation).Override().BodyInline(Code("Future { ??? }"))
        class.Add(method)
    }

    unit.AddDeclarations(class)

    return &gen.TextFile{
        Path:    filepath.Join(outPath, apiClassName+".scala"),
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

func generateApiController(api spec.Api, packageName string, outPath string) *gen.TextFile {
    unit := Unit(packageName)

    unit.
        Import("javax.inject._").
        Import("scala.util._").
        Import("scala.concurrent._").
        Import("play.api.mvc._").
        Import("controllers.PlayResultHelpers._").
        Import("models._").
        Import("json._").
        Import("services._")

    class :=
        Class(controllerType(api.Name)).Attribute("Singleton").Extends("AbstractController(cc)").
            Constructor(Constructor().
                Attribute("Inject()").
                Param("api", apiTraitType(api.Name)).
                Param("cc", "ControllerComponents").
                ImplicitParam("ec", "ExecutionContext"),
            )

    for _, operation := range api.Operations {
        class.Add(generateControllerMethod(operation))
    }

    unit.AddDeclarations(class)

    return &gen.TextFile{
        Path:    filepath.Join(outPath, controllerType(api.Name)+".scala"),
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
        If(operation.Body != nil).
            Then(Code("Action(parse.byteString).async ")).
            Else(Code("Action.async ")),
        Scope(
            Line("implicit request =>"),
            Block(
                If(len(parseParams) > 0).
                    Then(
                        Code("val params = Try "),
                        Scope(
                            addParamsParsing(operation.HeaderParams, "header", "request.headers.get"),
                            If(operation.Body != nil).
                                Then(Lazy(func() scala.Writable {
                                    return Line( "val body = Jsoner.read[%s](request.body.utf8String)", ScalaType(&operation.Body.Type.Definition))
                                })),
                            Line("(%s)", JoinParams(parseParams)),
                        ),
                        Code("params match "),
                        Scope(
                            Line("case Failure(ex) => Future { BadRequest }"),
                            Line("case Success(params) => "),
                            Block(
                                Line("val (%s) = params", JoinParams(parseParams)),
                                Line("val result = api.%s(%s)", operation.Name.CamelCase(), JoinParams(allParams)),
                                Line("result.map(_.toResult.toPlay).recover { case _: Exception => InternalServerError }"),
                            ),
                        ),
                    ).
                    Else(
                        Line("val result = api.%s(%s)", operation.Name.CamelCase(), JoinParams(allParams)),
                        Line("result.map(_.toResult.toPlay).recover { case _: Exception => InternalServerError }"),
                    ),
            ),
        ),
    )
    return method
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
            for _, param := range operation.QueryParams {
                paramDefinition := param.Name.Source+": "+ScalaType(&param.Type.Definition)
                if param.Default != nil {
                    paramDefinition += " ?= " + DefaultValue(&param.Type.Definition, *param.Default)
                }
                params = append(params, paramDefinition)
            }
            route := tail(operation.Endpoint.Method, 8) + " " + tail(routeUrl(operation), routeLength) + "   " + controllerEndpoint + "(" + JoinParams(params) + ")\n"
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
