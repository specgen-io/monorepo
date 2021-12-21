package genscala

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/genopenapi"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func GeneratePlayService(specification *spec.Spec, swaggerPath string, generatePath string, servicesPath string) (err error) {
	modelsPackage := NewPackage(generatePath, "", "models")
	controllersPackage := NewPackage(generatePath, "", "controllers")
	servicesPackage := NewPackage(generatePath, "", "services")
	appPackage := NewPackage(generatePath, "", "app")
	paramsPackage := controllersPackage

	sourcesOverwrite := []gen.TextFile{}
	sourcesScaffold := []gen.TextFile{}

	for _, version := range specification.Versions {
		versionServicesPackage := servicesPackage.Subpackage(version.Version.FlatCase())
		versionControllersPackage := controllersPackage.Subpackage(version.Version.FlatCase())
		versionModelsPackage := modelsPackage.Subpackage(version.Version.FlatCase())
		versionAppPackage := appPackage.Subpackage(version.Version.FlatCase())
		apisSourceManaged := generateApis(&version, versionServicesPackage, versionModelsPackage)
		sourcesOverwrite = append(sourcesOverwrite, apisSourceManaged...)
		apiControllerFile := generateApiControllers(&version, versionControllersPackage, versionServicesPackage, versionModelsPackage, modelsPackage, paramsPackage)
		sourcesOverwrite = append(sourcesOverwrite, *apiControllerFile)
		apiRoutersFile := generateRouter(&version, versionAppPackage, versionControllersPackage, versionModelsPackage, paramsPackage)
		sourcesOverwrite = append(sourcesOverwrite, *apiRoutersFile)
	}
	routesFile := generateMainRouter(specification.Versions, appPackage)
	sourcesOverwrite = append(sourcesOverwrite, *routesFile)

	scalaPlayParamsFile := generatePlayParams(paramsPackage)
	sourcesOverwrite = append(sourcesOverwrite, *scalaPlayParamsFile)
	scalaHttpStaticFile := generateStringParams(paramsPackage)
	sourcesOverwrite = append(sourcesOverwrite, *scalaHttpStaticFile)

	modelsFiles := GenerateCirceModels(specification, modelsPackage)
	sourcesOverwrite = append(sourcesOverwrite, modelsFiles...)

	if swaggerPath != "" {
		sourcesOverwrite = append(sourcesOverwrite, *genopenapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		servicesImplPackage := NewPackage(servicesPath, "services", "")
		for _, version := range specification.Versions {
			versionServicesImplPackage := servicesImplPackage.Subpackage(version.Version.FlatCase())
			versionModelsPackage := modelsPackage.Subpackage(version.Version.FlatCase())
			servicesSource := generateApisServices(&version, versionServicesImplPackage, versionModelsPackage)
			sourcesScaffold = append(sourcesScaffold, servicesSource...)
		}
	}

	err = gen.WriteFiles(sourcesScaffold, false)
	if err != nil {
		return
	}

	err = gen.WriteFiles(sourcesOverwrite, true)
	if err != nil {
		return
	}

	return
}

func generateApis(version *spec.Version, thepackage, modelsPackage Package) []gen.TextFile {
	sourceManaged := []gen.TextFile{}
	for _, api := range version.Http.Apis {
		apiTraitFile := generateApiInterface(&api, thepackage, modelsPackage)
		sourceManaged = append(sourceManaged, *apiTraitFile)
	}
	return sourceManaged
}

