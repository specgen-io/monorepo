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

var Chi = "chi"

type ChiGenerator struct {
	Types   *types.Types
	Models  models.Generator
	Modules *Modules
}

func NewChiGenerator(types *types.Types, models models.Generator, modules *Modules) *ChiGenerator {
	return &ChiGenerator{types, models, modules}
}

func (g *ChiGenerator) Routings(version *spec.Version) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *g.routing(&api))
	}
	return files
}

func (g *ChiGenerator) routing(api *spec.Api) *generator.CodeFile {
	w := writer.New(g.Modules.Routing(api.InHttp.InVersion), fmt.Sprintf("%s.go", api.Name.SnakeCase()))

	w.Imports.Add("github.com/go-chi/chi/v5")
	if walkers.ApiHasHasHeaderParams(api) {
		w.Imports.Add("github.com/go-chi/cors")
	}
	w.Imports.AddAliased("github.com/sirupsen/logrus", "log")
	w.Imports.Add("net/http")
	w.Imports.Add("fmt")
	if walkers.ApiHasBodyOfKind(api, spec.BodyString) {
		w.Imports.Add("io/ioutil")
	}
	if walkers.ApiHasBodyOfKind(api, spec.BodyJson) {
		w.Imports.Add("encoding/json")
	}
	if walkers.ApiHasBodyOfKind(api, spec.BodyJson) || walkers.ApiHasBodyOfKind(api, spec.BodyString) {
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
	w.Line(`func %s(router *chi.Mux, %s %s) {`, addRoutesMethodName(api), serviceInterfaceTypeVar(api), g.Modules.ServicesApi(api).Get(serviceInterfaceName))
	w.Indent()
	for _, operation := range api.Operations {
		url := g.getEndpointUrl(&operation)
		w.Line(`%s := log.Fields{"operationId": "%s.%s", "method": "%s", "url": "%s"}`, logFieldsName(&operation), operation.InApi.Name.Source, operation.Name.Source, casee.ToUpperCase(operation.Endpoint.Method), url)
		w.Line(`router.%s("%s", func(res http.ResponseWriter, req *http.Request) {`, casee.ToPascalCase(operation.Endpoint.Method), url)
		g.operation(w.Indented(), &operation)
		w.Line(`})`)
		if walkers.OperationHasHeaderParams(&operation) {
			g.addSetCors(w, &operation)
		}
		w.EmptyLine()
	}
	w.Unindent()
	w.Line(`}`)

	return w.ToCodeFile()
}

func (g *ChiGenerator) getEndpointUrl(operation *spec.NamedOperation) string {
	url := operation.FullUrl()
	if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
		for _, param := range operation.Endpoint.UrlParams {
			url = strings.Replace(url, spec.UrlParamStr(&param), fmt.Sprintf("{%s}", param.Name.Source), -1)
		}
	}
	return url
}

func (g *ChiGenerator) addSetCors(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line(`router.With(cors.Handler(cors.Options{`)
	params := []string{}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf(`"%s"`, param.Name.Source))
	}
	w.Line(`  AllowedHeaders: []string{%s},`, strings.Join(params, ", "))
	w.Line(`}))`)
}

func (g *ChiGenerator) parserParameterCall(param *spec.NamedParam, paramsParserName string) string {
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

func (g *ChiGenerator) headerParsing(w *writer.Writer, operation *spec.NamedOperation) {
	g.parametersParsing(w, operation, operation.HeaderParams, "header", "req.Header")
}

func (g *ChiGenerator) queryParsing(w *writer.Writer, operation *spec.NamedOperation) {
	g.parametersParsing(w, operation, operation.QueryParams, "query", "req.URL.Query()")
}

func (g *ChiGenerator) urlParamsParsing(w *writer.Writer, operation *spec.NamedOperation) {
	if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
		w.Line(`urlParams := paramsparser.NewUrlParser(req.Context(), false)`)
		for _, param := range operation.Endpoint.UrlParams {
			w.Line(`%s := %s`, param.Name.CamelCase(), g.parserParameterCall(&param, "urlParams"))
		}
		w.Line(`if len(urlParams.Errors) > 0 {`)
		respondNotFound(w.Indented(), operation, g.Types, fmt.Sprintf(`"Failed to parse url parameters"`))
		w.Line(`}`)
	}
}

