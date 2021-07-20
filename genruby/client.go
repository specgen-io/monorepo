package genruby

import (
	"fmt"
	"github.com/pinzolo/casee"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
	"strings"
)

func GenerateClient(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}

	gemName := specification.Name.SnakeCase() + "_client"
	moduleName := clientModuleName(specification.Name)
	libGemPath := filepath.Join(generatePath, gemName)
	models := generateModels(specification, moduleName, filepath.Join(libGemPath, "models.rb"))
	clients := generateClientApisClasses(specification, libGemPath)
	baseclient := generateBaseClient(moduleName, filepath.Join(libGemPath, "baseclient.rb"))
	clientroot := generateClientRoot(gemName, filepath.Join(generatePath, gemName+".rb"))

	sources := []gen.TextFile{*clientroot, *baseclient, *models, *clients}
	err = gen.WriteFiles(sources, true)
	return err
}

func clientClassName(apiName spec.Name) string {
	return apiName.PascalCase() + "Client"
}

func clientModuleName(serviceName spec.Name) string {
	return serviceName.PascalCase()
}

func generateClientApisClasses(specification *spec.Spec, generatePath string) *gen.TextFile {
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

	return &gen.TextFile{Path: filepath.Join(generatePath, "client.rb"), Content: w.String()}
}

func generateVersionClientModule(w *gen.Writer, version *spec.Version, moduleName string) {
	w.Line("module %s", versionedModule(moduleName, version.Version))
	for index, api := range version.Http.Apis {
		if index != 0 { w.EmptyLine() }
		generateClientApiClass(w.Indented(), api)
	}
	w.Line("end")
}

func operationResult(operation *spec.NamedOperation, response *spec.NamedResponse) string {
	body := "nil"
	if !response.Type.Definition.IsEmpty() {
		body = fmt.Sprintf("Jsoner.from_json(%s, response.body)", RubyType(&response.Type.Definition))
	}
	flags := ""
	for _, r := range operation.Responses {
		if r.Name.Source == response.Name.Source {
			flags += fmt.Sprintf(", :%s? => true", r.Name.Source)
		} else {
			flags += fmt.Sprintf(", :%s? => false", r.Name.Source)
		}
	}
	return fmt.Sprintf("OpenStruct.new(:%s => %s%s)", response.Name.Source, body, flags)
}

func generateClientOperation(w *gen.Writer, operation spec.NamedOperation) {
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

	addParamsWriting(w, operation.Endpoint.UrlParams, "url_params")
	addParamsWriting(w, operation.QueryParams, "query")
	addParamsWriting(w, operation.HeaderParams, "header")

	url_compose := "url = @base_uri"
	if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
		url_compose = url_compose + fmt.Sprintf(" + url_params.set_to_url('%s')", operation.FullUrl())
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

	if operation.Body != nil {
		bodyRubyType := RubyType(&operation.Body.Type.Definition)
		w.Line(fmt.Sprintf("body_json = Jsoner.to_json(%s, T.check_var('body', %s, body))", bodyRubyType, bodyRubyType))
		w.Line("request.body = body_json")
	}
	w.Line("response = @client.request(request)")
	w.Line("case response.code")

	for _, response := range operation.Responses {
		w.Line("when '%s'", spec.HttpStatusCode(response.Name))
		w.Line("  %s", operationResult(&operation, &response))
	}

	w.Line("else")
	w.Line(`  raise StandardError.new("Unexpected HTTP response code #{response.code}")`)
	w.Line("end")
	w.Unindent()
	w.Line("end")
}

func generateClientApiClass(w *gen.Writer, api spec.Api) {
	w.Line("class %s < Http::BaseClient", clientClassName(api.Name))
	for index, operation := range api.Operations {
		if index != 0 { w.EmptyLine() }
		generateClientOperation(w.Indented(), operation)
	}
	w.Line("end")
}

func addParams(params []spec.NamedParam) []string {
	args := []string{}
	for _, param := range params {
		arg := param.Name.SnakeCase()+":"
		if param.Default != nil {
			arg += " "+DefaultValue(&param.Type.Definition, *param.Default)
		}
		args = append(args, arg)
	}
	return args
}

func addParamsWriting(w *gen.Writer, params []spec.NamedParam, paramsName string) {
	if params != nil && len(params) > 0 {
		w.Line("%s = Http::StringParams.new", paramsName)
		for _, p := range params {
			w.Line("%s.set('%s', %s, %s)", paramsName, p.Name.Source, RubyType(&p.Type.Definition), p.Name.SnakeCase())
		}
	}
}
