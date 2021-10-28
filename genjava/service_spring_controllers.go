package genjava

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
)

func generateServicesControllers(version *spec.Version, thePackage Module, jsonPackage Module, modelsVersionPackage Module, serviceVersionPackage Module) []gen.TextFile {
	files := []gen.TextFile{}
	for _, api := range version.Http.Apis {
		serviceVersionSubpackage := serviceVersionPackage.Subpackage(api.Name.SnakeCase())
		files = append(files, generateController(version, &api, thePackage, jsonPackage, modelsVersionPackage, serviceVersionSubpackage)...)
	}
	return files
}

func generateController(version *spec.Version, api *spec.Api, apiPackage Module, jsonPackage Module, modelsVersionPackage Module, serviceVersionPackage Module) []gen.TextFile {
	files := []gen.TextFile{}
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
	w.Line(`@RestController("%s")`, versionControllerName(controllerName(api), version))
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
		generateMethod(w.Indented(), version, api, operation)
	}
	w.Line(`}`)

	files = append(files, gen.TextFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", controllerName(api))),
		Content: w.String(),
	})

	return files
}

func generateMethod(w *gen.Writer, version *spec.Version, api *spec.Api, operation spec.NamedOperation) {
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	w.Line(`@%sMapping("%s")`, ToPascalCase(methodName), url)
	w.Line(`public ResponseEntity<String> %s(%s) throws IOException {`, controllerMethodName(operation), JoinParams(addMethodParams(operation)))
	w.Line(`  logger.info("Received request, operationId: %s.%s, method: %s, url: %s");`, operation.Api.Name.Source, operation.Name.Source, methodName, url)
	w.Line(`  HttpHeaders headers = new HttpHeaders();`)
	w.Line(`  headers.add(CONTENT_TYPE, "application/json");`)
	w.EmptyLine()
	if operation.Body != nil {
		w.Line(`  Message requestBody;`)
		w.Line(`  try {`)
		w.Line(`    requestBody = objectMapper.readValue(jsonStr, %s.class);`, JavaType(&operation.Body.Type.Definition))
		w.Line(`  } catch (Exception e) {`)
		w.Line(`    logger.error("Completed request with status code: {}", HttpStatus.BAD_REQUEST);`)
		w.Line(`    return new ResponseEntity<>(headers, HttpStatus.BAD_REQUEST);`)
		w.Line(`  }`)
	}
	serviceCall := fmt.Sprintf(`%s.%s(%s)`, serviceVarName(api), operation.Name.CamelCase(), JoinParams(addServiceMethodParams(operation)))
	if len(operation.Responses) == 1 {
		for _, resp := range operation.Responses {
			if resp.Type.Definition.IsEmpty() {
				w.Line(`  %s;`, serviceCall)
				w.EmptyLine()
				w.Line(`  logger.info("Completed request with status code: {}", HttpStatus.%s);`, resp.Name.UpperCase())
				w.Line(`  return new ResponseEntity<>(HttpStatus.%s);`, resp.Name.UpperCase())
			} else {
				w.Line(`  var result = %s;`, serviceCall)
				w.Line(`  String responseJson = objectMapper.writeValueAsString(result);`)
				w.EmptyLine()
				w.Line(`  logger.info("Completed request with status code: {}", HttpStatus.%s);`, resp.Name.UpperCase())
				w.Line(`  return new ResponseEntity<>(responseJson, headers, HttpStatus.%s);`, resp.Name.UpperCase())
			}
		}
	}
	if len(operation.Responses) > 1 {
		w.Line(`  var result = %s;`, serviceCall)
		w.EmptyLine()
		for _, resp := range operation.Responses {
			w.Line(`  if (result instanceof %s) {`, serviceResponseImplName(operation, resp))
			w.Line(`    logger.info("Completed request with status code: {}", HttpStatus.%s);`, resp.Name.UpperCase())
			w.Line(`    return new ResponseEntity<>(HttpStatus.%s);`, resp.Name.UpperCase())
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

	return fmt.Sprintf(`@%s(%s)`, paramAnnotationName, JoinParams(annotationParams))
}

func generateMethodParam(namedParams []spec.NamedParam, paramAnnotationName string) []string {
	params := []string{}

	if namedParams != nil && len(namedParams) > 0 {
		for _, param := range namedParams {
			paramAnnotation := getSpringParameterAnnotation(paramAnnotationName, &param)
			paramType := fmt.Sprintf(`%s %s`, JavaType(&param.Type.Definition), param.Name.CamelCase())
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

func addMethodParams(operation spec.NamedOperation) []string {
	methodParams := []string{}

	if operation.Body != nil {
		methodParams = append(methodParams, "@RequestBody String jsonStr")
	}
	methodParams = append(methodParams, generateMethodParam(operation.QueryParams, "RequestParam")...)
	methodParams = append(methodParams, generateMethodParam(operation.HeaderParams, "RequestHeader")...)
	methodParams = append(methodParams, generateMethodParam(operation.Endpoint.UrlParams, "PathVariable")...)

	return methodParams
}

func addServiceMethodParams(operation spec.NamedOperation) []string {
	methodParams := []string{}
	if operation.Body != nil {
		methodParams = append(methodParams, "requestBody")
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