func generateApisServices(version *spec.Version, thepackage, modelsPackage Package) []gen.TextFile {
	source := []gen.TextFile{}
	for _, api := range version.Http.Apis {
		apiClassFile := generateApiClass(api, thepackage, modelsPackage)
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

func controllerMethodName(operation spec.NamedOperation) string {
	return operation.Name.CamelCase()
}

func operationSignature(operation spec.NamedOperation) string {
	params := []string{}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf(`%s: %s`, param.Name.CamelCase(), ScalaType(&param.Type.Definition)))
	}
	if operation.Body != nil {
		params = append(params, fmt.Sprintf(`body: %s`, ScalaType(&operation.Body.Type.Definition)))
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, fmt.Sprintf(`%s: %s`, param.Name.CamelCase(), ScalaType(&param.Type.Definition)))
	}
	for _, param := range operation.QueryParams {
		params = append(params, fmt.Sprintf(`%s: %s`, param.Name.CamelCase(), ScalaType(&param.Type.Definition)))
	}

	return fmt.Sprintf(`def %s(%s): Future[%s]`, controllerMethodName(operation), JoinParams(params), responseType(operation))
}

func generateApiInterface(api *spec.Api, thepackage, modelsPackage Package) *gen.TextFile {
	w := NewScalaWriter()
	w.Line(`package %s`, thepackage.PackageName)
	w.EmptyLine()
	w.Line(`import com.google.inject.ImplementedBy`)
	w.Line(`import scala.concurrent.Future`)
	w.Line(`import %s`, modelsPackage.PackageStar)

	apiTraitName := apiTraitType(api.Name)

	w.EmptyLine()
	w.Line(`@ImplementedBy(classOf[%s])`, apiClassType(api.Name))
	w.Line(`trait %s {`, apiTraitName)
	w.Line(`  import %s._`, apiTraitName)
	for _, operation := range api.Operations {
		w.Line(`  %s`, operationSignature(operation))
	}
	w.Line(`}`)

	w.EmptyLine()
	generateApiInterfaceResponse(w, api, apiTraitName)

	return &gen.TextFile{
		Path:    thepackage.GetPath(fmt.Sprintf("%s.scala", apiClassType(api.Name))),
		Content: w.String(),
	}
}

func generateApiClass(api spec.Api, thepackage, modelsPackage Package) *gen.TextFile {
	w := NewScalaWriter()
	w.Line(`package %s`, thepackage.PackageName)

	w.EmptyLine()
	w.Line(`import javax.inject._`)
	w.Line(`import scala.concurrent._`)
	w.Line(`import %s`, modelsPackage.PackageStar)

	apiClassName := apiClassType(api.Name)
	apiTraitName := apiTraitType(api.Name)

	w.EmptyLine()
	w.Line(`@Singleton`)
	w.Line(`class %s @Inject()()(implicit ec: ExecutionContext) extends %s {`, apiClassName, apiTraitName)
	w.Line(`  import %s._`, apiTraitName)

	for _, operation := range api.Operations {
		w.Line(`  override %s = Future { ??? }`, operationSignature(operation))
	}

	w.Line(`}`)

	return &gen.TextFile{
		Path:    thepackage.GetPath(fmt.Sprintf("%s.scala", apiClassName)),
		Content: w.String(),
	}
}

func addParamsParsing(w *gen.Writer, params []spec.NamedParam, paramsName string, values string) {
	if params != nil && len(params) > 0 {
		w.Line(`val %s = new StringParamsReader(%s)`, paramsName, values)
		for _, param := range params {
			paramBaseType := param.Type.Definition.BaseType()
			paramLine := fmt.Sprintf(`val %s = %s.read[%s]("%s")`, param.Name.CamelCase(), paramsName, ScalaType(paramBaseType), param.Name.Source)
			if paramBaseType.Node == spec.ArrayType {
				paramItemType := paramBaseType.Child
				paramLine = fmt.Sprintf(`val %s = %s.readList[%s]("%s")`, param.Name.CamelCase(), paramsName, ScalaType(paramItemType), param.Name.Source)
			}
			if !param.Type.Definition.IsNullable() {
				if param.Default != nil {
					paramLine += fmt.Sprintf(`.getOrElse(%s)`, DefaultValue(&param.Type.Definition, *param.Default))
				} else {
					paramLine += fmt.Sprintf(`.get`)
				}
			}
			w.Line(paramLine)
		}
	}
}

