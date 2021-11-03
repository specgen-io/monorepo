package gengo

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func generateRouting(version *spec.Version, versionModule module, modelsModule module) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", versionModule.Name)

	imports := Imports()
	imports.Add("encoding/json")
	imports.Add("github.com/husobee/vestigo")
	imports.AddAlias("github.com/sirupsen/logrus", "log")
	imports.Add("net/http")

	for _, api := range version.Http.Apis {
		apiModule := versionModule.Submodule(api.Name.SnakeCase())
		imports.Add(apiModule.Package)
	}
	imports.Add(modelsModule.Package)
	imports.Write(w)

	for _, api := range version.Http.Apis {
		generateApiRouter(w, &api)
	}

	return &gen.TextFile{
		Path:    versionModule.GetPath("routing.go"),
		Content: w.String(),
	}
}

func getVestigoUrl(operation *spec.NamedOperation) string {
	url := operation.FullUrl()
	if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
		for _, param := range operation.Endpoint.UrlParams {
			url = strings.Replace(url, spec.UrlParamStr(param.Name.Source), fmt.Sprintf(":%s", param.Name.Source), -1)
		}
	}
	return url
}

func addSetCors(w *gen.Writer, operation *spec.NamedOperation) {
	w.Line(`router.SetCors("%s", &vestigo.CorsAccessControl{`, getVestigoUrl(operation))
	params := []string{}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf(`"%s"`, param.Name.Source))
	}
	w.Line(`  AllowHeaders: []string{%s},`, JoinDelimParams(params))
	w.Line(`})`)
}

func logFieldsName(operation *spec.NamedOperation) string {
	return fmt.Sprintf("log%s", operation.Name.PascalCase())
}

func generateApiRouter(w *gen.Writer, api *spec.Api) {
	w.EmptyLine()
	w.Line(`func %s(router *vestigo.Router, %s %s) {`, addRoutes(api), serviceInterfaceTypeVar(api), fmt.Sprintf("%s.Service", api.Name.SnakeCase()))
	for i, operation := range api.Operations {
		if i != 0 {
			w.EmptyLine()
		}
		w.Indent()
		url := getVestigoUrl(&operation)
		w.Line(`%s := log.Fields{"operationId": "%s.%s", "method": "%s", "url": "%s"}`, logFieldsName(&operation), operation.Api.Name.Source, operation.Name.Source, ToUpperCase(operation.Endpoint.Method), url)
		w.Line(`router.%s("%s", func(res http.ResponseWriter, req *http.Request) {`, ToPascalCase(operation.Endpoint.Method), url)
		generateOperationMethod(w.Indented(), &operation)
		w.Line(`})`)
		if operation.HeaderParams != nil && len(operation.HeaderParams) > 0 {
			addSetCors(w, &operation)
		}
		w.Unindent()
	}
	w.Line(`}`)
}

func parserParameterCall(operation *spec.NamedOperation, param *spec.NamedParam, paramsParserName string) string {
	paramNameSource := param.Name.Source
	if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
		paramNameSource = ":" + paramNameSource
	}
	parserParams := []string{fmt.Sprintf(`"%s"`, paramNameSource)}
	methodName, defaultParam := parserDefaultName(param)
	isEnum := param.Type.Definition.Info.Model != nil && param.Type.Definition.Info.Model.IsEnum()
	enumModel := param.Type.Definition.Info.Model
	if isEnum {
		parserParams = append(parserParams, fmt.Sprintf("%s.%s", modelsPackage, enumValuesStrings(enumModel)))
	}
	if defaultParam != nil {
		parserParams = append(parserParams, *defaultParam)
	}
	call := fmt.Sprintf(`%s.%s(%s)`, paramsParserName, methodName, JoinDelimParams(parserParams))
	if isEnum {
		call = fmt.Sprintf(`%s.%s(%s)`, modelsPackage, enumModel.Name.PascalCase(), call)
	}
	return call
}

func generateOperationParametersParsing(w *gen.Writer, operation *spec.NamedOperation, namedParams []spec.NamedParam, paramsParserName string, paramName string) {
	if namedParams != nil && len(namedParams) > 0 {
		w.Line(`%s := NewParamsParser(%s)`, paramsParserName, paramName)
		for _, param := range namedParams {
			w.Line(`%s := %s`, param.Name.CamelCase(), parserParameterCall(operation, &param, paramsParserName))
		}
		w.Line(`if len(%s.Errors) > 0 {`, paramsParserName)
		w.Line(`  log.WithFields(%s).Warnf("Can't parse %s: %%s", %s.Errors)`, logFieldsName(operation), paramsParserName, paramsParserName)
		w.Line(`  res.WriteHeader(400)`)
		w.Line(`  log.WithFields(%s).WithField("status", 400).Info("Completed request")`, logFieldsName(operation))
		w.Line(`  return`)
		w.Line(`}`)
	}
}

