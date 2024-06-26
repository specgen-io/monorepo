package service

import (
	"fmt"
	"generator"
	"github.com/pinzolo/casee"
	"kotlin/models"
	"kotlin/packages"
	"kotlin/types"
	"kotlin/writer"
	"spec"
	"strings"
)

var Micronaut = "micronaut"

type MicronautGenerator struct {
	Types    *types.Types
	Models   models.Generator
	Packages *Packages
}

func NewMicronautGenerator(types *types.Types, models models.Generator, packages *Packages) *MicronautGenerator {
	return &MicronautGenerator{types, models, packages}
}

func (g *MicronautGenerator) ServiceImplAnnotation(api *spec.Api) (annotationImport, annotation string) {
	return `io.micronaut.context.annotation.Bean`, `Bean`
}

func (g *MicronautGenerator) ServicesControllers(version *spec.Version) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *g.serviceController(&api))
	}
	files = append(files, dateConverters(g.Packages.Converters)...)
	return files
}

func (g *MicronautGenerator) FilesImports() []string {
	return []string{
		`io.micronaut.http.multipart.CompletedFileUpload`,
		`io.micronaut.http.server.types.files.StreamedFile`,
	}
}

func (g *MicronautGenerator) ServiceImports() []string {
	return []string{
		`org.slf4j.*`,
		`io.micronaut.http.*`,
		`io.micronaut.http.annotation.*`,
		`jakarta.inject.Inject`,
	}
}

func (g *MicronautGenerator) ExceptionController(responses *spec.ErrorResponses) *generator.CodeFile {
	w := writer.New(g.Packages.RootControllers, `ExceptionController`)
	w.Imports.Add(g.ServiceImports()...)
	w.Imports.Add(`io.micronaut.http.annotation.Error`)
	w.Imports.PackageStar(g.Packages.Json)
	w.Imports.PackageStar(g.Packages.ErrorsModels)
	w.Imports.PackageStar(g.Packages.Errors)
	w.Line(`@Controller`)
	w.Line(`class [[.ClassName]](@Inject private val json: Json) {`)
	w.Line(`    private val logger = LoggerFactory.getLogger([[.ClassName]]::class.java)`)
	w.EmptyLine()
	g.errorHandler(w.Indented(), *responses)
	w.Line(`}`)
	return w.ToCodeFile()
}

