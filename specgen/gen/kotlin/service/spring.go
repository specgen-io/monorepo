package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen/kotlin/imports"
	"github.com/specgen-io/specgen/v2/gen/kotlin/models"
	"github.com/specgen-io/specgen/v2/gen/kotlin/modules"
	"github.com/specgen-io/specgen/v2/gen/kotlin/responses"
	"github.com/specgen-io/specgen/v2/gen/kotlin/types"
	"github.com/specgen-io/specgen/v2/gen/kotlin/writer"
	"strings"

	"github.com/pinzolo/casee"
	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/spec/v2"
)

var Spring = "spring"

type SpringGenerator struct {
	Types  *types.Types
	Models models.Generator
}

func NewSpringGenerator(types *types.Types, models models.Generator) *SpringGenerator {
	return &SpringGenerator{types, models}
}

func (g *SpringGenerator) ServiceImplAnnotation(api *spec.Api) (annotationImport, annotation string) {
	return `org.springframework.stereotype.Service`, fmt.Sprintf(`Service("%s")`, versionServiceName(serviceName(api), api.Http.Version))
}

func (g *SpringGenerator) ServicesControllers(version *spec.Version, mainPackage, thePackage, modelsVersionPackage, serviceVersionPackage modules.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		serviceVersionSubpackage := serviceVersionPackage.Subpackage(api.Name.SnakeCase())
		files = append(files, g.serviceController(&api, thePackage, modelsVersionPackage, serviceVersionSubpackage)...)
	}
	files = append(
		files,
		*g.errorsHelpers(thePackage, modelsVersionPackage),
		*g.Models.GenerateJsonParseException(thePackage, modelsVersionPackage),
		*contentTypeMismatchException(thePackage),
	)
	return files
}