func generateServiceCall(w *gen.Writer, operation *spec.NamedOperation) {
	singleEmptyResponse := len(operation.Responses) == 1  && operation.Responses[0].Type.Definition.IsEmpty()
	serviceCallResult := "response, err"
	if singleEmptyResponse {
		serviceCallResult = "err"
	}
	w.Line(`%s := %s.%s(%s)`, serviceCallResult, serviceInterfaceTypeVar(operation.Api), operation.Name.PascalCase(), JoinDelimParams(addOperationMethodParams(operation)))

	w.Line(`if err != nil {`)
	w.Line(`  log.WithFields(%s).Errorf("Error returned from service implementation: %%s", err.Error())`, logFieldsName(operation))
	w.Line(`  res.WriteHeader(500)`)
	w.Line(`  log.WithFields(%s).WithField("status", 500).Info("Completed request")`, logFieldsName(operation))
	w.Line(`  return`)
	w.Line(`}`)

	if !singleEmptyResponse {
		w.Line(`if response == nil {`)
		w.Line(`  log.WithFields(%s).Errorf("No result returned from service implementation")`, logFieldsName(operation))
		w.Line(`  res.WriteHeader(500)`)
		w.Line(`  log.WithFields(%s).WithField("status", 500).Info("Completed request")`, logFieldsName(operation))
		w.Line(`  return`)
		w.Line(`}`)
	}
}

func generateResponseWriting(w *gen.Writer, response *spec.NamedResponse, responseVar string) {
	w.Line(`res.WriteHeader(%s)`, spec.HttpStatusCode(response.Name))
	if !response.Type.Definition.IsEmpty() {
		w.Line(`json.NewEncoder(res).Encode(%s)`, responseVar)
	}
	w.Line(`log.WithFields(%s).WithField("status", %s).Info("Completed request")`, logFieldsName(response.Operation), spec.HttpStatusCode(response.Name))
	w.Line(`return`)
}

func generateOperationMethod(w *gen.Writer, operation *spec.NamedOperation) {
	w.Line(`log.WithFields(%s).Info("Received request")`, logFieldsName(operation))
	if operation.Body != nil {
		w.Line(`var body %s`, GoType(&operation.Body.Type.Definition))
		w.Line(`err := json.NewDecoder(req.Body).Decode(&body)`)
		w.Line(`if err != nil {`)
		w.Line(`  log.WithFields(%s).Warnf("Decoding body JSON failed: %%s", err.Error())`, logFieldsName(operation))
		w.Line(`  res.WriteHeader(400)`)
		w.Line(`  log.WithFields(%s).WithField("status", 400).Info("Completed request")`, logFieldsName(operation))
		w.Line(`  return`)
		w.Line(`}`)
	}
	generateOperationParametersParsing(w, operation, operation.QueryParams, "queryParams", "req.URL.Query()")
	generateOperationParametersParsing(w, operation, operation.HeaderParams, "headerParams", "req.Header")
	generateOperationParametersParsing(w, operation, operation.Endpoint.UrlParams, "urlParams", "req.URL.Query()")

	generateServiceCall(w, operation)

	if len(operation.Responses) == 1 {
		generateResponseWriting(w, &operation.Responses[0], "response")
	} else {
		for _, response := range operation.Responses {
			responseVar := fmt.Sprintf("response.%s", response.Name.PascalCase())
			w.Line(`if %s != nil {`, responseVar)
			generateResponseWriting(w.Indented(), &response, responseVar)
			w.Line(`}`)
		}
		w.Line(`log.WithFields(%s).Error("Result from service implementation does not have anything in it")`, logFieldsName(operation))
		w.Line(`res.WriteHeader(500)`)
		w.Line(`log.WithFields(%s).WithField("status", 500).Info("Completed request")`, logFieldsName(operation))
		w.Line(`return`)
	}
}

func addOperationMethodParams(operation *spec.NamedOperation) []string {
	urlParams := []string{}
	if operation.Body != nil {
		urlParams = append(urlParams, fmt.Sprintf("%s", "&body"))
	}
	for _, param := range operation.QueryParams {
		urlParams = append(urlParams, fmt.Sprintf("%s", param.Name.CamelCase()))
	}
	for _, param := range operation.HeaderParams {
		urlParams = append(urlParams, fmt.Sprintf("%s", param.Name.CamelCase()))
	}
	for _, param := range operation.Endpoint.UrlParams {
		urlParams = append(urlParams, fmt.Sprintf("%s", param.Name.CamelCase()))
	}
	return urlParams
}

func addRoutes(api *spec.Api) string {
	return fmt.Sprintf(`Add%sRoutes`, api.Name.PascalCase())
}
