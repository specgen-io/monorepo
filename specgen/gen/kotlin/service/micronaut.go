package service

import (
	"fmt"
	"github.com/pinzolo/casee"
	"github.com/specgen-io/specgen/spec/v2"
	"github.com/specgen-io/specgen/specgen/v2/gen/kotlin/imports"
	"github.com/specgen-io/specgen/specgen/v2/gen/kotlin/models"
	"github.com/specgen-io/specgen/specgen/v2/gen/kotlin/modules"
	"github.com/specgen-io/specgen/specgen/v2/gen/kotlin/responses"
	"github.com/specgen-io/specgen/specgen/v2/gen/kotlin/types"
	"github.com/specgen-io/specgen/specgen/v2/gen/kotlin/writer"
	"github.com/specgen-io/specgen/specgen/v2/generator"
	"strings"
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

func (g *MicronautGenerator) ServicesControllers(version *spec.Version, mainPackage, thePackage, jsonPackage, modelsVersionPackage, serviceVersionPackage modules.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		serviceVersionSubpackage := serviceVersionPackage.Subpackage(api.Name.SnakeCase())
		files = append(files, g.serviceController(&api, thePackage, modelsVersionPackage, serviceVersionSubpackage)...)
	}
	files = append(
		files,
		*g.errorsHelpers(thePackage, modelsVersionPackage),
		*g.Models.GenerateJsonParseException(thePackage, modelsVersionPackage),
		*g.contentTypeMismatchException(thePackage),
	)
	files = append(files, dateConverters(mainPackage)...)
	return files
}

func (g *MicronautGenerator) serviceController(api *spec.Api, apiPackage, modelsVersionPackage, serviceVersionPackage modules.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	w := writer.NewKotlinWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(`org.slf4j.*`)
	imports.Add(`io.micronaut.http.*`)
	imports.Add(`io.micronaut.http.annotation.*`)
	imports.Add(`jakarta.inject.Inject`)
	imports.Add(apiPackage.Subpackage("ErrorsHelpers").Get("getBadRequestError"))
	imports.Add(apiPackage.Subpackage("ErrorsHelpers").Get("getNotFoundError"))
	imports.Add(modelsVersionPackage.PackageStar)
	imports.Add(serviceVersionPackage.PackageStar)
	imports.Add(g.Models.ModelsUsageImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`@Controller`)
	className := controllerName(api)
	w.Line(`class %s(`, className)
	w.Line(`  @Inject private val %s: %s,`, serviceVarName(api), serviceInterfaceName(api))
	w.Line(`  @Inject %s`, g.Models.CreateJsonMapperField())
	w.Line(`) {`)
	w.Indent()
	w.Line(`private val logger = LoggerFactory.getLogger(%s::class.java)`, className)
	for _, operation := range api.Operations {
		w.EmptyLine()
		g.controllerMethod(w, &operation)
	}
	w.EmptyLine()
	g.checkContentType(w)
	w.EmptyLine()
	g.errorHandler(w, api.Http.Errors)
	w.Unindent()
	w.Line(`}`)

	files = append(files, generator.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.kt", className)),
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
	w.Line(`  logger.info("Received request, operationId: %s.%s, method: %s, url: %s")`, operation.Api.Name.Source, operation.Name.Source, methodName, url)
	w.Indent()
	g.parseBody(w, operation, "bodyStr", "requestBody")
	g.serviceCall(w, operation, "bodyStr", "requestBody", "result")
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
		requestBody, exception := g.Models.ReadJson(bodyStringVar, &operation.Body.Type.Definition)
		w.Line(`val %s = try {`, bodyJsonVar)
		w.Line(`  %s`, requestBody)
		w.Line(`} catch (exception: %s) {`, exception)
		w.Line(`  throw JsonParseException(exception)`)
		w.Line(`}`)
	}
}

func (g *MicronautGenerator) serviceCall(w *generator.Writer, operation *spec.NamedOperation, bodyStringVar, bodyJsonVar, resultVarName string) {
	serviceCall := fmt.Sprintf(`%s.%s(%s)`, serviceVarName(operation.Api), operation.Name.CamelCase(), joinParams(addServiceMethodParams(operation, bodyStringVar, bodyJsonVar)))
	if len(operation.Responses) == 1 && operation.Responses[0].BodyIs(spec.BodyEmpty) {
		w.Line(serviceCall)
	} else {
		w.Line(`val %s = %s`, resultVarName, serviceCall)
	}
}

