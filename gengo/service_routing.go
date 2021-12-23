package gengo

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func generateRoutings(version *spec.Version, versionModule module, modelsModule module) []gen.TextFile {
	files := []gen.TextFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *generateRouting(modelsModule, versionModule, &api))
	}
	return files
}

func generateRouting(modelsModule module, versionModule module, api *spec.Api) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", versionModule.Name)

	imports := Imports()
	imports.Add("encoding/json")
	imports.Add("github.com/husobee/vestigo")
	imports.AddAlias("github.com/sirupsen/logrus", "log")
	imports.Add("net/http")
	if bodyHasType(api, spec.TypeString) {
		imports.Add("io/ioutil")
	}
	apiModule := versionModule.Submodule(api.Name.SnakeCase())
	imports.Add(apiModule.Package)
	if paramsContainsModel(api) {
		imports.Add(modelsModule.Package)
	}
	imports.Write(w)

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

	return &gen.TextFile{
		Path:    versionModule.GetPath(fmt.Sprintf("%s_routing.go", api.Name.SnakeCase())),
		Content: w.String(),
	}
}

func getVestigoUrl(operation *spec.NamedOperation) string {
	url := operation.FullUrl()
	if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
		for _, param := range operation.Endpoint.UrlParams {
			url = strings.Replace(url, spec.UrlParamStr(&param), fmt.Sprintf(":%s", param.Name.Source), -1)
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

func parserParameterCall(isUrlParam bool, param *spec.NamedParam, paramsParserName string) string {
	paramNameSource := param.Name.Source
	if isUrlParam {
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

func generateOperationParametersParsing(w *gen.Writer, operation *spec.NamedOperation, namedParams []spec.NamedParam, isUrlParam bool, paramsParserName string, paramName string, parseCommaSeparatedArray bool) {
	if namedParams != nil && len(namedParams) > 0 {
		w.Line(`%s := NewParamsParser(%s, %t)`, paramsParserName, paramName, parseCommaSeparatedArray)
		for _, param := range namedParams {
			w.Line(`%s := %s`, param.Name.CamelCase(), parserParameterCall(isUrlParam, &param, paramsParserName))
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
	singleEmptyResponse := len(operation.Responses) == 1 && operation.Responses[0].Type.Definition.IsEmpty()
	serviceCall := fmt.Sprintf(`%s.%s(%s)`, serviceInterfaceTypeVar(operation.Api), operation.Name.PascalCase(), JoinDelimParams(addOperationMethodParams(operation)))
	if singleEmptyResponse {
		w.Line(`err = %s`, serviceCall)
	} else {
		w.Line(`response, err := %s`, serviceCall)
	}

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
		if response.Type.Definition.Plain == spec.TypeString {
			w.Line(`res.Write([]byte(*response))`)
		} else {
			w.Line(`json.NewEncoder(res).Encode(%s)`, responseVar)
		}
	}
	w.Line(`log.WithFields(%s).WithField("status", %s).Info("Completed request")`, logFieldsName(response.Operation), spec.HttpStatusCode(response.Name))
	w.Line(`return`)
}

func generateOperationMethod(w *gen.Writer, operation *spec.NamedOperation) {
	w.Line(`log.WithFields(%s).Info("Received request")`, logFieldsName(operation))
	w.Line(`var err error`)
	if operation.Body != nil {
		if operation.Body.Type.Definition.Plain == spec.TypeString {
			w.Line(`bodyData, err := ioutil.ReadAll(req.Body)`)
			w.Line(`if err != nil {`)
			w.Line(`  log.WithFields(%s).Warnf("%s", err.Error())`, logFieldsName(operation), `Reading request body failed: %s`)
			w.Line(`  res.WriteHeader(400)`)
			w.Line(`  log.WithFields(%s).WithField("status", 400).Info("Completed request")`, logFieldsName(operation))
			w.Line(`  return`)
			w.Line(`}`)
			w.Line(`body := string(bodyData)`)
		} else {
			w.Line(`var body %s`, GoType(&operation.Body.Type.Definition))
			w.Line(`err = json.NewDecoder(req.Body).Decode(&body)`)
			w.Line(`if err != nil {`)
			w.Line(`  log.WithFields(%s).Warnf("%s", err.Error())`, logFieldsName(operation), `Decoding body JSON failed: %s`)
			w.Line(`  res.WriteHeader(400)`)
			w.Line(`  log.WithFields(%s).WithField("status", 400).Info("Completed request")`, logFieldsName(operation))
			w.Line(`  return`)
			w.Line(`}`)
		}
	}
	generateOperationParametersParsing(w, operation, operation.QueryParams, false, "queryParams", "req.URL.Query()", false)
	generateOperationParametersParsing(w, operation, operation.HeaderParams, false, "headerParams", "req.Header", true)
	generateOperationParametersParsing(w, operation, operation.Endpoint.UrlParams, true, "urlParams", "req.URL.Query()", false)

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
		if operation.Body.Type.Definition.Plain == spec.TypeString {
			urlParams = append(urlParams, "body")
		} else {
			urlParams = append(urlParams, "&body")
		}
	}
	for _, param := range operation.QueryParams {
		urlParams = append(urlParams, param.Name.CamelCase())
	}
	for _, param := range operation.HeaderParams {
		urlParams = append(urlParams, param.Name.CamelCase())
	}
	for _, param := range operation.Endpoint.UrlParams {
		urlParams = append(urlParams, param.Name.CamelCase())
	}
	return urlParams
}

func addRoutes(api *spec.Api) string {
	return fmt.Sprintf(`Add%sRoutes`, api.Name.PascalCase())
}
