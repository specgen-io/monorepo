package gengo

import (
	"fmt"
	"github.com/pinzolo/casee"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
	"strings"
)

func GenerateService(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}
	files := []gen.TextFile{}

	for _, version := range specification.Versions {
		w := NewGoWriter()

		folder := "spec"
		if version.Version.Source != "" {
			folder += "_" + version.Version.FlatCase()
		}
		apiRoutersFile := generateRouting(w, &version, folder, generatePath)
		files = append(files, *apiRoutersFile)

	}
	return gen.WriteFiles(files, true)
}

func generateRouting(w *gen.Writer, version *spec.Version, packageName string, generatePath string) *gen.TextFile {
	w.Line("package %s", packageName)
	w.EmptyLine()
	w.Line("import (")
	w.Line(`  "encoding/json"`)
	w.Line(`  "fmt"`)
	w.Line(`  "github.com/husobee/vestigo"`)
	w.Line(`  "net/http"`)
	w.Line(`)`)
	generateCheckErrors(w)
	for _, api := range version.Http.Apis {
		generateApiRouter(w, api)
	}
	return &gen.TextFile{
		Path:    filepath.Join(generatePath, packageName, "routing.go"),
		Content: w.String(),
	}
}

func generateCheckErrors(w *gen.Writer) {
	w.Line(`
func checkErrors(params *ParamsParser, w http.ResponseWriter) bool {
	if len(params.Errors) > 0 {
		w.WriteHeader(400)
		fmt.Println(params.Errors)
		return false
		}
	return true
}
`)
}

func generateApiRouter(w *gen.Writer, api spec.Api) {
	apiName := api.Name.PascalCase()
	w.Line(`func Add%sRoutes(router *vestigo.Router, %sService I%sService) {`, apiName, api.Name.Source, apiName)
	for _, operation := range api.Operations {
		reminder := operation.FullUrl()
		urlParams := []string{}
		for _, param := range operation.Endpoint.UrlParams {
			paramName := param.Name.Source
			parts := strings.Split(reminder, spec.UrlParamStr(paramName))
			urlParams = append(urlParams, fmt.Sprintf("%s:%s", parts[0], paramName))
			reminder = parts[1]
		}
		w.Indent()
		if reminder != `` {
			w.Line(`router.%s("%s", func(w http.ResponseWriter, r *http.Request) {`, nameToPascalCase(operation.Endpoint.Method), reminder)
		} else {
			w.Line(`router.SetCors("/echo/header", &vestigo.CorsAccessControl{`)
			w.Line(`  AllowHeaders: []string{"Int-Header", "String-Header"},`)
			w.Line(`})`)

			w.Line(`router.%s("%s", func(w http.ResponseWriter, r *http.Request) {`, nameToPascalCase(operation.Endpoint.Method), strings.Join(urlParams, ""))
		}

		parts := strings.Split(operation.Name.Source, "_")
		operationName := parts[1]
		for _, response := range operation.Responses {
			if strings.Contains(operation.Endpoint.Method, "POST") {
				generatePostMethod(w, response, operationName, api.Name.Source, operation.Name.PascalCase())
			} else if strings.Contains(operation.Endpoint.Method, "GET") {
				generateGetMethod(w, response, operation.Name.Source, operationName, operation)
			}
		}

		w.Unindent()
		w.Line(`  })`)
	}
	w.Line(`}`)
	w.EmptyLine()
}

func nameToPascalCase(name string) string {
	return casee.ToPascalCase(name)
}

func generatePostMethod(w *gen.Writer, response spec.NamedResponse, operationName string, apiName string, operation string) {
	if !response.Type.Definition.IsEmpty() {
		w.Line("  var %s %s", operationName, GoType(&response.Type.Definition))
	}
	w.Line("  json.NewDecoder(r.Body).Decode(&%s)", operationName)
	w.Line("  response := %sService.%s(&%s)", apiName, operation, operationName)
	w.Line("  w.WriteHeader(%s)", spec.HttpStatusCode(response.Name))
	w.Line("  json.NewEncoder(w).Encode(response.%s)", response.Name.PascalCase())
}

func generateGetMethod(w *gen.Writer, response spec.NamedResponse, namedOperation string, operationName string, operation spec.NamedOperation) {
	var paramsName string
	urlParams := []string{}
	if strings.Contains(namedOperation, "query") {
		w.Line("  %s := NewParamsParser(r.URL.Query())", operationName)
		for _, params := range operation.QueryParams {
			queryParams := params.Name.Source
			paramsName = params.Name.CamelCase()
			w.Line(`  %s := %s.%s("%s")`, paramsName, operationName, nameToPascalCase(params.Type.Definition.Name), queryParams)
			urlParams = append(urlParams, fmt.Sprintf("%s", paramsName))
		}
		w.Line("  if checkErrors(%s, w) {", operationName)
	} else if strings.Contains(namedOperation, "header") {
		w.Line("  %s := NewParamsParser(r.URL.Query())", operationName)
		for _, params := range operation.HeaderParams {
			headerParams := params.Name.Source
			paramsName = params.Name.CamelCase()
			w.Line(`  %s := %s.%s("%s")`, paramsName, operationName, nameToPascalCase(params.Type.Definition.Name), headerParams)
			urlParams = append(urlParams, fmt.Sprintf("%s", paramsName))
		}
		w.Line("  if checkErrors(%s, w) {", operationName)
	} else if strings.Contains(namedOperation, "url_params") {
		w.Line("  %s := NewParamsParser(r.URL.Query())", "query")
		for _, param := range operation.Endpoint.UrlParams {
			queryParams := param.Name.Source
			paramsName = param.Name.CamelCase()
			w.Line(`  %s := %s.%s(":%s")`, paramsName, "query", nameToPascalCase(param.Type.Definition.Name), queryParams)
			urlParams = append(urlParams, fmt.Sprintf("%s", paramsName))
		}
		w.Line("  if checkErrors(%s, w) {", "query")
	}
	w.Line("    response := echoService.%s(%s)", operation.Name.PascalCase(), strings.Join(urlParams, ", "))
	w.Line("  w.WriteHeader(%s)", spec.HttpStatusCode(response.Name))
	w.Line("  json.NewEncoder(w).Encode(response.%s)", response.Name.PascalCase())
	w.Line("  }")
}
