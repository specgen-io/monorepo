package service

import (
	"fmt"
	"generator"
	"github.com/pinzolo/casee"
	"golang/models"
	"golang/types"
	"golang/walkers"
	"golang/writer"
	"spec"
	"strings"
)

var HttpRouter = "httprouter"

type HttpRouterGenerator struct {
	Types   *types.Types
	Models  models.Generator
	Modules *Modules
}

func NewHttpRouterGenerator(types *types.Types, models models.Generator, modules *Modules) *HttpRouterGenerator {
	return &HttpRouterGenerator{types, models, modules}
}

func (g *HttpRouterGenerator) Routings(version *spec.Version) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *g.routing(&api))
	}
	return files
}

func (g *HttpRouterGenerator) routing(api *spec.Api) *generator.CodeFile {
	w := writer.New(g.Modules.Routing(api.InHttp.InVersion), fmt.Sprintf("%s.go", api.Name.SnakeCase()))

	w.Imports.Add("github.com/julienschmidt/httprouter")
	w.Imports.AddAliased("github.com/sirupsen/logrus", "log")
	w.Imports.Add("net/http")
	w.Imports.Add("fmt")
	if walkers.ApiHasBodyOfKind(api, spec.RequestBodyString) {
		w.Imports.Add("io/ioutil")
	}
	if walkers.ApiHasBodyOfKind(api, spec.RequestBodyJson) {
		w.Imports.Add("encoding/json")
	}
	if walkers.ApiHasBodyOfKind(api, spec.RequestBodyJson) || walkers.ApiHasBodyOfKind(api, spec.RequestBodyString) {
		w.Imports.Module(g.Modules.ContentType)
	}
	w.Imports.Module(g.Modules.ServicesApi(api))
	w.Imports.Module(g.Modules.HttpErrors)
	w.Imports.Module(g.Modules.HttpErrorsModels)
	if walkers.ApiIsUsingModels(api) {
		w.Imports.Module(g.Modules.Models(api.InHttp.InVersion))
	}
	if operationHasParams(api) {
		w.Imports.Module(g.Modules.ParamsParser)
	}
	w.Imports.Module(g.Modules.Respond)
	w.Line(`func %s(router *httprouter.Router, %s %s) {`, addRoutesMethodName(api), serviceInterfaceTypeVar(api), g.Modules.ServicesApi(api).Get(serviceInterfaceName))
	w.Indent()
	for _, operation := range api.Operations {
		url := g.getEndpointUrl(&operation)
		w.Line(`%s := log.Fields{"operationId": "%s.%s", "method": "%s", "url": "%s"}`, logFieldsName(&operation), operation.InApi.Name.Source, operation.Name.Source, casee.ToUpperCase(operation.Endpoint.Method), url)
		w.Line(`router.%s("%s", func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {`, casee.ToUpperCase(operation.Endpoint.Method), url)
		g.operation(w.Indented(), &operation)
		w.Line(`})`)
		w.EmptyLine()
	}
	w.Unindent()
	w.Line(`}`)

	return w.ToCodeFile()
}

func (g *HttpRouterGenerator) getEndpointUrl(operation *spec.NamedOperation) string {
	url := operation.FullUrl()
	if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
		for _, param := range operation.Endpoint.UrlParams {
			url = strings.Replace(url, spec.UrlParamStr(&param), fmt.Sprintf(":%s", param.Name.Source), -1)
		}
	}
	return url
}

func (g *HttpRouterGenerator) addSetCors(w *writer.Writer, operation *spec.NamedOperation) {
	params := []string{}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf(`%s`, param.Name.Source))
	}
	w.Line(`res.Header().Set("Access-Control-Allow-Headers", "%s")`, strings.Join(params, ", "))
}

