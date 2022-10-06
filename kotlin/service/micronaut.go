package service

import (
	"fmt"
	"strings"

	"generator"
	"github.com/pinzolo/casee"
	"kotlin/imports"
	"kotlin/models"
	"kotlin/packages"
	"kotlin/types"
	"kotlin/writer"
	"spec"
)

var Micronaut = "micronaut"

type MicronautGenerator struct {
	Types  *types.Types
	Models models.Generator
}

func NewMicronautGenerator(types *types.Types, models models.Generator) *MicronautGenerator {
	return &MicronautGenerator{types, models}
}

func (g *MicronautGenerator) ServiceImplAnnotation(api *spec.Api) (annotationImport, annotation string) {
	return `io.micronaut.context.annotation.Bean`, `Bean`
}

func (g *MicronautGenerator) ServicesControllers(version *spec.Version, mainPackage, thePackage, contentTypePackage, jsonPackage, modelsVersionPackage, errorModelsPackage, serviceVersionPackage packages.Package) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		serviceVersionSubpackage := serviceVersionPackage.Subpackage(api.Name.SnakeCase())
		files = append(files, g.serviceController(&api, thePackage, contentTypePackage, jsonPackage, modelsVersionPackage, errorModelsPackage, serviceVersionSubpackage)...)
	}
	files = append(files, dateConverters(mainPackage)...)
	return files
}

func (g *MicronautGenerator) ServiceImports() []string {
	return []string{
		`org.slf4j.*`,
		`io.micronaut.http.*`,
		`io.micronaut.http.annotation.*`,
		`jakarta.inject.Inject`,
	}
}

func (g *MicronautGenerator) ExceptionController(responses *spec.Responses, thePackage, errorsPackage, errorsModelsPackage, jsonPackage packages.Package) *generator.CodeFile {
	w := writer.NewKotlinWriter()
	w.Line(`package %s`, thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.ServiceImports()...)
	imports.Add(`io.micronaut.http.annotation.Error`)
	imports.Add(jsonPackage.Get("Json"))
	imports.Add(errorsModelsPackage.PackageStar)
	imports.Add(errorsPackage.Subpackage("ErrorsHelpers").Get("getBadRequestError"))
	imports.Add(errorsPackage.Subpackage("ErrorsHelpers").Get("getNotFoundError"))
	imports.Write(w)
	w.EmptyLine()
	w.Line(`@Controller`)
	className := `ExceptionController`
	w.Line(`class %s(@Inject private val json: Json) {`, className)
	w.Line(`  private val logger = LoggerFactory.getLogger(%s::class.java)`, className)
	w.EmptyLine()
	g.errorHandler(w.Indented(), *responses)
	w.Line(`}`)

	return &generator.CodeFile{
		Path:    thePackage.GetPath(fmt.Sprintf("%s.kt", className)),
		Content: w.String(),
	}
}

func (g *MicronautGenerator) errorHandler(w *generator.Writer, errors spec.Responses) {
	notFoundError := errors.GetByStatusName(spec.HttpStatusNotFound)
	badRequestError := errors.GetByStatusName(spec.HttpStatusBadRequest)
	internalServerError := errors.GetByStatusName(spec.HttpStatusInternalServerError)
	w.Line(`@Error(global = true, exception = Throwable::class)`)
	w.Line(`fun error(request: HttpRequest<Any>, exception: Throwable): HttpResponse<*> {`)
	w.Line(`  val notFoundError = getNotFoundError(exception)`)
	w.Line(`  if (notFoundError != null) {`)
	g.processResponse(w.IndentedWith(2), notFoundError, "notFoundError")
	w.Line(`  }`)
	w.Line(`  val badRequestError = getBadRequestError(exception)`)
	w.Line(`  if (badRequestError != null) {`)
	g.processResponse(w.IndentedWith(2), badRequestError, "badRequestError")
	w.Line(`  }`)
	w.Line(`  val internalServerError = InternalServerError(exception.message ?: "Unknown error")`)
	g.processResponse(w.IndentedWith(1), internalServerError, "internalServerError")
	w.Line(`}`)
}

