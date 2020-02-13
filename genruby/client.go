package genruby

import (
	"fmt"
	"github.com/ModaOperandi/spec"
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
	moduleName := specification.ServiceName.PascalCase()
	libGemPath := filepath.Join(generatePath, "lib", gemName)
	models := GenerateModels(specification, libGemPath)
	clients := generateClientApisClasses(specification, libGemPath)

	data := static.RubyGem{GemName: gemName, ModuleName: moduleName}
	sources, err := static.RenderTemplate("genruby", generatePath, data)
	if err != nil { return err }

	sources = append(sources, *models, *clients)
	err = gen.WriteFiles(sources, true)
	return err
}

func clientClassName(apiName spec.Name) string {
	return apiName.PascalCase() + "Client"
}

func generateClientApisClasses(specification *spec.Spec, generatePath string) *gen.TextFile {
	gemName := specification.ServiceName.SnakeCase()+"_client"
	moduleName := specification.ServiceName.PascalCase()

	module := ruby.Module(moduleName)
	for _, api := range specification.Apis {
		apiClass := generateClientApiClass(api)
		module.AddDeclarations(apiClass)
	}

	unit := ruby.Unit()
	unit.Require("net/http")
	unit.Require("net/https")
	unit.Require("uri")
	unit.Require(gemName+"/type")
	unit.Require(gemName+"/jsoner")
	unit.Require(gemName+"/models")
	unit.Require(gemName+"/baseclient")
	unit.AddDeclarations(module)

	return &gen.TextFile{
		Path:    filepath.Join(generatePath, "client.rb"),
		Content: unit.Code(),
	}
}

func generateClientApiClass(api spec.Api) *ruby.ClassDeclaration {
	apiClassName := clientClassName(api.Name)
	apiClass := ruby.Class(apiClassName).Inherits("BaseClient")
	for _, operation := range api.Operations {
		method := apiClass.Def(operation.Name.SnakeCase())
		methodBody := method.Body()

		addParams(method, operation.HeaderParams)
		if operation.Body != nil {
			method.KeywordArg("body")
		}
		for _, param := range operation.Endpoint.UrlParams {
			method.KeywordArg(param.Name.SnakeCase())
		}
		addParams(method, operation.QueryParams)

		addParamsTypeCheck(methodBody, operation.HeaderParams)
		addParamsTypeCheck(methodBody, operation.Endpoint.UrlParams)
		addParamsTypeCheck(methodBody, operation.QueryParams)
		if operation.Body != nil {
			methodBody.AddLn(fmt.Sprintf("T.check_var('body', %s, body)", RubyType(&operation.Body.Type.Definition)))
		}

		httpMethod := casee.ToPascalCase(operation.Endpoint.Method)

		if operation.QueryParams != nil && len(operation.QueryParams) > 0 {
			methodBody.AddLn("query = StringParams.new")
			for _, p := range operation.QueryParams {
				methodBody.AddLn(fmt.Sprintf("query['%s'] = %s", p.Name.SnakeCase(), p.Name.SnakeCase()))
			}
		}

		if operation.QueryParams != nil && len(operation.QueryParams) > 0 {
			methodBody.AddLn(fmt.Sprintf("url = @base_uri + '%s' + query.query_str", operation.Endpoint.Url))
		} else {
			methodBody.AddLn(fmt.Sprintf("url = @base_uri + '%s'", operation.Endpoint.Url))
		}

		methodBody.AddLn(fmt.Sprintf("request = Net::HTTP::%s.new(url)", httpMethod))
		if operation.Body != nil {
			methodBody.AddLn("body_json = Jsoner.to_json(Message, body)")
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
		methodBody.Scope().Add("raise StandardError.new('Unexpected HTTP response code')")
		methodBody.AddLn("end")
	}
	return apiClass
}

func addParams(method *ruby.MethodDeclaration, params []spec.NamedParam) {
	for _, param := range params {
		arg := method.KeywordArg(param.Name.SnakeCase())
		if param.Default != nil {
			arg.Default(ruby.Code(DefaultValue(&param.Type.Definition, *param.Default)))
		}
	}
}

func addParamsTypeCheck(methodBody *ruby.StatementsDeclaration, params []spec.NamedParam) {
	for _, param := range params {
		methodBody.AddLn(fmt.Sprintf("T.check_var('%s', %s, %s)", param.Name.SnakeCase(), RubyType(&param.Type.Definition), param.Name.SnakeCase()))
	}
}