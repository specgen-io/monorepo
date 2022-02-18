package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen/ts/modules"
	"github.com/specgen-io/specgen/v2/gen/ts/types"
	"github.com/specgen-io/specgen/v2/gen/ts/validations"
	"github.com/specgen-io/specgen/v2/gen/ts/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

type expressGenerator struct {
	validation validations.Validation
}

func (g *expressGenerator) SpecRouter(specification *spec.Spec, rootModule modules.Module, module modules.Module) *sources.CodeFile {
	w := writer.NewTsWriter()
	w.Line("import {Router} from 'express'")
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			versionModule := rootModule.Submodule(version.Version.FlatCase())
			apiModule := versionModule.Submodule(serviceName(&api)) //TODO: This logic is repeated here, it also exists where api module is created
			routerModule := versionModule.Submodule("routing")      //TODO: This logic is repeated here, it also exists where router module is created
			w.Line("import {%s as %s} from '%s'", serviceInterfaceName(&api), serviceInterfaceNameVersioned(&api), apiModule.GetImport(module))
			w.Line("import {%s as %s} from '%s'", apiRouterName(&api), apiRouterNameVersioned(&api), routerModule.GetImport(module))
		}
	}

	routerParams := []string{}
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			routerParams = append(routerParams, fmt.Sprintf("%s: %s", apiServiceParamName(&api), serviceInterfaceNameVersioned(&api)))
		}
	}

	w.EmptyLine()
	w.Line("export let specRouter = (%s) => {", strings.Join(routerParams, ", "))
	w.Line("  let router = Router()")
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			w.Line("  router.use('%s', %s(%s))", expressVersionUrl(&version), apiRouterNameVersioned(&api), apiServiceParamName(&api))
		}
	}
	w.Line("  return router")
	w.Line("}")

	return &sources.CodeFile{module.GetPath(), w.String()}
}

func expressVersionUrl(version *spec.Version) string {
	url := version.Http.GetUrl()
	if url == "" {
		return "/"
	}
	return url
}

func (g *expressGenerator) VersionRouting(version *spec.Version, validationModule, paramsModule, module modules.Module) *sources.CodeFile {
	w := writer.NewTsWriter()

	w.Line("import {Router} from 'express'")
	w.Line("import {Request, Response} from 'express'") //TODO: Join with above
	w.Line(`import {zipHeaders} from '%s'`, paramsModule.GetImport(module))
	w.Line(`import * as t from '%s'`, validationModule.GetImport(module))
	w.Line("import * as models from './models'")

	for _, api := range version.Http.Apis {
		w.Line("import {%s} from './%s'", serviceInterfaceName(&api), serviceName(&api))
	}

	for _, api := range version.Http.Apis {
		for _, operation := range api.Operations {
			g.validation.WriteParamsType(w, paramsTypeName(&operation, "HeaderParams"), operation.HeaderParams)
			g.validation.WriteParamsType(w, paramsTypeName(&operation, "UrlParams"), operation.Endpoint.UrlParams)
			g.validation.WriteParamsType(w, paramsTypeName(&operation, "QueryParams"), operation.QueryParams)
		}

		g.apiRouting(w, &api)
	}

	return &sources.CodeFile{module.GetPath(), w.String()}
}

func (g *expressGenerator) apiRouting(w *sources.Writer, api *spec.Api) {
	w.EmptyLine()
	w.Line("export let %s = (service: %s) => {", apiRouterName(api), serviceInterfaceName(api))
	w.Line("  let router = Router()")
	for _, operation := range api.Operations {
		w.EmptyLine()
		g.operationRouting(w.Indented(), &operation)
	}
	w.EmptyLine()
	w.Line("  return router")
	w.Line("}")
}

func getExpressUrl(endpoint spec.Endpoint) string {
	url := endpoint.Url
	for _, param := range endpoint.UrlParams {
		url = strings.Replace(url, spec.UrlParamStr(&param), ":"+param.Name.Source, -1)
	}
	return url
}

func (g *expressGenerator) response(w *sources.Writer, response *spec.NamedResponse, dataParam string) {
	if response.BodyIs(spec.BodyEmpty) {
		w.Line("response.status(%s).send()", spec.HttpStatusCode(response.Name))
		w.Line("return")
	}
	if response.BodyIs(spec.BodyString) {
		w.Line("response.status(%s).type('text').send(%s)", spec.HttpStatusCode(response.Name), dataParam)
		w.Line("return")
	}
	if response.BodyIs(spec.BodyJson) {
		w.Line("response.status(%s).type('json').send(JSON.stringify(t.encode(%s, %s)))", spec.HttpStatusCode(response.Name), g.validation.RuntimeTypeFromPackage(types.ModelsPackage, &response.Type.Definition), dataParam)
		w.Line("return")
	}
}

func (g *expressGenerator) operationRouting(w *sources.Writer, operation *spec.NamedOperation) {
	w.Line("router.%s('%s', async (request: Request, response: Response) => {", strings.ToLower(operation.Endpoint.Method), getExpressUrl(operation.Endpoint))
	w.Indent()

	if operation.BodyIs(spec.BodyString) {
		w.Line(`if (!request.is('text/plain')) {`)
		w.Line(`  response.status(400).send()`)
		w.Line(`  return`)
		w.Line(`}`)
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`if (!request.is('application/json')) {`)
		w.Line(`  response.status(400).send()`)
		w.Line(`  return`)
		w.Line(`}`)
	}

	generateParametersParsing(w, operation, "zipHeaders(request.rawHeaders)", "request.params", "request.query", "response.status(400).send()")
	generateBodyParsing(w, g.validation, operation, "request.body", "request.body", "response.status(400).send()")

	w.Line("try {")
	w.Line("  %s", serviceCall(operation, getApiCallParamsObject(operation)))
	if len(operation.Responses) == 1 {
		g.response(w.IndentedWith(1), &operation.Responses[0], "result")
	} else {
		w.Line("  switch (result.status) {")
		for _, response := range operation.Responses {
			w.Line("    case '%s':", response.Name.SnakeCase())
			g.response(w.IndentedWith(3), &response, "result.data")
		}
		w.Line("  }")
	}
	w.Line("} catch (error) {")
	w.Line("  response.status(500).send()")
	w.Line("}")
	w.Unindent()
	w.Line("})")
}
