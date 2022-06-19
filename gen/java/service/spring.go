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

func (g *SpringGenerator) ServicesControllers(version *spec.Version, mainPackage, thePackage, jsonPackage, modelsVersionPackage, serviceVersionPackage packages.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		serviceVersionSubpackage := serviceVersionPackage.Subpackage(api.Name.SnakeCase())
		files = append(files, g.serviceController(&api, thePackage, jsonPackage, modelsVersionPackage, serviceVersionSubpackage)...)
	}
	return files
}

func (g *SpringGenerator) serviceController(api *spec.Api, apiPackage, jsonPackage, modelsVersionPackage, serviceVersionPackage packages.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	w := writer.NewJavaWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(`org.apache.logging.log4j.*`)
	imports.Add(`org.springframework.beans.factory.annotation.Autowired`)
	imports.Add(`org.springframework.format.annotation.DateTimeFormat`)
	imports.Add(`org.springframework.http.*`)
	imports.Add(`org.springframework.web.bind.annotation.*`)
	imports.Add(modelsVersionPackage.PackageStar)
	imports.Add(serviceVersionPackage.PackageStar)
	imports.Add(g.Models.ModelsUsageImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Add(`static org.apache.tomcat.util.http.fileupload.FileUploadBase.CONTENT_TYPE`)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`@RestController("%s")`, versionControllerName(controllerName(api), api.Http.Version))
	className := controllerName(api)
	w.Line(`public class %s {`, className)
	w.Line(`  private static final Logger logger = LogManager.getLogger(%s.class);`, className)
	w.EmptyLine()
	w.Line(`  @Autowired`)
	w.Line(`  private %s %s;`, serviceInterfaceName(api), serviceVarName(api))
	w.EmptyLine()
	g.Models.CreateJsonMapperField(w.Indented(), "@Autowired")
	for _, operation := range api.Operations {
		w.EmptyLine()
		g.controllerMethod(w.Indented(), &operation)
	}
	w.Line(`}`)

	files = append(files, sources.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", className)),
		Content: w.String(),
	})

	return files
}

func (g *SpringGenerator) controllerMethod(w *sources.Writer, operation *spec.NamedOperation) {
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	w.Line(`@%sMapping("%s")`, casee.ToPascalCase(methodName), url)
	w.Line(`public ResponseEntity<String> %s(%s) {`, controllerMethodName(operation), joinParams(addSpringMethodParams(operation, g.Types)))
	w.Indent()
	w.Line(`logger.info("Received request, operationId: %s.%s, method: %s, url: %s");`, operation.Api.Name.Source, operation.Name.Source, methodName, url)
	w.EmptyLine()
	if operation.BodyIs(spec.BodyJson) {
		requestBody, exception := g.Models.ReadJson("bodyStr", &operation.Body.Type.Definition)
		w.Line(`%s requestBody;`, g.Types.Java(&operation.Body.Type.Definition))
		w.Line(`try {`)
		w.Line(`  requestBody = %s;`, requestBody)
		w.Line(`} catch (%s e) {`, exception)
		g.badRequest(w.Indented(), operation, `"Failed to deserialize request body {}", e.getMessage()`)
		w.Line(`}`)
	}
	serviceCall := fmt.Sprintf(`%s.%s(%s)`, serviceVarName(operation.Api), operation.Name.CamelCase(), joinParams(addServiceMethodParams(operation, "bodyStr", "requestBody")))
	if len(operation.Responses) == 1 && operation.Responses[0].BodyIs(spec.BodyEmpty) {
		w.Line(`%s;`, serviceCall)
	} else {
		w.Line(`var result = %s;`, serviceCall)
		w.Line(`if (result == null) {`)
		g.internalServerError(w.Indented(), operation, `"Service implementation returned nil"`)
		w.Line(`}`)
	}
	if len(operation.Responses) == 1 {
		g.processResponse(w, &operation.Responses[0].Response, "result")
	}
	if len(operation.Responses) > 1 {
		for _, response := range operation.Responses {
			w.Line(`if (result instanceof %s.%s) {`, responses.InterfaceName(operation), response.Name.PascalCase())
			g.processResponse(w.Indented(), &response.Response, responses.GetBody(&response, "result"))
			w.Line(`}`)
		}
		w.EmptyLine()
		g.internalServerError(w, operation, `"No result returned from service implementation"`)
	}
	w.Unindent()
	w.Line(`}`)
}

func (g *SpringGenerator) processResponse(w *sources.Writer, response *spec.Response, result string) {
	if response.BodyIs(spec.BodyEmpty) {
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s);`, response.Name.UpperCase())
		w.Line(`return new ResponseEntity<>(HttpStatus.%s);`, response.Name.UpperCase())
	}
	if response.BodyIs(spec.BodyString) {
		w.Line(`HttpHeaders headers = new HttpHeaders();`)
		w.Line(`headers.add(CONTENT_TYPE, "text/plain");`)
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s);`, response.Name.UpperCase())
		w.Line(`return new ResponseEntity<>(%s, headers, HttpStatus.%s);`, result, response.Name.UpperCase())
	}
	if response.BodyIs(spec.BodyJson) {
		responseWrite, exception := g.Models.WriteJson(result, &response.Type.Definition)
		w.Line(`String responseJson = "";`)
		w.Line(`try {`)
		w.Line(`  responseJson = %s;`, responseWrite)
		w.Line(`} catch (%s e) {`, exception)
		w.Line(`  logger.error("Failed to serialize response body: {}", e.getMessage());`)
		w.Line(`}`)
		w.Line(`HttpHeaders headers = new HttpHeaders();`)
		w.Line(`headers.add(CONTENT_TYPE, "application/json");`)
		w.Line(`logger.info("Completed request with status code: {}", HttpStatus.%s);`, response.Name.UpperCase())
		w.Line(`return new ResponseEntity<>(responseJson, headers, HttpStatus.%s);`, response.Name.UpperCase())
	}
}

func (g *SpringGenerator) badRequest(w *sources.Writer, operation *spec.NamedOperation, message string) {
	w.Line(`logger.error(%s);`, message)
	w.Line(`logger.info("Completed request with status code: {}", HttpStatus.BAD_REQUEST);`)
	w.Line(`return new ResponseEntity<>(HttpStatus.BAD_REQUEST);`)
}

func (g *SpringGenerator) internalServerError(w *sources.Writer, operation *spec.NamedOperation, message string) {
	w.Line(`logger.error(%s);`, message)
	w.Line(`logger.info("Completed request with status code: {}", HttpStatus.INTERNAL_SERVER_ERROR);`)
	w.Line(`return new ResponseEntity<>(HttpStatus.INTERNAL_SERVER_ERROR);`)
}

func addSpringMethodParams(operation *spec.NamedOperation, types *types.Types) []string {
	methodParams := []string{}

	if operation.Body != nil {
		methodParams = append(methodParams, "@RequestBody String bodyStr")
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
			paramType := fmt.Sprintf(`%s %s`, types.Java(&param.Type.Definition), param.Name.CamelCase())
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