func (g *MicronautGenerator) errorHandler(w *writer.Writer, errors spec.ErrorResponses) {
	notFoundError := errors.GetByStatusName(spec.HttpStatusNotFound)
	badRequestError := errors.GetByStatusName(spec.HttpStatusBadRequest)
	internalServerError := errors.GetByStatusName(spec.HttpStatusInternalServerError)
	w.Line(`@Error(global = true, exception = Throwable::class)`)
	w.Line(`fun error(request: HttpRequest<Any>, exception: Throwable): HttpResponse<*> {`)
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

func (g *MicronautGenerator) serviceController(api *spec.Api) *generator.CodeFile {
	w := writer.New(g.Packages.Controllers(api.InHttp.InVersion), controllerName(api))
	w.Imports.Add(g.ServiceImports()...)
	w.Imports.Add(g.FilesImports()...)
	w.Imports.PackageStar(g.Packages.ContentType)
	w.Imports.PackageStar(g.Packages.Json)
	w.Imports.PackageStar(g.Packages.ErrorsModels)
	w.Imports.PackageStar(g.Packages.Models(api.InHttp.InVersion))
	w.Imports.PackageStar(g.Packages.ServicesApi(api))
	w.Imports.Add(g.Types.Imports()...)
	w.Line(`@Controller`)
	w.Line(`class [[.ClassName]](`)
	w.Line(`    @Inject private val %s: %s,`, serviceVarName(api), serviceInterfaceName(api))
	w.Line(`    @Inject private val json: Json`)
	w.Line(`) {`)
	w.Line(`    private val logger = LoggerFactory.getLogger([[.ClassName]]::class.java)`)
	for _, operation := range api.Operations {
		w.EmptyLine()
		g.controllerMethod(w.Indented(), &operation)
	}
	w.Line(`}`)
	return w.ToCodeFile()
}

func (g *MicronautGenerator) controllerMethod(w *writer.Writer, operation *spec.NamedOperation) {
	if !operation.Body.IsEmpty() {
		w.Line(`@Consumes(%s)`, g.requestContentType(operation))
	}
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	w.Line(`@%s("%s")`, casee.ToPascalCase(methodName), url)
	w.Line(`fun %s(%s): HttpResponse<*> {`, controllerMethodName(operation), strings.Join(micronautMethodParams(operation, g.Types), ", "))
	w.Line(`    logger.info("Received request, operationId: %s.%s, method: %s, url: %s")`, operation.InApi.Name.Source, operation.Name.Source, methodName, url)
	w.Indent()
	bodyStringVar := "bodyStr"
	if operation.Body.IsJson() {
		bodyStringVar += ".reader()"
	}
	g.parseBody(w, operation, bodyStringVar, "requestBody")
	serviceCall(w, operation, bodyStringVar, "requestBody", "file", "result", false)
	g.processResponses(w, operation, "result")
	w.Unindent()
	w.Line(`}`)
}

func (g *MicronautGenerator) parseBody(w *writer.Writer, operation *spec.NamedOperation, bodyStringVar, bodyJsonVar string) {
	if !operation.Body.IsEmpty() {
		w.Line(`checkContentType(request, %s)`, g.requestContentType(operation))
	}
	if operation.Body.IsJson() {
		typ := g.Types.Kotlin(&operation.Body.Type.Definition)
		w.Line(`val %s: %s = json.%s`, bodyJsonVar, typ, g.Models.ReadJson(bodyStringVar, &operation.Body.Type.Definition))
	}
}

func (g *MicronautGenerator) requestContentType(operation *spec.NamedOperation) string {
	switch operation.Body.Kind() {
	case spec.BodyEmpty:
		return ""
	case spec.BodyText:
		return `MediaType.TEXT_PLAIN`
	case spec.BodyJson:
		return `MediaType.APPLICATION_JSON`
	case spec.BodyBinary:
		return `MediaType.APPLICATION_OCTET_STREAM`
	case spec.BodyFormData:
		return `MediaType.MULTIPART_FORM_DATA`
	case spec.BodyFormUrlEncoded:
		return `MediaType.APPLICATION_FORM_URLENCODED`
	default:
		panic(fmt.Sprintf("Unknown Content Type"))
	}
}

func (g *MicronautGenerator) responseContentType(response *spec.Response) string {
	switch response.Body.Kind() {
	case spec.BodyEmpty:
		return ""
	case spec.BodyText:
		return `MediaType.TEXT_PLAIN`
	case spec.BodyJson:
		return `MediaType.APPLICATION_JSON`
	case spec.BodyBinary:
		return `MediaType.APPLICATION_OCTET_STREAM`
	case spec.BodyFile:
		return ""
	default:
		panic(fmt.Sprintf("Unknown Content Type"))
	}
}

func (g *MicronautGenerator) processResponses(w *writer.Writer, operation *spec.NamedOperation, resultVarName string) {
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

func (g *MicronautGenerator) processResponse(w *writer.Writer, response *spec.Response, bodyVar string) {
	if response.Body.IsEmpty() {
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s)`, response.Name.UpperCase())
		w.Line(`return HttpResponse.status<Any>(HttpStatus.%s)`, response.Name.UpperCase())
	} else {
		if response.Body.IsJson() {
			w.Line(`val bodyJson = json.%s`, g.Models.WriteJson(bodyVar, &response.Body.Type.Definition))
			bodyVar = "bodyJson"
		}
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s)`, response.Name.UpperCase())
		if response.Body.IsFile() {
			w.Line(`return HttpResponse.status<Any>(HttpStatus.%s).body(%s)`, response.Name.UpperCase(), bodyVar)
		} else {
			w.Line(`return HttpResponse.status<Any>(HttpStatus.%s).body(%s).contentType(%s)`, response.Name.UpperCase(), bodyVar, g.responseContentType(response))
		}
	}
}

func (g *MicronautGenerator) ContentType() []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *g.contentTypeMismatchException())
	files = append(files, *g.checkContentType())
	return files
}

