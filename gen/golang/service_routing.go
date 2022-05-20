package golang

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func generateRoutings(version *spec.Version, versionModule module, modelsModule module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *generateRouting(modelsModule, versionModule, &api))
	}
	return files
}

func generateRouting(modelsModule module, versionModule module, api *spec.Api) *sources.CodeFile {
	w := NewGoWriter()
	w.Line("package %s", versionModule.Name)

	imports := Imports()
	imports.Add("encoding/json")
	imports.Add("github.com/husobee/vestigo")
	imports.AddAlias("github.com/sirupsen/logrus", "log")
	imports.Add("net/http")
	imports.Add("fmt")
	imports.Add("strings")
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

	return &sources.CodeFile{
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

func addSetCors(w *sources.Writer, operation *spec.NamedOperation) {
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

func generateOperationParametersParsing(w *sources.Writer, operation *spec.NamedOperation, namedParams []spec.NamedParam, isUrlParam bool, paramsParserName string, paramName string, parseCommaSeparatedArray bool) {
	if namedParams != nil && len(namedParams) > 0 {
		w.Line(`%s := NewParamsParser(%s, %t)`, paramsParserName, paramName, parseCommaSeparatedArray)
		for _, param := range namedParams {
			w.Line(`%s := %s`, param.Name.CamelCase(), parserParameterCall(isUrlParam, &param, paramsParserName))
		}

		w.Line(`if len(%s.Errors) > 0 {`, paramsParserName)
		generateBadRequestResponse(w.Indented(), operation, fmt.Sprintf(`"Can't parse %s"`, paramsParserName))
		w.Line(`}`)
	}
}

func generateServiceCall(w *sources.Writer, operation *spec.NamedOperation) {
	singleEmptyResponse := len(operation.Responses) == 1 && operation.Responses[0].Type.Definition.IsEmpty()
	serviceCall := fmt.Sprintf(`%s.%s(%s)`, serviceInterfaceTypeVar(operation.Api), operation.Name.PascalCase(), JoinDelimParams(addOperationMethodParams(operation)))
	if singleEmptyResponse {
		w.Line(`err = %s`, serviceCall)
	} else {
		w.Line(`response, err := %s`, serviceCall)
	}

	w.Line(`if err != nil {`)
	generateInternalServerErrorResponse(w.Indented(), operation, genFmtSprintf("Error returned from service implementation: %s", `err.Error()`))
	w.Line(`}`)

	if !singleEmptyResponse {
		w.Line(`if response == nil {`)
		generateInternalServerErrorResponse(w.Indented(), operation, `"No result returned from service implementation"`)
		w.Line(`}`)
	}
}

func generateResponseWriting(w *sources.Writer, response *spec.NamedResponse, responseVar string) {
	if response.BodyIs(spec.BodyString) {
		w.Line(`res.Header().Set("Content-Type", "text/plain")`)
		w.Line(`res.WriteHeader(%s)`, spec.HttpStatusCode(response.Name))
		w.Line(`res.Write([]byte(*response))`)
	} else if response.BodyIs(spec.BodyJson) {
		w.Line(`res.Header().Set("Content-Type", "application/json")`)
		w.Line(`res.WriteHeader(%s)`, spec.HttpStatusCode(response.Name))
		w.Line(`json.NewEncoder(res).Encode(%s)`, responseVar)
	} else {
		w.Line(`res.WriteHeader(%s)`, spec.HttpStatusCode(response.Name))
	}
	w.Line(`log.WithFields(%s).WithField("status", %s).Info("Completed request")`, logFieldsName(response.Operation), spec.HttpStatusCode(response.Name))
	w.Line(`return`)
}

func checkRequestContentType(w *sources.Writer, operation *spec.NamedOperation, contentType string) {
	w.Line(`contentType := req.Header.Get("Content-Type")`)
	w.Line(`if !strings.Contains(contentType, "%s") {`, contentType)
	generateBadRequestResponse(w.Indented(), operation, genFmtSprintf("Wrong Content-type: %s", "contentType"))
	w.Line(`}`)
}

func generateOperationMethod(w *sources.Writer, operation *spec.NamedOperation) {
	w.Line(`log.WithFields(%s).Info("Received request")`, logFieldsName(operation))
	w.Line(`var err error`)
	if operation.BodyIs(spec.BodyString) {
		checkRequestContentType(w, operation, "text/plain")
		w.Line(`bodyData, err := ioutil.ReadAll(req.Body)`)
		w.Line(`if err != nil {`)
		generateBadRequestResponse(w.Indented(), operation, genFmtSprintf(`Reading request body failed: %s`, `err.Error()`))
		w.Line(`}`)
		w.Line(`body := string(bodyData)`)
	}
	if operation.BodyIs(spec.BodyJson) {
		checkRequestContentType(w, operation, "application/json")
		w.Line(`var body %s`, GoType(&operation.Body.Type.Definition))
		w.Line(`err = json.NewDecoder(req.Body).Decode(&body)`)
		w.Line(`if err != nil {`)
		generateBadRequestResponse(w.Indented(), operation, genFmtSprintf(`Decoding body JSON failed: %s`, `err.Error()`))
		w.Line(`}`)
	}
	generateOperationParametersParsing(w, operation, operation.QueryParams, false, "queryParams", "req.URL.Query()", true)
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
		generateInternalServerErrorResponse(w, operation, `"Result from service implementation does not have anything in it"`)
	}
}

func addOperationMethodParams(operation *spec.NamedOperation) []string {
	urlParams := []string{}
	if operation.BodyIs(spec.BodyString) {
		urlParams = append(urlParams, "body")
	}
	if operation.BodyIs(spec.BodyJson) {
		urlParams = append(urlParams, "&body")
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

func generateBadRequestResponse(w *sources.Writer, operation *spec.NamedOperation, message string) {
	badRequest := operation.Errors.Get(spec.HttpStatusBadRequest)
	w.Line(`message := %s`, message)
	w.Line(`log.WithFields(%s).Warn(message)`, logFieldsName(operation))
	w.Line(`errorResponse := %s{message, nil}`, GoType(&badRequest.Type.Definition))
	generateResponseWriting(w, badRequest, `errorResponse`)
}

func generateInternalServerErrorResponse(w *sources.Writer, operation *spec.NamedOperation, message string) {
	internalServerError := operation.Errors.Get(spec.HttpStatusInternalServerError)
	w.Line(`message := %s`, message)
	w.Line(`log.WithFields(%s).Error(message)`, logFieldsName(operation))
	w.Line(`errorResponse := %s{message}`, GoType(&internalServerError.Type.Definition))
	generateResponseWriting(w, internalServerError, `errorResponse`)
}

func genFmtSprintf(format string, args ...string) string {
	if len(args) > 0 {
		return fmt.Sprintf(`fmt.Sprintf("%s", %s)`, format, strings.Join(args, ", "))
	} else {
		return format
	}
}
