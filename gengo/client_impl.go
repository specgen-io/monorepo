package gengo

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
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
		Add("encoding/json").
		AddAlias("github.com/sirupsen/logrus", "log")
	if apiHasBody(api) {
		imports.Add("bytes")
	}
	imports.AddApiTypes(api)
	imports.Add(modelsModule.Package)
	imports.Write(w)

	w.EmptyLine()
	w.Line(`type EmptyDef struct{}`)
	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			w.EmptyLine()
			generateOperationResponseStruct(w, operation)
		}
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
	w.Line(`var %s = log.Fields{"operationId": "%s.%s", "method": "%s", "url": "%s"}`, logFieldsName(&operation), operation.Api.Name.Source, operation.Name.Source, ToUpperCase(operation.Endpoint.Method), operation.FullUrl())
	w.Line(`func (client *%s) %s(%s) %s {`, clientTypeName(), operation.Name.PascalCase(), JoinDelimParams(addMethodParams(operation)), operationReturn(&operation))
	var body = "nil"
	if operation.Body != nil {
		w.Line(`  bodyJSON, err := json.Marshal(body)`)
		body = "bytes.NewBuffer(bodyJSON)"
	}
	w.Line(`  req, err := http.NewRequest("%s", client.baseUrl+%s, %s)`, operation.Endpoint.Method, addRequestUrlParams(operation), body)
	w.Line(`  if err != nil {`)
	w.Line(`    log.WithFields(%s).Error("Failed to create HTTP request", err.Error())`, logFieldsName(&operation))
	w.Line(`    %s`, returnErr(&operation))
	w.Line(`  }`)
	w.EmptyLine()
	parseParams(w, operation)
	w.Line(`  log.WithFields(%s).Info("Sending request")`, logFieldsName(&operation))
	w.Line(`  resp, err := http.DefaultClient.Do(req)`)
	w.Line(`  if err != nil {`)
	w.Line(`    log.WithFields(%s).Error("Request failed", err.Error())`, logFieldsName(&operation))
	w.Line(`    %s`, returnErr(&operation))
	w.Line(`  }`)
	w.Indent()
	addClientResponses(w, operation)
	w.Unindent()
	w.EmptyLine()
	w.Line(`  msg := fmt.Sprintf("Unexpected status code received: %s", resp.StatusCode)`, "%d")
	w.Line(`  log.WithFields(%s).Error(msg)`, logFieldsName(&operation))
	w.Line(`  err = errors.New(msg)`)
	w.Line(`  %s`, returnErr(&operation))
	w.Line(`}`)
}

func returnErr(operation *spec.NamedOperation) string {
	if len(operation.Responses) == 1 && operation.Responses[0].Type.Definition.IsEmpty() {
		return `return err`
	}
	return `return nil, err`
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

func addClientResponses(w *gen.Writer, operation spec.NamedOperation) {
	for _, response := range operation.Responses {
		w.EmptyLine()
		w.Line(`if resp.StatusCode == %s {`, spec.HttpStatusCode(response.Name))
		w.Line(`  log.WithFields(%s).WithField("status", %s).Info("Received response")`, logFieldsName(&operation), spec.HttpStatusCode(response.Name))
		if !response.Type.Definition.IsEmpty() {
			w.Line(`  responseBody, err := ioutil.ReadAll(resp.Body)`)
			w.Line(`  err = resp.Body.Close()`)
			w.EmptyLine()
			if operation.Body == nil {
				w.Line(`  var body *%s.%s`, modelsPackage, ToPascalCase(response.Type.Definition.Name))
			}
			w.Line(`  err = json.Unmarshal(responseBody, &body)`)
			w.Line(`  if err != nil {`)
			w.Line(`    log.WithFields(%s).Error("Failed to parse response JSON", err.Error())`, logFieldsName(&operation))
			w.Line(`    return nil, err`)
			w.Line(`  }`)
		}
		w.Line(`  %s`, returnStatement(&response))
		w.Line(`}`)
	}
}

func returnStatement(response *spec.NamedResponse) string {
	operation := response.Operation
	if len(operation.Responses) == 1 {
		if response.Type.Definition.IsEmpty() {
			return `return nil`
		}
		return fmt.Sprintf(`return body, nil`)
	} else {
		if response.Type.Definition.IsEmpty() {
			return fmt.Sprintf(`return &%s{%s: &%s{}}, nil`, responseTypeName(operation), response.Name.PascalCase(), GoType(&response.Type.Definition))
		}
		return fmt.Sprintf(`return &%s{%s: body}, nil`, responseTypeName(operation), response.Name.PascalCase())
	}
}
