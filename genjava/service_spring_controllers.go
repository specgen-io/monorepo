package genjava

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func (g *Generator) generateServicesControllers(version *spec.Version, thePackage Module, jsonPackage Module, modelsVersionPackage Module, serviceVersionPackage Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		serviceVersionSubpackage := serviceVersionPackage.Subpackage(api.Name.SnakeCase())
		files = append(files, g.generateController(&api, thePackage, jsonPackage, modelsVersionPackage, serviceVersionSubpackage)...)
	}
	return files
}

func (g *Generator) generateController(api *spec.Api, apiPackage Module, jsonPackage Module, modelsVersionPackage Module, serviceVersionPackage Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	w := NewJavaWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	w.Line(`import java.math.BigDecimal;`)
	w.Line(`import java.io.IOException;`)
	w.Line(`import java.time.*;`)
	w.Line(`import java.util.UUID;`)
	w.EmptyLine()
	w.Line(`import org.apache.logging.log4j.*;`)
	w.Line(`import org.springframework.beans.factory.annotation.Autowired;`)
	w.Line(`import org.springframework.format.annotation.DateTimeFormat;`)
	w.Line(`import org.springframework.http.*;`)
	w.Line(`import org.springframework.web.bind.annotation.*;`)
	w.Line(`import com.fasterxml.jackson.databind.ObjectMapper;`)
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
		g.generateMethod(w.Indented(), &operation)
	}
	w.Line(`}`)

	files = append(files, sources.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", className)),
		Content: w.String(),
	})

	return files
}

func (g *Generator) generateMethod(w *sources.Writer, operation *spec.NamedOperation) {
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	w.Line(`@%sMapping("%s")`, ToPascalCase(methodName), url)
	w.Line(`public ResponseEntity<String> %s(%s) throws IOException {`, controllerMethodName(operation), JoinDelimParams(addMethodParams(operation, g.Types)))
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
			w.Line(`  Message requestBody;`)
			w.Line(`  try {`)
			w.Line(`    requestBody = %s;`, g.Models.ReadJson("bodyStr", g.Types.JavaType(&operation.Body.Type.Definition)))
			w.Line(`  } catch (Exception e) {`)
			w.Line(`    logger.error("Completed request with status code: {}", HttpStatus.BAD_REQUEST);`)
			w.Line(`    return new ResponseEntity<>(headers, HttpStatus.BAD_REQUEST);`)
			w.Line(`  }`)
		}
	}
	serviceCall := fmt.Sprintf(`%s.%s(%s)`, serviceVarName(operation.Api), operation.Name.CamelCase(), JoinDelimParams(addServiceMethodParams(operation)))
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
			w.Line(`  if (result instanceof %s.%s) {`, serviceResponseInterfaceName(operation), resp.Name.PascalCase())
			if !resp.Type.Definition.IsEmpty() {
				responseWrite := g.Models.WriteJson(fmt.Sprintf(`((%s.%s) result).body`, serviceResponseInterfaceName(operation), resp.Name.PascalCase()))
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

	return fmt.Sprintf(`@%s(%s)`, paramAnnotationName, JoinDelimParams(annotationParams))
}

func generateMethodParam(namedParams []spec.NamedParam, paramAnnotationName string, types *Types) []string {
	params := []string{}

	if namedParams != nil && len(namedParams) > 0 {
		for _, param := range namedParams {
			paramAnnotation := getSpringParameterAnnotation(paramAnnotationName, &param)
			paramType := fmt.Sprintf(`%s %s`, types.JavaType(&param.Type.Definition), param.Name.CamelCase())
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

func addMethodParams(operation *spec.NamedOperation, types *Types) []string {
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