func (g *MicronautGenerator) contentTypeMismatchException() *generator.CodeFile {
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

func (g *MicronautGenerator) checkContentType() *generator.CodeFile {
	w := writer.New(g.Packages.ContentType, `CheckContentType`)
	w.Lines(`
import io.micronaut.http.HttpRequest

fun checkContentType(request: HttpRequest<*>, expectedContentType: String) {
	val contentType = request.headers.contentType
	if (!(contentType.isPresent && contentType.get().contains(expectedContentType))) {
		throw ContentTypeMismatchException(expectedContentType, if (contentType.isPresent) contentType.get() else null )
	}
}
`)
	return w.ToCodeFile()
}

func (g *MicronautGenerator) ErrorsHelpers() *generator.CodeFile {
	w := writer.New(g.Packages.Errors, `ErrorsHelpers`)
	w.Template(
		map[string]string{
			`ContentTypePackage`:  g.Packages.ContentType.PackageName,
			`ErrorsModelsPackage`: g.Packages.ErrorsModels.PackageName,
			`ErrorsPackage`:       g.Packages.Errors.PackageName,
			`JsonPackage`:         g.Packages.Json.PackageName,
		}, `
import io.micronaut.core.annotation.AnnotationValue
import io.micronaut.core.convert.exceptions.ConversionErrorException
import io.micronaut.core.type.Argument
import io.micronaut.web.router.exceptions.*
import [[.ContentTypePackage]].*
import [[.ErrorsModelsPackage]].*
import [[.ErrorsPackage]].ValidationErrorsHelpers.extractValidationErrors
import [[.JsonPackage]].*
import java.util.*
import javax.validation.ConstraintViolationException

const val NOT_FOUND_ERROR = "Failed to parse url parameters"

fun getNotFoundError(exception: Throwable?): NotFoundError? {
	if (exception is UnsatisfiedPathVariableRouteException) {
		return NotFoundError(NOT_FOUND_ERROR)
	}
	if (exception is ConversionErrorException) {
		val annotation = exception.argument.annotationMetadata.findDeclaredAnnotation<Annotation>("io.micronaut.http.annotation.PathVariable")
		if (annotation.isPresent) {
			return NotFoundError(NOT_FOUND_ERROR)
		}
	}
	return null
}

private fun getLocation(argument: Argument<*>): ErrorLocation {
	val query = argument.annotationMetadata.findDeclaredAnnotation<Annotation>("io.micronaut.http.annotation.QueryValue")
	if (query.isPresent) {
		return ErrorLocation.QUERY
	}
	val header = argument.annotationMetadata.findDeclaredAnnotation<Annotation>("io.micronaut.http.annotation.Headers")
	if (header.isPresent) {
		return ErrorLocation.HEADER
	}
	val part = argument.annotationMetadata.findDeclaredAnnotation<Annotation>("io.micronaut.http.annotation.Part")
	if (part.isPresent) {
		return ErrorLocation.BODY
	}
	return ErrorLocation.BODY
}

private fun getParameterName(argument: Argument<*>): String {
	val query = argument.annotationMetadata.findDeclaredAnnotation<Annotation>("io.micronaut.http.annotation.QueryValue")
	if (query.isPresent) {
		return query.get().values["value"].toString()
	}
	val header = argument.annotationMetadata.findDeclaredAnnotation<Annotation>("io.micronaut.http.annotation.Headers")
	if (header.isPresent) {
		val annotationValues = header.get().values["value"] as Array<AnnotationValue<*>>?
		return annotationValues!![0].values["value"].toString()
	}
	val part = argument.annotationMetadata.findDeclaredAnnotation<Annotation>("io.micronaut.http.annotation.Part")
	if (part.isPresent) {
		return part.get().values["value"].toString()
	}
	return "unknown"
}

private fun argumentBadRequestError(arg: Argument<*>, errorMessage: String?, code: String): BadRequestError {
	val location = getLocation(arg)
	val parameterName = getParameterName(arg)
	val validation = ValidationError(parameterName, code, errorMessage)
	val message = String.format("Failed to parse %s", location.name.lowercase(Locale.getDefault()))
	return BadRequestError(message, location, listOf(validation))
}

fun getBadRequestError(exception: Throwable): BadRequestError? {
	if (exception is UnsatisfiedPartRouteException) {
		val errors = listOf(ValidationError(exception.partName, "missing", exception.message))
		return BadRequestError("Failed to parse body", ErrorLocation.BODY, errors)
	}
	if (exception is JsonParseException) {
		val errors = extractValidationErrors(exception)
		return BadRequestError("Failed to parse body", ErrorLocation.BODY, errors)
	}
	if (exception is ContentTypeMismatchException) {
		val error = ValidationError("Content-Type", "missing", exception.message)
		return BadRequestError("Failed to parse header", ErrorLocation.HEADER, listOf(error))
	}
	if (exception is UnsatisfiedRouteException) {
		return argumentBadRequestError(exception.argument, exception.message, "missing")
	}
	if (exception is ConversionErrorException) {
		return argumentBadRequestError(exception.argument, exception.message, "parsing_failed")
	}
	if (exception is ConstraintViolationException) {
		val message = "Failed to parse body"
		return BadRequestError(message, ErrorLocation.BODY, null)
	}
	return null
}
`)
	return w.ToCodeFile()
}

func micronautMethodParams(operation *spec.NamedOperation, types *types.Types) []string {
	methodParams := []string{"request: HttpRequest<*>"}
	
	if operation.Body.IsText() || operation.Body.IsJson() {
		methodParams = append(methodParams, "@Body bodyStr: String")
	}
	if operation.Body.IsBinary() {
		methodParams = append(methodParams, fmt.Sprintf("@Body file: %s", types.RequestBodyKotlinType(&operation.Body)))
	}
	methodParams = append(methodParams, generateMicronautMethodParam(operation.Body.FormData, "Part", false, types)...)
	methodParams = append(methodParams, generateMicronautMethodParam(operation.Body.FormUrlEncoded, "Part", false, types)...)
	methodParams = append(methodParams, generateMicronautMethodParam(operation.QueryParams, "QueryValue", true, types)...)
	methodParams = append(methodParams, generateMicronautMethodParam(operation.HeaderParams, "Header", true, types)...)
	methodParams = append(methodParams, generateMicronautMethodParam(operation.Endpoint.UrlParams, "PathVariable", true, types)...)
	
	return methodParams
}

func generateMicronautMethodParam(namedParams []spec.NamedParam, paramAnnotation string, isSupportDefaulted bool, types *types.Types) []string {
	params := []string{}
	
	if namedParams != nil && len(namedParams) > 0 {
		for _, param := range namedParams {
			paramType := fmt.Sprintf(`%s: %s`, param.Name.CamelCase(), types.ParamKotlinType(&param))
			if !isSupportDefaulted && param.DefinitionDefault.Default != nil {
				paramType = fmt.Sprintf(`%s?`, paramType)
			}
			params = append(params, fmt.Sprintf(`%s %s`, getMicronautParameterAnnotation(paramAnnotation, &param, isSupportDefaulted), paramType))
		}
	}
	
	return params
}

func getMicronautParameterAnnotation(paramAnnotation string, param *spec.NamedParam, isSupportDefaulted bool) string {
	annotationParams := []string{fmt.Sprintf(`value = "%s"`, param.Name.Source)}
	
	if param.DefinitionDefault.Default != nil && isSupportDefaulted {
		annotationParams = append(annotationParams, fmt.Sprintf(`defaultValue = "%s"`, *param.DefinitionDefault.Default))
	}
	
	return fmt.Sprintf(`@%s(%s)`, paramAnnotation, strings.Join(annotationParams, ", "))
}

func dateConverters(convertersPackage packages.Package) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *localDateConverter(convertersPackage))
	files = append(files, *localDateTimeConverter(convertersPackage))
	return files
}

