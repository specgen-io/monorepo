package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/genjava/packages"
	"github.com/specgen-io/specgen/v2/genjava/responses"
	"github.com/specgen-io/specgen/v2/genjava/types"
	"github.com/specgen-io/specgen/v2/genjava/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func (g *Generator) ServicesControllers(version *spec.Version, thePackage packages.Module, jsonPackage packages.Module, modelsVersionPackage packages.Module, serviceVersionPackage packages.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		serviceVersionSubpackage := serviceVersionPackage.Subpackage(api.Name.SnakeCase())
		files = append(files, g.serviceController(&api, thePackage, jsonPackage, modelsVersionPackage, serviceVersionSubpackage)...)
	}
	return files
}

func (g *Generator) serviceController(api *spec.Api, apiPackage packages.Module, jsonPackage packages.Module, modelsVersionPackage packages.Module, serviceVersionPackage packages.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	w := writer.NewJavaWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	w.Line(`import java.math.BigDecimal;`)
	w.Line(`import java.io.IOException;`)
	w.Line(`import java.time.*;`)
	w.Line(`import java.util.*;`)
	w.EmptyLine()
	w.Line(`import org.apache.logging.log4j.*;`)
	w.Line(`import org.springframework.beans.factory.annotation.Autowired;`)
	w.Line(`import org.springframework.format.annotation.DateTimeFormat;`)
	w.Line(`import org.springframework.http.*;`)
	w.Line(`import org.springframework.web.bind.annotation.*;`)
	w.Line(`import com.fasterxml.jackson.databind.ObjectMapper;`)
	w.Line(`import com.fasterxml.jackson.core.type.TypeReference;`)
	w.EmptyLine()
	w.Line(`import static org.apache.tomcat.util.http.fileupload.FileUploadBase.CONTENT_TYPE;`)
	w.EmptyLine()
	w.Line(`import %s.Json;`, jsonPackage.PackageName)
	w.Line(`import %s;`, modelsVersionPackage.PackageStar)
	w.Line(`import %s;`, serviceVersionPackage.PackageStar)
	w.EmptyLine()
	w.Line(`@RestController("%s")`, versionControllerName(controllerName(api), api.Apis.Version))
	className := controllerName(api)
	w.Line(`public class %s {`, className)
	w.Line(`  private static final Logger logger = LogManager.getLogger(%s.class);`, className)
	w.EmptyLine()
	w.Line(`  @Autowired`)
	w.Line(`  private %s %s;`, serviceInterfaceName(api), serviceVarName(api))
	w.EmptyLine()
	w.Line(`  @Autowired`)
	w.Line(`  private ObjectMapper objectMapper;`)
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

func (g *Generator) controllerMethod(w *sources.Writer, operation *spec.NamedOperation) {
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	w.Line(`@%sMapping("%s")`, ToPascalCase(methodName), url)
	w.Line(`public ResponseEntity<String> %s(%s) throws IOException {`, controllerMethodName(operation), joinParams(addMethodParams(operation, g.Types)))
	w.Line(`  logger.info("Received request, operationId: %s.%s, method: %s, url: %s");`, operation.Api.Name.Source, operation.Name.Source, methodName, url)
	w.Line(`  HttpHeaders headers = new HttpHeaders();`)
	contentType := "application/json"
	if operation.Body != nil && operation.Body.Type.Definition.Plain == spec.TypeString {
		contentType = "text/plain"
	}
	w.Line(`  headers.add(CONTENT_TYPE, "%s");`, contentType)
	w.EmptyLine()
	if operation.Body != nil {
		if operation.Body.Type.Definition.Plain != spec.TypeString {
			bodyJavaType := g.Types.Java(&operation.Body.Type.Definition)
			w.Line(`  %s requestBody;`, bodyJavaType)
			valueTypeName := fmt.Sprintf("%s.class", bodyJavaType)
			if operation.Body.Type.Definition.Node == spec.MapType {
				valueTypeName = "typeRef"
				w.Line(`  TypeReference<%s> typeRef = new TypeReference<%s>() {};`, bodyJavaType, bodyJavaType)
			}
			w.Line(`  try {`)
			w.Line(`    requestBody = %s;`, g.Models.ReadJson("bodyStr", valueTypeName))
			w.Line(`  } catch (Exception e) {`)
			w.Line(`    logger.error("Completed request with status code: {}", HttpStatus.BAD_REQUEST);`)
			w.Line(`    return new ResponseEntity<>(headers, HttpStatus.BAD_REQUEST);`)
			w.Line(`  }`)
		}
	}
	serviceCall := fmt.Sprintf(`%s.%s(%s)`, serviceVarName(operation.Api), operation.Name.CamelCase(), joinParams(addServiceMethodParams(operation)))
	if len(operation.Responses) == 1 {
		for _, resp := range operation.Responses {
			if resp.Type.Definition.IsEmpty() {
				w.Line(`  %s;`, serviceCall)
				w.EmptyLine()
				w.Line(`  logger.info("Completed request with status code: {}", HttpStatus.%s);`, resp.Name.UpperCase())
				w.Line(`  return new ResponseEntity<>(HttpStatus.%s);`, resp.Name.UpperCase())
			} else {
				w.Line(`  var result = %s;`, serviceCall)
				w.Line(`  if (result == null) {`)
				w.Line(`    logger.error("Completed request with status code: {}", HttpStatus.INTERNAL_SERVER_ERROR);`)
				w.Line(`    return new ResponseEntity<>(HttpStatus.INTERNAL_SERVER_ERROR);`)
				w.Line(`  }`)
				responseVar := "result"
				if resp.Type.Definition.Plain != spec.TypeString {
					responseVar = "responseJson"
					w.Line(`  String %s = %s;`, responseVar, g.Models.WriteJson("result"))
				}
				w.EmptyLine()
				w.Line(`  logger.info("Completed request with status code: {}", HttpStatus.%s);`, resp.Name.UpperCase())
				w.Line(`  return new ResponseEntity<>(%s, headers, HttpStatus.%s);`, responseVar, resp.Name.UpperCase())
			}
		}
	}
	if len(operation.Responses) > 1 {
		w.Line(`  var result = %s;`, serviceCall)
		w.Line(`  if (result == null) {`)
		w.Line(`    logger.error("Completed request with status code: {}", HttpStatus.INTERNAL_SERVER_ERROR);`)
		w.Line(`    return new ResponseEntity<>(HttpStatus.INTERNAL_SERVER_ERROR);`)
		w.Line(`  }`)
		for _, resp := range operation.Responses {
			w.EmptyLine()
			w.Line(`  if (result instanceof %s.%s) {`, responses.InterfaceName(operation), resp.Name.PascalCase())
			if !resp.Type.Definition.IsEmpty() {
				responseWrite := g.Models.WriteJson(fmt.Sprintf(`((%s.%s) result).body`, responses.InterfaceName(operation), resp.Name.PascalCase()))
				w.Line(`    String responseJson = %s;`, responseWrite)
				w.Line(`    logger.info("Completed request with status code: {}", HttpStatus.%s);`, resp.Name.UpperCase())
				w.Line(`    return new ResponseEntity<>(responseJson, headers, HttpStatus.%s);`, resp.Name.UpperCase())
			} else {
				w.Line(`    logger.info("Completed request with status code: {}", HttpStatus.%s);`, resp.Name.UpperCase())
				w.Line(`    return new ResponseEntity<>(HttpStatus.%s);`, resp.Name.UpperCase())
			}
			w.Line(`  }`)
		}
		w.EmptyLine()
		w.Line(`  logger.error("Completed request with status code: {}", HttpStatus.INTERNAL_SERVER_ERROR);`)
		w.Line(`  return new ResponseEntity<>(HttpStatus.INTERNAL_SERVER_ERROR);`)
	}
	w.Line(`}`)
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

func generateMethodParam(namedParams []spec.NamedParam, paramAnnotationName string, types *types.Types) []string {
	params := []string{}

	if namedParams != nil && len(namedParams) > 0 {
		for _, param := range namedParams {
			paramAnnotation := getSpringParameterAnnotation(paramAnnotationName, &param)
			paramType := fmt.Sprintf(`%s %s`, types.Java(&param.Type.Definition), param.Name.CamelCase())
			dateFormatAnnotation := dateFormatAnnotation(&param.Type.Definition)
			if dateFormatAnnotation != "" {
				params = append(params, fmt.Sprintf(`%s %s %s`, paramAnnotation, dateFormatAnnotation, paramType))
			} else {
				params = append(params, fmt.Sprintf(`%s %s`, paramAnnotation, paramType))
			}
		}
	}

	return params
}

func addMethodParams(operation *spec.NamedOperation, types *types.Types) []string {
	methodParams := []string{}

	if operation.Body != nil {
		methodParams = append(methodParams, "@RequestBody String bodyStr")
	}
	methodParams = append(methodParams, generateMethodParam(operation.QueryParams, "RequestParam", types)...)
	methodParams = append(methodParams, generateMethodParam(operation.HeaderParams, "RequestHeader", types)...)
	methodParams = append(methodParams, generateMethodParam(operation.Endpoint.UrlParams, "PathVariable", types)...)

	return methodParams
}

func addServiceMethodParams(operation *spec.NamedOperation) []string {
	methodParams := []string{}
	if operation.Body != nil {
		if operation.Body.Type.Definition.Plain == spec.TypeString {
			methodParams = append(methodParams, "bodyStr")
		} else {
			methodParams = append(methodParams, "requestBody")
		}
	}
	for _, param := range operation.QueryParams {
		methodParams = append(methodParams, param.Name.CamelCase())
	}
	for _, param := range operation.HeaderParams {
		methodParams = append(methodParams, param.Name.CamelCase())
	}
	for _, param := range operation.Endpoint.UrlParams {
		methodParams = append(methodParams, param.Name.CamelCase())
	}
	return methodParams
}

func dateFormatAnnotation(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return dateFormatAnnotationPlain(typ.Plain)
	case spec.NullableType:
		return dateFormatAnnotation(typ.Child)
	case spec.ArrayType:
		return dateFormatAnnotation(typ.Child)
	case spec.MapType:
		return dateFormatAnnotation(typ.Child)
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func dateFormatAnnotationPlain(typ string) string {
	switch typ {
	case spec.TypeDate:
		return "@DateTimeFormat(iso = DateTimeFormat.ISO.DATE)"
	case spec.TypeDateTime:
		return "@DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME)"
	default:
		return ""
	}
}
