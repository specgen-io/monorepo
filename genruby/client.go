package genruby

import (
	"fmt"
	"github.com/pinzolo/casee"
	spec "github.com/specgen-io/spec.v2"
	"github.com/vsapronov/gopoetry/ruby"
	"path/filepath"
	"specgen/gen"
)

func GenerateClient(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil { return err }

	gemName := specification.Name.SnakeCase()+"_client"
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
	rootModule := ruby.Module(moduleName)

	for _, version := range specification.Versions {
		module := rootModule
		if version.Version.Source != "" {
			module = ruby.Module(version.Version.PascalCase())
			rootModule.AddDeclarations(module)
		}
		for _, api := range version.Http.Apis {
			apiClass := generateClientApiClass(api)
			module.AddDeclarations(apiClass)
		}
	}

	unit := ruby.Unit()
	unit.Require("net/http")
	unit.Require("net/https")
	unit.Require("uri")
	unit.Require("emery")
	unit.AddDeclarations(rootModule)

	return &gen.TextFile{
		Path:    filepath.Join(generatePath, "client.rb"),
		Content: unit.Code(),
	}
}

func operationResult(operation *spec.NamedOperation, response *spec.NamedResponse) *ruby.WritableCode {
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
	return ruby.Code(fmt.Sprintf("OpenStruct.new(:%s => %s%s)", response.Name.Source, body, flags))
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
		url_compose = url_compose + fmt.Sprintf(" + url_params.set_to_url('%s')", operation.FullUrl())
	} else {
		url_compose = url_compose + fmt.Sprintf(" + '%s'", operation.FullUrl())
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
		methodBody.Scope().AddCode(operationResult(&operation, &response))
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
