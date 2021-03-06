package genruby

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/vsapronov/casee"
	"github.com/vsapronov/gopoetry/ruby"
	"path/filepath"
	"specgen/gen"
	"specgen/static"
)

func GenerateClient(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil { return err }

	gemName := specification.ServiceName.SnakeCase()+"_client"
	moduleName := clientModuleName(specification.ServiceName)
	libGemPath := filepath.Join(generatePath, gemName)
	models := generateModels(specification.ResolvedModels, moduleName, filepath.Join(libGemPath, "models.rb"))
	clients := generateClientApisClasses(specification, libGemPath)

	data := static.RubyClient{GemName: gemName, ModuleName: moduleName}
	sources, err := static.RenderTemplate("ruby-client", generatePath, data)
	if err != nil { return err }

	sources = append(sources, *models, *clients)
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
	moduleName := clientModuleName(specification.ServiceName)
	module := ruby.Module(moduleName)

	for _, api := range specification.Apis {
		apiClass := generateClientApiClass(api)
		module.AddDeclarations(apiClass)
	}

	if len(specification.Apis) > 1 {
		module.AddDeclarations(generateClientSuperClass(specification))
	}

	unit := ruby.Unit()
	unit.Require("net/http")
	unit.Require("net/https")
	unit.Require("uri")
	unit.Require("emery")
	unit.AddDeclarations(module)

	return &gen.TextFile{
		Path:    filepath.Join(generatePath, "client.rb"),
		Content: unit.Code(),
	}
}

func generateClientOperation(operation spec.NamedOperation) *ruby.MethodDeclaration {
	method := ruby.Method(operation.Name.SnakeCase())
	methodBody := method.Body()

	addParams(method, operation.HeaderParams)
	if operation.Body != nil {
		method.KeywordArg("body")
	}
	addParams(method, operation.Endpoint.UrlParams)
	addParams(method, operation.QueryParams)

	httpMethod := casee.ToPascalCase(operation.Endpoint.Method)

	addParamsWriting(methodBody, operation.Endpoint.UrlParams, "url_params")
	addParamsWriting(methodBody, operation.QueryParams, "query")
	addParamsWriting(methodBody, operation.HeaderParams, "header")

	url_compose := "url = @base_uri"
	if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
		url_compose = url_compose + fmt.Sprintf(" + url_params.set_to_url('%s')", operation.Endpoint.Url)
	} else {
		url_compose = url_compose + fmt.Sprintf(" + '%s'", operation.Endpoint.Url)
	}
	if operation.QueryParams != nil && len(operation.QueryParams) > 0 {
		url_compose = url_compose + " + query.query_str"
	}
	methodBody.AddLn(url_compose)

	methodBody.AddLn(fmt.Sprintf("request = Net::HTTP::%s.new(url)", httpMethod))

	if operation.HeaderParams != nil && len(operation.HeaderParams) > 0 {
		methodBody.AddLn("header.params.each { |name, value| request.add_field(name, value) }")
	}

	if operation.Body != nil {
		boydRubyType := RubyType(&operation.Body.Type.Definition)
		methodBody.AddLn(fmt.Sprintf("body_json = Jsoner.to_json(%s, T.check_var('body', %s, body))", boydRubyType, boydRubyType))
		methodBody.AddLn("request.body = body_json")
	}
	methodBody.AddLn("response = @client.request(request)")
	methodBody.AddLn("case response.code")

	for _, response := range operation.Responses {
		methodBody.AddLn(fmt.Sprintf("when '%s'", spec.HttpStatusCode(response.Name)))
		if response.Type.Definition.IsEmpty() {
			methodBody.Scope().Add("nil")
		} else {
			methodBody.Scope().Add(fmt.Sprintf("Jsoner.from_json(%s, response.body)", RubyType(&response.Type.Definition)))
		}
	}

	methodBody.AddLn("else")
	methodBody.Scope().Add("raise StandardError.new(\"Unexpected HTTP response code #{response.code}\")")
	methodBody.AddLn("end")
	return method
}

func generateClientApiClass(api spec.Api) *ruby.ClassDeclaration {
	apiClassName := clientClassName(api.Name)
	apiClass := ruby.Class(apiClassName).Inherits("BaseClient")
	for _, operation := range api.Operations {
		method := generateClientOperation(operation)
		apiClass.AddMembers(method)
	}
	return apiClass
}

func generateClientSuperClass(specification *spec.Spec) *ruby.ClassDeclaration {
	class := ruby.Class(specification.ServiceName.PascalCase() + "Client")
	for _, api := range specification.Apis {
		class.AddMembers(ruby.Code(fmt.Sprintf("attr_reader :%s_client", api.Name.SnakeCase())))
	}
	init := class.Initialize()
	init.Arg("base_url")
	initBody := init.Body()
	for _, api := range specification.Apis {
		initBody.AddLn(fmt.Sprintf("@%s_client = %s.new(base_url)", api.Name.SnakeCase(), clientClassName(api.Name)))
	}
	return class
}

func addParams(method *ruby.MethodDeclaration, params []spec.NamedParam) {
	for _, param := range params {
		arg := method.KeywordArg(param.Name.SnakeCase())
		if param.Default != nil {
			arg.Default(ruby.Code(DefaultValue(&param.Type.Definition, *param.Default)))
		}
	}
}

func addParamsWriting(code *ruby.StatementsDeclaration, params []spec.NamedParam, paramsName string) {
	if params != nil && len(params) > 0 {
		code.AddLn(paramsName + " = StringParams.new")
		for _, p := range params {
			code.AddLn(fmt.Sprintf("%s.set('%s', %s, %s)", paramsName, p.Name.Source, RubyType(&p.Type.Definition), p.Name.SnakeCase()))
		}
	}
}
