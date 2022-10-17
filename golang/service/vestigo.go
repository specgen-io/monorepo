package service

import (
	"fmt"
	"github.com/pinzolo/casee"
	"strings"

	"generator"
	"golang/imports"
	"golang/models"
	"golang/module"
	"golang/types"
	"golang/writer"
	"spec"
)

func generateRoutings(version *spec.Version, versionModule, routingModule, contentTypeModule, errorsModule, errorsModelsModule, modelsModule, paramsParserModule, respondModule module.Module, modelsGenerator models.Generator) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *generateRouting(&api, versionModule, routingModule, contentTypeModule, errorsModule, errorsModelsModule, modelsModule, paramsParserModule, respondModule, modelsGenerator))
	}
	return files
}

func signatureAddRouting(api *spec.Api) string {
	fullServiceInterfaceName := fmt.Sprintf("%s.%s", api.Name.SnakeCase(), serviceInterfaceName)
	return fmt.Sprintf(`%s(router *vestigo.Router, %s %s)`, addRoutesMethodName(api), serviceInterfaceTypeVar(api), fullServiceInterfaceName)
}

func callAddRouting(api *spec.Api, serviceVar string) string {
	return fmt.Sprintf(`%s(router, %s)`, addRoutesMethodName(api), serviceVar)
}

func generateRouting(api *spec.Api, versionModule, routingModule, contentTypeModule, errorsModule, errorsModelsModule, modelsModule, paramsParserModule, respondModule module.Module, modelsGenerator models.Generator) *generator.CodeFile {
	apiModule := versionModule.Submodule(api.Name.SnakeCase())

	w := writer.New(routingModule, fmt.Sprintf("%s.go", api.Name.SnakeCase()))

	imports := imports.New()
	if types.ApiHasBody(api) {
		imports.Add("encoding/json")
	}
	imports.Add("github.com/husobee/vestigo")
	imports.AddAliased("github.com/sirupsen/logrus", "log")
	imports.Add("net/http")
	imports.Add("fmt")
	if types.BodyHasType(api, spec.TypeString) {
		imports.Add("io/ioutil")
	}
	if hasNonEmptyBody(api) {
		imports.Module(contentTypeModule)
	}
	imports.Module(apiModule)
	imports.Module(errorsModule)
	imports.Module(errorsModelsModule)
	if isRouterUsingModels(api) {
		imports.Module(modelsModule)
	}
	if operationHasParams(api) {
		imports.Module(paramsParserModule)
	}
	imports.Module(respondModule)
	imports.Write(w)

	w.EmptyLine()
	w.Line(`func %s {`, signatureAddRouting(api))
	w.Indent()
	for _, operation := range api.Operations {
		url := getEndpointUrl(&operation)
		w.Line(`%s := log.Fields{"operationId": "%s.%s", "method": "%s", "url": "%s"}`, logFieldsName(&operation), operation.InApi.Name.Source, operation.Name.Source, casee.ToUpperCase(operation.Endpoint.Method), url)
		w.Line(`router.%s("%s", func(res http.ResponseWriter, req *http.Request) {`, casee.ToPascalCase(operation.Endpoint.Method), url)
		generateOperationMethod(w.Indented(), &operation, modelsGenerator)
		w.Line(`})`)
		if operation.HeaderParams != nil && len(operation.HeaderParams) > 0 {
			addSetCors(w, &operation)
		}
		w.EmptyLine()
	}
	w.Unindent()
	w.Line(`}`)

	return w.ToCodeFile()
}