func (g *MicronautGenerator) serviceController(api *spec.Api, thePackage, contentTypePackage, jsonPackage, modelsVersionPackage, errorModelsPackage, serviceVersionPackage packages.Package) []generator.CodeFile {
	files := []generator.CodeFile{}
	w := writer.NewKotlinWriter()
	w.Line(`package %s`, thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.ServiceImports()...)
	imports.Add(contentTypePackage.PackageStar)
	imports.Add(jsonPackage.Get("Json"))
	imports.Add(modelsVersionPackage.PackageStar)
	imports.Add(errorModelsPackage.PackageStar)
	imports.Add(serviceVersionPackage.PackageStar)
	imports.Add(g.Models.ModelsUsageImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`@Controller`)
	className := controllerName(api)
	w.Line(`class %s(`, className)
	w.Line(`  @Inject private val %s: %s,`, serviceVarName(api), serviceInterfaceName(api))
	w.Line(`  @Inject private val json: Json`)
	w.Line(`) {`)
	w.Line(`  private val logger = LoggerFactory.getLogger(%s::class.java)`, className)

	for _, operation := range api.Operations {
		w.EmptyLine()
		g.controllerMethod(w.Indented(), &operation)
	}
	w.Line(`}`)

	files = append(files, generator.CodeFile{
		Path:    thePackage.GetPath(fmt.Sprintf("%s.kt", className)),
		Content: w.String(),
	})

	return files
}

func (g *MicronautGenerator) controllerMethod(w *generator.Writer, operation *spec.NamedOperation) {
	if operation.BodyIs(spec.BodyString) {
		w.Line(`@Consumes(MediaType.TEXT_PLAIN)`)
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`@Consumes(MediaType.APPLICATION_JSON)`)
	}
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	w.Line(`@%s("%s")`, casee.ToPascalCase(methodName), url)
	w.Line(`fun %s(%s): HttpResponse<*> {`, controllerMethodName(operation), joinParams(micronautMethodParams(operation, g.Types)))
	w.Line(`  logger.info("Received request, operationId: %s.%s, method: %s, url: %s")`, operation.InApi.Name.Source, operation.Name.Source, methodName, url)
	w.Indent()
	g.parseBody(w, operation, "bodyStr", "requestBody")
	serviceCall(w, operation, "bodyStr", "requestBody", "result")
	g.processResponses(w, operation, "result")
	w.Unindent()
	w.Line(`}`)
}

func (g *MicronautGenerator) parseBody(w *generator.Writer, operation *spec.NamedOperation, bodyStringVar, bodyJsonVar string) {
	if operation.BodyIs(spec.BodyString) {
		w.Line(`checkContentType(request, MediaType.TEXT_PLAIN)`)
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`checkContentType(request, MediaType.APPLICATION_JSON)`)
		typ := g.Types.Kotlin(&operation.Body.Type.Definition)
		w.Line(`val %s: %s = json.read(%s);`, bodyJsonVar, typ, g.Models.JsonRead(bodyStringVar, &operation.Body.Type.Definition))
	}
}

func (g *MicronautGenerator) processResponses(w *generator.Writer, operation *spec.NamedOperation, resultVarName string) {
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

func (g *MicronautGenerator) processResponse(w *generator.Writer, response *spec.Response, bodyVar string) {
	if response.BodyIs(spec.BodyEmpty) {
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s)`, response.Name.UpperCase())
		w.Line(`return HttpResponse.status<Any>(HttpStatus.%s)`, response.Name.UpperCase())
	}
	if response.BodyIs(spec.BodyString) {
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s)`, response.Name.UpperCase())
		w.Line(`return HttpResponse.status<Any>(HttpStatus.%s).body(%s).contentType("text/plain")`, response.Name.UpperCase(), bodyVar)
	}
	if response.BodyIs(spec.BodyJson) {
		w.Line(`val bodyJson = json.write(%s)`, g.Models.JsonWrite(bodyVar, &response.Type.Definition))
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s)`, response.Name.UpperCase())
		w.Line(`return HttpResponse.status<Any>(HttpStatus.%s).body(bodyJson).contentType("application/json")`, response.Name.UpperCase())
	}
}

func (g *MicronautGenerator) ContentType(thePackage packages.Package) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *contentTypeMismatchException(thePackage))
	files = append(files, *g.checkContentType(thePackage))
	return files
}

func (g *MicronautGenerator) checkContentType(thePackage packages.Package) *generator.CodeFile {
	code := `
