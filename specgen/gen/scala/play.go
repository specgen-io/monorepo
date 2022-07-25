package scala

import (
	"fmt"
	"github.com/specgen-io/specgen/spec/v2"
	"github.com/specgen-io/specgen/v2/gen/openapi"
	"github.com/specgen-io/specgen/v2/generator"
	"strings"
)

func GeneratePlayService(specification *spec.Spec, swaggerPath string, generatePath string, servicesPath string) *generator.Sources {
	rootPackage := NewPackage(generatePath, "", "")
	implRootPackage := NewPackage(servicesPath, "", "")

	appPackage := rootPackage.Subpackage("app")
	jsonPackage := rootPackage.Subpackage("json")
	paramsPackage := rootPackage.Subpackage("params")
	exceptionsPackage := rootPackage.Subpackage("exceptions")

	sources := generator.NewSources()

	playParamsFile := generatePlayParams(paramsPackage)
	scalaHttpStaticFile := generateStringParams(paramsPackage)
	exceptionsFile := generateExceptions(exceptionsPackage)
	sources.AddGenerated(scalaHttpStaticFile, playParamsFile, exceptionsFile)
	jsonHelpers := generateJson(jsonPackage)
	taggedUnion := generateTaggedUnion(jsonPackage)
	sources.AddGenerated(taggedUnion, jsonHelpers)

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
			apiController := generateApiController(&api, controllersPackage, apiPackage, modelsPackage, jsonPackage, paramsPackage, exceptionsPackage)
			apiRouter := generateApiRouter(&api, routersPackage, controllersPackage, modelsPackage, jsonPackage, paramsPackage)
			sources.AddGenerated(apiRouter, apiController, apiTrait)
		}
		versionModels := generateCirceModels(&version, modelsPackage, jsonPackage)
		sources.AddGenerated(versionModels)
	}
	routesFile := generateMainRouter(specification.Versions, appPackage)
	sources.AddGenerated(routesFile)

	if swaggerPath != "" {
		sources.AddGenerated(openapi.GenerateOpenapi(specification, swaggerPath))
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
				sources.AddScaffolded(apiImpl)
			}
		}
	}

	return sources
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

func generateApiTrait(api *spec.Api, thepackage, modelsPackage, servicesImplPackage Package) *generator.CodeFile {
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

	return &generator.CodeFile{
		Path:    thepackage.GetPath(fmt.Sprintf("%s.scala", apiClassType(api.Name))),
		Content: w.String(),
	}
}

func generateApiImpl(api spec.Api, thepackage, apiPackage, modelsPackage Package) *generator.CodeFile {
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

	return &generator.CodeFile{
		Path:    thepackage.GetPath(fmt.Sprintf("%s.scala", apiClassName)),
		Content: w.String(),
	}
}

func addParamsParsing(w *generator.Writer, params []spec.NamedParam, location string, paramsVar string, valuesVar string) {
	if params != nil && len(params) > 0 {
		w.Line(`val %s = new StringParamsReader(%s, %s)`, paramsVar, location, valuesVar)
		for _, param := range params {
			w.Line(`val %s = %s.%s`, param.Name.CamelCase(), paramsVar, addParamParsingCall(&param))
		}
	}
}

func addParamParsingCall(param *spec.NamedParam) string {
	paramBaseType := param.Type.Definition.BaseType()
	if paramBaseType.Node == spec.ArrayType {
		paramItemType := paramBaseType.Child
		if param.Type.Definition.IsNullable() {
			return fmt.Sprintf(`getOptionList[%s]("%s")`, ScalaType(paramItemType), param.Name.Source)
		} else {
			if param.Default != nil {
				return fmt.Sprintf(`getList[%s]("%s", Some(%s))`, ScalaType(paramItemType), param.Name.Source, DefaultValue(&param.Type.Definition, *param.Default))
			} else {
				return fmt.Sprintf(`getList[%s]("%s")`, ScalaType(paramItemType), param.Name.Source)
			}
		}
	} else {
		if param.Type.Definition.IsNullable() {
			return fmt.Sprintf(`getOption[%s]("%s")`, ScalaType(paramBaseType), param.Name.Source)
		} else {
			if param.Default != nil {
				return fmt.Sprintf(`get[%s]("%s", Some(%s))`, ScalaType(paramBaseType), param.Name.Source, DefaultValue(&param.Type.Definition, *param.Default))
			} else {
				return fmt.Sprintf(`get[%s]("%s")`, ScalaType(paramBaseType), param.Name.Source)
			}
		}
	}
}