func generateApiControllers(version *spec.Version, thepackage, servicesPackage, modelsPackage, jsonPackage, paramsPackage Package) *gen.TextFile {
	w := NewScalaWriter()
	w.Line(`package %s`, thepackage.PackageName)
	w.EmptyLine()
	w.Line(`import javax.inject._`)
	w.Line(`import scala.util._`)
	w.Line(`import scala.concurrent._`)
	w.Line(`import play.api.mvc._`)
	w.Line(`import %s.ParamsTypesBindings._`, paramsPackage.PackageName)
	w.Line(`import %s.Jsoner`, jsonPackage.PackageName)
	w.Line(`import %s`, servicesPackage.PackageStar)
	w.Line(`import %s`, modelsPackage.PackageStar)

	for _, api := range version.Http.Apis {
		w.EmptyLine()
		w.Line(`@Singleton`)
		w.Line(`class %s @Inject()(api: %s, cc: ControllerComponents)(implicit ec: ExecutionContext) extends AbstractController(cc) {`, controllerType(api.Name), apiTraitType(api.Name))
		w.Line(`  import %s._`, apiTraitType(api.Name))
		for _, operation := range api.Operations {
			generateControllerMethod(w.Indented(), operation)
		}
		w.Line(`}`)
	}

	return &gen.TextFile{
		Path:    thepackage.GetPath("Controllers.scala"),
		Content: w.String(),
	}
}

func generateControllerMethod(w *gen.Writer, operation spec.NamedOperation) {
	parseParams := getParsedOperationParams(operation)
	allParams := getOperationCallParams(operation)

	params := []string{}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, fmt.Sprintf(`%s: %s`, param.Name.CamelCase(), ScalaType(&param.Type.Definition)))
	}
	for _, param := range operation.QueryParams {
		params = append(params, fmt.Sprintf(`%s: %s`, param.Name.Source, ScalaType(&param.Type.Definition)))
	}

	payload := ``
	if operation.Body != nil {
		payload = `(parse.byteString)`
	}

	w.Line(`def %s(%s) = Action%s.async {`, operation.Name.CamelCase(), JoinParams(params), payload)
	w.Line(`  implicit request =>`)
	w.IndentWith(2)
	if len(parseParams) > 0 {
		w.Line(`val params = Try {`)
		addParamsParsing(w.Indented(), operation.HeaderParams, "header", "request.headers.toMap")
		if operation.Body != nil {
			if operation.Body.Type.Definition.Plain == spec.TypeString {
				w.Line(`  val body = request.body.utf8String`)
			} else {
				w.Line(`  val body = Jsoner.readThrowing[%s](request.body.utf8String)`, ScalaType(&operation.Body.Type.Definition))
			}
		}
		w.Line(`  (%s)`, JoinParams(parseParams))
		w.Line(`}`)
		w.Line(`params match {`)
		w.Line(`  case Failure(ex) => Future { BadRequest }`)
		w.Line(`  case Success(params) => `)
		w.Line(`    val (%s) = params`, JoinParams(parseParams))
		w.Line(`    val result = api.%s(%s)`, operation.Name.CamelCase(), JoinParams(allParams))
		w.Line(`    val response = result.map {`)
		genResponseCases(w.IndentedWith(3), operation)
		w.Line(`    }`)
		w.Line(`    response.recover { case _: Exception => InternalServerError }`)
		w.Line(`}`)
	} else {
		w.Line(`val result = api.%s(%s)`, operation.Name.CamelCase(), JoinParams(allParams))
		w.Line(`val response = result.map {`)
		genResponseCases(w.Indented(), operation)
		w.Line(`}`)
		w.Line(`response.recover { case _: Exception => InternalServerError }`)
	}
	w.UnindentWith(2)
	w.Line(`}`)
}

