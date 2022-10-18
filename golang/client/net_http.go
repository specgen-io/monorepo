package client

import (
	"fmt"
	"generator"
	"github.com/pinzolo/casee"
	"golang/common"
	"golang/imports"
	"golang/module"
	"golang/types"
	"golang/writer"
	"spec"
	"strings"
)

func NewNetHttpGenerator(types *types.Types) *NetHttpGenerator {
	return &NetHttpGenerator{types}
}

type NetHttpGenerator struct {
	Types *types.Types
}

func (g *NetHttpGenerator) GenerateClientsImplementations(version *spec.Version, versionModule, convertModule, emptyModule, errorsModule, errorsModelsModule, modelsModule, respondModule module.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		apiModule := versionModule.Submodule(api.Name.SnakeCase())
		files = append(files, *g.generateClientImplementation(&api, apiModule, convertModule, emptyModule, errorsModule, errorsModelsModule, modelsModule, respondModule))
	}
	return files
}

func (g *NetHttpGenerator) generateClientImplementation(api *spec.Api, versionModule, convertModule, emptyModule, errorsModule, errorsModelsModule, modelsModule, responseModule module.Module) *generator.CodeFile {
	w := writer.New(versionModule, "client.go")

	imports := imports.New().
		Add("fmt").
		Add("errors").
		Add("net/http").
		Add("encoding/json").
		AddAliased("github.com/sirupsen/logrus", "log")
	if types.ApiHasBody(api) {
		imports.Add("bytes")
	}
	if types.ApiHasUrlParams(api) {
		imports.Module(convertModule)
	}
	if types.ApiHasType(api, spec.TypeEmpty) {
		imports.Module(emptyModule)
	}
	imports.Module(errorsModule)
	imports.Module(errorsModelsModule)
	imports.AddApiTypes(api)
	imports.Module(modelsModule)
	imports.Module(responseModule)
	imports.Write(w)

	for _, operation := range api.Operations {
		if common.ResponsesNumber(&operation) > 1 {
			w.EmptyLine()
			generateResponseStruct(w, g.Types, &operation)
		}
	}
	w.EmptyLine()
	g.generateClientWithCtor(w)
	for _, operation := range api.Operations {
		w.EmptyLine()
		g.generateClientFunction(w, &operation)
	}

	return w.ToCodeFile()
}

func (g *NetHttpGenerator) generateClientWithCtor(w generator.Writer) {
	w.Line(`type %s struct {`, clientTypeName())
	w.Line(`  baseUrl string`)
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`func New%s(baseUrl string) *%s {`, casee.ToPascalCase(clientTypeName()), clientTypeName())
	w.Line(`  return &%s{baseUrl}`, clientTypeName())
	w.Line(`}`)
}

func (g *NetHttpGenerator) generateClientFunction(w generator.Writer, operation *spec.NamedOperation) {
	w.Line(`func (client *%s) %s {`, clientTypeName(), operationSignature(operation, g.Types, nil))
	w.Line(`  var %s = log.Fields{"operationId": "%s.%s", "method": "%s", "url": "%s"}`, logFieldsName(operation), operation.InApi.Name.Source, operation.Name.Source, casee.ToUpperCase(operation.Endpoint.Method), operation.FullUrl())
	body := "nil"
	if operation.BodyIs(spec.BodyString) {
		w.Line(`  bodyData := []byte(body)`)
		body = "bytes.NewBuffer(bodyData)"
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`  bodyData, err := json.Marshal(body)`)
		generateErrHandler(w)
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
	generateErrHandler(w)
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
		if types.IsEnumModel(&param.Type.Definition) || g.Types.GoType(&param.Type.Definition) == "string" {
			urlParams = append(urlParams, param.Name.CamelCase())
		} else {
			urlParams = append(urlParams, callRawConvert(&param.Type.Definition, param.Name.CamelCase()))
		}
	}
	return urlParams
}

func (g *NetHttpGenerator) parseParams(w generator.Writer, operation *spec.NamedOperation) {
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

func (g *NetHttpGenerator) addParsedParams(w generator.Writer, namedParams []spec.NamedParam, paramsConverterName string, paramsParserName string) {
	w.Line(`  %s := convert.NewParamsConverter(%s)`, paramsConverterName, paramsParserName)
	for _, param := range namedParams {
		w.Line(`  %s.%s`, paramsConverterName, callConverter(&param.Type.Definition, param.Name.Source, param.Name.CamelCase()))
	}
}

func (g *NetHttpGenerator) addClientResponses(w generator.Writer, operation *spec.NamedOperation) {
	for _, response := range operation.Responses {
		w.EmptyLine()
		g.generateResponse(w, operation, response)

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

func (g *NetHttpGenerator) generateResponse(w generator.Writer, operation *spec.NamedOperation, response spec.OperationResponse) {
	w.Line(`if resp.StatusCode == %s {`, spec.HttpStatusCode(response.Name))
	if response.BodyIs(spec.BodyString) {
		w.Line(`  responseBody, err := response.Text(%s, resp)`, logFieldsName(operation))
		generateErrHandler(w)
		w.Line(`  result := string(responseBody)`)
	}
	if response.BodyIs(spec.BodyJson) {
		w.Line(`  var result %s`, g.Types.GoType(&response.Type.Definition))
		w.Line(`  err := response.Json(%s, resp, &result)`, logFieldsName(operation))
		generateErrHandler(w)
	}
	if response.BodyIs(spec.BodyEmpty) {
		w.Line(`  response.Empty(%s, resp)`, logFieldsName(operation))
	}
}

func newResponse(response *spec.OperationResponse, body string) string {
	return fmt.Sprintf(`%s{%s: &%s}`, responseTypeName(response.Operation), response.Name.PascalCase(), body)
}

func generateErrHandler(w generator.Writer) {
	w.Line(`  if err != nil {`)
	w.Line(`    return nil, err`)
	w.Line(`  }`)
}

func clientTypeName() string {
	return `Client`
}

func logFieldsName(operation *spec.NamedOperation) string {
	return fmt.Sprintf("log%s", operation.Name.PascalCase())
}
