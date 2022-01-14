package gents

import (
	"fmt"
	"github.com/pinzolo/casee"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

var koa = "koa"

func generateKoaSpecRouter(specification *spec.Spec, rootModule module, module module) *sources.CodeFile {
	w := NewTsWriter()
	w.Line("import Router from '@koa/router'")
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
	w.Line("  let router = new Router()")
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			apiRouterNameConst := "the" + casee.ToPascalCase(apiVersionedRouterName(&api))
			w.Line("  const %s = %s(%s)", apiRouterNameConst, apiRouterNameVersioned(&api), apiServiceParamName(&api))
			if version.Http.GetUrl() == "" {
				w.Line("  router.use(%s.routes(), %s.allowedMethods())", apiRouterNameConst, apiRouterNameConst)
			} else {
				w.Line("  router.use('%s', %s.routes(), %s.allowedMethods())", version.Http.GetUrl(), apiRouterNameConst, apiRouterNameConst)
			}
		}
	}
	w.Line("  return router")
	w.Line("}")

	return &sources.CodeFile{module.GetPath(), w.String()}
}

func generateKoaVersionRouting(version *spec.Version, validation string, validationModule, paramsModule, module module) *sources.CodeFile {
	w := NewTsWriter()

	w.Line(`import Router from '@koa/router'`)
	w.Line(`import {zipHeaders} from '%s'`, paramsModule.GetImport(module))
	w.Line(`import * as t from '%s'`, validationModule.GetImport(module))
	w.Line(`import * as models from './models'`)
	for _, api := range version.Http.Apis {
		w.Line("import {%s} from './%s'", serviceInterfaceName(&api), serviceName(&api))
	}

	for _, api := range version.Http.Apis {
		for _, operation := range api.Operations {
			generateParams(w, paramsTypeName(&operation, "HeaderParams"), operation.HeaderParams, validation)
			generateParams(w, paramsTypeName(&operation, "UrlParams"), operation.Endpoint.UrlParams, validation)
			generateParams(w, paramsTypeName(&operation, "QueryParams"), operation.QueryParams, validation)
		}

		generateKoaApiRouting(w, &api, validation)
	}

	return &sources.CodeFile{module.GetPath(), w.String()}
}

func generateKoaApiRouting(w *sources.Writer, api *spec.Api, validation string) {
	w.EmptyLine()
	w.Line("export let %s = (service: %s) => {", apiRouterName(api), serviceInterfaceName(api))
	w.Line("  let router = new Router()")
	for _, operation := range api.Operations {
		w.EmptyLine()
		generateKoaOperationRouting(w.Indented(), &operation, validation)
	}
	w.EmptyLine()
	w.Line("  return router")
	w.Line("}")
}

func getKoaUrl(endpoint spec.Endpoint) string {
	url := endpoint.Url
	for _, param := range endpoint.UrlParams {
		url = strings.Replace(url, spec.UrlParamStr(&param), ":"+param.Name.Source, -1)
	}
	return url
}

func generateKoaResponse(w *sources.Writer, response *spec.NamedResponse, validation string, dataParam string) {
	w.Line("ctx.status = %s", spec.HttpStatusCode(response.Name))
	if !response.Type.Definition.IsEmpty() {
		if response.Type.Definition.Plain == spec.TypeString {
			w.Line("ctx.body = %s", dataParam)
		} else {
			w.Line("ctx.body = t.encode(%s.%s, %s)", modelsPackage, runtimeType(validation, &response.Type.Definition), dataParam)
		}
	}
	w.Line("return")
}

func generateKoaOperationRouting(w *sources.Writer, operation *spec.NamedOperation, validation string) {
	w.Line("router.%s('%s', async (ctx) => {", strings.ToLower(operation.Endpoint.Method), getKoaUrl(operation.Endpoint))
	w.Indent()

	if operation.Body != nil {
		if operation.Body.Type.Definition.Plain == spec.TypeString {
			w.Line(`if (ctx.request.type == 'text/plain') {`)
			w.Line(`  ctx.throw(400)`)
			w.Line(`  return`)
			w.Line(`}`)
		} else {
			w.Line(`if (ctx.request.type == 'application/json') {`)
			w.Line(`  ctx.throw(400)`)
			w.Line(`  return`)
			w.Line(`}`)
		}
	}

	generateParametersParsing(w, operation, "zipHeaders(ctx.req.rawHeaders)", "ctx.params", "ctx.request.query", "ctx.throw(400)")
	generateBodyParsing(w, validation, operation, "ctx.request.body", "ctx.request.rawBody", "ctx.throw(400)")

	w.Line("try {")
	w.Line("  %s", serviceCall(operation, getApiCallParamsObject(operation)))
	if len(operation.Responses) == 1 {
		generateKoaResponse(w.IndentedWith(1), &operation.Responses[0], validation, "result")
	} else {
		w.Line("  switch (result.status) {")
		for _, response := range operation.Responses {
			w.Line("    case '%s':", response.Name.FlatCase())
			generateKoaResponse(w.IndentedWith(3), &response, validation, "result.data")
		}
		w.Line("  }")
	}
	w.Line("} catch (error) {")
	w.Line("  ctx.throw(500)")
	w.Line("}")
	w.Unindent()
	w.Line("})")
}
