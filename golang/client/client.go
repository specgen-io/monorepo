package client

import (
	"fmt"
	"generator"
	"github.com/pinzolo/casee"
	"golang/common"
	"golang/imports"
	"golang/models"
	"golang/module"
	"golang/responses"
	"golang/types"
	"golang/writer"
	"spec"
	"strings"
)

var ToPascalCase = casee.ToPascalCase
var ToUpperCase = casee.ToUpperCase

func GenerateClient(specification *spec.Spec, moduleName string, generatePath string) *generator.Sources {
	sources := generator.NewSources()

	modules := models.NewModules(moduleName, generatePath, specification)
	modelsGenerator := models.NewGenerator(modules)

	rootModule := module.New(moduleName, generatePath)

	sources.AddGenerated(modelsGenerator.GenerateEnumsHelperFunctions())

	emptyModule := rootModule.Submodule("empty")
	sources.AddGenerated(types.GenerateEmpty(emptyModule))

	convertModule := rootModule.Submodule("convert")
	sources.AddGenerated(generateConverter(convertModule))

	responseModule := rootModule.Submodule("response")
	sources.AddGenerated(generateResponseFunctions(responseModule))

	errorsModule := rootModule.Submodule("httperrors")
	errorsModelsModule := errorsModule.SubmoduleAliased("models", types.ErrorsModelsPackage)
	sources.AddGenerated(modelsGenerator.GenerateErrorModels(specification.HttpErrors))
	sources.AddGenerated(httpErrors(errorsModule, errorsModelsModule, &specification.HttpErrors.Responses))

	for _, version := range specification.Versions {
		versionModule := rootModule.Submodule(version.Name.FlatCase())
		modelsModule := versionModule.Submodule(types.VersionModelsPackage)
		sources.AddGenerated(modelsGenerator.GenerateVersionModels(&version))
		sources.AddGeneratedAll(generateClientsImplementations(&version, versionModule, convertModule, emptyModule, errorsModule, errorsModelsModule, modelsModule, responseModule))
	}
	return sources
}

func generateClientsImplementations(version *spec.Version, versionModule, convertModule, emptyModule, errorsModule, errorsModelsModule, modelsModule, respondModule module.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		apiModule := versionModule.Submodule(api.Name.SnakeCase())
		files = append(files, *generateClientImplementation(&api, apiModule, convertModule, emptyModule, errorsModule, errorsModelsModule, modelsModule, respondModule))
	}
	return files
}

func generateClientImplementation(api *spec.Api, versionModule, convertModule, emptyModule, errorsModule, errorsModelsModule, modelsModule, responseModule module.Module) *generator.CodeFile {
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
	imports.ModuleAliased(errorsModelsModule)
	imports.AddApiTypes(api)
	imports.Module(modelsModule)
	imports.Module(responseModule)
	imports.Write(w)

	for _, operation := range api.Operations {
		if common.ResponsesNumber(&operation) > 1 {
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

	return w.ToCodeFile()
}

func generateClientWithCtor(w generator.Writer) {
	w.Line(`type %s struct {`, clientTypeName())
	w.Line(`  baseUrl string`)
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`func New%s(baseUrl string) *%s {`, ToPascalCase(clientTypeName()), clientTypeName())
	w.Line(`  return &%s{baseUrl}`, clientTypeName())
	w.Line(`}`)
}

func generateClientFunction(w generator.Writer, operation *spec.NamedOperation) {
	w.Line(`func (client *%s) %s {`, clientTypeName(), OperationSignature(operation, nil))
	w.Line(`  var %s = log.Fields{"operationId": "%s.%s", "method": "%s", "url": "%s"}`, logFieldsName(operation), operation.InApi.Name.Source, operation.Name.Source, ToUpperCase(operation.Endpoint.Method), operation.FullUrl())
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
	w.Line(`  req, err := http.NewRequest("%s", client.baseUrl+%s, %s)`, operation.Endpoint.Method, addRequestUrlParams(operation), body)
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
	parseParams(w, operation)
	w.Line(`  log.WithFields(%s).Info("Sending request")`, logFieldsName(operation))
	w.Line(`  resp, err := http.DefaultClient.Do(req)`)
	w.Line(`  if err != nil {`)
	w.Line(`    log.WithFields(%s).Error("Request failed", err.Error())`, logFieldsName(operation))
	w.Line(`    return nil, err`)
	w.Line(`  }`)
	addClientResponses(w.Indented(), operation)
	w.EmptyLine()
	w.Line(`  msg := fmt.Sprintf("Unexpected status code received: %s", resp.StatusCode)`, "%d")
	w.Line(`  log.WithFields(%s).Error(msg)`, logFieldsName(operation))
	w.Line(`  err = errors.New(msg)`)
	w.Line(`  return nil, err`)
	w.Line(`}`)
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
		if paramIsEnumModel(&param.Type.Definition) || types.GoType(&param.Type.Definition) == "string" {
			urlParams = append(urlParams, param.Name.CamelCase())
		} else {
			urlParams = append(urlParams, callRawConvert(&param.Type.Definition, param.Name.CamelCase()))
		}
	}
	return urlParams
}