func generateApiController(api *spec.Api, thepackage, apiPackage, modelsPackage, jsonPackage, paramsPackage, exceptionsPackage Package) *generator.CodeFile {
	w := NewScalaWriter()
	w.Line(`package %s`, thepackage.PackageName)
	w.EmptyLine()
	w.Line(`import javax.inject._`)
	w.Line(`import scala.util._`)
	w.Line(`import scala.concurrent._`)
	w.Line(`import play.api.mvc._`)
	w.Line(`import io.circe._`)
	w.Line(`import %s`, exceptionsPackage.PackageStar)
	w.Line(`import %s.ParamsTypesBindings._`, paramsPackage.PackageName)
	w.Line(`import %s`, paramsPackage.PackageStar)
	w.Line(`import %s.Jsoner`, jsonPackage.PackageName)
	w.Line(`import %s`, apiPackage.PackageStar)
	w.Line(`import %s`, modelsPackage.PackageStar)

	w.EmptyLine()
	w.Line(`@Singleton`)
	w.Line(`class %s @Inject()(api: %s, cc: ControllerComponents)(implicit ec: ExecutionContext) extends AbstractController(cc) {`, controllerType(api.Name), apiTraitType(api.Name))
	w.Indent()
	w.Lines(`
def error(e: Throwable): Result = {
  e match {
    case ex: ContentTypeMismatchException =>
      val validationError = ValidationError("Content-Type", "missing", Some(ex.getMessage))
      val body = BadRequestError("Failed to parse header", ErrorLocation.Header, Some(List(validationError)))
      new Status(400)(Jsoner.write(body)).as("application/json")
    case ex: DecodingFailure =>
      val path = CursorOp.opsToPath(ex.history)
      val validationError = ValidationError(if (path.startsWith(".")) path.substring(1, path.length) else path, "parsing_failed", Some(ex.getMessage))
      val body = BadRequestError("Failed to parse body", ErrorLocation.Body, Some(List(validationError)))
      new Status(400)(Jsoner.write(body)).as("application/json")
    case ex: ParamReadException =>
      val validationError = ValidationError(ex.paramName, ex.code, Some(ex.getMessage))
      val location = ex.location match {
        case _: ParamLocation.Query => ErrorLocation.Query
        case _: ParamLocation.Header => ErrorLocation.Header
      }
      val locationStr = ex.location match {
        case _: ParamLocation.Query => "query"
        case _: ParamLocation.Header => "header"
      }
      val body = BadRequestError(s"Failed to parse $locationStr", location, Some(List(validationError)))
      new Status(400)(Jsoner.write(body)).as("application/json")
    case _ => InternalServerError
  }
}
`)
	w.Lines(`
def checkContentType(request: RequestHeader, expectedContentType: String) = {
  if (request.contentType != Some(expectedContentType)) {
    throw new ContentTypeMismatchException(expectedContentType, request.contentType)
  }
}
`)
	for _, operation := range api.Operations {
		generateControllerMethod(w, &operation)
	}
	w.Unindent()
	w.Line(`}`)

	return &generator.CodeFile{
		Path:    thepackage.GetPath(fmt.Sprintf("%s.scala", controllerType(api.Name))),
		Content: w.String(),
	}
}

func generateControllerMethod(w *generator.Writer, operation *spec.NamedOperation) {
	params := getControllerParams(operation, true)

	payload := ``
	if operation.Body != nil {
		payload = `(parse.byteString)`
	}

	w.Line(`def %s(%s) = Action%s.async {`, operation.Name.CamelCase(), strings.Join(params, ", "), payload)
	w.Line(`  implicit request =>`)
	generateControllerMethodRequest(w.IndentedWith(2), operation)
	w.Line(`}`)
}