package [[.PackageName]]

import io.micronaut.http.HttpRequest

fun checkContentType(request: HttpRequest<*>, expectedContentType: String) {
	val contentType = request.headers.contentType
	if (!(contentType.isPresent && contentType.get().contains(expectedContentType))) {
		throw ContentTypeMismatchException(expectedContentType, if (contentType.isPresent) contentType.get() else null )
	}
}
`
	code, _ = generator.ExecuteTemplate(code, struct {
		PackageName string
	}{
		thePackage.PackageName,
	})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("CheckContentType.kt"),
		Content: strings.TrimSpace(code),
	}
}

func (g *MicronautGenerator) Errors(thePackage, errorsModelsPackage, contentTypePackage, jsonPackage packages.Package) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *g.errorsHelpers(thePackage, errorsModelsPackage, contentTypePackage, jsonPackage))
	files = append(files, *g.Models.ValidationErrorsHelpers())
	return files
}

func (g *MicronautGenerator) errorsHelpers(thePackage, errorsModelsPackage, contentTypePackage, jsonPackage packages.Package) *generator.CodeFile {
	code := `
package [[.PackageName]]

import io.micronaut.core.annotation.AnnotationValue
import io.micronaut.core.convert.exceptions.ConversionErrorException
import io.micronaut.core.type.Argument
import io.micronaut.web.router.exceptions.*
import [[.ContentTypePackage]].*
import [[.ErrorsModelsPackage]].*
import [[.PackageName]].ValidationErrorsHelpers.extractValidationErrors
import [[.JsonPackage]].*
import java.util.*
import javax.validation.ConstraintViolationException

object ErrorsHelpers {
    private val NOT_FOUND_ERROR = NotFoundError("Failed to parse url parameters")

    fun getNotFoundError(exception: Throwable?): NotFoundError? {
        if (exception is UnsatisfiedPathVariableRouteException) {
            return NOT_FOUND_ERROR
        }
        if (exception is UnsatisfiedPartRouteException) {
            return NOT_FOUND_ERROR
        }
        if (exception is ConversionErrorException) {
            val annotation =
                exception.argument.annotationMetadata.findDeclaredAnnotation<Annotation>("io.micronaut.http.annotation.PathVariable")
            if (annotation.isPresent) {
                return NOT_FOUND_ERROR
            }
        }
        return null
    }

    private fun getLocation(argument: Argument<*>): ErrorLocation {
        val query =
            argument.annotationMetadata.findDeclaredAnnotation<Annotation>("io.micronaut.http.annotation.QueryValue")
        if (query.isPresent) {
            return ErrorLocation.QUERY
        }
        val header =
            argument.annotationMetadata.findDeclaredAnnotation<Annotation>("io.micronaut.http.annotation.Headers")
        return if (header.isPresent) {
            ErrorLocation.HEADER
        } else ErrorLocation.BODY
    }