func (g *ChiGenerator) parametersParsing(w *writer.Writer, operation *spec.NamedOperation, namedParams []spec.NamedParam, paramsParserName string, paramsValuesVar string) {
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

func (g *ChiGenerator) serviceCallAndResponseCheck(w *writer.Writer, operation *spec.NamedOperation, responseVar string) {
	singleEmptyResponse := len(operation.Responses) == 1 && operation.Responses[0].Type.Definition.IsEmpty()
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

func (g *ChiGenerator) operation(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line(`log.WithFields(%s).Info("Received request")`, logFieldsName(operation))
	w.Line(`var err error`)
	g.urlParamsParsing(w, operation)
	g.headerParsing(w, operation)
	g.queryParsing(w, operation)
	g.bodyParsing(w, operation)
	g.serviceCallAndResponseCheck(w, operation, `response`)
	g.response(w, operation, `response`)
}

func (g *ChiGenerator) response(w *writer.Writer, operation *spec.NamedOperation, responseVar string) {
	if len(operation.Responses) == 1 {
		writeResponse(w, logFieldsName(operation), &operation.Responses[0].Response, responseVar)
	} else {
		for _, response := range operation.Responses {
			responseVar := fmt.Sprintf("%s.%s", responseVar, response.Name.PascalCase())
			w.Line(`if %s != nil {`, responseVar)
			writeResponse(w.Indented(), logFieldsName(operation), &response.Response, responseVar)
			w.Line(`  return`)
			w.Line(`}`)
		}
		respondInternalServerError(w, operation, g.Types, `"Result from service implementation does not have anything in it"`)
	}
}

func (g *ChiGenerator) bodyParsing(w *writer.Writer, operation *spec.NamedOperation) {
	if operation.BodyIs(spec.BodyString) {
		w.Line(`if !%s {`, callCheckContentType(logFieldsName(operation), `"text/plain"`, "req", "res"))
		w.Line(`  return`)
		w.Line(`}`)
		w.Line(`bodyData, err := ioutil.ReadAll(req.Body)`)
		w.Line(`if err != nil {`)
		respondBadRequest(w.Indented(), operation, g.Types, "body", genFmtSprintf(`Reading request body failed: %s`, `err.Error()`), "nil")
		w.Line(`}`)
		w.Line(`body := string(bodyData)`)
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`if !%s {`, callCheckContentType(logFieldsName(operation), `"application/json"`, "req", "res"))
		w.Line(`  return`)
		w.Line(`}`)
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
}

func (g *ChiGenerator) RootRouting(specification *spec.Spec) *generator.CodeFile {
	w := writer.New(g.Modules.Root, "spec.go")

	w.Imports.Add("github.com/go-chi/chi/v5")
	for _, version := range specification.Versions {
		w.Imports.ModuleAliased(g.Modules.Routing(&version).Aliased(routingPackageAlias(&version)))
		for _, api := range version.Http.Apis {
			w.Imports.ModuleAliased(g.Modules.ServicesApi(&api).Aliased(apiPackageAlias(&api)))
		}
	}
	w.EmptyLine()
	w.Line(`type Services struct {`)
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			apiModule := g.Modules.ServicesApi(&api).Aliased(apiPackageAlias(&api))
			w.LineAligned(`  %s %s`, serviceApiPublicNameVersioned(&api), apiModule.Get(serviceInterfaceName))
		}
	}
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`func AddRoutes(router *chi.Mux, services Services) {`)
	for _, version := range specification.Versions {
		routingModule := g.Modules.Routing(&version).Aliased(routingPackageAlias(&version))
		for _, api := range version.Http.Apis {
			w.Line(`  %s(router, services.%s)`, routingModule.Get(addRoutesMethodName(&api)), serviceApiPublicNameVersioned(&api))
		}
	}
	w.Line(`}`)

	return w.ToCodeFile()
}

func (g *ChiGenerator) GenerateUrlParamsCtor() *generator.CodeFile {
	w := writer.New(g.Modules.ParamsParser, `url_parser.go`)

	w.Lines(`
import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/url"
)

func NewUrlParser(context context.Context, parseCommaSeparatedArray bool) *ParamsParser {
	values := make(url.Values)
	urlParams := chi.RouteContext(context).URLParams
	for i := 0; i < len(urlParams.Keys); i++ {
		values.Add(urlParams.Keys[i], urlParams.Values[i])
	}
	return &ParamsParser{values, parseCommaSeparatedArray, []ParsingError{}}
}
`)

	return w.ToCodeFile()
}