func (g *MicronautGenerator) processResponses(w *generator.Writer, operation *spec.NamedOperation, resultVarName string) {
	if len(operation.Responses) == 1 {
		g.processResponse(w, &operation.Responses[0].Response, resultVarName)
	}
	if len(operation.Responses) > 1 {
		for _, response := range operation.Responses {
			w.Line(`if (%s is %s.%s) {`, resultVarName, responses.InterfaceName(operation), response.Name.PascalCase())
			g.processResponse(w.Indented(), &response.Response, responses.GetBody(resultVarName))
			w.Line(`}`)
		}
		w.Line(`throw RuntimeException("Service implementation didn't return any value'")`)
	}
}

func (g *MicronautGenerator) processResponse(w *generator.Writer, response *spec.Response, result string) {
	if response.BodyIs(spec.BodyEmpty) {
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s)`, response.Name.UpperCase())
		w.Line(`return HttpResponse.status<Any>(HttpStatus.%s)`, response.Name.UpperCase())
	}
	if response.BodyIs(spec.BodyString) {
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s)`, response.Name.UpperCase())
		w.Line(`return HttpResponse.status<Any>(HttpStatus.%s).body(%s).contentType("text/plain")`, response.Name.UpperCase(), result)
	}
	if response.BodyIs(spec.BodyJson) {
		responseWrite, _ := g.Models.WriteJson(result, &response.Type.Definition)
		w.Line(`val bodyJson = %s`, responseWrite)
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s)`, response.Name.UpperCase())
		w.Line(`return HttpResponse.status<Any>(HttpStatus.%s).body(bodyJson).contentType("application/json")`, response.Name.UpperCase())
	}
}

func (g *MicronautGenerator) checkContentType(w *generator.Writer) {
	w.Lines(`
private fun checkContentType(request: HttpRequest<*>, expectedContentType: String) {
	val contentType = request.headers.contentType
	if (!(contentType.isPresent && contentType.get().contains(expectedContentType))) {
		throw ContentTypeMismatchException(expectedContentType, if (contentType.isPresent) contentType.get() else null )
	}
}
`)
}

func (g *MicronautGenerator) errorHandler(w *generator.Writer, errors spec.Responses) {
	notFoundError := errors.GetByStatusName(spec.HttpStatusNotFound)
	badRequestError := errors.GetByStatusName(spec.HttpStatusBadRequest)
	internalServerError := errors.GetByStatusName(spec.HttpStatusInternalServerError)
	w.Line(`@Error`)
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

func (g *MicronautGenerator) contentTypeMismatchException(thePackage modules.Module) *generator.CodeFile {
	code := `
package [[.PackageName]];

class ContentTypeMismatchException(expected: String, actual: String?) :
    RuntimeException(
        String.format(
            "Expected Content-Type header: '%s' was not provided, found: '%s'",
            expected,
            actual
        )
    )
`
	code, _ = generator.ExecuteTemplate(code, struct {
		PackageName string
	}{
		thePackage.PackageName,
	})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("ContentTypeMismatchException.kt"),
		Content: strings.TrimSpace(code),
	}
}

func (g *MicronautGenerator) errorsHelpers(thePackage, modelsPackage modules.Module) *generator.CodeFile {
	code := `
package [[.PackageName]]

import io.micronaut.core.annotation.AnnotationValue
import io.micronaut.core.convert.exceptions.ConversionErrorException
import io.micronaut.core.type.Argument
import io.micronaut.web.router.exceptions.*

import java.util.*
import javax.validation.ConstraintViolationException;

import [[.ModelsPackage]].*

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
            val parameterName = query.get().values["value"]
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
            return BadRequestError("Failed to parse body", ErrorLocation.BODY, exception.errors)
        }
        if (exception is ContentTypeMismatchException) {
            val error = ValidationError("Content-Type", "missing", exception.message)
            return BadRequestError("Failed to parse header", ErrorLocation.HEADER, listOf(error))
        }
        if (exception is UnsatisfiedRouteException) {
            val e = exception
            return argumentBadRequestError(e.argument, e.message, "missing")
        }
        if (exception is ConversionErrorException) {
            val e = exception
            return argumentBadRequestError(e.argument, e.message, "parsing_failed")
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
		PackageName   string
		ModelsPackage string
	}{thePackage.PackageName, modelsPackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("ErrorsHelpers.kt"),
		Content: strings.TrimSpace(code),
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

func dateConverters(thePackage modules.Module) []generator.CodeFile {
	convertersPackage := thePackage.Subpackage("converters")

	files := []generator.CodeFile{}
	files = append(files, *localDateConverter(convertersPackage))
	files = append(files, *localDateTimeConverter(convertersPackage))
	return files
}

func localDateConverter(thePackage modules.Module) *generator.CodeFile {
	code := `
package [[.PackageName]];

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

func localDateTimeConverter(thePackage modules.Module) *generator.CodeFile {
	code := `
package [[.PackageName]];

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
