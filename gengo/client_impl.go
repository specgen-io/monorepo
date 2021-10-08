package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"strings"
)

func generateClientsImplementations(version *spec.Version, versionModule module, modelsModule module) []gen.TextFile {
	files := []gen.TextFile{}
	for _, api := range version.Http.Apis {
		apiModule := versionModule.Submodule(api.Name.SnakeCase())
		files = append(files, *generateConverter(apiModule))
		files = append(files, *generateClientImplementation(&api, apiModule, modelsModule))
	}
	return files
}

func generateClientImplementation(api *spec.Api, versionModule module, modelsModule module) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", versionModule.Name)

	imports := Imports().
		Add("fmt").
		Add("errors").
		Add("io/ioutil").
		Add("net/http").
		Add("encoding/json")
	if apiHasBody(api) {
		imports.Add("bytes")
	}
	imports.AddApiTypes(api)
	imports.Add(modelsModule.Package)
	imports.Write(w)

	w.EmptyLine()
	w.Line(`type EmptyDef struct{}`)
	w.EmptyLine()
	w.Line(`var Empty = EmptyDef{}`)

	for _, operation := range api.Operations {
		w.EmptyLine()
		generateOperationResponseStruct(w, operation)
	}

	w.EmptyLine()
	generateClientWithCtor(w)
	for _, operation := range api.Operations {
		w.EmptyLine()
		generateClientFunction(w, operation)
	}

	return &gen.TextFile{
		Path:    versionModule.GetPath("client.go"),
		Content: w.String(),
	}
}

func generateClientWithCtor(w *gen.Writer) {
	w.Line(`type %s struct {`, clientTypeName())
	w.Line(`  baseUrl string`)
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`func New%s(baseUrl string) *%s {`, ToPascalCase(clientTypeName()), clientTypeName())
	w.Line(`  return &%s{baseUrl}`, clientTypeName())
	w.Line(`}`)
}

func generateClientFunction(w *gen.Writer, operation spec.NamedOperation) {
	w.Line(`func (client *%s) %s(%s) (*%s, error) {`, clientTypeName(), operation.Name.PascalCase(), JoinDelimParams(addMethodParams(operation)), responseTypeName(&operation))
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
					w.Line(`  var body *%s.%s`, modelsPackage, ToPascalCase(response.Type.Definition.Name))
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
