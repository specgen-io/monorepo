package service

import (
	"fmt"
	"generator"
	"github.com/pinzolo/casee"
	"kotlin/models"
	"kotlin/types"
	"kotlin/writer"
	"spec"
	"strings"
)

var Spring = "spring"

type SpringGenerator struct {
	Types    *types.Types
	Models   models.Generator
	Packages *Packages
}

func NewSpringGenerator(types *types.Types, models models.Generator, packages *Packages) *SpringGenerator {
	return &SpringGenerator{types, models, packages}
}

func (g *SpringGenerator) ServiceImplAnnotation(api *spec.Api) (annotationImport, annotation string) {
	return `org.springframework.stereotype.Service`, fmt.Sprintf(`Service("%s")`, versionServiceName(serviceName(api), api.InHttp.InVersion))
}

func (g *SpringGenerator) ServicesControllers(version *spec.Version) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *g.serviceController(&api))
	}
	return files
}

func (g *SpringGenerator) ServiceImports() []string {
	return []string{
		`org.apache.logging.log4j.*`,
		`org.springframework.beans.factory.annotation.Autowired`,
		`org.springframework.format.annotation.*`,
		`org.springframework.http.*`,
		`org.springframework.web.bind.annotation.*`,
	}
}

func (g *SpringGenerator) ExceptionController(responses *spec.ErrorResponses) *generator.CodeFile {
	w := writer.New(g.Packages.RootControllers, `ExceptionController`)
	w.Imports.Add(g.ServiceImports()...)
	w.Imports.Add(`javax.servlet.http.HttpServletRequest`)
	w.Imports.Add(`org.apache.tomcat.util.http.fileupload.FileUploadBase.CONTENT_TYPE`)
	w.Imports.PackageStar(g.Packages.Json)
	w.Imports.PackageStar(g.Packages.ErrorsModels)
	w.Imports.PackageStar(g.Packages.Errors)
	w.Line(`@ControllerAdvice`)
	w.Line(`class [[.ClassName]](@Autowired private val json: Json) {`)
	w.Line(`    private val logger = LogManager.getLogger([[.ClassName]]::class.java)`)
	w.EmptyLine()
	g.errorHandler(w.Indented(), *responses)
	w.Line(`}`)
	return w.ToCodeFile()
}

func (g *SpringGenerator) errorHandler(w *writer.Writer, errors spec.ErrorResponses) {
	notFoundError := errors.GetByStatusName(spec.HttpStatusNotFound)
	badRequestError := errors.GetByStatusName(spec.HttpStatusBadRequest)
	internalServerError := errors.GetByStatusName(spec.HttpStatusInternalServerError)
	w.Line(`@ExceptionHandler(Throwable::class)`)
	w.Line(`fun error(request: HttpServletRequest, exception: Throwable): ResponseEntity<String> {`)
	w.Line(`    val notFoundError = getNotFoundError(exception)`)
	w.Line(`    if (notFoundError != null) {`)
	g.processResponse(w.IndentedWith(2), &notFoundError.Response, "notFoundError")
	w.Line(`    }`)
	w.Line(`    val badRequestError = getBadRequestError(exception)`)
	w.Line(`    if (badRequestError != null) {`)
	g.processResponse(w.IndentedWith(2), &badRequestError.Response, "badRequestError")
	w.Line(`    }`)
	w.Line(`    val internalServerError = InternalServerError(exception.message ?: "Unknown error")`)
	g.processResponse(w.IndentedWith(1), &internalServerError.Response, "internalServerError")
	w.Line(`}`)
}

