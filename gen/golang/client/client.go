package client

import (
	"fmt"
	"github.com/pinzolo/casee"
	"github.com/specgen-io/specgen/v2/gen/golang/common"
	"github.com/specgen-io/specgen/v2/gen/golang/imports"
	"github.com/specgen-io/specgen/v2/gen/golang/models"
	"github.com/specgen-io/specgen/v2/gen/golang/module"
	"github.com/specgen-io/specgen/v2/gen/golang/responses"
	"github.com/specgen-io/specgen/v2/gen/golang/types"
	"github.com/specgen-io/specgen/v2/gen/golang/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

var ToPascalCase = casee.ToPascalCase
var ToUpperCase = casee.ToUpperCase

func GenerateClient(specification *spec.Spec, moduleName string, generatePath string) *sources.Sources {
	sources := sources.NewSources()

	rootModule := module.New(moduleName, generatePath)

	emptyModule := rootModule.Submodule("empty")
	sources.AddGenerated(types.GenerateEmpty(emptyModule))

	for _, version := range specification.Versions {
		versionModule := rootModule.Submodule(version.Version.FlatCase())
		modelsModule := versionModule.Submodule(types.ModelsPackage)

		sources.AddGeneratedAll(models.GenerateVersionModels(&version, modelsModule))
		sources.AddGeneratedAll(generateClientsImplementations(&version, versionModule, modelsModule, emptyModule))
	}
	return sources
}

func generateClientsImplementations(version *spec.Version, versionModule, modelsModule, emptyModule module.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		apiModule := versionModule.Submodule(api.Name.SnakeCase())
		files = append(files, *generateConverter(apiModule))
		files = append(files, *generateClientImplementation(&api, apiModule, modelsModule, emptyModule))
	}
	return files
}

func generateClientImplementation(api *spec.Api, versionModule, modelsModule, emptyModule module.Module) *sources.CodeFile {
	w := writer.NewGoWriter()
	w.Line("package %s", versionModule.Name)

	imports := imports.Imports().
		Add("fmt").
		Add("errors").
		Add("io/ioutil").
		Add("net/http").
		Add("encoding/json").
		AddAlias("github.com/sirupsen/logrus", "log")
	if types.ApiHasBody(api) {
		imports.Add("bytes")
	}
	if types.ApiHasType(api, spec.TypeEmpty) {
		imports.Add(emptyModule.Package)
	}
	imports.AddApiTypes(api)
	imports.Add(modelsModule.Package)
	imports.Write(w)

	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			w.EmptyLine()
			responses.GenerateOperationResponseStruct(w, &operation)
		}
	}
	w.EmptyLine()
	generateClientWithCtor(w)
	for _, operation := range api.Operations {
		w.EmptyLine()
		generateClientFunction(w, &operation)
	}

	return &sources.CodeFile{
		Path:    versionModule.GetPath("client.go"),
		Content: w.String(),
	}
}

func generateClientWithCtor(w *sources.Writer) {
	w.Line(`type %s struct {`, clientTypeName())
	w.Line(`  baseUrl string`)
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`func New%s(baseUrl string) *%s {`, ToPascalCase(clientTypeName()), clientTypeName())
	w.Line(`  return &%s{baseUrl}`, clientTypeName())
	w.Line(`}`)
}

func generateClientFunction(w *sources.Writer, operation *spec.NamedOperation) {
	w.Line(`var %s = log.Fields{"operationId": "%s.%s", "method": "%s", "url": "%s"}`, logFieldsName(operation), operation.Api.Name.Source, operation.Name.Source, ToUpperCase(operation.Endpoint.Method), operation.FullUrl())
	w.Line(`func (client *%s) %s {`, clientTypeName(), common.OperationSignature(operation, nil))
	body := "nil"
	if operation.BodyIs(spec.BodyString) {
		w.Line(`  bodyData := []byte(body)`)
		body = "bytes.NewBuffer(bodyData)"
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`  bodyData, err := json.Marshal(body)`)
		body = "bytes.NewBuffer(bodyData)"
	}
	w.Line(`  req, err := http.NewRequest("%s", client.baseUrl+%s, %s)`, operation.Endpoint.Method, addRequestUrlParams(operation), body)
	w.Line(`  if err != nil {`)
	w.Line(`    log.WithFields(%s).Error("Failed to create HTTP request", err.Error())`, logFieldsName(operation))
	w.Line(`    %s`, returnErr(operation))
	w.Line(`  }`)
	if operation.BodyIs(spec.BodyString) {
		w.Line(`  req.Header.Set("Content-Type", "text/plain")`)
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`  req.Header.Set("Content-Type", "application/json")`)
	}
	w.EmptyLine()
	parseParams(w, operation)
	w.Line(`  log.WithFields(%s).Info("Sending request")`, logFieldsName(operation))
	w.Line(`  resp, err := http.DefaultClient.Do(req)`)
	w.Line(`  if err != nil {`)
	w.Line(`    log.WithFields(%s).Error("Request failed", err.Error())`, logFieldsName(operation))
	w.Line(`    %s`, returnErr(operation))
	w.Line(`  }`)
	addClientResponses(w.Indented(), operation)
	w.EmptyLine()
	w.Line(`  msg := fmt.Sprintf("Unexpected status code received: %s", resp.StatusCode)`, "%d")
	w.Line(`  log.WithFields(%s).Error(msg)`, logFieldsName(operation))
	w.Line(`  err = errors.New(msg)`)
	w.Line(`  %s`, returnErr(operation))
	w.Line(`}`)
}

