package client

import (
	"fmt"
	"generator"
	"github.com/pinzolo/casee"
	"golang/common"
	"golang/types"
	"golang/walkers"
	"golang/writer"
	"spec"
	"strings"
)

func NewNetHttpGenerator(modules *Modules, types *types.Types) *NetHttpGenerator {
	return &NetHttpGenerator{modules, types}
}

type NetHttpGenerator struct {
	Modules *Modules
	Types   *types.Types
}

func (g *NetHttpGenerator) Clients(version *spec.Version) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *g.client(&api))
	}
	return files
}

func (g *NetHttpGenerator) client(api *spec.Api) *generator.CodeFile {
	w := writer.New(g.Modules.Client(api), "client.go")

	w.Imports.Add("fmt")
	w.Imports.Add("errors")
	w.Imports.Add("net/http")
	w.Imports.Add("encoding/json")
	w.Imports.AddAliased("github.com/sirupsen/logrus", "log")
	if walkers.ApiHasBodyOfKind(api, spec.BodyJson) || walkers.ApiHasBodyOfKind(api, spec.BodyString) {
		w.Imports.Add("bytes")
	}
	if walkers.ApiHasUrlParams(api) {
		w.Imports.Module(g.Modules.Convert)
	}
	if walkers.ApiHasType(api, spec.TypeEmpty) {
		w.Imports.Module(g.Modules.Empty)
	}
	w.Imports.Module(g.Modules.HttpErrors)
	w.Imports.AddApiTypes(api)
	w.Imports.Module(g.Modules.Models(api.InHttp.InVersion))
	w.Imports.Module(g.Modules.Response)

	for _, operation := range api.Operations {
		if common.ResponsesNumber(&operation) > 1 {
			w.EmptyLine()
			responseStruct(w, g.Types, &operation)
		}
	}
	w.EmptyLine()
	g.clientWithCtor(w)
	for _, operation := range api.Operations {
		w.EmptyLine()
		g.operation(w, &operation)
	}

	return w.ToCodeFile()
}

func (g *NetHttpGenerator) clientWithCtor(w *writer.Writer) {
	w.Line(`type %s struct {`, clientTypeName())
	w.Line(`  baseUrl string`)
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`func New%s(baseUrl string) *%s {`, casee.ToPascalCase(clientTypeName()), clientTypeName())
	w.Line(`  return &%s{baseUrl}`, clientTypeName())
	w.Line(`}`)
}

