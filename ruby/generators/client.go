package generators

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pinzolo/casee"
	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/spec/v2"
)

func GenerateClient(specification *spec.Spec, generatePath string) *generator.Sources {
	sources := generator.NewSources()
	gemName := specification.Name.SnakeCase() + "_client"
	moduleName := clientModuleName(specification.Name)
	libGemPath := filepath.Join(generatePath, gemName)
	models := generateModels(specification, moduleName, filepath.Join(libGemPath, "models.rb"))
	clients := generateClientApisClasses(specification, libGemPath)
	baseclient := generateBaseClient(moduleName, filepath.Join(libGemPath, "baseclient.rb"))
	clientroot := generateClientRoot(gemName, filepath.Join(generatePath, gemName+".rb"))

	sources.AddGenerated(clientroot, baseclient, models, clients)
	return sources
}

func clientClassName(apiName spec.Name) string {
	return apiName.PascalCase() + "Client"
}

func clientModuleName(serviceName spec.Name) string {
	return serviceName.PascalCase()
}

func generateClientApisClasses(specification *spec.Spec, generatePath string) *generator.CodeFile {
	moduleName := clientModuleName(specification.Name)

	w := NewRubyWriter()
	w.Line(`require "net/http"`)
	w.Line(`require "net/https"`)
	w.Line(`require "uri"`)
	w.Line(`require "emery"`)

	for _, version := range specification.Versions {
		if version.Version.Source == "" {
			w.EmptyLine()
			generateVersionClientModule(w, &version, moduleName)
		}
	}

	for _, version := range specification.Versions {
		if version.Version.Source != "" {
			w.EmptyLine()
			generateVersionClientModule(w, &version, moduleName)
		}
	}

	return &generator.CodeFile{Path: filepath.Join(generatePath, "client.rb"), Content: w.String()}
}

func generateVersionClientModule(w *generator.Writer, version *spec.Version, moduleName string) {
	w.Line("module %s", versionedModule(moduleName, version.Version))
	for index, api := range version.Http.Apis {
		if index != 0 {
			w.EmptyLine()
		}
		generateClientApiClass(w.Indented(), moduleName, &api)
	}
	w.Line("end")
}

func operationResult(operation *spec.NamedOperation, response *spec.Response) string {
	flags := ""
	for _, r := range operation.Responses {
		if r.Name.Source == response.Name.Source {
			flags += fmt.Sprintf(", :%s? => true", r.Name.Source)
		} else {
			flags += fmt.Sprintf(", :%s? => false", r.Name.Source)
		}
	}

	if response.BodyIs(spec.BodyString) {
		body := "response.body"
		return fmt.Sprintf("OpenStruct.new(:%s => %s%s)", response.Name.Source, body, flags)
	} else if response.BodyIs(spec.BodyJson) {
		body := fmt.Sprintf("Jsoner.from_json(%s, response.body)", RubyType(&response.Type.Definition))
		return fmt.Sprintf("OpenStruct.new(:%s => %s%s)", response.Name.Source, body, flags)
	} else {
		body := "nil"
		return fmt.Sprintf("OpenStruct.new(:%s => %s%s)", response.Name.Source, body, flags)
	}
}

func generateClientOperation(w *generator.Writer, moduleName string, operation *spec.NamedOperation) {
	args := []string{}
	args = append(args, addParams(operation.HeaderParams)...)
	if operation.Body != nil {
		args = append(args, "body:")
	}
	args = append(args, addParams(operation.Endpoint.UrlParams)...)
	args = append(args, addParams(operation.QueryParams)...)

	w.Line("def %s(%s)", operation.Name.SnakeCase(), strings.Join(args, ", "))
	w.Indent()

	httpMethod := casee.ToPascalCase(operation.Endpoint.Method)

	addParamsWriting(w, moduleName, operation.QueryParams, "query")
	addParamsWriting(w, moduleName, operation.HeaderParams, "header")

	url_compose := "url = @base_uri"
	if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
		w.Line("url_params = {")
		for _, p := range operation.Endpoint.UrlParams {
			w.Line("  '%s' => T.check(%s, %s),", p.Name.Source, RubyType(&p.Type.Definition), p.Name.SnakeCase())
		}
		w.Line("}")
		url_compose = url_compose + fmt.Sprintf(" + Stringify::set_params_to_url('%s', url_params)", operation.FullUrl())
	} else {
		url_compose = url_compose + fmt.Sprintf(" + '%s'", operation.FullUrl())
	}
	if operation.QueryParams != nil && len(operation.QueryParams) > 0 {
		url_compose = url_compose + " + query.query_str"
	}
	w.Line(url_compose)

	w.Line(fmt.Sprintf("request = Net::HTTP::%s.new(url)", httpMethod))

	if operation.HeaderParams != nil && len(operation.HeaderParams) > 0 {
		w.Line("header.params.each { |name, value| request.add_field(name, value) }")
	}

	if operation.BodyIs(spec.BodyString) {
		w.Line(`request.add_field('Content-Type', 'text/plain')`)
		w.Line("request.body = T.check_var('body', String, body)")
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`request.add_field('Content-Type', 'application/json')`)
		bodyRubyType := RubyType(&operation.Body.Type.Definition)
		w.Line(fmt.Sprintf("body_json = Jsoner.to_json(%s, T.check_var('body', %s, body))", bodyRubyType, bodyRubyType))
		w.Line("request.body = body_json")
	}
	w.Line("response = @client.request(request)")
	w.Line("case response.code")

	for _, response := range operation.Responses {
		w.Line("when '%s'", spec.HttpStatusCode(response.Name))
		w.Line("  %s", operationResult(operation, &response.Response))
	}

	w.Line("else")
	w.Line(`  raise StandardError.new("Unexpected HTTP response code #{response.code}")`)
	w.Line("end")
	w.Unindent()
	w.Line("end")
}

func generateClientApiClass(w *generator.Writer, moduleName string, api *spec.Api) {
	w.Line("class %s < %s::BaseClient", clientClassName(api.Name), moduleName)
	for index, operation := range api.Operations {
		if index != 0 {
			w.EmptyLine()
		}
		generateClientOperation(w.Indented(), moduleName, &operation)
	}
	w.Line("end")
}

func addParams(params []spec.NamedParam) []string {
	args := []string{}
	for _, param := range params {
		arg := param.Name.SnakeCase() + ":"
		if param.Default != nil {
			arg += " " + DefaultValue(&param.Type.Definition, *param.Default)
		}
		args = append(args, arg)
	}
	return args
}

func addParamsWriting(w *generator.Writer, moduleName string, params []spec.NamedParam, paramsName string) {
	if params != nil && len(params) > 0 {
		w.Line("%s = %s::StringParams.new", paramsName, moduleName)
		for _, p := range params {
			w.Line("%s.set('%s', %s, %s)", paramsName, p.Name.Source, RubyType(&p.Type.Definition), p.Name.SnakeCase())
		}
	}
}