func operationHasParams(api *spec.Api) bool {
	for _, operation := range api.Operations {
		for _, param := range operation.QueryParams {
			if &param != nil {
				return true
			}
		}
		for _, param := range operation.HeaderParams {
			if &param != nil {
				return true
			}
		}
		for _, param := range operation.Endpoint.UrlParams {
			if &param != nil {
				return true
			}
		}
	}
	return false
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

func addSetCors(w generator.Writer, operation *spec.NamedOperation) {
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

func parserParameterCall(isUrlParam bool, param *spec.NamedParam, paramsParserName string, modelsGenerator models.Generator) string {
	paramNameSource := param.Name.Source
	if isUrlParam {
		paramNameSource = ":" + paramNameSource
	}
	parserParams := []string{fmt.Sprintf(`"%s"`, paramNameSource)}
	methodName, defaultParam := parserDefaultName(param)
	isEnum := param.Type.Definition.Info.Model != nil && param.Type.Definition.Info.Model.IsEnum()
	enumModel := param.Type.Definition.Info.Model
	if isEnum {
		parserParams = append(parserParams, fmt.Sprintf("%s.%s", types.VersionModelsPackage, modelsGenerator.EnumValuesStrings(enumModel)))
	}
	if defaultParam != nil {
		parserParams = append(parserParams, *defaultParam)
	}
	call := fmt.Sprintf(`%s.%s(%s)`, paramsParserName, methodName, strings.Join(parserParams, ", "))
	if isEnum {
		call = fmt.Sprintf(`%s.%s(%s)`, types.VersionModelsPackage, enumModel.Name.PascalCase(), call)
	}
	return call
}

func generateHeaderParsing(w generator.Writer, operation *spec.NamedOperation, modelsGenerator models.Generator) {
	generateParametersParsing(w, operation, operation.HeaderParams, "header", "req.Header", modelsGenerator)
}

func generateQueryParsing(w generator.Writer, operation *spec.NamedOperation, modelsGenerator models.Generator) {
	generateParametersParsing(w, operation, operation.QueryParams, "query", "req.URL.Query()", modelsGenerator)
}

func generateUrlParamsParsing(w generator.Writer, operation *spec.NamedOperation, modelsGenerator models.Generator) {
	if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
		w.Line(`urlParams := paramsparser.New(req.URL.Query(), false)`)
		for _, param := range operation.Endpoint.UrlParams {
			w.Line(`%s := %s`, param.Name.CamelCase(), parserParameterCall(true, &param, "urlParams", modelsGenerator))
		}
		w.Line(`if len(urlParams.Errors) > 0 {`)
		respondNotFound(w.Indented(), operation, fmt.Sprintf(`"Failed to parse url parameters"`))
		w.Line(`}`)
	}
}

func generateParametersParsing(w generator.Writer, operation *spec.NamedOperation, namedParams []spec.NamedParam, paramsParserName string, paramsValuesVar string, modelsGenerator models.Generator) {
	if namedParams != nil && len(namedParams) > 0 {
		w.Line(`%s := paramsparser.New(%s, true)`, paramsParserName, paramsValuesVar)
		for _, param := range namedParams {
			w.Line(`%s := %s`, param.Name.CamelCase(), parserParameterCall(false, &param, paramsParserName, modelsGenerator))
		}

		w.Line(`if len(%s.Errors) > 0 {`, paramsParserName)
		respondBadRequest(w.Indented(), operation, paramsParserName, fmt.Sprintf(`"Failed to parse %s"`, paramsParserName), fmt.Sprintf(`httperrors.Convert(%s.Errors)`, paramsParserName))
		w.Line(`}`)
	}
}

func generateServiceCall(w generator.Writer, operation *spec.NamedOperation, responseVar string) {
	singleEmptyResponse := len(operation.Responses) == 1 && operation.Responses[0].Type.Definition.IsEmpty()
	serviceCall := serviceCall(serviceInterfaceTypeVar(operation.InApi), operation)
	if singleEmptyResponse {
		w.Line(`err = %s`, serviceCall)
	} else {
		w.Line(`%s, err := %s`, responseVar, serviceCall)
	}

	w.Line(`if err != nil {`)
	respondInternalServerError(w.Indented(), operation, genFmtSprintf("Error returned from service implementation: %s", `err.Error()`))
	w.Line(`}`)

	if !singleEmptyResponse {
		w.Line(`if response == nil {`)
		respondInternalServerError(w.Indented(), operation, `"Service implementation returned nil"`)
		w.Line(`}`)
	}
}

func generateResponseWriting(w generator.Writer, logFieldsName string, response *spec.Response, responseVar string) {
	if response.BodyIs(spec.BodyEmpty) {
		w.Line(respondEmpty(logFieldsName, `res`, spec.HttpStatusCode(response.Name)))
	}
	if response.BodyIs(spec.BodyString) {
		w.Line(respondText(logFieldsName, `res`, spec.HttpStatusCode(response.Name), `*`+responseVar))
	}
	if response.BodyIs(spec.BodyJson) {
		w.Line(respondJson(logFieldsName, `res`, spec.HttpStatusCode(response.Name), responseVar))
	}
}

func generateOperationMethod(w generator.Writer, operation *spec.NamedOperation, modelsGenerator models.Generator) {
	w.Line(`log.WithFields(%s).Info("Received request")`, logFieldsName(operation))
	w.Line(`var err error`)
	generateUrlParamsParsing(w, operation, modelsGenerator)
	generateHeaderParsing(w, operation, modelsGenerator)
	generateQueryParsing(w, operation, modelsGenerator)
	generateBodyParsing(w, operation)
	generateServiceCall(w, operation, `response`)
	generateResponse(w, operation, `response`)
}

func generateResponse(w generator.Writer, operation *spec.NamedOperation, responseVar string) {
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
		respondInternalServerError(w, operation, `"Result from service implementation does not have anything in it"`)
	}
}