func (g *HttpRouterGenerator) parserParameterCall(param *spec.NamedParam, paramsParserName string) string {
	parserParams := []string{fmt.Sprintf(`"%s"`, param.Name.Source)}
	methodName, defaultParam := parserDefaultName(param)
	isEnum := param.Type.Definition.Info.Model != nil && param.Type.Definition.Info.Model.IsEnum()
	enumModel := param.Type.Definition.Info.Model
	if isEnum {
		parserParams = append(parserParams, fmt.Sprintf("%s.%s", types.VersionModelsPackage, g.Models.EnumValuesStrings(enumModel)))
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

func (g *HttpRouterGenerator) headerParsing(w *writer.Writer, operation *spec.NamedOperation) {
	g.parametersParsing(w, operation, operation.HeaderParams, "header", "req.Header")
}

func (g *HttpRouterGenerator) queryParsing(w *writer.Writer, operation *spec.NamedOperation) {
	g.parametersParsing(w, operation, operation.QueryParams, "query", "req.URL.Query()")
}

func (g *HttpRouterGenerator) urlParamsParsing(w *writer.Writer, operation *spec.NamedOperation) {
	if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
		w.Line(`urlParams := paramsparser.NewUrlParser(params, false)`)
		for _, param := range operation.Endpoint.UrlParams {
			w.Line(`%s := %s`, param.Name.CamelCase(), g.parserParameterCall(&param, "urlParams"))
		}
		w.Line(`if len(urlParams.Errors) > 0 {`)
		respondNotFound(w.Indented(), operation, g.Types, fmt.Sprintf(`"Failed to parse url parameters"`))
		w.Line(`}`)
	}
}

func (g *HttpRouterGenerator) parametersParsing(w *writer.Writer, operation *spec.NamedOperation, namedParams []spec.NamedParam, paramsParserName string, paramsValuesVar string) {
	if namedParams != nil && len(namedParams) > 0 {
		w.Line(`%s := paramsparser.New(%s, true)`, paramsParserName, paramsValuesVar)
		for _, param := range namedParams {
			w.Line(`%s := %s`, param.Name.CamelCase(), g.parserParameterCall(&param, paramsParserName))
		}
		w.Line(`if len(%s.Errors) > 0 {`, paramsParserName)
		respondBadRequest(w.Indented(), operation, g.Types, paramsParserName, fmt.Sprintf(`"Failed to parse %s"`, paramsParserName), fmt.Sprintf(`httperrors.Convert(%s.Errors)`, paramsParserName))
		w.Line(`}`)
	}
}

func (g *HttpRouterGenerator) serviceCallAndResponseCheck(w *writer.Writer, operation *spec.NamedOperation, responseVar string) {
	singleEmptyResponse := len(operation.Responses) == 1 && operation.Responses[0].Body.Type.Definition.IsEmpty()
	serviceCall := serviceCall(serviceInterfaceTypeVar(operation.InApi), operation)
	if singleEmptyResponse {
		w.Line(`err = %s`, serviceCall)
	} else {
		w.Line(`%s, err := %s`, responseVar, serviceCall)
	}
	w.Line(`if err != nil {`)
	respondInternalServerError(w.Indented(), operation, g.Types, genFmtSprintf("Error returned from service implementation: %s", `err.Error()`))
	w.Line(`}`)
	if !singleEmptyResponse {
		w.Line(`if response == nil {`)
		respondInternalServerError(w.Indented(), operation, g.Types, `"Service implementation returned nil"`)
		w.Line(`}`)
	}
}

func (g *HttpRouterGenerator) operation(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line(`log.WithFields(%s).Info("Received request")`, logFieldsName(operation))
	w.Line(`var err error`)
	g.urlParamsParsing(w, operation)
	g.headerParsing(w, operation)
	g.queryParsing(w, operation)
	g.bodyParsing(w, operation)
	g.serviceCallAndResponseCheck(w, operation, `response`)
	g.response(w, operation, `response`)
}

func (g *HttpRouterGenerator) response(w *writer.Writer, operation *spec.NamedOperation, responseVar string) {
	if len(operation.Responses) == 1 {
		if walkers.OperationHasHeaderParams(operation) {
			g.addSetCors(w, operation)
		}
		writeResponse(w, logFieldsName(operation), &operation.Responses[0].Response, responseVar)
	} else {
		for _, response := range operation.Responses {
			responseVar := fmt.Sprintf("%s.%s", responseVar, response.Name.PascalCase())
			w.Line(`if %s != nil {`, responseVar)
			writeResponse(w.Indented(), logFieldsName(operation), &response.Response, responseVar)
			w.Line(`  return`)
			w.Line(`}`)
		}
		if walkers.OperationHasHeaderParams(operation) {
			g.addSetCors(w, operation)
		}
		respondInternalServerError(w, operation, g.Types, `"Result from service implementation does not have anything in it"`)
	}
}

func (g *HttpRouterGenerator) bodyParsing(w *writer.Writer, operation *spec.NamedOperation) {
	if operation.Body != nil {
		w.Line(`if !%s {`, callCheckContentType(logFieldsName(operation), fmt.Sprintf(`"%s"`, ContentType(operation)), "req", "res"))
		w.Line(`  return`)
		w.Line(`}`)
	}
	if operation.BodyIs(spec.RequestBodyString) {
		w.Line(`bodyData, err := ioutil.ReadAll(req.Body)`)
		w.Line(`if err != nil {`)
		respondBadRequest(w.Indented(), operation, g.Types, "body", genFmtSprintf(`Reading request body failed: %s`, `err.Error()`), "nil")
		w.Line(`}`)
		w.Line(`body := string(bodyData)`)
	}
	if operation.BodyIs(spec.RequestBodyJson) {
		w.Line(`var body %s`, g.Types.GoType(&operation.Body.Type.Definition))
		w.Line(`err = json.NewDecoder(req.Body).Decode(&body)`)
		w.Line(`if err != nil {`)
		w.Line(`  var errors []errmodels.ValidationError = nil`)
		w.Line(`  if unmarshalError, ok := err.(*json.UnmarshalTypeError); ok {`)
		w.Line(`    message := fmt.Sprintf("Failed to parse JSON, field: PERCENT_s", unmarshalError.Field)`)
		w.Line(`    errors = []errmodels.ValidationError{{Path: unmarshalError.Field, Code: "parsing_failed", Message: &message}}`)
		w.Line(`  }`)
		respondBadRequest(w.Indented(), operation, g.Types, "body", `"Failed to parse body"`, "errors")
		w.Line(`}`)
	}
	if operation.BodyIs(spec.RequestBodyFormData) || operation.BodyIs(spec.RequestBodyFormUrlEncoded) {
		w.Line(`formBody, err := paramsparser.New%sParser(req, true)`, casee.ToPascalCase(formBodyTypeName(operation)))
		w.Line(`if err != nil {`)
		respondBadRequest(w.Indented(), operation, g.Types, "body", `"Failed to parse body"`, fmt.Sprintf(`[]errmodels.ValidationError{{Path: "", Code: "%s_parse_failed"}}`, formBodyTypeName(operation)))
		w.Line(`}`)
		for _, param := range operation.Body.FormData {
			w.Line(`%s := %s`, param.Name.CamelCase(), g.parserParameterCall(&param, "formBody"))
		}
		for _, param := range operation.Body.FormUrlEncoded {
			w.Line(`%s := %s`, param.Name.CamelCase(), g.parserParameterCall(&param, "formBody"))
		}
		w.Line(`if len(formBody.Errors) > 0 {`)
		respondBadRequest(w.Indented(), operation, g.Types, "body", fmt.Sprintf(`"Failed to parse body"`), fmt.Sprintf(`httperrors.Convert(formBody.Errors)`))
		w.Line(`}`)
	}
}

func (g *HttpRouterGenerator) RootRouting(specification *spec.Spec) *generator.CodeFile {
	w := writer.New(g.Modules.Root, "spec.go")

	w.Imports.Add("github.com/julienschmidt/httprouter")
	for _, version := range specification.Versions {
		w.Imports.ModuleAliased(g.Modules.Routing(&version).Aliased(routingPackageAlias(&version)))
		for _, api := range version.Http.Apis {
			w.Imports.ModuleAliased(g.Modules.ServicesApi(&api).Aliased(apiPackageAlias(&api)))
		}
	}
	w.Line(`type Services struct {`)
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			apiModule := g.Modules.ServicesApi(&api).Aliased(apiPackageAlias(&api))
			w.LineAligned(`  %s %s`, serviceApiPublicNameVersioned(&api), apiModule.Get(serviceInterfaceName))
		}
	}
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`func AddRoutes(router *httprouter.Router, services Services) {`)
	for _, version := range specification.Versions {
		routingModule := g.Modules.Routing(&version).Aliased(routingPackageAlias(&version))
		for _, api := range version.Http.Apis {
			w.Line(`  %s(router, services.%s)`, routingModule.Get(addRoutesMethodName(&api)), serviceApiPublicNameVersioned(&api))
		}
	}
	w.Line(`}`)

	return w.ToCodeFile()
}

func (g *HttpRouterGenerator) GenerateUrlParamsCtor() *generator.CodeFile {
	w := writer.New(g.Modules.ParamsParser, `url_parser.go`)

	w.Lines(`
import (
	"github.com/julienschmidt/httprouter"
	"net/url"
)

func NewUrlParser(params httprouter.Params, parseCommaSeparatedArray bool) *ParamsParser {
	values := make(url.Values)
	for _, param := range params {
		values.Add(param.Key, param.Value)
	}
	return &ParamsParser{values, parseCommaSeparatedArray, []ParsingError{}}
}
`)

	return w.ToCodeFile()
}