func generateControllerMethodRequest(w *generator.Writer, operation *spec.NamedOperation) {
	parseParams := getParsedOperationParams(operation)
	allParams := getOperationCallParams(operation)

	if len(parseParams) > 0 {
		w.Line(`val params = Try {`)
		addParamsParsing(w.Indented(), operation.HeaderParams, "ParamLocation.HEADER", "header", "request.headers.toMap")
		addParamsParsing(w.Indented(), operation.QueryParams, "ParamLocation.QUERY", "query", "request.queryString")
		genBodyParsing(w.Indented(), operation)
		w.Line(`  (%s)`, JoinParams(parseParams))
		w.Line(`}`)
		w.Line(`params match {`)
		w.Line(`  case Failure(ex) => Future { error(ex) }`)
		w.Line(`  case Success(params) =>`)
		w.Line(`    val (%s) = params`, JoinParams(parseParams))
		w.Line(`    val result = api.%s(%s)`, operation.Name.CamelCase(), JoinParams(allParams))
		genMatchResult(w.IndentedWith(2), operation, "result")
		w.Line(`}`)
	} else {
		w.Line(`val result = api.%s(%s)`, operation.Name.CamelCase(), JoinParams(allParams))
		genMatchResult(w, operation, "result")
	}
}

func genBodyParsing(w *generator.Writer, operation *spec.NamedOperation) {
	if operation.BodyIs(spec.BodyString) {
		w.Line(`checkContentType(request, "text/plain")`)
		w.Line(`val body = request.body.utf8String`)
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`checkContentType(request, "application/json")`)
		w.Line(`val body = Jsoner.readThrowing[%s](request.body.utf8String)`, ScalaType(&operation.Body.Type.Definition))
	}
}

func getPlayStatus(response *spec.Response) string {
	if response.BodyIs(spec.BodyEmpty) {
		return fmt.Sprintf(`new Status(%s)`, spec.HttpStatusCode(response.Name))
	}
	if response.BodyIs(spec.BodyString) {
		return fmt.Sprintf(`new Status(%s)(body)`, spec.HttpStatusCode(response.Name))
	} else {
		return fmt.Sprintf(`new Status(%s)(Jsoner.write(body)).as("application/json")`, spec.HttpStatusCode(response.Name))
	}
}

func genMatchResult(w *generator.Writer, operation *spec.NamedOperation, resultVarName string) {
	w.Line(`%s.map {`, resultVarName)
	w.Indent()
	if len(operation.Responses) == 1 {
		r := operation.Responses[0]
		if !r.Type.Definition.IsEmpty() {
			w.Line(`body => %s`, getPlayStatus(&r.Response))
		} else {
			w.Line(`_ => %s`, getPlayStatus(&r.Response))
		}
	} else {
		for _, r := range operation.Responses {
			if !r.Type.Definition.IsEmpty() {
				w.Line(`case %s.%s(body) => %s`, responseType(operation), r.Name.PascalCase(), getPlayStatus(&r.Response))
			} else {
				w.Line(`case %s.%s() => %s`, responseType(operation), r.Name.PascalCase(), getPlayStatus(&r.Response))
			}
		}
	}
	w.Unindent()
	w.Line(`} recover { case ex: Exception => error(ex) }`)
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
			params = append(params, param.Name.CamelCase())
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
	if operation.QueryParams != nil {
		for _, param := range operation.QueryParams {
			params = append(params, param.Name.CamelCase())
		}
	}
	if operation.Body != nil {
		params = append(params, "body")
	}
	return params
}

func generateApiRouter(api *spec.Api, thepackage, controllersPackage, modelsPackage, jsonPackage, paramsPackage Package) *generator.CodeFile {
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
	w.Line(`import %s`, jsonPackage.PackageStar)

	w.EmptyLine()
	generateApiRouterClass(w, api)

	return &generator.CodeFile{
		Path:    thepackage.GetPath(fmt.Sprintf("%s.scala", routerType(api.Name))),
		Content: w.String(),
	}
}

