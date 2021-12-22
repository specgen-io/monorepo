package genscala

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/genopenapi"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func GeneratePlayService(specification *spec.Spec, swaggerPath string, generatePath string, servicesPath string) (err error) {
	rootPackage := NewPackage(generatePath, "", "")
	implRootPackage := NewPackage(servicesPath, "", "")

	appPackage := rootPackage.Subpackage("app")
	jsonPackage := rootPackage.Subpackage("json")
	paramsPackage := rootPackage.Subpackage("params")

	sourcesOverwrite := []gen.TextFile{}
	sourcesScaffold := []gen.TextFile{}

	scalaPlayParamsFile := generatePlayParams(paramsPackage)
	scalaHttpStaticFile := generateStringParams(paramsPackage)
	sourcesOverwrite = append(sourcesOverwrite, *scalaHttpStaticFile, *scalaPlayParamsFile)
	jsonHelpers := generateJson(jsonPackage)
	taggedUnion := generateTaggedUnion(jsonPackage)
	sourcesOverwrite = append(sourcesOverwrite, *taggedUnion, *jsonHelpers)

	for _, version := range specification.Versions {
		versionPackage := rootPackage.Subpackage(version.Version.FlatCase())

		routersPackage := versionPackage.Subpackage("routers")
		controllersPackage := versionPackage.Subpackage("controllers")
		servicesPackage := versionPackage.Subpackage("services")
		modelsPackage := versionPackage.Subpackage("models")

		servicesImplPackage := implRootPackage.Subpackage(version.Version.FlatCase()).Subpackage("services")
		for _, api := range version.Http.Apis {
			apiPackage := servicesPackage.Subpackage(api.Name.FlatCase())
			apiTrait := generateApiTrait(&api, apiPackage, modelsPackage, servicesImplPackage)
			apiController := generateApiController(&api, controllersPackage, apiPackage, modelsPackage, jsonPackage, paramsPackage)
			apiRouter := generateApiRouter(&api, routersPackage, controllersPackage, modelsPackage, paramsPackage)
			sourcesOverwrite = append(sourcesOverwrite, *apiRouter, *apiController, *apiTrait)
		}
		versionModels := generateCirceModels(&version, modelsPackage, jsonPackage)
		sourcesOverwrite = append(sourcesOverwrite, *versionModels)
	}
	routesFile := generateMainRouter(specification.Versions, appPackage)
	sourcesOverwrite = append(sourcesOverwrite, *routesFile)

	if swaggerPath != "" {
		sourcesOverwrite = append(sourcesOverwrite, *genopenapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		for _, version := range specification.Versions {
			servicesImplPackage := implRootPackage.Subpackage(version.Version.FlatCase()).Subpackage("services")
			versionPackage := rootPackage.Subpackage(version.Version.FlatCase())
			modelsPackage := versionPackage.Subpackage("models")
			servicesPackage := versionPackage.Subpackage("services")
			for _, api := range version.Http.Apis {
				apiPackage := servicesPackage.Subpackage(api.Name.FlatCase())
				apiImpl := generateApiImpl(api, servicesImplPackage, apiPackage, modelsPackage)
				sourcesScaffold = append(sourcesScaffold, *apiImpl)
			}
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

func controllerType(apiName spec.Name) string {
	return apiName.PascalCase() + "Controller"
}

func apiTraitType(apiName spec.Name) string {
	return "I" + apiName.PascalCase() + "Service"
}

func apiClassType(apiName spec.Name) string {
	return apiName.PascalCase() + "Service"
}

func controllerMethodName(operation *spec.NamedOperation) string {
	return operation.Name.CamelCase()
}

func operationSignature(operation *spec.NamedOperation) string {
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

func generateApiTrait(api *spec.Api, thepackage, modelsPackage, servicesImplPackage Package) *gen.TextFile {
	w := NewScalaWriter()
	w.Line(`package %s`, thepackage.PackageName)
	w.EmptyLine()
	w.Line(`import com.google.inject.ImplementedBy`)
	w.Line(`import scala.concurrent.Future`)
	w.Line(`import %s`, modelsPackage.PackageStar)

	apiTraitName := apiTraitType(api.Name)

	w.EmptyLine()
	w.Line(`@ImplementedBy(classOf[%s.%s])`, servicesImplPackage.PackageName, apiClassType(api.Name))
	w.Line(`trait %s {`, apiTraitName)
	for _, operation := range api.Operations {
		w.Line(`  %s`, operationSignature(&operation))
	}
	w.Line(`}`)

	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			w.EmptyLine()
			generateResponse(w, &operation)
		}
	}

	return &gen.TextFile{
		Path:    thepackage.GetPath(fmt.Sprintf("%s.scala", apiClassType(api.Name))),
		Content: w.String(),
	}
}

func generateApiImpl(api spec.Api, thepackage, apiPackage, modelsPackage Package) *gen.TextFile {
	w := NewScalaWriter()
	w.Line(`package %s`, thepackage.PackageName)

	w.EmptyLine()
	w.Line(`import javax.inject._`)
	w.Line(`import scala.concurrent._`)
	w.Line(`import %s`, apiPackage.PackageStar)
	w.Line(`import %s`, modelsPackage.PackageStar)

	apiClassName := apiClassType(api.Name)
	apiTraitName := apiTraitType(api.Name)

	w.EmptyLine()
	w.Line(`@Singleton`)
	w.Line(`class %s @Inject()()(implicit ec: ExecutionContext) extends %s {`, apiClassName, apiTraitName)
	for _, operation := range api.Operations {
		w.Line(`  override %s = Future { ??? }`, operationSignature(&operation))
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

func generateApiController(api *spec.Api, thepackage, apiPackage, modelsPackage, jsonPackage, paramsPackage Package) *gen.TextFile {
	w := NewScalaWriter()
	w.Line(`package %s`, thepackage.PackageName)
	w.EmptyLine()
	w.Line(`import javax.inject._`)
	w.Line(`import scala.util._`)
	w.Line(`import scala.concurrent._`)
	w.Line(`import play.api.mvc._`)
	w.Line(`import %s.ParamsTypesBindings._`, paramsPackage.PackageName)
	w.Line(`import %s.Jsoner`, jsonPackage.PackageName)
	w.Line(`import %s`, apiPackage.PackageStar)
	w.Line(`import %s`, modelsPackage.PackageStar)

	w.EmptyLine()
	w.Line(`@Singleton`)
	w.Line(`class %s @Inject()(api: %s, cc: ControllerComponents)(implicit ec: ExecutionContext) extends AbstractController(cc) {`, controllerType(api.Name), apiTraitType(api.Name))
	for _, operation := range api.Operations {
		generateControllerMethod(w.Indented(), &operation)
	}
	w.Line(`}`)

	return &gen.TextFile{
		Path:    thepackage.GetPath(fmt.Sprintf("%s.scala", controllerType(api.Name))),
		Content: w.String(),
	}
}

func generateControllerMethod(w *gen.Writer, operation *spec.NamedOperation) {
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

func genResponseCases(w *gen.Writer, operation *spec.NamedOperation) {
	if len(operation.Responses) == 1 {
		r := operation.Responses[0]
		if !r.Type.Definition.IsEmpty() {
			body := `body`
			if r.Type.Definition.Plain != spec.TypeString {
				body = `Jsoner.write(body)`
			}
			w.Line(`body => new Status(%s)(%s)`, spec.HttpStatusCode(r.Name), body)
		} else {
			w.Line(`_ => new Status(%s)`, spec.HttpStatusCode(r.Name))
		}
	} else {
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
}

func getOperationCallParams(operation *spec.NamedOperation) []string {
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

func getParsedOperationParams(operation *spec.NamedOperation) []string {
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

func generateApiRouter(api *spec.Api, thepackage, controllersPackage, modelsPackage, paramsPackage Package) *gen.TextFile {
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

	w.EmptyLine()
	generateApiRouterClass(w, api)

	return &gen.TextFile{
		Path:    thepackage.GetPath(fmt.Sprintf("%s.scala", routerType(api.Name))),
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
			params = append(params, fmt.Sprintf(`%s: %s`, apiParamName, routerTypeName(&api)))
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

func generateApiRouterClass(w *gen.Writer, api *spec.Api) {
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
		arguments := JoinParams(getControllerParams(&operation))
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
			w.Line(`        case Right((%s)) => controller.%s(%s)`, arguments, controllerMethodName(&operation), arguments)
			w.Line(`      }`)
		} else {
			w.Line(`      controller.%s(%s)`, controllerMethodName(&operation), arguments)
		}
	}
	w.Line(`  }`)
	w.Line(`}`)
}

func getControllerParams(operation *spec.NamedOperation) []string {
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

func routerTypeName(api *spec.Api) string {
	typeName := fmt.Sprintf(`routers.%s`, routerType(api.Name))
	if api.Apis.Version.Version.Source != "" {
		typeName = fmt.Sprintf(`%s.%s`, api.Apis.Version.Version.FlatCase(), typeName)
	}
	return typeName
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