func paramIsEnumModel(typ *spec.TypeDef) bool {
	if types.IsModel(typ) {
		if typ.Info.Model.IsEnum() {
			return true
		}
	}
	return false
}

func parseParams(w generator.Writer, operation *spec.NamedOperation) {
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

func addParsedParams(w generator.Writer, namedParams []spec.NamedParam, paramsConverterName string, paramsParserName string) {
	w.Line(`  %s := convert.NewParamsConverter(%s)`, paramsConverterName, paramsParserName)
	for _, param := range namedParams {
		w.Line(`  %s.%s`, paramsConverterName, callConverter(&param.Type.Definition, param.Name.Source, param.Name.CamelCase()))
	}
}

func addClientResponses(w generator.Writer, operation *spec.NamedOperation) {
	for _, response := range operation.Responses {
		w.EmptyLine()
		generateResponse(w, operation, response)

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
				w.Line(`  return &%s, nil`, responses.NewResponse(&response, fmt.Sprintf(`%s{}`, types.EmptyType)))
			} else {
				w.Line(`  return &%s, nil`, responses.NewResponse(&response, `result`))
			}
		}
		w.Line(`}`)
	}

	for _, response := range operation.InApi.InHttp.InVersion.InSpec.HttpErrors.Responses {
		w.EmptyLine()
		generateErrorResponse(w, operation, response)
	}
}

func generateResponse(w generator.Writer, operation *spec.NamedOperation, response spec.OperationResponse) {
	w.Line(`if resp.StatusCode == %s {`, spec.HttpStatusCode(response.Name))
	if response.BodyIs(spec.BodyString) {
		w.Line(`  responseBody, err := response.Text(%s, resp)`, logFieldsName(operation))
		generateErrHandler(w)
		w.Line(`  result := string(responseBody)`)
	}
	if response.BodyIs(spec.BodyJson) {
		w.Line(`  var result %s`, types.GoType(&response.Type.Definition))
		w.Line(`  err := response.Json(%s, resp, &result)`, logFieldsName(operation))
		generateErrHandler(w)
	}
	if response.BodyIs(spec.BodyEmpty) {
		w.Line(`  response.Empty(%s, resp)`, logFieldsName(operation))
	}
}

func generateErrorResponse(w generator.Writer, operation *spec.NamedOperation, response spec.Response) {
	w.Line(`if resp.StatusCode == %s {`, spec.HttpStatusCode(response.Name))
	w.Line(`  var result %s`, types.GoType(&response.Type.Definition))
	w.Line(`  err := response.Json(%s, resp, &result)`, logFieldsName(operation))
	generateErrHandler(w)
	w.Line(`  return nil, %s`, responses.NewErrorResponse(response, `result`))
	w.Line(`}`)
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
