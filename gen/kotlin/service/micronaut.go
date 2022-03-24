package service

import (
	"fmt"
	"github.com/pinzolo/casee"
	"github.com/specgen-io/specgen/v2/gen/kotlin/imports"
	"github.com/specgen-io/specgen/v2/gen/kotlin/models"
	"github.com/specgen-io/specgen/v2/gen/kotlin/modules"
	"github.com/specgen-io/specgen/v2/gen/kotlin/responses"
	"github.com/specgen-io/specgen/v2/gen/kotlin/types"
	"github.com/specgen-io/specgen/v2/gen/kotlin/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
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

func (m *MicronautGenerator) ServiceImplAnnotation(api *spec.Api) (annotationImport, annotation string) {
	return `io.micronaut.context.annotation.Bean`, `Bean`
}

func (m *MicronautGenerator) ServicesControllers(version *spec.Version, mainPackage, thePackage, jsonPackage, modelsVersionPackage, serviceVersionPackage modules.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		serviceVersionSubpackage := serviceVersionPackage.Subpackage(api.Name.SnakeCase())
		files = append(files, m.serviceController(&api, thePackage, jsonPackage, modelsVersionPackage, serviceVersionSubpackage)...)
	}
	files = append(files, dateConverters(mainPackage)...)
	return files
}

func (m *MicronautGenerator) serviceController(api *spec.Api, apiPackage, jsonPackage, modelsVersionPackage, serviceVersionPackage modules.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	w := writer.NewKotlinWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(`org.slf4j.*`)
	imports.Add(`io.micronaut.http.*`)
	imports.Add(`io.micronaut.http.annotation.*`)
	imports.Add(`jakarta.inject.Inject`)
	imports.Add(modelsVersionPackage.PackageStar)
	imports.Add(serviceVersionPackage.PackageStar)
	imports.Add(m.Models.JsonImports()...)
	imports.Add(m.Types.Imports()...)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`@Controller`)
	className := controllerName(api)
	w.Line(`class %s(`, className)
	w.Line(`  @Inject private val %s: %s,`, serviceVarName(api), serviceInterfaceName(api))
	w.Line(`  @Inject %s`, m.Models.CreateJsonMapperField())
	w.Line(`) {`)
	w.Line(`  private val logger = LoggerFactory.getLogger(%s::class.java)`, className)
	for _, operation := range api.Operations {
		w.EmptyLine()
		m.controllerMethod(w.Indented(), &operation)
	}
	w.Line(`}`)

	files = append(files, sources.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.kt", className)),
		Content: w.String(),
	})

	return files
}

func (m *MicronautGenerator) controllerMethod(w *sources.Writer, operation *spec.NamedOperation) {
	if operation.BodyIs(spec.BodyString) {
		w.Line(`@Consumes(MediaType.TEXT_PLAIN)`)
	}
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	w.Line(`@%s("%s")`, casee.ToPascalCase(methodName), url)
	w.Line(`fun %s(%s): MutableHttpResponse<*> {`, controllerMethodName(operation), joinParams(addMicronautMethodParams(operation, m.Types)))
	w.Line(`  logger.info("Received request, operationId: %s.%s, method: %s, url: %s")`, operation.Api.Name.Source, operation.Name.Source, methodName, url)
	w.EmptyLine()
	if operation.BodyIs(spec.BodyJson) {
		requestBody, exception := m.Models.ReadJson("bodyStr", &operation.Body.Type.Definition)
		w.Line(`  val requestBody = try {`)
		w.Line(`    %s`, requestBody)
		w.Line(`  } catch (e: %s) {`, exception)
		w.Line(`    logger.error("Completed request with status code: {}", HttpStatus.BAD_REQUEST)`)
		w.Line(`    return HttpResponse.status<Any>(HttpStatus.BAD_REQUEST)`)
		w.Line(`  }`)
	}
	serviceCall := fmt.Sprintf(`%s.%s(%s)`, serviceVarName(operation.Api), operation.Name.CamelCase(), joinParams(addServiceMethodParams(operation, "bodyStr", "requestBody!!")))
	if len(operation.Responses) == 1 && operation.Responses[0].BodyIs(spec.BodyEmpty) {
		w.Line(`  %s`, serviceCall)
	} else {
		w.Line(`  val result = %s`, serviceCall)
	}
	if len(operation.Responses) == 1 {
		m.processResponse(w.Indented(), &operation.Responses[0], "result")
	}
	if len(operation.Responses) > 1 {
		for _, response := range operation.Responses {
			w.Line(`  if (result is %s.%s) {`, responses.InterfaceName(operation), response.Name.PascalCase())
			m.processResponse(w.IndentedWith(2), &response, "result")
			w.Line(`  }`)
		}
		w.EmptyLine()
		w.Line(`  logger.error("Completed request with status code: {}", HttpStatus.INTERNAL_SERVER_ERROR)`)
		w.Line(`  return HttpResponse.status<Any>(HttpStatus.INTERNAL_SERVER_ERROR)`)
	}
	w.Line(`}`)
}

func (m *MicronautGenerator) processResponse(w *sources.Writer, response *spec.NamedResponse, result string) {
	if len(response.Operation.Responses) > 1 {
		result = fmt.Sprintf(`%s.body`, result)
	}
	if response.BodyIs(spec.BodyEmpty) {
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s)`, response.Name.UpperCase())
		w.Line(`return HttpResponse.status<Any>(HttpStatus.%s)`, response.Name.UpperCase())
	}
	if response.BodyIs(spec.BodyString) {
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s)`, response.Name.UpperCase())
		w.Line(`return HttpResponse.status<Any>(HttpStatus.%s).body(%s).contentType("text/plain")`, response.Name.UpperCase(), result)
	}
	if response.BodyIs(spec.BodyJson) {
		responseWrite, _ := m.Models.WriteJson(result, &response.Type.Definition)
		w.Line(`val responseJson = %s;`, responseWrite)
		w.EmptyLine()
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s)`, response.Name.UpperCase())
		w.Line(`return HttpResponse.status<Any>(HttpStatus.%s).body(responseJson).contentType("application/json")`, response.Name.UpperCase())
	}
}

func addMicronautMethodParams(operation *spec.NamedOperation, types *types.Types) []string {
	methodParams := []string{}

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

func dateConverters(thePackage modules.Module) []sources.CodeFile {
	convertersPackage := thePackage.Subpackage("converters")

	files := []sources.CodeFile{}
	files = append(files, *localDateConverter(convertersPackage))
	files = append(files, *localDateTimeConverter(convertersPackage))
	return files
}

func localDateConverter(thePackage modules.Module) *sources.CodeFile {
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

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("LocalDateConverter.kt"),
		Content: strings.TrimSpace(code),
	}
}

func localDateTimeConverter(thePackage modules.Module) *sources.CodeFile {
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

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("LocalDateTimeConverter.kt"),
		Content: strings.TrimSpace(code),
	}
}