func localDateConverter(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `LocalDateConverter`)
	w.Lines(`
import io.micronaut.core.convert.*
import jakarta.inject.Singleton

import java.time.format.*
import java.time.LocalDate
import java.util.*

@Singleton
class LocalDateConverter : TypeConverter<String, LocalDate> {
    override fun convert(
        value: String,
        targetType: Class<LocalDate>,
        context: ConversionContext
    ): Optional<LocalDate> {
        return try {
            val formatter = DateTimeFormatter.ofPattern("yyyy-MM-dd", context.locale)
            val result = LocalDate.parse(value, formatter)
            Optional.of(result)
        } catch (e: DateTimeParseException) {
            context.reject(value, e)
            Optional.empty()
        }
    }
}`)
	return w.ToCodeFile()
}

func localDateTimeConverter(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `LocalDateTimeConverter`)
	w.Lines(`
import io.micronaut.core.convert.*
import jakarta.inject.Singleton

import java.time.format.*
import java.time.LocalDateTime
import java.util.*

@Singleton
class LocalDateTimeConverter : TypeConverter<String, LocalDateTime> {
    override fun convert(
        value: String,
        targetType: Class<LocalDateTime>,
        context: ConversionContext
    ): Optional<LocalDateTime> {
        return try {
            val formatter = DateTimeFormatter.ofPattern("yyyy-MM-dd'T'HH:mm:ss", context.locale)
            val result = LocalDateTime.parse(value, formatter)
            Optional.of(result)
        } catch (e: DateTimeParseException) {
            context.reject(value, e)
            Optional.empty()
        }
    }
}
`)
	return w.ToCodeFile()
}
