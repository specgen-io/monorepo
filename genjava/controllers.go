package genjava

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
)

func generateServicesControllers(version *spec.Version, thePackage Module, modelsPackage Module, modelsVersionPackage Module, serviceVersionPackage Module) []gen.TextFile {
	files := []gen.TextFile{}
	for _, api := range version.Http.Apis {
		serviceVersionSubpackage := serviceVersionPackage.Subpackage(api.Name.SnakeCase())
		files = append(files, generateController(version, &api, thePackage, modelsPackage, modelsVersionPackage, serviceVersionSubpackage)...)
	}
	return files
}

func generateController(version *spec.Version, api *spec.Api, apiPackage Module, modelsPackage Module, modelsVersionPackage Module, serviceVersionPackage Module) []gen.TextFile {
	files := []gen.TextFile{}
	w := NewJavaWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	w.Line(`import com.fasterxml.jackson.databind.ObjectMapper;`)
	w.Line(`import java.math.BigDecimal;`)
	w.Line(`import java.io.IOException;`)
	w.Line(`import java.time.*;`)
	w.Line(`import java.util.UUID;`)
	w.Line(`import org.springframework.format.annotation.DateTimeFormat;`)
	w.Line(`import org.springframework.http.*;`)
	w.Line(`import org.springframework.web.bind.annotation.*;`)
	w.EmptyLine()
	w.Line(`import static org.apache.tomcat.util.http.fileupload.FileUploadBase.CONTENT_TYPE;`)
	w.EmptyLine()
	w.Line(`import %s.Jsoner;`, modelsPackage.PackageName)
	w.Line(`import %s;`, modelsVersionPackage.PackageStar)
	w.Line(`import %s;`, serviceVersionPackage.PackageStar)
	w.EmptyLine()
	w.Line(`@RestController("%s")`, versionControllerName(controllerName(api), version))
	w.Line(`public class %s {`, controllerName(api))
	w.Line(`  final %s %s;`, serviceInterfaceName(api), serviceVarName(api))
	w.EmptyLine()
	w.Line(`  public %s(%s %s) {`, controllerName(api), serviceInterfaceName(api), serviceVarName(api))
	w.Line(`    this.%s = %s;`, serviceVarName(api), serviceVarName(api))
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  ObjectMapper objectMapper = new ObjectMapper();`)
	for _, operation := range api.Operations {
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
	w.EmptyLine()
	w.Line(`@%sMapping("%s")`, ToPascalCase(operation.Endpoint.Method), versionUrl(version, operation.Endpoint.Url))
	w.Line(`public ResponseEntity<String> %s(%s) throws IOException {`, controllerMethodName(operation), JoinParams(addMethodParams(operation)))
	if operation.Body != nil {
		w.Line(`  var requestBody = Jsoner.deserialize(objectMapper, jsonStr, %s.class);`, JavaType(&operation.Body.Type.Definition))
	}
	response := collectResponses(operation)
	if len(response) == 1 {
		for _, resp := range operation.Responses {
			if resp.Type.Definition.IsEmpty() {
				w.Line(`  %s.%s(%s);`, serviceVarName(api), operation.Name.CamelCase(), JoinParams(addServiceMethodParams(operation)))
				w.EmptyLine()
				w.Line(`  return new ResponseEntity<>(HttpStatus.OK);`)
			} else {
				w.Line(`  var result = %s.%s(%s);`, serviceVarName(api), operation.Name.CamelCase(), JoinParams(addServiceMethodParams(operation)))
				w.EmptyLine()
				w.Line(`  HttpHeaders headers = new HttpHeaders();`)
				w.Line(`  headers.add(CONTENT_TYPE, "application/json");`)
				w.Line(`  String responseJson = Jsoner.serialize(objectMapper, result);`)
				w.EmptyLine()
				w.Line(`  return new ResponseEntity<>(responseJson, headers, HttpStatus.%s);`, resp.Name.UpperCase())
			}
		}
	}
	if len(response) > 1 {
		w.Line(`  var result = %s.%s(%s);`, serviceVarName(api), operation.Name.CamelCase(), JoinParams(addServiceMethodParams(operation)))
		w.EmptyLine()
		for _, resp := range operation.Responses {
			w.Line(`  if (result instanceof %s) {`, serviceResponseImplName(operation, resp))
			w.Line(`    return new ResponseEntity<>(HttpStatus.%s);`, resp.Name.UpperCase())
			w.Line(`  }`)
		}
		w.EmptyLine()
		w.Line(`  return new ResponseEntity<>(HttpStatus.INTERNAL_SERVER_ERROR);`)
	}
	w.Line(`}`)
}

func addMethodParams(operation spec.NamedOperation) []string {
	methodParams := []string{}
	if operation.Body != nil {
		methodParams = append(methodParams, "@RequestBody String jsonStr")
	}
	for _, param := range operation.QueryParams {
		methodParams = append(methodParams, fmt.Sprintf(`@RequestParam("%s")%s %s %s`, param.Name.Source, addDateFormatAnnotation(&param.Type.Definition), JavaType(&param.Type.Definition), param.Name.CamelCase()))
	}
	for _, param := range operation.HeaderParams {
		methodParams = append(methodParams, fmt.Sprintf(`@RequestHeader("%s")%s %s %s`, param.Name.Source, addDateFormatAnnotation(&param.Type.Definition), JavaType(&param.Type.Definition), param.Name.CamelCase()))
	}
	for _, param := range operation.Endpoint.UrlParams {
		methodParams = append(methodParams, fmt.Sprintf(`@PathVariable("%s")%s %s %s`, param.Name.Source, addDateFormatAnnotation(&param.Type.Definition), JavaType(&param.Type.Definition), param.Name.CamelCase()))
	}
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