func returnErr(operation *spec.NamedOperation) string {
	if len(operation.Responses) == 1 && operation.Responses[0].Type.Definition.IsEmpty() {
		return `return err`
	}
	return `return nil, err`
}

func getUrl(operation *spec.NamedOperation) []string {
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

func addRequestUrlParams(operation *spec.NamedOperation) string {
	if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
		return fmt.Sprintf(`fmt.Sprintf("%s", %s)`, strings.Join(getUrl(operation), ""), strings.Join(addUrlParam(operation), ", "))
	} else {
		return fmt.Sprintf(`"%s"`, operation.FullUrl())
	}
}

func addUrlParam(operation *spec.NamedOperation) []string {
	urlParams := []string{}
	for _, param := range operation.Endpoint.UrlParams {
		if types.GoType(&param.Type.Definition) != "string" {
			urlParams = append(urlParams, fmt.Sprintf("convert%s(%s)", converterMethodName(&param.Type.Definition), param.Name.CamelCase()))
		} else {
			urlParams = append(urlParams, param.Name.CamelCase())
		}
	}
	return urlParams
}

func parseParams(w *sources.Writer, operation *spec.NamedOperation) {
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

func addParsedParams(w *sources.Writer, namedParams []spec.NamedParam, paramsConverterName string, paramsParserName string) {
	w.Line(`  %s := NewParamsConverter(%s)`, paramsConverterName, paramsParserName)
	for _, param := range namedParams {
		w.Line(`  %s.%s("%s", %s)`, paramsConverterName, converterMethodName(&param.Type.Definition), param.Name.Source, param.Name.CamelCase())
	}
}

func addClientResponses(w *sources.Writer, operation *spec.NamedOperation) {
	for _, response := range operation.Responses {
		w.EmptyLine()
		w.Line(`if resp.StatusCode == %s {`, spec.HttpStatusCode(response.Name))
		w.Line(`  log.WithFields(%s).WithField("status", %s).Info("Received response")`, logFieldsName(operation), spec.HttpStatusCode(response.Name))
		if response.BodyIs(spec.BodyString) {
			w.Line(`  responseBody, err := ioutil.ReadAll(resp.Body)`)
			w.Line(`  err = resp.Body.Close()`)
			w.Line(`  if err != nil {`)
			w.Line(`    log.WithFields(%s).Error("%s", err.Error())`, logFieldsName(operation), `Reading request body failed`)
			w.Line(`    return nil, err`)
			w.Line(`  }`)
			w.Line(`  result := string(responseBody)`)
		}
		if response.BodyIs(spec.BodyJson) {
			w.Line(`  responseBody, err := ioutil.ReadAll(resp.Body)`)
			w.Line(`  err = resp.Body.Close()`)
			w.Line(`  var result %s`, types.GoType(&response.Type.Definition))
			w.Line(`  err = json.Unmarshal(responseBody, &result)`)
			w.Line(`  if err != nil {`)
			w.Line(`    log.WithFields(%s).Error("%s", err.Error())`, logFieldsName(operation), `Failed to parse response JSON`)
			w.Line(`    return nil, err`)
			w.Line(`  }`)
		}

		if len(operation.Responses) == 1 {
			if response.Type.Definition.IsEmpty() {
				w.Line(`  return nil`)
			} else {
				w.Line(`  return &result, nil`)
			}
		} else {
			if response.Type.Definition.IsEmpty() {
				w.Line(`  return &%s, nil`, responses.NewResponse(&response, fmt.Sprintf(`%s{}`, types.EmptyType)))
			} else {
				w.Line(`  return &%s, nil`, responses.NewResponse(&response, `result`))
			}
		}
		w.Line(`}`)
	}
}

func clientTypeName() string {
	return `Client`
}

func logFieldsName(operation *spec.NamedOperation) string {
	return fmt.Sprintf("log%s", operation.Name.PascalCase())
}
