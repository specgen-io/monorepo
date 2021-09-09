package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"specgen/gen"
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

	imports := []string{
		`"fmt"`,
		`"errors"`,
		`"io/ioutil"`,
		`"net/http"`,
		`"encoding/json"`,
	}
	if apiHasBody(api) {
		imports = append(imports, `"bytes"`)
	}
	imports = generateImports(api, imports)
	w.EmptyLine()
	for _, imp := range imports {
		w.Line(`import %s`, imp)
	}

	w.EmptyLine()
	generateClientTypeStruct(w, *api)
	w.EmptyLine()
	generateNewClient(w, *api)
	for _, operation := range api.Operations {
		w.EmptyLine()
		generateClientFunction(w, *api, operation)
	}

	return &gen.TextFile{
		Path:    filepath.Join(generatePath, fmt.Sprintf("%s_client.go", api.Name.SnakeCase())),
		Content: w.String(),
	}
}

func generateClientTypeStruct(w *gen.Writer, api spec.Api) {
	w.Line(`type %s struct {`, clientTypeName(&api))
	w.Line(`  baseUrl string`)
	w.Line(`}`)
}

func generateNewClient(w *gen.Writer, api spec.Api) {
	w.Line(`func New%s(baseUrl string) *%s {`, ToPascalCase(clientTypeName(&api)), clientTypeName(&api))
	w.Line(`  return &%s{baseUrl}`, clientTypeName(&api))
	w.Line(`}`)
}

func generateClientFunction(w *gen.Writer, api spec.Api, operation spec.NamedOperation) {
	w.Line(`func (client *%s) %s(%s) (*%s, error) {`, clientTypeName(&api), operation.Name.PascalCase(), JoinDelimParams(addMethodParams(operation)), responseTypeName(&operation))
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
		return fmt.Sprintf(`fmt.Sprintf("%s", %s)`, JoinParams(getUrl(operation)), JoinDelimParams(addUrlParam(operation)))
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
		w.Line(`  return &%s{%s: body}, nil`, responseTypeName(&operation), response.Name.PascalCase())
		w.Line(`}`)
	}
}
