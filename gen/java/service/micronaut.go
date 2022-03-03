package service

import (
	"fmt"
	"github.com/pinzolo/casee"
	"github.com/specgen-io/specgen/v2/gen/java/imports"
	"github.com/specgen-io/specgen/v2/gen/java/models"
	"github.com/specgen-io/specgen/v2/gen/java/packages"
	"github.com/specgen-io/specgen/v2/gen/java/responses"
	"github.com/specgen-io/specgen/v2/gen/java/types"
	"github.com/specgen-io/specgen/v2/gen/java/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
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

func (m *MicronautGenerator) ServicesControllers(version *spec.Version, thePackage packages.Module, jsonPackage packages.Module, modelsVersionPackage packages.Module, serviceVersionPackage packages.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		serviceVersionSubpackage := serviceVersionPackage.Subpackage(api.Name.SnakeCase())
		files = append(files, m.serviceController(&api, thePackage, jsonPackage, modelsVersionPackage, serviceVersionSubpackage)...)
	}
	return files
}

func (m *MicronautGenerator) serviceController(api *spec.Api, apiPackage packages.Module, jsonPackage packages.Module, modelsVersionPackage packages.Module, serviceVersionPackage packages.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	w := writer.NewJavaWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(`org.slf4j.*`)
	imports.Add(`io.micronaut.core.convert.format.Format`)
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
	w.Line(`public class %s {`, className)
	w.Line(`  private static final Logger logger = LoggerFactory.getLogger(%s.class);`, className)
	w.EmptyLine()
	w.Line(`  @Inject`)
	w.Line(`  private %s %s;`, serviceInterfaceName(api), serviceVarName(api))
	w.EmptyLine()
	w.Line(`  @Inject`)
	m.Models.CreateJsonMapperField(w.Indented())
	for _, operation := range api.Operations {
		w.EmptyLine()
		m.controllerMethod(w.Indented(), &operation)
	}
	w.Line(`}`)

	files = append(files, sources.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", className)),
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
	w.Line(`public MutableHttpResponse<?> %s(%s) throws IOException {`, controllerMethodName(operation), joinParams(addMicronautMethodParams(operation, m.Types)))
	w.Line(`  logger.info("Received request, operationId: %s.%s, method: %s, url: %s");`, operation.Api.Name.Source, operation.Name.Source, methodName, url)
	w.EmptyLine()
	if operation.BodyIs(spec.BodyJson) {
		bodyJavaType := m.Types.Java(&operation.Body.Type.Definition)
		requestBody, exception := m.Models.ReadJson("bodyStr", &operation.Body.Type.Definition)
		w.Line(`  %s requestBody;`, bodyJavaType)
		w.Line(`  try {`)
		w.Line(`    requestBody = %s;`, requestBody)
		w.Line(`  } catch (%s e) {`, exception)
		w.Line(`    logger.error("Completed request with status code: {}", HttpStatus.BAD_REQUEST);`)
		w.Line(`    return HttpResponse.status(HttpStatus.BAD_REQUEST).contentType("application/json");`)
		w.Line(`  }`)
	}
	serviceCall := fmt.Sprintf(`%s.%s(%s)`, serviceVarName(operation.Api), operation.Name.CamelCase(), joinParams(addServiceMethodParams(operation, "bodyStr", "requestBody")))
	if len(operation.Responses) == 1 && operation.Responses[0].BodyIs(spec.BodyEmpty) {
		w.Line(`  %s;`, serviceCall)
	} else {
		w.Line(`  var result = %s;`, serviceCall)
		w.Line(`  if (result == null) {`)
		w.Line(`    logger.error("Completed request with status code: {}", HttpStatus.INTERNAL_SERVER_ERROR);`)
		w.Line(`    return HttpResponse.status(HttpStatus.INTERNAL_SERVER_ERROR);`)
		w.Line(`  }`)
	}
	if len(operation.Responses) == 1 {
		m.processResponse(w.Indented(), &operation.Responses[0], "result")
	}
	if len(operation.Responses) > 1 {
		for _, response := range operation.Responses {
			w.Line(`  if (result instanceof %s.%s) {`, responses.InterfaceName(operation), response.Name.PascalCase())
			m.processResponse(w.IndentedWith(2), &response, "result")
			w.Line(`  }`)
		}
		w.EmptyLine()
		w.Line(`  logger.error("Completed request with status code: {}", HttpStatus.INTERNAL_SERVER_ERROR);`)
		w.Line(`  return HttpResponse.status(HttpStatus.INTERNAL_SERVER_ERROR);`)
	}
	w.Line(`}`)
}

func (m *MicronautGenerator) processResponse(w *sources.Writer, response *spec.NamedResponse, result string) {
	if len(response.Operation.Responses) > 1 {
		result = fmt.Sprintf(`((%s.%s) %s).body`, responses.InterfaceName(response.Operation), response.Name.PascalCase(), result)
	}
	if response.BodyIs(spec.BodyEmpty) {
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s);`, response.Name.UpperCase())
		w.Line(`return HttpResponse.status(HttpStatus.%s);`, response.Name.UpperCase())
	}
	if response.BodyIs(spec.BodyString) {
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s);`, response.Name.UpperCase())
		w.Line(`return HttpResponse.status(HttpStatus.%s).body(%s).contentType("text/plain");`, response.Name.UpperCase(), result)
	}
	if response.BodyIs(spec.BodyJson) {
		responseWrite, _ := m.Models.WriteJson(result, &response.Type.Definition)
		w.Line(`String responseJson = %s;`, responseWrite)
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s);`, response.Name.UpperCase())
		w.Line(`return HttpResponse.status(HttpStatus.%s).body(responseJson).contentType("application/json");`, response.Name.UpperCase())
	}
}

func addMicronautMethodParams(operation *spec.NamedOperation, types *types.Types) []string {
	methodParams := []string{}

	if operation.Body != nil {
		methodParams = append(methodParams, "@Body String bodyStr")
	}

	methodParams = append(methodParams, generateMicronautMethodParam(operation.QueryParams, "QueryValue", "value", types)...)
	methodParams = append(methodParams, generateMicronautMethodParam(operation.HeaderParams, "Header", "name", types)...)
	methodParams = append(methodParams, generateMicronautMethodParam(operation.Endpoint.UrlParams, "PathVariable", "name", types)...)

	return methodParams
}

func generateMicronautMethodParam(namedParams []spec.NamedParam, paramAnnotation string, paramAnnotationField string, types *types.Types) []string {
	params := []string{}

	if namedParams != nil && len(namedParams) > 0 {
		for _, param := range namedParams {
			paramAnnotation := getMicronautParameterAnnotation(paramAnnotation, paramAnnotationField, &param)
			paramType := fmt.Sprintf(`%s %s`, types.Java(&param.Type.Definition), param.Name.CamelCase())
			dateFormatAnnotation := dateFormatMicronautAnnotation(&param.Type.Definition)
			if dateFormatAnnotation != "" {
				params = append(params, fmt.Sprintf(`%s %s %s`, paramAnnotation, dateFormatAnnotation, paramType))
			} else {
				params = append(params, fmt.Sprintf(`%s %s`, paramAnnotation, paramType))
			}
		}
	}

	return params
}

func getMicronautParameterAnnotation(paramAnnotation string, paramAnnotationField string, param *spec.NamedParam) string {
	annotationParams := []string{fmt.Sprintf(`%s = "%s"`, paramAnnotationField, param.Name.Source)}

	if param.DefinitionDefault.Default != nil {
		annotationParams = append(annotationParams, fmt.Sprintf(`defaultValue = "%s"`, *param.DefinitionDefault.Default))
	}

	return fmt.Sprintf(`@%s(%s)`, paramAnnotation, joinParams(annotationParams))
}

func dateFormatMicronautAnnotation(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return dateFormatMicronautAnnotationPlain(typ.Plain)
	case spec.NullableType:
		return dateFormatMicronautAnnotation(typ.Child)
	case spec.ArrayType:
		return dateFormatMicronautAnnotation(typ.Child)
	case spec.MapType:
		return dateFormatMicronautAnnotation(typ.Child)
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func dateFormatMicronautAnnotationPlain(typ string) string {
	switch typ {
	case spec.TypeDate:
		return `@Format("yyyy-MM-dd")`
	case spec.TypeDateTime:
		return `@Format("yyyy-MM-dd'T'HH:mm:ss")`
	default:
		return ``
	}
}
