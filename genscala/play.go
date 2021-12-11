package genscala

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/genopenapi"
	"github.com/specgen-io/specgen/v2/spec"
	"path/filepath"
	"strings"
)

func GeneratePlayService(specification *spec.Spec, swaggerPath string, generatePath string, servicesPath string) (err error) {
	modelsPackage := "models"
	controllersPackage := "controllers"
	servicesPackage := "services"

	sourcesOverwrite := []gen.TextFile{}
	sourcesScaffold := []gen.TextFile{}

	for _, version := range specification.Versions {
		apisSourceManaged := generateApis(&version, servicesPackage, controllersPackage, generatePath)
		sourcesOverwrite = append(sourcesOverwrite, apisSourceManaged...)
		apiControllerFile := generateApiControllers(&version, controllersPackage, generatePath)
		sourcesOverwrite = append(sourcesOverwrite, *apiControllerFile)
		apiRoutersFile := generateRouter(&version, "app", generatePath)
		sourcesOverwrite = append(sourcesOverwrite, *apiRoutersFile)
	}
	routesFile := generateMainRouter(specification.Versions, "app", generatePath)
	sourcesOverwrite = append(sourcesOverwrite, *routesFile)

	scalaPlayParamsFile := generatePlayParams("spec.controllers", filepath.Join(generatePath, "PlayParamsTypesBindings.scala"))
	sourcesOverwrite = append(sourcesOverwrite, *scalaPlayParamsFile)
	scalaHttpStaticFile := generateStringParams("spec.controllers", filepath.Join(generatePath, "StringParams.scala"))
	sourcesOverwrite = append(sourcesOverwrite, *scalaHttpStaticFile)

	scalaCirceFile := generateJson("spec.models", filepath.Join(generatePath, "Json.scala"))
	sourcesOverwrite = append(sourcesOverwrite, *scalaCirceFile)
	modelsFiles := GenerateCirceModels(specification, modelsPackage, generatePath)
	sourcesOverwrite = append(sourcesOverwrite, modelsFiles...)

	if swaggerPath != "" {
		sourcesOverwrite = append(sourcesOverwrite, *genopenapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		for _, version := range specification.Versions {
			servicesSource := generateApisServices(&version, servicesPackage, servicesPath)
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

func generateApis(version *spec.Version, servicesPackage string, controllersPackage string, generatePath string) []gen.TextFile {
	sourceManaged := []gen.TextFile{}
	for _, api := range version.Http.Apis {
		apiTraitFile := generateApiInterface(api, servicesPackage, generatePath)
		sourceManaged = append(sourceManaged, *apiTraitFile)
	}
	return sourceManaged
}

func generateApisServices(version *spec.Version, servicesPackage string, servicesPath string) []gen.TextFile {
	source := []gen.TextFile{}
	for _, api := range version.Http.Apis {
		apiClassFile := generateApiClass(api, servicesPackage, servicesPath)
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

func generateApiInterface(api spec.Api, packageName string, outPath string) *gen.TextFile {
	version := api.Apis.Version.Version

	w := NewScalaWriter()
	w.Line(`package %s`, versionedPackage(version, packageName))
	modelsPackage := versionedPackage(version, "models")

	w.EmptyLine()
	w.Line(`import com.google.inject.ImplementedBy`)
	w.Line(`import scala.concurrent.Future`)
	w.Line(`import %s._`, modelsPackage)

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
		Path:    filepath.Join(outPath, fmt.Sprintf("%s%s.scala", apiTraitName, version.PascalCase())),
		Content: w.String(),
	}
}

func generateApiClass(api spec.Api, packageName string, outPath string) *gen.TextFile {
	version := api.Apis.Version.Version

	w := NewScalaWriter()
	w.Line(`package %s`, versionedPackage(version, packageName))

	modelsPackage := versionedPackage(version, "models")

	w.EmptyLine()
	w.Line(`import javax.inject._`)
	w.Line(`import scala.concurrent._`)
	w.Line(`import %s._`, modelsPackage)

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
		Path:    filepath.Join(outPath, version.FlatCase(), fmt.Sprintf("%s.scala", apiClassName)),
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

func generateApiControllers(version *spec.Version, packageName string, outPath string) *gen.TextFile {
	modelsPackage := versionedPackage(version.Version, "models")
	servicePackage := versionedPackage(version.Version, "services")
	w := NewScalaWriter()
	w.Line(`package %s`, versionedPackage(version.Version, packageName))
	w.EmptyLine()
	w.Line(`import javax.inject._`)
	w.Line(`import scala.util._`)
	w.Line(`import scala.concurrent._`)
	w.Line(`import play.api.mvc._`)
	w.Line(`import spec.controllers.ParamsTypesBindings._`)
	w.Line(`import spec.models.Jsoner`)
	w.Line(`import %s._`, servicePackage)
	w.Line(`import %s._`, modelsPackage)

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
		Path:    filepath.Join(outPath, fmt.Sprintf("%sControllers.scala", version.Version.PascalCase())),
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

func generateRouter(version *spec.Version, packageName string, outPath string) *gen.TextFile {
	packageName = versionedPackage(version.Version, packageName)
	controllersPackage := versionedPackage(version.Version, "controllers")
	modelsPackage := versionedPackage(version.Version, "models")

	w := NewScalaWriter()
	w.Line(`package %s`, packageName)

	w.EmptyLine()
	w.Line(`import javax.inject._`)
	w.Line(`import play.api.mvc._`)
	w.Line(`import play.api.routing._`)
	w.Line(`import play.core.routing._`)
	w.Line(`import spec.controllers.ParamsTypesBindings._`)
	w.Line(`import spec.controllers.PlayParamsTypesBindings._`)
	w.Line(`import %s._`, controllersPackage)
	w.Line(`import %s._`, modelsPackage)

	for _, api := range version.Http.Apis {
		w.EmptyLine()
		generateApiRouter(w, api)
	}

	return &gen.TextFile{
		Path:    filepath.Join(outPath, fmt.Sprintf("%sRouters.scala", version.Version.PascalCase())),
		Content: w.String(),
	}
}

func generateMainRouter(versions []spec.Version, packageName string, outPath string) *gen.TextFile {
	w := NewScalaWriter()
	w.Line(`package %s`, packageName)

	w.EmptyLine()
	w.Line(`import javax.inject._`)
	w.Line(`import play.api.routing._`)

	generateSpecRouterMainClass(w, versions)

	return &gen.TextFile{
		Path:    filepath.Join(outPath, "SpecRouter.scala"),
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
		reminder := operation.FullUrl()
		for _, param := range operation.Endpoint.UrlParams {
			parts := strings.Split(reminder, spec.UrlParamStr(param.Name.Source))
			w.Line(`    StaticPart("%s"),`, parts[0])
			w.Line(`    DynamicPart("%s", """[^/]+""", true),`, param.Name.Source)
			reminder = parts[1]
		}
		if reminder != `` {
			w.Line(`    StaticPart("%s"),`, reminder)
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

func generatePlayParams(packageName string, path string) *gen.TextFile {
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
	code, _ = gen.ExecuteTemplate(code, struct{ PackageName string }{packageName})
	return &gen.TextFile{path, strings.TrimSpace(code)}
}