func generateBodyParsing(w generator.Writer, operation *spec.NamedOperation) {
	if operation.BodyIs(spec.BodyString) {
		w.Line(`if !%s {`, callCheckContentType(logFieldsName(operation), `"text/plain"`, "req", "res"))
		w.Line(`  return`)
		w.Line(`}`)
		w.Line(`bodyData, err := ioutil.ReadAll(req.Body)`)
		w.Line(`if err != nil {`)
		respondBadRequest(w.Indented(), operation, "body", genFmtSprintf(`Reading request body failed: %s`, `err.Error()`), "nil")
		w.Line(`}`)
		w.Line(`body := string(bodyData)`)
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`if !%s {`, callCheckContentType(logFieldsName(operation), `"application/json"`, "req", "res"))
		w.Line(`  return`)
		w.Line(`}`)
		w.Line(`var body %s`, types.GoType(&operation.Body.Type.Definition))
		w.Line(`err = json.NewDecoder(req.Body).Decode(&body)`)
		w.Line(`if err != nil {`)
		w.Line(`  var errors []errmodels.ValidationError = nil`)
		w.Line(`  if unmarshalError, ok := err.(*json.UnmarshalTypeError); ok {`)
		w.Line(`    message := fmt.Sprintf("Failed to parse JSON, field: PERCENT_s", unmarshalError.Field)`)
		w.Line(`    errors = []errmodels.ValidationError{{Path: unmarshalError.Field, Code: "parsing_failed", Message: &message}}`)
		w.Line(`  }`)
		respondBadRequest(w.Indented(), operation, "body", `"Failed to parse body"`, "errors")
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

func generateSpecRouting(specification *spec.Spec, rootModule module.Module) *generator.CodeFile {
	w := writer.New(rootModule, "spec.go")

	imports := imports.New()
	imports.Add("github.com/husobee/vestigo")
	for _, version := range specification.Versions {
		versionModule := rootModule.Submodule(version.Name.FlatCase())
		routingModule := versionModule.SubmoduleAliased("routing", routingPackageAlias(&version))
		imports.ModuleAliased(routingModule)
		for _, api := range version.Http.Apis {
			apiModule := versionModule.SubmoduleAliased(api.Name.SnakeCase(), apiPackageAlias(&api))
			imports.ModuleAliased(apiModule)
		}
	}
	imports.Write(w)

	w.EmptyLine()
	routesParams := []string{}
	for _, version := range specification.Versions {
		versionModule := rootModule.Submodule(version.Name.FlatCase())
		for _, api := range version.Http.Apis {
			apiModule := versionModule.SubmoduleAliased(api.Name.SnakeCase(), apiPackageAlias(&api))
			routesParams = append(routesParams, fmt.Sprintf(`%s %s`, serviceApiNameVersioned(&api), apiModule.Get(serviceInterfaceName)))
		}
	}
	w.Line(`func AddRoutes(router *vestigo.Router, %s) {`, strings.Join(routesParams, ", "))
	for _, version := range specification.Versions {
		versionModule := rootModule.Submodule(version.Name.FlatCase())
		routingModule := versionModule.SubmoduleAliased("routing", routingPackageAlias(&version))
		for _, api := range version.Http.Apis {
			w.Line(`  %s(router, %s)`, routingModule.Get(addRoutesMethodName(&api)), serviceApiNameVersioned(&api))
		}
	}
	w.Line(`}`)

	return w.ToCodeFile()
}

func routingPackageAlias(version *spec.Version) string {
	if version.Name.Source != "" {
		return fmt.Sprintf(`%s`, version.Name.FlatCase())
	} else {
		return fmt.Sprintf(`root`)
	}
}

func apiPackageAlias(api *spec.Api) string {
	version := api.InHttp.InVersion.Name
	if version.Source != "" {
		return api.Name.CamelCase() + version.PascalCase()
	}
	return api.Name.CamelCase()
}

func serviceApiNameVersioned(api *spec.Api) string {
	return fmt.Sprintf(`%sService%s`, api.Name.Source, api.InHttp.InVersion.Name.PascalCase())
}

func checkContentType(contentTypeModule, errorsModule, errorsModelsModule module.Module) *generator.CodeFile {
	w := writer.New(contentTypeModule, `check.go`)
	w.Template(
		map[string]string{
			`ErrorsPackage`:       errorsModule.Package,
			`ErrorsModelsPackage`: errorsModelsModule.Package,
		}, `
import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"[[.ErrorsPackage]]"
	"[[.ErrorsModelsPackage]]"
)

func Check(logFields log.Fields, expectedContentType string, req *http.Request, res http.ResponseWriter) bool {
	contentType := req.Header.Get("Content-Type")
	if !strings.Contains(contentType, expectedContentType) {
		message := fmt.Sprintf("Expected Content-Type header: '%s' was not provided, found: '%s'", expectedContentType, contentType)
		httperrors.RespondBadRequest(logFields, res, &errmodels.BadRequestError{Location: "header", Message: "Failed to parse header", Errors: []errmodels.ValidationError{{Path: "Content-Type", Code: "missing", Message: &message}}})
		return false
	}
	return true
}
`)
	return w.ToCodeFile()
}

func httpErrors(converterModule, errorsModelsModule, paramsParserModule, respondModule module.Module, responses *spec.Responses) []generator.CodeFile {
	files := []generator.CodeFile{}

	files = append(files, *errorsModelsConverter(converterModule, errorsModelsModule, paramsParserModule))
	files = append(files, *generateErrors(converterModule, errorsModelsModule, respondModule, responses))

	return files
}

func errorsModelsConverter(converterModule, errorsModelsModule, paramsParserModule module.Module) *generator.CodeFile {
	w := writer.New(converterModule, `converter.go`)
	w.Template(
		map[string]string{
			`ErrorsModelsPackage`: errorsModelsModule.Package,
			`ParamsParserModule`:  paramsParserModule.Package,
		}, `
import (
	"[[.ErrorsModelsPackage]]"
	"[[.ParamsParserModule]]"
)

func Convert(parsingErrors []paramsparser.ParsingError) []errmodels.ValidationError {
	var validationErrors []errmodels.ValidationError

	for _, parsingError := range parsingErrors {
		validationError := errmodels.ValidationError(parsingError)
		validationErrors = append(validationErrors, validationError)
	}

	return validationErrors
}
`)
	return w.ToCodeFile()
}

func hasNonEmptyBody(api *spec.Api) bool {
	hasNonEmptyBody := false
	walk := spec.NewWalker().
		OnOperation(func(operation *spec.NamedOperation) {
			if operation.BodyIs(spec.BodyJson) || operation.BodyIs(spec.BodyString) {
				hasNonEmptyBody = true
			}
		})
	walk.Api(api)
	return hasNonEmptyBody
}

func isRouterUsingModels(api *spec.Api) bool {
	usingModels := false
	walk := spec.NewWalker().
		OnOperation(func(operation *spec.NamedOperation) {
			if operation.Body != nil {
				if types.IsModel(&operation.Body.Type.Definition) {
					usingModels = true
				}
			}
		}).
		OnParam(func(param *spec.NamedParam) {
			if param.Type.Definition.Info.Model != nil {
				usingModels = true
			}
		})
	walk.Api(api)
	return usingModels
}
