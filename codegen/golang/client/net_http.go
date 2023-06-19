package client

import (
	"fmt"
	"generator"
	"github.com/pinzolo/casee"
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
	if walkers.ApiHasBodyOfKind(api, spec.RequestBodyJson) || walkers.ApiHasBodyOfKind(api, spec.RequestBodyString) {
		w.Imports.Add("bytes")
	}
	if walkers.ApiHasUrlParams(api) {
		w.Imports.Module(g.Modules.Convert)
	}
	if walkers.ApiHasMultiSuccessResponsesWithEmptyBody(api) {
		w.Imports.Module(g.Modules.Empty)
	}
	if walkers.ApiHasType(api, spec.TypeDate) {
		w.Imports.Add("cloud.google.com/go/civil")
	}
	if walkers.ApiHasType(api, spec.TypeJson) {
		w.Imports.Add("encoding/json")
	}
	if walkers.ApiHasType(api, spec.TypeUuid) {
		w.Imports.Add("github.com/google/uuid")
	}
	if walkers.ApiHasType(api, spec.TypeDecimal) {
		w.Imports.Add("github.com/shopspring/decimal")
	}
	w.Imports.Module(g.Modules.HttpErrors)
	if walkers.ApiIsUsingErrorModels(api) {
		w.Imports.Module(g.Modules.HttpErrorsModels)
	}
	w.Imports.Module(g.Modules.Models(api.InHttp.InVersion))
	w.Imports.Module(g.Modules.Response)

	for _, operation := range api.Operations {
		responseStruct(w, g.Types, &operation)
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
	w.Line(`func (client *%s) %s {`, clientTypeName(), operationSignature(g.Types, operation))
	w.Line(`  var %s = log.Fields{"operationId": "%s.%s", "method": "%s", "url": "%s"}`, logFieldsName(operation), operation.InApi.Name.Source, operation.Name.Source, casee.ToUpperCase(operation.Endpoint.Method), operation.FullUrl())
	g.createRequest(w, operation, `req`)
	g.addQueryParams(w, operation, `req`)
	g.addHeaderParams(w, operation, `req`)
	g.sendRequest(w, operation, `req`, `resp`)
	g.processResponses(w.Indented(), operation)
	w.Line(`}`)
}

func (g *NetHttpGenerator) createRequest(w *writer.Writer, operation *spec.NamedOperation, requestVar string) {
	if operation.BodyIs(spec.RequestBodyString) {
		w.Line(`  bodyData := []byte(body)`)
	}
	if operation.BodyIs(spec.RequestBodyJson) {
		w.Line(`  bodyData, err := json.Marshal(body)`)
		w.Line(`  if err != nil {`)
		w.Line(`    return %s`, operationError(operation, `err`))
		w.Line(`  }`)
	}
	body := "nil"
	if !operation.BodyIs(spec.RequestBodyEmpty) {
		body = "bytes.NewBuffer(bodyData)"
	}
	w.Line(`  %s, err := http.NewRequest("%s", client.baseUrl+%s, %s)`, requestVar, operation.Endpoint.Method, g.addRequestUrlParams(operation), body)
	w.Line(`  if err != nil {`)
	w.Line(`    log.WithFields(%s).Error("Failed to create HTTP request", err.Error())`, logFieldsName(operation))
	w.Line(`    return %s`, operationError(operation, `err`))
	w.Line(`  }`)
	if operation.BodyIs(spec.RequestBodyString) {
		w.Line(`  %s.Header.Set("Content-Type", "text/plain")`, requestVar)
	}
	if operation.BodyIs(spec.RequestBodyJson) {
		w.Line(`  %s.Header.Set("Content-Type", "application/json")`, requestVar)
	}
	w.EmptyLine()
}

func (g *NetHttpGenerator) sendRequest(w *writer.Writer, operation *spec.NamedOperation, requestVar, responseVar string) {
	w.Line(`  log.WithFields(%s).Info("Sending request")`, logFieldsName(operation))
	w.Line(`  %s, err := http.DefaultClient.Do(%s)`, responseVar, requestVar)
	w.Line(`  if err != nil {`)
	w.Line(`    log.WithFields(%s).Error("Request failed", err.Error())`, logFieldsName(operation))
	w.Line(`    return %s`, operationError(operation, `err`))
	w.Line(`  }`)
	w.Line(`  log.WithFields(%s).WithField("status", %s.StatusCode).Info("Received response")`, logFieldsName(operation), responseVar)
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

func (g *NetHttpGenerator) addQueryParams(w *writer.Writer, operation *spec.NamedOperation, requestVar string) {
	if operation.QueryParams != nil && len(operation.QueryParams) > 0 {
		w.Line(`  query := %s.URL.Query()`, requestVar)
		w.Line(`  q := convert.NewParamsConverter(query)`)
		for _, param := range operation.QueryParams {
			w.Line(`  q.%s`, callConverter(&param.Type.Definition, param.Name.Source, param.Name.CamelCase()))
		}
		w.Line(`  %s.URL.RawQuery = query.Encode()`, requestVar)
		w.EmptyLine()
	}
}

func (g *NetHttpGenerator) addHeaderParams(w *writer.Writer, operation *spec.NamedOperation, requestVar string) {
	if operation.HeaderParams != nil && len(operation.HeaderParams) > 0 {
		w.Line(`  h := convert.NewParamsConverter(%s.Header)`, requestVar)
		for _, param := range operation.HeaderParams {
			w.Line(`  h.%s`, callConverter(&param.Type.Definition, param.Name.Source, param.Name.CamelCase()))
		}
		w.EmptyLine()
	}
}

func (g *NetHttpGenerator) processResponses(w *writer.Writer, operation *spec.NamedOperation) {
	w.EmptyLine()
	w.Line(`switch resp.StatusCode {`)
	for _, response := range operation.Responses {
		w.Line(`case %s:`, spec.HttpStatusCode(response.Name))
		if response.BodyIs(spec.ResponseBodyString) {
			w.Line(`  result, err := response.Text(resp)`)
			w.Line(`  if err != nil {`)
			w.Line(`    return %s`, operationError(response.Operation, `err`))
			w.Line(`  }`)
		}
		if response.BodyIs(spec.ResponseBodyJson) {
			w.Line(`  var result %s`, g.Types.GoType(&response.ResponseBody.Type.Definition))
			w.Line(`  err := response.Json(resp, &result)`)
			w.Line(`  if err != nil {`)
			w.Line(`    return %s`, operationError(response.Operation, `err`))
			w.Line(`  }`)
		}

		if response.IsSuccess() {
			w.Line(`  return %s`, resultSuccess(&response, `result`))
		} else {
			w.Line(`  return %s`, resultError(&response, g.Modules.HttpErrors, `result`))
		}
	}
	w.Line(`default:`)
	w.Line(`  err = httperrors.HandleErrors(resp)`)
	w.Line(`  if err != nil {`)
	w.Line(`    return %s`, operationError(operation, `err`))
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  msg := fmt.Sprintf("Unexpected status code received: %s", resp.StatusCode)`, "%d")
	w.Line(`  log.WithFields(%s).Error(msg)`, logFieldsName(operation))
	w.Line(`  return %s`, operationError(operation, `errors.New(msg)`))
	w.Line(`}`)
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

func responseStruct(w *writer.Writer, types *types.Types, operation *spec.NamedOperation) {
	if len(operation.Responses.Success()) > 1 {
		w.EmptyLine()
		w.Line(`type %s struct {`, responseTypeName(operation))
		w.Indent()
		for _, response := range operation.Responses.Success() {
			w.LineAligned(`%s %s`, response.Name.PascalCase(), types.GoType(spec.Nullable(&response.ResponseBody.Type.Definition)))
		}
		w.Unindent()
		w.Line(`}`)
	}
}

func responseTypeName(operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%sResponse`, operation.Name.PascalCase())
}