func (g *SpringGenerator) serviceController(api *spec.Api) *generator.CodeFile {
	w := writer.New(g.Packages.Controllers(api.InHttp.InVersion), controllerName(api))
	w.Imports.Add(g.ServiceImports()...)
	w.Imports.Add(`javax.servlet.http.HttpServletRequest`)
	w.Imports.Add(`org.apache.tomcat.util.http.fileupload.FileUploadBase.CONTENT_TYPE`)
	w.Imports.PackageStar(g.Packages.ContentType)
	w.Imports.PackageStar(g.Packages.Json)
	w.Imports.PackageStar(g.Packages.Models(api.InHttp.InVersion))
	w.Imports.PackageStar(g.Packages.ErrorsModels)
	w.Imports.PackageStar(g.Packages.ServicesApi(api))
	w.Imports.Add(g.Models.ModelsUsageImports()...)
	w.Imports.Add(g.Types.Imports()...)
	w.Line(`@RestController("%s")`, versionControllerName(controllerName(api), api.InHttp.InVersion))
	w.Line(`class [[.ClassName]](`)
	w.Line(`    @Autowired private val %s: %s,`, serviceVarName(api), serviceInterfaceName(api))
	w.Line(`    @Autowired private val json: Json`)
	w.Line(`) {`)
	w.Line(`    private val logger = LogManager.getLogger([[.ClassName]]::class.java)`)
	for _, operation := range api.Operations {
		w.EmptyLine()
		g.controllerMethod(w.Indented(), &operation)
	}
	w.Line(`}`)
	return w.ToCodeFile()
}

func (g *SpringGenerator) controllerMethod(w *writer.Writer, operation *spec.NamedOperation) {
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	w.Line(`@%sMapping("%s")`, casee.ToPascalCase(methodName), url)
	w.Line(`fun %s(%s): ResponseEntity<String> {`, controllerMethodName(operation), strings.Join(springMethodParams(operation, g.Types), ", "))
	w.Indent()
	w.Line(`logger.info("Received request, operationId: %s.%s, method: %s, url: %s")`, operation.InApi.Name.Source, operation.Name.Source, methodName, url)
	bodyStringVar := "bodyStr"
	if operation.Body.IsJson() {
		bodyStringVar += ".reader()"
	}
	g.parseBody(w, operation, bodyStringVar, "requestBody")
	serviceCall(w, operation, bodyStringVar, "requestBody", "result", true)
	g.processResponses(w, operation, "result")
	w.Unindent()
	w.Line(`}`)
}

func (g *SpringGenerator) parseBody(w *writer.Writer, operation *spec.NamedOperation, bodyStringVar, bodyJsonVar string) {
	if !operation.Body.IsEmpty() {
		w.Line(`checkContentType(request, %s)`, g.contentType(operation))
	}
	if operation.Body.IsJson() {
		typ := g.Types.Kotlin(&operation.Body.Type.Definition)
		w.Line(`val %s: %s = json.%s`, bodyJsonVar, typ, g.Models.ReadJson(bodyStringVar, &operation.Body.Type.Definition))
	}
}

func (g *SpringGenerator) contentType(operation *spec.NamedOperation) string {
	if operation.Body.IsEmpty() {
		return ""
	} else if operation.Body.IsText() {
		return `MediaType.TEXT_PLAIN`
	} else if operation.Body.IsJson() {
		return `MediaType.APPLICATION_JSON`
	} else if operation.Body.IsBodyFormData() {
		return `MediaType.MULTIPART_FORM_DATA`
	} else if operation.Body.IsBodyFormUrlEncoded() {
		return `MediaType.APPLICATION_FORM_URLENCODED`
	} else {
		panic(fmt.Sprintf("Unknown Contet Type"))
	}
}

func (g *SpringGenerator) processResponses(w *writer.Writer, operation *spec.NamedOperation, resultVarName string) {
	if len(operation.Responses) == 1 {
		g.processResponse(w, &operation.Responses[0].Response, resultVarName)
	}
	if len(operation.Responses) > 1 {
		for _, response := range operation.Responses {
			w.Line(`if (%s is %s.%s) {`, resultVarName, responseInterfaceName(operation), response.Name.PascalCase())
			g.processResponse(w.Indented(), &response.Response, getResponseBody(resultVarName))
			w.Line(`}`)
		}
		w.Line(`throw RuntimeException("Service implementation didn't return any value'")`)
	}
}

