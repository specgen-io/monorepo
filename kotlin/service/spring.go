package service

import (
	"fmt"
	"strings"

	"kotlin/imports"
	"kotlin/models"
	"kotlin/modules"
	"kotlin/responses"
	"kotlin/types"
	"kotlin/writer"

	"generator"
	"github.com/pinzolo/casee"
	"spec"
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

func (g *SpringGenerator) ServicesControllers(version *spec.Version, mainPackage, thePackage, contentTypePackage, jsonPackage, modelsVersionPackage, serviceVersionPackage modules.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		serviceVersionSubpackage := serviceVersionPackage.Subpackage(api.Name.SnakeCase())
		files = append(files, g.serviceController(&api, thePackage, contentTypePackage, jsonPackage, modelsVersionPackage, serviceVersionSubpackage)...)
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

func (g *SpringGenerator) ExceptionController(responses *spec.Responses, thePackage, errorsPackage, errorsModelsPackage, jsonPackage modules.Module) *generator.CodeFile {
	w := writer.NewKotlinWriter()
	w.Line(`package %s`, thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.ServiceImports()...)
	imports.Add(`javax.servlet.http.HttpServletRequest`)
	imports.Add(`org.apache.tomcat.util.http.fileupload.FileUploadBase.CONTENT_TYPE`)
	imports.Add(jsonPackage.PackageStar)
	imports.Add(errorsModelsPackage.PackageStar)
	imports.Add(errorsPackage.Subpackage("ErrorsHelpers").Get("getBadRequestError"))
	imports.Add(errorsPackage.Subpackage("ErrorsHelpers").Get("getNotFoundError"))
	imports.Write(w)
	w.EmptyLine()
	w.Line(`@ControllerAdvice`)
	className := `ExceptionController`
	w.Line(`class %s(@Autowired private val json: Json) {`, className)
	w.Line(`  private val logger = LogManager.getLogger(%s::class.java)`, className)
	w.EmptyLine()
	g.errorHandler(w.Indented(), *responses)
	w.Line(`}`)

	return &generator.CodeFile{
		Path:    thePackage.GetPath(fmt.Sprintf("%s.kt", className)),
		Content: w.String(),
	}
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

func (g *SpringGenerator) serviceController(api *spec.Api, thePackage, contentTypePackage, jsonPackage, modelsVersionPackage, serviceVersionPackage modules.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	w := writer.NewKotlinWriter()
	w.Line(`package %s`, thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.ServiceImports()...)
	imports.Add(`javax.servlet.http.HttpServletRequest`)
	imports.Add(`org.apache.tomcat.util.http.fileupload.FileUploadBase.CONTENT_TYPE`)
	imports.Add(contentTypePackage.PackageStar)
	imports.Add(jsonPackage.Get("Json"))
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
	w.Line(`  @Autowired private val json: Json`)
	w.Line(`) {`)
	w.Line(`  private val logger = LogManager.getLogger(%s::class.java)`, className)

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
		w.Line(`checkContentType(request, MediaType.TEXT_PLAIN)`)
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`checkContentType(request, MediaType.APPLICATION_JSON)`)
		typ := g.Types.Kotlin(&operation.Body.Type.Definition)
		w.Line(`val %s: %s = json.read(%s)`, bodyJsonVar, typ, g.Models.JsonRead(bodyStringVar, &operation.Body.Type.Definition))
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

func (g *SpringGenerator) processResponse(w *generator.Writer, response *spec.Response, bodyVar string) {
	if response.BodyIs(spec.BodyEmpty) {
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s)`, response.Name.UpperCase())
		w.Line(`return ResponseEntity(HttpStatus.%s)`, response.Name.UpperCase())
	}
	if response.BodyIs(spec.BodyString) {
		w.Line(`val headers = HttpHeaders()`)
		w.Line(`headers.add(CONTENT_TYPE, "text/plain")`)
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s)`, response.Name.UpperCase())
		w.Line(`return ResponseEntity(%s, headers, HttpStatus.%s)`, bodyVar, response.Name.UpperCase())
	}
	if response.BodyIs(spec.BodyJson) {
		w.Line(`val bodyJson = json.write(%s)`, g.Models.JsonWrite(bodyVar, &response.Type.Definition))
		w.Line(`val headers = HttpHeaders()`)
		w.Line(`headers.add(CONTENT_TYPE, "application/json")`)
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s)`, response.Name.UpperCase())
		w.Line(`return ResponseEntity(bodyJson, headers, HttpStatus.%s)`, response.Name.UpperCase())
	}
}

func (g *SpringGenerator) ContentType(thePackage modules.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *contentTypeMismatchException(thePackage))
	files = append(files, *g.checkContentType(thePackage))
	return files
}

func (g *SpringGenerator) checkContentType(thePackage modules.Module) *generator.CodeFile {
	code := `
package [[.PackageName]]

import javax.servlet.http.HttpServletRequest
import org.springframework.http.MediaType

fun checkContentType(request: HttpServletRequest, expectedContentType: MediaType) {
	val contentType = request.getHeader("Content-Type")
	if (contentType == null || !contentType.contains(expectedContentType.toString())) {
		throw ContentTypeMismatchException(expectedContentType.toString(), contentType)
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

func (g *SpringGenerator) Errors(thePackage, errorsModelsPackage, contentTypePackage, jsonPackage modules.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *g.errorsHelpers(thePackage, errorsModelsPackage, contentTypePackage, jsonPackage))
	files = append(files, *g.Models.ValidationErrorsHelpers(thePackage, errorsModelsPackage, jsonPackage))
	return files
}

func (g *SpringGenerator) errorsHelpers(thePackage, errorsModelsPackage, contentTypePackage, jsonPackage modules.Module) *generator.CodeFile {
	code := `
package [[.PackageName]]

import org.springframework.web.bind.*
import org.springframework.web.bind.annotation.*
import org.springframework.web.method.annotation.MethodArgumentTypeMismatchException
import [[.ErrorsModelsPackage]].*
import [[.ContentTypePackage]].*
import [[.JsonPackage]].*
import [[.PackageName]].ValidationErrorsHelpers.extractValidationErrors

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
			val errors = extractValidationErrors(exception)
            return BadRequestError("Failed to parse body", ErrorLocation.BODY, errors)
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
		PackageName         string
		ErrorsModelsPackage string
		ContentTypePackage  string
		JsonPackage         string
	}{
		thePackage.PackageName,
		errorsModelsPackage.PackageName,
		contentTypePackage.PackageName,
		jsonPackage.PackageName,
	})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("ErrorsHelpers.kt"),
		Content: strings.TrimSpace(code),
	}
}

func (g *SpringGenerator) JsonHelpers(thePackage modules.Module) []generator.CodeFile {
	files := []generator.CodeFile{}

	files = append(files, *g.Json(thePackage))
	files = append(files, *jsonParseException(thePackage))
	files = append(files, g.Models.SetupLibrary(thePackage)...)

	return files
}

func (g *SpringGenerator) Json(thePackage modules.Module) *generator.CodeFile {
	w := writer.NewKotlinWriter()
	w.Line(`package %s`, thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(`org.springframework.beans.factory.annotation.Autowired`)
	imports.Add(`org.springframework.stereotype.Component`)
	imports.Add(g.Models.ModelsUsageImports()...)
	imports.Add(`java.io.IOException`)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`@Component`)
	className := `Json`
	w.Line(`class %s(%s) {`, className, g.Models.CreateJsonMapperField("Autowired"))
	w.Line(g.Models.JsonHelpersMethods())
	w.Line(`}`)

	return &generator.CodeFile{
		Path:    thePackage.GetPath(fmt.Sprintf("%s.kt", className)),
		Content: w.String(),
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