func generateMainRouter(versions []spec.Version, thepackage Package) *generator.CodeFile {
	w := NewScalaWriter()
	w.Line(`package %s`, thepackage.PackageName)

	w.EmptyLine()
	w.Line(`import javax.inject._`)
	w.Line(`import play.api.routing._`)

	generateSpecRouterMainClass(w, versions)

	return &generator.CodeFile{
		Path:    thepackage.GetPath("SpecRouter.scala"),
		Content: w.String(),
	}
}

func generateSpecRouterMainClass(w *generator.Writer, versions []spec.Version) {
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

func generateApiRouterClass(w *generator.Writer, api *spec.Api) {
	w.Line(`class %s @Inject()(Action: DefaultActionBuilder, controller: %s) extends SimpleRouter {`, routerType(api.Name), controllerType(api.Name))

	for _, operation := range api.Operations {
		w.Line(`  lazy val %s = Route("%s", PathPattern(List(`, routeName(operation.Name), operation.Endpoint.Method)
		if operation.Api.Http.GetUrl() != "" {
			w.Line(`    StaticPart("%s"),`, operation.Api.Http.GetUrl())
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
		controllerParams := getControllerParams(&operation, false)
		controllerParamsStr := strings.Join(controllerParams, ", ")
		w.Line(`    case %s(params@_) =>`, routeName(operation.Name))
		if len(controllerParams) > 0 {
			w.Line(`      val arguments =`)
			w.Line(`        for {`)
			for _, p := range operation.Endpoint.UrlParams {
				w.Line(`          %s <- params.fromPath[%s]("%s").value`, p.Name.CamelCase(), ScalaType(&p.Type.Definition), p.Name.Source)
			}
			//for _, p := range operation.QueryParams {
			//	defaultValue := `None`
			//	if p.Default != nil {
			//		defaultValue = fmt.Sprintf(`Some(%s)`, DefaultValue(&p.Type.Definition, *p.Default))
			//	}
			//	w.Line(`          %s <- params.fromQuery[%s]("%s", %s).value`, p.Name.CamelCase(), ScalaType(&p.Type.Definition), p.Name.Source, defaultValue)
			//}
			w.Line(`        }`)
			w.Line(`        yield (%s)`, controllerParamsStr)
			w.Line(`      arguments match{`)
			w.Line(`        case Left(_) => Action { Results.NotFound(Jsoner.write(NotFoundError("Failed to parse url parameters"))).as("application/json") }`)
			w.Line(`        case Right((%s)) => controller.%s(%s)`, controllerParamsStr, controllerMethodName(&operation), controllerParamsStr)
			w.Line(`      }`)
		} else {
			w.Line(`      controller.%s(%s)`, controllerMethodName(&operation), controllerParamsStr)
		}
	}
	w.Line(`  }`)
	w.Line(`}`)
}

func getControllerParams(operation *spec.NamedOperation, withTypes bool) []string {
	params := []string{}
	for _, p := range operation.Endpoint.UrlParams {
		param := p.Name.CamelCase()
		if withTypes {
			param += ": " + ScalaType(&p.Type.Definition)
		}
		params = append(params, param)
	}
	//if operation.QueryParams != nil {
	//	for _, p := range operation.QueryParams {
	//		param := p.Name.CamelCase()
	//		if withTypes {
	//			param += ": " + ScalaType(&p.Type.Definition)
	//		}
	//		params = append(params, p)
	//	}
	//}
	return params
}

func routerTypeName(api *spec.Api) string {
	typeName := fmt.Sprintf(`routers.%s`, routerType(api.Name))
	if api.Http.Version.Version.Source != "" {
		typeName = fmt.Sprintf(`%s.%s`, api.Http.Version.Version.FlatCase(), typeName)
	}
	return typeName
}

func generatePlayParams(thepackage Package) *generator.CodeFile {
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
	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thepackage.PackageName})
	return &generator.CodeFile{
		Path:    thepackage.GetPath("PlayParamsTypesBindings.scala"),
		Content: strings.TrimSpace(code),
	}
}