func (g *SpringGenerator) serviceController(api *spec.Api, apiPackage, modelsVersionPackage, serviceVersionPackage modules.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	w := writer.NewKotlinWriter()
	w.Line(`package %s`, apiPackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(`org.apache.logging.log4j.LogManager`)
	imports.Add(`org.apache.tomcat.util.http.fileupload.FileUploadBase.CONTENT_TYPE`)
	imports.Add(`org.springframework.beans.factory.annotation.Autowired`)
	imports.Add(`org.springframework.format.annotation.DateTimeFormat`)
	imports.Add(`org.springframework.http.*`)
	imports.Add(`org.springframework.web.bind.annotation.*`)
	imports.Add(`javax.servlet.http.HttpServletRequest`)
	imports.Add(apiPackage.Subpackage("ErrorsHelpers").Get("getBadRequestError"))
	imports.Add(apiPackage.Subpackage("ErrorsHelpers").Get("getNotFoundError"))
	imports.Add(modelsVersionPackage.PackageStar)
	imports.Add(serviceVersionPackage.PackageStar)
	imports.Add(g.Models.ModelsUsageImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`@RestController("%s")`, versionControllerName(controllerName(api), api.Http.Version))
	className := controllerName(api)
	w.Line(`class %s(`, className)
	w.Line(`  @Autowired private val %s: %s,`, serviceVarName(api), serviceInterfaceName(api))
	w.Line(`  @Autowired %s`, g.Models.CreateJsonMapperField())
	w.Line(`) {`)
	w.Indent()
	w.Line(`private val logger = LogManager.getLogger(%s::class.java)`, className)
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

func (g *SpringGenerator) controllerMethod(w *generator.Writer, operation *spec.NamedOperation) {
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	w.Line(`@%sMapping("%s")`, casee.ToPascalCase(methodName), url)
	w.Line(`fun %s(%s): ResponseEntity<String> {`, controllerMethodName(operation), strings.Join(springMethodParams(operation, g.Types), ", "))
	w.Indent()
	w.Line(`logger.info("Received request, operationId: %s.%s, method: %s, url: %s")`, operation.Api.Name.Source, operation.Name.Source, methodName, url)
	g.parseBody(w, operation, "bodyStr", "requestBody")
	serviceCall(w, operation, "bodyStr", "requestBody", "result")
	g.processResponses(w, operation, "result")
	w.Unindent()
	w.Line(`}`)
}

func (g *SpringGenerator) parseBody(w *generator.Writer, operation *spec.NamedOperation, bodyStringVar, bodyJsonVar string) {
	if operation.BodyIs(spec.BodyString) {
		w.Line(`checkContentType(request, MediaType.TEXT_PLAIN);`)
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`checkContentType(request, MediaType.APPLICATION_JSON);`)
		requestBody, exception := g.Models.ReadJson(bodyStringVar, &operation.Body.Type.Definition)
		w.Line(`val %s = try {`, bodyJsonVar)
		w.Line(`  %s`, requestBody)
		w.Line(`} catch (exception: %s) {`, exception)
		w.Line(`  throw JsonParseException(exception)`)
		w.Line(`}`)
	}
}

func (g *SpringGenerator) processResponses(w *generator.Writer, operation *spec.NamedOperation, resultVarName string) {
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

func (g *SpringGenerator) processResponse(w *generator.Writer, response *spec.Response, result string) {
	if response.BodyIs(spec.BodyEmpty) {
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s)`, response.Name.UpperCase())
		w.Line(`return ResponseEntity(HttpStatus.%s);`, response.Name.UpperCase())
	}
	if response.BodyIs(spec.BodyString) {
		w.Line(`val headers = HttpHeaders()`)
		w.Line(`headers.add(CONTENT_TYPE, "text/plain")`)
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s)`, response.Name.UpperCase())
		w.Line(`return ResponseEntity(%s, headers, HttpStatus.%s)`, result, response.Name.UpperCase())
	}
	if response.BodyIs(spec.BodyJson) {
		responseWrite, _ := g.Models.WriteJson(result, &response.Type.Definition)
		w.Line(`val bodyJson = %s;`, responseWrite)
		w.Line(`val headers = HttpHeaders()`)
		w.Line(`headers.add(CONTENT_TYPE, "application/json")`)
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s)`, response.Name.UpperCase())
		w.Line(`return ResponseEntity(bodyJson, headers, HttpStatus.%s)`, response.Name.UpperCase())
	}
}

func (g *SpringGenerator) checkContentType(w *generator.Writer) {
	w.Lines(`
private fun checkContentType(request: HttpServletRequest, expectedContentType: MediaType) {
	val contentType = request.getHeader("Content-Type")
	if (contentType == null || !contentType.contains(expectedContentType.toString())) {
		throw ContentTypeMismatchException(expectedContentType.toString(), contentType)
	}
}
`)
}

func (g *SpringGenerator) errorHandler(w *generator.Writer, errors spec.Responses) {
	notFoundError := errors.GetByStatusName(spec.HttpStatusNotFound)
	badRequestError := errors.GetByStatusName(spec.HttpStatusBadRequest)
	internalServerError := errors.GetByStatusName(spec.HttpStatusInternalServerError)
	w.Line(`@ExceptionHandler(Throwable::class)`)
	w.Line(`fun error(request: HttpServletRequest, exception: Throwable): ResponseEntity<String> {`)
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

func (g *SpringGenerator) errorsHelpers(thePackage, modelsPackage modules.Module) *generator.CodeFile {
	code := `
package [[.PackageName]]

import org.springframework.web.bind.*
import org.springframework.web.bind.annotation.*
import org.springframework.web.method.annotation.MethodArgumentTypeMismatchException
import [[.ModelsPackage]].*

object ErrorsHelpers {
    private val NOT_FOUND_ERROR = NotFoundError("Failed to parse url parameters")
    fun getNotFoundError(exception: Throwable?): NotFoundError? {
        if (exception is MethodArgumentTypeMismatchException) {
            if (exception.parameter.hasParameterAnnotation(PathVariable::class.java)) {
                return NOT_FOUND_ERROR
            }
        }
        return null
    }

    fun getBadRequestError(exception: Throwable): BadRequestError? {
        if (exception is JsonParseException) {
            return BadRequestError("Failed to parse body", ErrorLocation.BODY, exception.errors)
        }
        if (exception is ContentTypeMismatchException) {
            val error = ValidationError("Content-Type", "missing", exception.message)
            return BadRequestError("Failed to parse header", ErrorLocation.HEADER, listOf(error))
        }
        if (exception is MissingServletRequestParameterException) {
            val message = String.format("Failed to parse query")
            val validation = ValidationError(exception.parameterName, "missing", exception.message)
            return BadRequestError(message, ErrorLocation.QUERY, listOf(validation))
        }
        if (exception is MethodArgumentTypeMismatchException) {
            val message = String.format("Failed to parse query")
            val validation = ValidationError(exception.name, "parsing_failed", exception.message)
            if (exception.parameter.hasParameterAnnotation(RequestParam::class.java)) {
                return BadRequestError(message, ErrorLocation.QUERY, listOf(validation))
            } else if (exception.parameter.hasParameterAnnotation(RequestHeader::class.java)) {
                return BadRequestError(message, ErrorLocation.HEADER, listOf(validation))
            }
        }
        if (exception is MissingRequestHeaderException) {
            val message = String.format("Failed to parse header")
            val validation = ValidationError(exception.headerName, "missing", exception.message)
            return BadRequestError(message, ErrorLocation.HEADER, listOf(validation))
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

func springMethodParams(operation *spec.NamedOperation, types *types.Types) []string {
	methodParams := []string{"request: HttpServletRequest"}

	if operation.Body != nil {
		methodParams = append(methodParams, "@RequestBody bodyStr: String")
	}
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
