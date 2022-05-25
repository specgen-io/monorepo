package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen/golang/client"
	"github.com/specgen-io/specgen/v2/gen/golang/imports"
	"github.com/specgen-io/specgen/v2/gen/golang/models"
	"github.com/specgen-io/specgen/v2/gen/golang/module"
	"github.com/specgen-io/specgen/v2/gen/golang/types"
	"github.com/specgen-io/specgen/v2/gen/golang/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func generateRoutings(version *spec.Version, versionModule module.Module, modelsModule module.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *generateRouting(modelsModule, versionModule, &api))
	}
	files = append(files, *generateErrors(version, versionModule, modelsModule))
	return files
}

func signatureAddRouting(api *spec.Api) string {
	fullServiceInterfaceName := fmt.Sprintf("%s.%s", api.Name.SnakeCase(), serviceInterfaceName)
	return fmt.Sprintf(`%s(router *vestigo.Router, %s %s)`, addRoutesMethodName(api), serviceInterfaceTypeVar(api), fullServiceInterfaceName)
}

func callAddRouting(api *spec.Api, serviceVar string) string {
	return fmt.Sprintf(`%s(router, %s)`, addRoutesMethodName(api), serviceVar)
}

func generateRouting(modelsModule module.Module, versionModule module.Module, api *spec.Api) *sources.CodeFile {
	w := writer.NewGoWriter()
	w.Line("package %s", versionModule.Name)

	imports := imports.Imports()
	imports.Add("encoding/json")
	imports.Add("github.com/husobee/vestigo")
	imports.AddAlias("github.com/sirupsen/logrus", "log")
	imports.Add("net/http")
	imports.Add("fmt")
	imports.Add("strings")
	if types.BodyHasType(api, spec.TypeString) {
		imports.Add("io/ioutil")
	}
	apiModule := versionModule.Submodule(api.Name.SnakeCase())
	imports.Add(apiModule.Package)
	if paramsContainsModel(api) {
		imports.Add(modelsModule.Package)
	}
	imports.Write(w)

	w.EmptyLine()
	w.Line(`func %s {`, signatureAddRouting(api))
	for i, operation := range api.Operations {
		if i != 0 {
			w.EmptyLine()
		}
		w.Indent()
		url := getEndpointUrl(&operation)
		w.Line(`%s := log.Fields{"operationId": "%s.%s", "method": "%s", "url": "%s"}`, logFieldsName(&operation), operation.Api.Name.Source, operation.Name.Source, client.ToUpperCase(operation.Endpoint.Method), url)
		w.Line(`router.%s("%s", func(res http.ResponseWriter, req *http.Request) {`, client.ToPascalCase(operation.Endpoint.Method), url)
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

func getEndpointUrl(operation *spec.NamedOperation) string {
	url := operation.FullUrl()
	if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
		for _, param := range operation.Endpoint.UrlParams {
			url = strings.Replace(url, spec.UrlParamStr(&param), fmt.Sprintf(":%s", param.Name.Source), -1)
		}
	}
	return url
}

func addSetCors(w *sources.Writer, operation *spec.NamedOperation) {
	w.Line(`router.SetCors("%s", &vestigo.CorsAccessControl{`, getEndpointUrl(operation))
	params := []string{}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf(`"%s"`, param.Name.Source))
	}
	w.Line(`  AllowHeaders: []string{%s},`, strings.Join(params, ", "))
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
		parserParams = append(parserParams, fmt.Sprintf("%s.%s", types.ModelsPackage, models.EnumValuesStrings(enumModel)))
	}
	if defaultParam != nil {
		parserParams = append(parserParams, *defaultParam)
	}
	call := fmt.Sprintf(`%s.%s(%s)`, paramsParserName, methodName, strings.Join(parserParams, ", "))
	if isEnum {
		call = fmt.Sprintf(`%s.%s(%s)`, types.ModelsPackage, enumModel.Name.PascalCase(), call)
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

func generateServiceCall(w *sources.Writer, operation *spec.NamedOperation, responseVar string) {
	singleEmptyResponse := len(operation.Responses) == 1 && operation.Responses[0].Type.Definition.IsEmpty()
	serviceCall := serviceCall(serviceInterfaceTypeVar(operation.Api), operation)
	if singleEmptyResponse {
		w.Line(`err = %s`, serviceCall)
	} else {
		w.Line(`%s, err := %s`, responseVar, serviceCall)
	}

	w.Line(`if err != nil {`)
	generateInternalServerErrorResponse(w.Indented(), operation, genFmtSprintf("Error returned from service implementation: %s", `err.Error()`))
	w.Line(`}`)

	if !singleEmptyResponse {
		w.Line(`if response == nil {`)
		generateInternalServerErrorResponse(w.Indented(), operation, `"Service implementation returned nil"`)
		w.Line(`}`)
	}
}

func generateResponseWriting(w *sources.Writer, logFieldsName string, response *spec.Response, responseVar string) {
	if response.BodyIs(spec.BodyEmpty) {
		w.Line(`res.WriteHeader(%s)`, spec.HttpStatusCode(response.Name))
	}
	if response.BodyIs(spec.BodyString) {
		w.Line(`res.Header().Set("Content-Type", "text/plain")`)
		w.Line(`res.WriteHeader(%s)`, spec.HttpStatusCode(response.Name))
		w.Line(`res.Write([]byte(*response))`)
	}
	if response.BodyIs(spec.BodyJson) {
		w.Line(`res.Header().Set("Content-Type", "application/json")`)
		w.Line(`res.WriteHeader(%s)`, spec.HttpStatusCode(response.Name))
		w.Line(`json.NewEncoder(res).Encode(%s)`, responseVar)
	}
	w.Line(`log.WithFields(%s).WithField("status", %s).Info("Completed request")`, logFieldsName, spec.HttpStatusCode(response.Name))
}

func generateOperationMethod(w *sources.Writer, operation *spec.NamedOperation) {
	w.Line(`log.WithFields(%s).Info("Received request")`, logFieldsName(operation))
	w.Line(`var err error`)
	generateBodyParsing(w, operation)
	generateOperationParametersParsing(w, operation, operation.QueryParams, false, "queryParams", "req.URL.Query()", true)
	generateOperationParametersParsing(w, operation, operation.HeaderParams, false, "headerParams", "req.Header", true)
	generateOperationParametersParsing(w, operation, operation.Endpoint.UrlParams, true, "urlParams", "req.URL.Query()", false)
	generateServiceCall(w, operation, `response`)
	generateReponse(w, operation, `response`)
}

func generateReponse(w *sources.Writer, operation *spec.NamedOperation, responseVar string) {
	if len(operation.Responses) == 1 {
		generateResponseWriting(w, logFieldsName(operation), &operation.Responses[0].Response, responseVar)
	} else {
		for _, response := range operation.Responses {
			responseVar := fmt.Sprintf("%s.%s", responseVar, response.Name.PascalCase())
			w.Line(`if %s != nil {`, responseVar)
			generateResponseWriting(w.Indented(), logFieldsName(operation), &response.Response, responseVar)
			w.Line(`  return`)
			w.Line(`}`)
		}
		generateInternalServerErrorResponse(w, operation, `"Result from service implementation does not have anything in it"`)
	}
}

func generateBodyParsing(w *sources.Writer, operation *spec.NamedOperation) {
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
		w.Line(`var body %s`, types.GoType(&operation.Body.Type.Definition))
		w.Line(`err = json.NewDecoder(req.Body).Decode(&body)`)
		w.Line(`if err != nil {`)
		generateBadRequestResponse(w.Indented(), operation, genFmtSprintf(`Decoding body JSON failed: %s`, `err.Error()`))
		w.Line(`}`)
	}
}