func (g *NetHttpGenerator) operation(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line(`func (client *%s) %s {`, clientTypeName(), operationSignature(operation, g.Types, nil))
	w.Line(`  var %s = log.Fields{"operationId": "%s.%s", "method": "%s", "url": "%s"}`, logFieldsName(operation), operation.InApi.Name.Source, operation.Name.Source, casee.ToUpperCase(operation.Endpoint.Method), operation.FullUrl())
	body := "nil"
	if operation.BodyIs(spec.BodyString) {
		w.Line(`  bodyData := []byte(body)`)
		body = "bytes.NewBuffer(bodyData)"
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`  bodyData, err := json.Marshal(body)`)
		w.Line(`  if err != nil {`)
		w.Line(`    return nil, err`)
		w.Line(`  }`)
		body = "bytes.NewBuffer(bodyData)"
	}
	w.Line(`  req, err := http.NewRequest("%s", client.baseUrl+%s, %s)`, operation.Endpoint.Method, g.addRequestUrlParams(operation), body)
	w.Line(`  if err != nil {`)
	w.Line(`    log.WithFields(%s).Error("Failed to create HTTP request", err.Error())`, logFieldsName(operation))
	w.Line(`    return nil, err`)
	w.Line(`  }`)
	if operation.BodyIs(spec.BodyString) {
		w.Line(`  req.Header.Set("Content-Type", "text/plain")`)
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`  req.Header.Set("Content-Type", "application/json")`)
	}
	w.EmptyLine()
	g.parseParams(w, operation)
	w.Line(`  log.WithFields(%s).Info("Sending request")`, logFieldsName(operation))
	w.Line(`  resp, err := http.DefaultClient.Do(req)`)
	w.Line(`  if err != nil {`)
	w.Line(`    log.WithFields(%s).Error("Request failed", err.Error())`, logFieldsName(operation))
	w.Line(`    return nil, err`)
	w.Line(`  }`)
	g.addClientResponses(w.Indented(), operation)
	w.EmptyLine()
	w.Line(`  err = httperrors.HandleErrors(resp, %s)`, logFieldsName(operation))
	w.Line(`  if err != nil {`)
	w.Line(`    return nil, err`)
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  msg := fmt.Sprintf("Unexpected status code received: %s", resp.StatusCode)`, "%d")
	w.Line(`  log.WithFields(%s).Error(msg)`, logFieldsName(operation))
	w.Line(`  err = errors.New(msg)`)
	w.Line(`  return nil, err`)
	w.Line(`}`)
}

func (g *NetHttpGenerator) getUrl(operation *spec.NamedOperation) []string {
	reminder := operation.FullUrl()
	urlParams := []string{}
	if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
		for _, param := range operation.Endpoint.UrlParams {
			parts := strings.Split(reminder, spec.UrlParamStr(&param))
			urlParams = append(urlParams, fmt.Sprintf("%s%s", parts[0], "%s"))
			reminder = parts[1]
		}
	}
	urlParams = append(urlParams, fmt.Sprintf("%s", reminder))
	return urlParams
}

func (g *NetHttpGenerator) addRequestUrlParams(operation *spec.NamedOperation) string {
	if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
		return fmt.Sprintf(`fmt.Sprintf("%s", %s)`, strings.Join(g.getUrl(operation), ""), strings.Join(g.addUrlParam(operation), ", "))
	} else {
		return fmt.Sprintf(`"%s"`, operation.FullUrl())
	}
}

func (g *NetHttpGenerator) addUrlParam(operation *spec.NamedOperation) []string {
	urlParams := []string{}
	for _, param := range operation.Endpoint.UrlParams {
		if param.Type.Definition.IsEnum() || g.Types.GoType(&param.Type.Definition) == "string" {
			urlParams = append(urlParams, param.Name.CamelCase())
		} else {
			urlParams = append(urlParams, callRawConvert(&param.Type.Definition, param.Name.CamelCase()))
		}
	}
	return urlParams
}

func (g *NetHttpGenerator) parseParams(w *writer.Writer, operation *spec.NamedOperation) {
	if operation.QueryParams != nil && len(operation.QueryParams) > 0 {
		w.Line(`  query := req.URL.Query()`)
		g.addParsedParams(w, operation.QueryParams, "q", "query")
		w.Line(`  req.URL.RawQuery = %s.Encode()`, "query")
		w.EmptyLine()
	}
	if operation.HeaderParams != nil && len(operation.HeaderParams) > 0 {
		w.Line(`  header := req.Header`)
		g.addParsedParams(w, operation.HeaderParams, "h", "header")
		w.EmptyLine()
	}
}

func (g *NetHttpGenerator) addParsedParams(w *writer.Writer, namedParams []spec.NamedParam, paramsConverterName string, paramsParserName string) {
	w.Line(`  %s := convert.NewParamsConverter(%s)`, paramsConverterName, paramsParserName)
	for _, param := range namedParams {
		w.Line(`  %s.%s`, paramsConverterName, callConverter(&param.Type.Definition, param.Name.Source, param.Name.CamelCase()))
	}
}

func (g *NetHttpGenerator) addClientResponses(w *writer.Writer, operation *spec.NamedOperation) {
	for _, response := range operation.Responses {
		w.EmptyLine()
		g.response(w, operation, response)

		if common.ResponsesNumber(operation) == 1 {
			if response.Type.Definition.IsEmpty() {
				if common.IsSuccessfulStatusCode(spec.HttpStatusCode(response.Name)) {
					w.Line(`  return &%s{}, nil`, types.EmptyType)
				} else {
					w.Line(`  return nil, err`)
				}
			} else {
				w.Line(`  return &result, nil`)
			}
		} else {
			if response.Type.Definition.IsEmpty() {
				w.Line(`  return &%s, nil`, newResponse(&response, fmt.Sprintf(`%s{}`, types.EmptyType)))
			} else {
				w.Line(`  return &%s, nil`, newResponse(&response, `result`))
			}
		}
		w.Line(`}`)
	}
}

func (g *NetHttpGenerator) response(w *writer.Writer, operation *spec.NamedOperation, response spec.OperationResponse) {
	w.Line(`if resp.StatusCode == %s {`, spec.HttpStatusCode(response.Name))
	if response.BodyIs(spec.BodyString) {
		w.Line(`  responseBody, err := response.Text(%s, resp)`, logFieldsName(operation))
		w.Line(`  if err != nil {`)
		w.Line(`    return nil, err`)
		w.Line(`  }`)
		w.Line(`  result := string(responseBody)`)
	}
	if response.BodyIs(spec.BodyJson) {
		w.Line(`  var result %s`, g.Types.GoType(&response.Type.Definition))
		w.Line(`  err := response.Json(%s, resp, &result)`, logFieldsName(operation))
		w.Line(`  if err != nil {`)
		w.Line(`    return nil, err`)
		w.Line(`  }`)
	}
	if response.BodyIs(spec.BodyEmpty) {
		w.Line(`  response.Empty(%s, resp)`, logFieldsName(operation))
	}
}

func newResponse(response *spec.OperationResponse, body string) string {
	return fmt.Sprintf(`%s{%s: &%s}`, responseTypeName(response.Operation), response.Name.PascalCase(), body)
}

func clientTypeName() string {
	return `Client`
}

func logFieldsName(operation *spec.NamedOperation) string {
	return fmt.Sprintf("log%s", operation.Name.PascalCase())
}
