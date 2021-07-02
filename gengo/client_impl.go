package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
	"strings"
)

func generateClientsImplementations(version *spec.Version, packageName string, generatePath string) []gen.TextFile {
	files := []gen.TextFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *generateClientImplementation(&api, packageName, generatePath))
	}
	return files
}

func generateClientImplementation(api *spec.Api, packageName string, generatePath string) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", packageName)

	imports := []string{}
	imports = append(imports, `import "fmt"`)
	imports = append(imports, `import "errors"`)
	imports = append(imports, `import "io/ioutil"`)
	imports = append(imports, `import "net/http"`)
	imports = append(imports, `import "encoding/json"`)

	if bodyHasType(api) {
		imports = append(imports, `import "bytes"`)
	}
	if operationHasType(api, spec.TypeDate) {
		imports = append(imports, fmt.Sprintf(`import "cloud.google.com/go/civil"`))
	}
	if operationHasType(api, spec.TypeUuid) {
		imports = append(imports, fmt.Sprintf(`import "github.com/google/uuid"`))
	}
	if operationHasType(api, spec.TypeDecimal) {
		imports = append(imports, fmt.Sprintf(`import "github.com/shopspring/decimal"`))
	}

	//TODO: Imports will never be empty according to code above
	if len(imports) > 0 {
		w.EmptyLine()
		w.Line(`%s`, strings.Join(imports, "\n"))
	}

	w.EmptyLine()
	generateClientTypeStruct(w, *api)
	w.EmptyLine()
	generateNewClient(w, *api)
	for _, operation := range api.Operations {
		w.EmptyLine()
		generateClientFunction(w, api.Name, operation)
	}

	return &gen.TextFile{
		Path:    filepath.Join(generatePath, fmt.Sprintf("%s_client.go", api.Name.SnakeCase())),
		Content: w.String(),
	}
}

func generateClientTypeStruct(w *gen.Writer, api spec.Api) {
	w.Line(`type %sClient struct {`, api.Name.SnakeCase())
	w.Line(`  baseUrl string`)
	w.Line(`}`)
}

func generateNewClient(w *gen.Writer, api spec.Api) {
	w.Line(`func New%sClient(baseUrl string) *%sClient {`, api.Name.PascalCase(), api.Name.SnakeCase())
	w.Line(`  return &%sClient{baseUrl}`, api.Name.SnakeCase())
	w.Line(`}`)
}

func generateClientFunction(w *gen.Writer, apiName spec.Name, operation spec.NamedOperation) {
	w.Line(`func (client *%sClient) %s(%s) (*%sResponse, error) {`, apiName.SnakeCase(), operation.Name.PascalCase(), strings.Join(addParams(operation), ", "), operation.Name.PascalCase())
	var body = "nil"
	if operation.Body != nil {
		w.Line(`  bodyJSON, err := json.Marshal(body)`)
		body = "bytes.NewBuffer(bodyJSON)"
	}
	w.Line(`  req, err := http.NewRequest("%s", client.baseUrl+%s, %s)`, operation.Endpoint.Method, addRequestUrlParams(operation), body)
	addCheckError(w)
	w.EmptyLine()
	parseParams(w, operation)
	w.Line(`  resp, err := http.DefaultClient.Do(req)`)
	addCheckError(w)
	w.Indent()
	addResponse(w, operation)
	w.Unindent()
	w.EmptyLine()
	w.Line(`  return nil, errors.New(fmt.Sprintf("Unexpected status code received: %s", resp.StatusCode))`, "%d")
	w.Line(`}`)
}

func addCheckError(w *gen.Writer) {
	w.Line(`  if err != nil { return nil, err }`)
}

func getUrl(operation spec.NamedOperation) []string {
	reminder := operation.FullUrl()
	urlParams := []string{}
	if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
		for _, param := range operation.Endpoint.UrlParams {
			parts := strings.Split(reminder, spec.UrlParamStr(param.Name.Source))
			urlParams = append(urlParams, fmt.Sprintf("%s%s", parts[0], "%s"))
			reminder = parts[1]
		}
	}
	urlParams = append(urlParams, fmt.Sprintf("%s", reminder))
	return urlParams
}

func addRequestUrlParams(operation spec.NamedOperation) string {
	if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
		return fmt.Sprintf(`fmt.Sprintf("%s", %s)`, strings.Join(getUrl(operation), ""), strings.Join(addUrlParam(operation), ", "))
	} else {
		return fmt.Sprintf(`"%s"`, operation.Endpoint.Url)
	}
}

func addUrlParam(operation spec.NamedOperation) []string {
	urlParams := []string{}
	for _, param := range operation.Endpoint.UrlParams {
		if GoType(&param.Type.Definition) != "string" {
			urlParams = append(urlParams, fmt.Sprintf("convert%s(%s)", parserMethodName(&param.Type.Definition), param.Name.CamelCase()))
		} else {
			urlParams = append(urlParams, fmt.Sprintf("%s", param.Name.CamelCase()))
		}
	}
	return urlParams
}

func parseParams(w *gen.Writer, operation spec.NamedOperation) {
	if operation.QueryParams != nil && len(operation.QueryParams) > 0 {
		w.Line(`  query := req.URL.Query()`)
		addParsedParams(w, operation.QueryParams, "q", "query")
		w.Line(`  req.URL.RawQuery = %s.Encode()`, "query")
		w.EmptyLine()
	}
	if operation.HeaderParams != nil && len(operation.HeaderParams) > 0 {
		w.Line(`  header := req.Header`)
		addParsedParams(w, operation.HeaderParams, "h", "header")
		w.EmptyLine()
	}
}

func addParsedParams(w *gen.Writer, namedParams []spec.NamedParam, paramsConverterName string, paramsParserName string) {
	w.Line(`  %s := NewParamsConverter(%s)`, paramsConverterName, paramsParserName)
	for _, param := range namedParams {
		w.Line(`  %s.%s("%s", %s)`, paramsConverterName, parserMethodName(&param.Type.Definition), param.Name.Source, param.Name.CamelCase())
	}
}

func addResponse(w *gen.Writer, operation spec.NamedOperation) {
	for _, response := range operation.Responses {
		w.EmptyLine()
		w.Line(`if resp.StatusCode == %s {`, spec.HttpStatusCode(response.Name))
		if !response.Type.Definition.IsEmpty() {
			if spec.HttpStatusCode(response.Name) == "200" {
				w.Line(`  responseBody, err := ioutil.ReadAll(resp.Body)`)
				w.Line(`  err = resp.Body.Close()`)
				w.EmptyLine()
				if operation.Body == nil {
					w.Line(`  var body *%s`, ToPascalCase(response.Type.Definition.Name))
				}
				w.Line(`  err = json.Unmarshal(responseBody, &body)`)
				addCheckError(w)
			}
		} else {
			w.Line(`  body := &%s`, ToPascalCase(response.Type.Definition.Name))
		}
		w.EmptyLine()
		w.Line(`  return &%sResponse{%s: body}, nil`, operation.Name.PascalCase(), response.Name.PascalCase())
		w.Line(`}`)
	}
}