    private fun getParameterName(argument: Argument<*>): String {
        val query =
            argument.annotationMetadata.findDeclaredAnnotation<Annotation>("io.micronaut.http.annotation.QueryValue")
        if (query.isPresent) {
            return query.get().values["value"].toString()
        }
        val header =
            argument.annotationMetadata.findDeclaredAnnotation<Annotation>("io.micronaut.http.annotation.Headers")
        if (header.isPresent) {
            val annotationValues = header.get().values["value"] as Array<AnnotationValue<*>>?
            return annotationValues!![0].values["value"].toString()
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
}
`

	code, _ = generator.ExecuteTemplate(code, struct {
		PackageName         string
		ErrorsModelsPackage string
		ContentTypePackage  string
		JsonPackage         string
	}{thePackage.PackageName,
		errorsModelsPackage.PackageName,
		contentTypePackage.PackageName,
		jsonPackage.PackageName,
	})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("ErrorsHelpers.kt"),
		Content: strings.TrimSpace(code),
	}
}

func (g *MicronautGenerator) JsonHelpers(thePackage packages.Package) []generator.CodeFile {
	files := []generator.CodeFile{}

	files = append(files, *g.Json(thePackage))
	files = append(files, *jsonParseException(thePackage))
	files = append(files, g.Models.SetupLibrary()...)

	return files
}

func (g *MicronautGenerator) Json(thePackage packages.Package) *generator.CodeFile {
	w := writer.NewKotlinWriter()
	w.Line(`package %s`, thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.Models.ModelsUsageImports()...)
	imports.Add(`jakarta.inject.*`)
	imports.Add(`java.io.IOException`)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`@Singleton`)
	className := `Json`
	w.Line(`class %s(%s) {`, className, g.Models.CreateJsonMapperField("Inject"))
	w.Line(g.Models.JsonHelpersMethods())
	w.Line(`}`)

	return &generator.CodeFile{
		Path:    thePackage.GetPath(fmt.Sprintf("%s.kt", className)),
		Content: w.String(),
	}
}

func micronautMethodParams(operation *spec.NamedOperation, types *types.Types) []string {
	methodParams := []string{"request: HttpRequest<*>"}

	if operation.Body != nil {
		methodParams = append(methodParams, "@Body bodyStr: String")
	}

	methodParams = append(methodParams, generateMicronautMethodParam(operation.QueryParams, "QueryValue", types)...)
	methodParams = append(methodParams, generateMicronautMethodParam(operation.HeaderParams, "Header", types)...)
	methodParams = append(methodParams, generateMicronautMethodParam(operation.Endpoint.UrlParams, "PathVariable", types)...)

	return methodParams
}

func generateMicronautMethodParam(namedParams []spec.NamedParam, paramAnnotation string, types *types.Types) []string {
	params := []string{}

	if namedParams != nil && len(namedParams) > 0 {
		for _, param := range namedParams {
			paramType := fmt.Sprintf(`%s: %s`, param.Name.CamelCase(), types.Kotlin(&param.Type.Definition))
			params = append(params, fmt.Sprintf(`%s %s`, getMicronautParameterAnnotation(paramAnnotation, &param), paramType))
		}
	}

	return params
}

func getMicronautParameterAnnotation(paramAnnotation string, param *spec.NamedParam) string {
	annotationParams := []string{fmt.Sprintf(`value = "%s"`, param.Name.Source)}

	if param.DefinitionDefault.Default != nil {
		annotationParams = append(annotationParams, fmt.Sprintf(`defaultValue = "%s"`, *param.DefinitionDefault.Default))
	}

	return fmt.Sprintf(`@%s(%s)`, paramAnnotation, joinParams(annotationParams))
}

func dateConverters(thePackage packages.Package) []generator.CodeFile {
	convertersPackage := thePackage.Subpackage("converters")

	files := []generator.CodeFile{}
	files = append(files, *localDateConverter(convertersPackage))
	files = append(files, *localDateTimeConverter(convertersPackage))
	return files
}

func localDateConverter(thePackage packages.Package) *generator.CodeFile {
	code := `
package [[.PackageName]]

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
}`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("LocalDateConverter.kt"),
		Content: strings.TrimSpace(code),
	}
}

func localDateTimeConverter(thePackage packages.Package) *generator.CodeFile {
	code := `
package [[.PackageName]]

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
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("LocalDateTimeConverter.kt"),
		Content: strings.TrimSpace(code),
	}
}