func serviceCall(serviceVar string, operation *spec.NamedOperation) string {
	params := []string{}
	if operation.BodyIs(spec.BodyString) {
		params = append(params, "body")
	}
	if operation.BodyIs(spec.BodyJson) {
		params = append(params, "&body")
	}
	for _, param := range operation.QueryParams {
		params = append(params, param.Name.CamelCase())
	}
	for _, param := range operation.HeaderParams {
		params = append(params, param.Name.CamelCase())
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, param.Name.CamelCase())
	}

	return fmt.Sprintf(`%s.%s(%s)`, serviceVar, operation.Name.PascalCase(), strings.Join(params, ", "))
}

func addRoutesMethodName(api *spec.Api) string {
	return fmt.Sprintf(`Add%sRoutes`, api.Name.PascalCase())
}

func genFmtSprintf(format string, args ...string) string {
	if len(args) > 0 {
		return fmt.Sprintf(`fmt.Sprintf("%s", %s)`, format, strings.Join(args, ", "))
	} else {
		return format
	}
}

func serviceInterfaceTypeVar(api *spec.Api) string {
	return fmt.Sprintf(`%sService`, api.Name.Source)
}

func generateSpecRouting(specification *spec.Spec, module module.Module) *sources.CodeFile {
	w := writer.NewGoWriter()
	w.Line("package %s", module.Name)

	imports := imports.Imports()
	imports.Add("github.com/husobee/vestigo")
	for _, version := range specification.Versions {
		versionModule := module.Submodule(version.Version.FlatCase())
		if version.Version.Source != "" {
			imports.Add(versionModule.Package)
		}
		for _, api := range version.Http.Apis {
			apiModule := versionModule.Submodule(api.Name.SnakeCase())
			imports.AddAlias(apiModule.Package, versionedApiImportAlias(&api))
		}
	}
	imports.Write(w)

	w.EmptyLine()
	routesParams := []string{}
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			routesParams = append(routesParams, fmt.Sprintf(`%s %s.%s`, serviceApiNameVersioned(&api), versionedApiImportAlias(&api), serviceInterfaceName))
		}
	}
	w.Line(`func AddRoutes(router *vestigo.Router, %s) {`, strings.Join(routesParams, ", "))
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			w.Line(`  %s%s`, packageFrom(&version), callAddRouting(&api, serviceApiNameVersioned(&api)))
		}
	}
	w.Line(`}`)

	return &sources.CodeFile{
		Path:    module.GetPath("spec_routing.go"),
		Content: w.String(),
	}
}

func versionedApiImportAlias(api *spec.Api) string {
	version := api.Apis.Version.Version
	if version.Source != "" {
		return api.Name.CamelCase() + version.PascalCase()
	}
	return api.Name.CamelCase()
}

func serviceApiNameVersioned(api *spec.Api) string {
	return fmt.Sprintf(`%sService%s`, api.Name.Source, api.Apis.Version.Version.PascalCase())
}

func packageFrom(version *spec.Version) string {
	if version.Version.Source != "" {
		return fmt.Sprintf(`%s.`, version.Version.FlatCase())
	}
	return ""
}