func (g *SpringGenerator) processResponse(w *writer.Writer, response *spec.Response, bodyVar string) {
	if response.Body.Is(spec.ResponseBodyEmpty) {
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s)`, response.Name.UpperCase())
		w.Line(`return ResponseEntity(HttpStatus.%s)`, response.Name.UpperCase())
	}
	if response.Body.Is(spec.ResponseBodyString) {
		w.Line(`val headers = HttpHeaders()`)
		w.Line(`headers.add(CONTENT_TYPE, "text/plain")`)
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s)`, response.Name.UpperCase())
		w.Line(`return ResponseEntity(%s, headers, HttpStatus.%s)`, bodyVar, response.Name.UpperCase())
	}
	if response.Body.Is(spec.ResponseBodyJson) {
		w.Line(`val bodyJson = json.%s`, g.Models.WriteJson(bodyVar, &response.Body.Type.Definition))
		w.Line(`val headers = HttpHeaders()`)
		w.Line(`headers.add(CONTENT_TYPE, "application/json")`)
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s)`, response.Name.UpperCase())
		w.Line(`return ResponseEntity(bodyJson, headers, HttpStatus.%s)`, response.Name.UpperCase())
	}
}

func (g *SpringGenerator) ContentType() []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *g.contentTypeMismatchException())
	files = append(files, *g.checkContentType())
	return files
}

func (g *SpringGenerator) contentTypeMismatchException() *generator.CodeFile {
	w := writer.New(g.Packages.ContentType, `ContentTypeMismatchException`)
	w.Lines(`
class ContentTypeMismatchException(expected: String, actual: String?) :
    RuntimeException(
        String.format(
            "Expected Content-Type header: '%s' was not provided, found: '%s'",
            expected,
            actual
        )
    )
`)
	return w.ToCodeFile()
}

func (g *SpringGenerator) checkContentType() *generator.CodeFile {
	w := writer.New(g.Packages.ContentType, `CheckContentType`)
	w.Lines(`
import javax.servlet.http.HttpServletRequest
import org.springframework.http.MediaType

fun checkContentType(request: HttpServletRequest, expectedContentType: MediaType) {
	val contentType = request.getHeader("Content-Type")
	if (contentType == null || !contentType.contains(expectedContentType.toString())) {
		throw ContentTypeMismatchException(expectedContentType.toString(), contentType)
	}
}
`)
	return w.ToCodeFile()
}

func (g *SpringGenerator) ErrorsHelpers() *generator.CodeFile {
	w := writer.New(g.Packages.Errors, `ErrorsHelpers`)
	w.Template(
		map[string]string{
			`ErrorsModelsPackage`: g.Packages.ErrorsModels.PackageName,
			`ContentTypePackage`:  g.Packages.ContentType.PackageName,
			`JsonPackage`:         g.Packages.Json.PackageName,
			`ErrorsPackage`:       g.Packages.Errors.PackageName,
		}, `
import org.springframework.web.bind.*
import org.springframework.web.bind.annotation.*
import org.springframework.web.method.annotation.MethodArgumentTypeMismatchException
import [[.ErrorsModelsPackage]].*
import [[.ContentTypePackage]].*
import [[.JsonPackage]].*
import [[.ErrorsPackage]].ValidationErrorsHelpers.extractValidationErrors

fun getNotFoundError(exception: Throwable?): NotFoundError? {
	if (exception is MethodArgumentTypeMismatchException) {
		if (exception.parameter.hasParameterAnnotation(PathVariable::class.java)) {
			return NotFoundError("Failed to parse url parameters")
		}
	}
	return null
}