func (g *NetHttpGenerator) ResponseHelperFunctions() *generator.CodeFile {
	w := writer.New(g.Modules.Response, `response.go`)
	w.Lines(`
import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Json(resp *http.Response, result any) error {
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return err
	}
	return nil
}

func Text(resp *http.Response) (string, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = resp.Body.Close()
	if err != nil {
		return "", err
	}
	return string(body), nil
}
`)
	return w.ToCodeFile()
}

func (g *NetHttpGenerator) ErrorsHandler(errors spec.ErrorResponses) *generator.CodeFile {
	w := writer.New(g.Modules.HttpErrors, `errors_handler.go`)

	w.Imports.Add("net/http")
	w.Imports.Module(g.Modules.HttpErrorsModels)
	w.Imports.Module(g.Modules.Response)

	w.EmptyLine()
	w.Line(`func HandleErrors(resp *http.Response) error {`)
	w.Indent()
	w.Line(`switch resp.StatusCode {`)
	for _, response := range errors.Required() {
		w.Line(`case %s:`, spec.HttpStatusCode(response.Name))
		if response.BodyIs(spec.ResponseBodyString) {
			w.Line(`  result, err := response.Text(resp)`)
			w.Line(`  if err != nil {`)
			w.Line(`    return err`)
			w.Line(`  }`)
		}
		if response.BodyIs(spec.ResponseBodyJson) {
			w.Line(`  var result %s`, g.Types.GoType(&response.ResponseBody.Type.Definition))
			w.Line(`  err := response.Json(resp, &result)`)
			w.Line(`  if err != nil {`)
			w.Line(`    return err`)
			w.Line(`  }`)
		}

		if response.ResponseBody.Type.Definition.IsEmpty() {
			w.Line(`  return &%s{}`, response.Name.PascalCase())
		} else {
			w.Line(`  return &%s{result}`, response.Name.PascalCase())
		}
	}
	w.Line(`default:`)
	w.Line(`  return nil`)
	w.Line(`}`)
	w.Unindent()
	w.Line(`}`)

	return w.ToCodeFile()
}