func genResponseCases(w *gen.Writer, operation spec.NamedOperation) {
	for _, r := range operation.Responses {
		if !r.Type.Definition.IsEmpty() {
			body := `body`
			if r.Type.Definition.Plain != spec.TypeString {
				body = `Jsoner.write(body)`
			}
			w.Line(`case %s.%s(body) => new Status(%s)(%s)`, responseType(operation), r.Name.PascalCase(), spec.HttpStatusCode(r.Name), body)
		} else {
			w.Line(`case %s.%s() => new Status(%s)`, responseType(operation), r.Name.PascalCase(), spec.HttpStatusCode(r.Name))
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

func generateRouter(version *spec.Version, thepackage, controllersPackage, modelsPackage, paramsPackage Package) *gen.TextFile {
	w := NewScalaWriter()
	w.Line(`package %s`, thepackage.PackageName)

	w.EmptyLine()
	w.Line(`import javax.inject._`)
	w.Line(`import play.api.mvc._`)
	w.Line(`import play.api.routing._`)
	w.Line(`import play.core.routing._`)
	w.Line(`import %s.ParamsTypesBindings._`, paramsPackage.PackageName)
	w.Line(`import %s.PlayParamsTypesBindings._`, paramsPackage.PackageName)
	w.Line(`import %s`, controllersPackage.PackageStar)
	w.Line(`import %s`, modelsPackage.PackageStar)

	for _, api := range version.Http.Apis {
		w.EmptyLine()
		generateApiRouter(w, api)
	}

	return &gen.TextFile{
		Path:    thepackage.GetPath("Routers.scala"),
		Content: w.String(),
	}
}

func generateMainRouter(versions []spec.Version, thepackage Package) *gen.TextFile {
	w := NewScalaWriter()
	w.Line(`package %s`, thepackage.PackageName)

	w.EmptyLine()
	w.Line(`import javax.inject._`)
	w.Line(`import play.api.routing._`)

	generateSpecRouterMainClass(w, versions)

	return &gen.TextFile{
		Path:    thepackage.GetPath("SpecRouter.scala"),
		Content: w.String(),
	}
}

func generateSpecRouterMainClass(w *gen.Writer, versions []spec.Version) {
	params := []string{}
	for _, version := range versions {
		for _, api := range version.Http.Apis {
			apiParamName := api.Name.CamelCase() + version.Version.PascalCase()
			apiTypeName := versionedTypeName(version.Version, routerType(api.Name))
			params = append(params, fmt.Sprintf(`%s: %s`, apiParamName, apiTypeName))
		}
	}

	w.EmptyLine()
	w.Line(`class SpecRouter @Inject()(%s) extends SimpleRouter {`, JoinParams(params))

	w.Line(`  def routes: Router.Routes =`)
	w.Line(`    Seq(`)
	for _, version := range versions {
		for _, api := range version.Http.Apis {
			apiParamName := api.Name.CamelCase() + version.Version.PascalCase()
			w.Line(`      %s.routes,`, apiParamName)
		}
	}
	w.Line(`    ).reduce { (r1, r2) => r1.orElse(r2) }`)
	w.Line(`}`)
}

func routerType(apiName spec.Name) string {
	return fmt.Sprintf("%sRouter", apiName.PascalCase())
}

func routeName(operationName spec.Name) string {
	return fmt.Sprintf("route%s", operationName.PascalCase())
}

func generateApiRouter(w *gen.Writer, api spec.Api) {
	w.Line(`class %s @Inject()(Action: DefaultActionBuilder, controller: %s) extends SimpleRouter {`, routerType(api.Name), controllerType(api.Name))

	for _, operation := range api.Operations {
		w.Line(`  lazy val %s = Route("%s", PathPattern(List(`, routeName(operation.Name), operation.Endpoint.Method)
		if operation.Api.Apis.GetUrl() != "" {
			w.Line(`    StaticPart("%s"),`, operation.Api.Apis.GetUrl())
		}
		for _, part := range operation.Endpoint.UrlParts {
			if part.Param != nil {
				w.Line(`    DynamicPart("%s", """[^/]+""", true),`, part.Param.Name.Source)
			} else {
				w.Line(`    StaticPart("%s"),`, part.Part)
			}
		}
		w.Line(`  )))`)
	}

	w.Line(`  def routes: Router.Routes = {`)
	for _, operation := range api.Operations {
		arguments := JoinParams(getControllerParams(operation))
		w.Line(`    case %s(params@_) =>`, routeName(operation.Name))
		if len(arguments) > 0 {
			w.Line(`      val arguments =`)
			w.Line(`        for {`)
			for _, p := range operation.Endpoint.UrlParams {
				w.Line(`          %s <- params.fromPath[%s]("%s").value`, p.Name.CamelCase(), ScalaType(&p.Type.Definition), p.Name.Source)
			}
			for _, p := range operation.QueryParams {
				defaultValue := `None`
				if p.Default != nil {
					defaultValue = fmt.Sprintf(`Some(%s)`, DefaultValue(&p.Type.Definition, *p.Default))
				}
				w.Line(`          %s <- params.fromQuery[%s]("%s", %s).value`, p.Name.CamelCase(), ScalaType(&p.Type.Definition), p.Name.Source, defaultValue)
			}
			w.Line(`        }`)
			w.Line(`        yield (%s)`, arguments)
			w.Line(`      arguments match{`)
			w.Line(`        case Left(_) => Action { Results.BadRequest }`)
			w.Line(`        case Right((%s)) => controller.%s(%s)`, arguments, controllerMethodName(operation), arguments)
			w.Line(`      }`)
		} else {
			w.Line(`      controller.%s(%s)`, controllerMethodName(operation), arguments)
		}
	}
	w.Line(`  }`)
	w.Line(`}`)
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

func versionedTypeName(version spec.Name, typeName string) string {
	if version.Source != "" {
		return fmt.Sprintf("%s.%s", version.FlatCase(), typeName)
	} else {
		return typeName
	}
}

func generatePlayParams(thepackage Package) *gen.TextFile {
	code := `
package [[.PackageName]]

import play.api.mvc.{PathBindable, QueryStringBindable}

object PlayParamsTypesBindings {
  implicit def queryBindableParser[T](implicit stringBinder: QueryStringBindable[String], codec: Codec[T]): QueryStringBindable[T] = new QueryStringBindable[T] {
    override def bind(key: String, params: Map[String, Seq[String]]): Option[Either[String, T]] =
      for {
        valueStr <- stringBinder.bind(key, params)
      } yield {
        valueStr match {
          case Right(valueStr) =>
            try {
              Right(codec.decode(valueStr))
            } catch {
              case t: Throwable => Left(s"Unable to bind from key: $key, error: ${t.getMessage}")
            }
          case _ => Left(s"Unable to bind from key: $key")
        }
      }

    override def unbind(key: String, value: T): String = stringBinder.unbind(key, codec.encode(value))
  }

  implicit def pathBindableParser[T](implicit stringBinder: PathBindable[String], codec: Codec[T]): PathBindable[T] = new PathBindable[T] {
    override def bind(key: String, value: String): Either[String, T] = {
      val valueStr = stringBinder.bind(key, value)
      valueStr match {
        case Right(valueStr) =>
          try {
            Right(codec.decode(valueStr))
          } catch {
            case t: Throwable => Left(s"Unable to bind from key: $key, error: ${t.getMessage}")
          }
        case _ => Left(s"Unable to bind from key: $key")
      }
    }

    override def unbind(key: String, value: T): String = stringBinder.unbind(key, codec.encode(value))
  }
}`
	code, _ = gen.ExecuteTemplate(code, struct{ PackageName string }{thepackage.PackageName})
	return &gen.TextFile{
		Path:    thepackage.GetPath("PlayParamsTypesBindings.scala"),
		Content: strings.TrimSpace(code),
	}
}