fun getBadRequestError(exception: Throwable): BadRequestError? {
	if (exception is JsonParseException) {
		val errors = extractValidationErrors(exception)
		return BadRequestError("Failed to parse body", ErrorLocation.BODY, errors)
	}
	if (exception is ContentTypeMismatchException) {
		val error = ValidationError("Content-Type", "missing", exception.message)
		return BadRequestError("Failed to parse header", ErrorLocation.HEADER, listOf(error))
	}
	if (exception is MissingServletRequestParameterException) {
		val message = String.format("Failed to parse parameters")
		val validation = ValidationError(exception.parameterName, "missing", exception.message)
		return BadRequestError(message, ErrorLocation.PARAMETERS, listOf(validation))
	}
	if (exception is MethodArgumentTypeMismatchException) {
		val validation = ValidationError(exception.name, "parsing_failed", exception.message)
		if (exception.parameter.hasParameterAnnotation(RequestParam::class.java)) {
			return BadRequestError("Failed to parse parameters", ErrorLocation.PARAMETERS, listOf(validation))
		} else if (exception.parameter.hasParameterAnnotation(RequestHeader::class.java)) {
			return BadRequestError("Failed to parse header", ErrorLocation.HEADER, listOf(validation))
		}
	}
	if (exception is MissingRequestHeaderException) {
		val message = String.format("Failed to parse header")
		val validation = ValidationError(exception.headerName, "missing", exception.message)
		return BadRequestError(message, ErrorLocation.HEADER, listOf(validation))
	}
	return null
}
`)
	return w.ToCodeFile()
}

func springMethodParams(operation *spec.NamedOperation, types *types.Types) []string {
	methodParams := []string{"request: HttpServletRequest"}

	if operation.Body.IsText() || operation.Body.IsJson() {
		methodParams = append(methodParams, "@RequestBody bodyStr: String")
	}
	methodParams = append(methodParams, generateSpringMethodParam(operation.Body.FormData, "RequestParam", types)...)
	methodParams = append(methodParams, generateSpringMethodParam(operation.Body.FormUrlEncoded, "RequestParam", types)...)
	methodParams = append(methodParams, generateSpringMethodParam(operation.QueryParams, "RequestParam", types)...)
	methodParams = append(methodParams, generateSpringMethodParam(operation.HeaderParams, "RequestHeader", types)...)
	methodParams = append(methodParams, generateSpringMethodParam(operation.Endpoint.UrlParams, "PathVariable", types)...)

	return methodParams
}

func generateSpringMethodParam(namedParams []spec.NamedParam, paramAnnotationName string, types *types.Types) []string {
	params := []string{}

	if namedParams != nil && len(namedParams) > 0 {
		for _, param := range namedParams {
			paramAnnotation := getSpringParameterAnnotation(paramAnnotationName, &param)
			paramType := fmt.Sprintf(`%s: %s`, param.Name.CamelCase(), types.Kotlin(&param.Type.Definition))
			dateFormatAnnotation := dateFormatSpringAnnotation(&param.Type.Definition)
			if dateFormatAnnotation != "" {
				params = append(params, fmt.Sprintf(`%s %s %s`, paramAnnotation, dateFormatAnnotation, paramType))
			} else {
				params = append(params, fmt.Sprintf(`%s %s`, paramAnnotation, paramType))
			}
		}
	}

	return params
}

func getSpringParameterAnnotation(paramAnnotationName string, param *spec.NamedParam) string {
	annotationParams := []string{fmt.Sprintf(`name = "%s"`, param.Name.Source)}

	if param.Type.Definition.IsNullable() {
		annotationParams = append(annotationParams, `required = false`)
	}
	if param.DefinitionDefault.Default != nil {
		annotationParams = append(annotationParams, fmt.Sprintf(`defaultValue = "%s"`, *param.DefinitionDefault.Default))
	}

	return fmt.Sprintf(`@%s(%s)`, paramAnnotationName, joinParams(annotationParams))
}

func dateFormatSpringAnnotation(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return dateFormatSpringAnnotationPlain(typ.Plain)
	case spec.NullableType:
		return dateFormatSpringAnnotation(typ.Child)
	case spec.ArrayType:
		return dateFormatSpringAnnotation(typ.Child)
	case spec.MapType:
		return dateFormatSpringAnnotation(typ.Child)
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func dateFormatSpringAnnotationPlain(typ string) string {
	switch typ {
	case spec.TypeDate:
		return `@DateTimeFormat(iso = DateTimeFormat.ISO.DATE)`
	case spec.TypeDateTime:
		return `@DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME)`
	default:
		return ``
	}
}
