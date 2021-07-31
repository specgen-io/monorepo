package gents

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
	"strings"
)

var koa = "koa"

func generateKoaSpecRouter(specification *spec.Spec, generatePath string) *gen.TextFile {
	w := NewTsWriter()
	w.Line("import Router from '@koa/router'")
	for _, version := range specification.Versions {
		w.Line("import * as %s from './%s'", versionModule(&version, "services"), versionFilename(&version, "services", ""))
		w.Line("import * as %s from './%s'", versionModule(&version, "routing"), versionFilename(&version, "routing", ""))
	}

	routerParams := []string{}
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			routerParams = append(routerParams, fmt.Sprintf("%s: %s.%s", apiServiceParamName(&api), versionModule(&version, "services"), serviceInterfaceName(&api)))
		}
	}

	w.EmptyLine()
	w.Line("export let specRouter = (%s) => {", strings.Join(routerParams, ", "))
	w.Line("  let router = new Router()")
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			apiRouterNameConst := apiVersionedRouterName(&api)
			w.Line("  const %s = %s.%s(%s)", apiRouterNameConst, versionModule(&version, "routing"), apiRouterName(&api), apiServiceParamName(&api))
			if version.Http.GetUrl() == "" {
				w.Line("  router.use(%s.routes(), %s.allowedMethods())", apiRouterNameConst, apiRouterNameConst)
			} else {
				w.Line("  router.use('%s', %s.routes(), %s.allowedMethods())", version.Http.GetUrl(), apiRouterNameConst, apiRouterNameConst)
			}
		}
	}
	w.Line("  return router")
	w.Line("}")

	return &gen.TextFile{filepath.Join(generatePath, "spec_router.ts"), w.String()}
}

func generateKoaVersionRouting(version *spec.Version, validation string, generatePath string) *gen.TextFile {
	w := NewTsWriter()

	w.Line("import Router from '@koa/router'")
	w.Line(importEncoding(validation))
	w.Line("import * as %s from './%s'", modelsPackage, versionFilename(version, "models", ""))
	w.Line("import * as services from './%s'", versionFilename(version, "services", ""))

	for _, api := range version.Http.Apis {
		for _, operation := range api.Operations {
			generateParams(w, paramsTypeName(&operation, "HeaderParams"), true, operation.HeaderParams, validation)
			generateParams(w, paramsTypeName(&operation, "UrlParams"), false, operation.Endpoint.UrlParams, validation)
			generateParams(w, paramsTypeName(&operation, "QueryParams"), false, operation.QueryParams, validation)
		}

		generateKoaApiRouting(w, &api, validation)
	}

	filename := versionFilename(version, "routing", "ts")
	return &gen.TextFile{filepath.Join(generatePath, filename), w.String()}
}

func generateKoaApiRouting(w *gen.Writer, api *spec.Api, validation string) {
	w.EmptyLine()
	w.Line("export let %s = (service: services.%s) => {", apiRouterName(api), serviceInterfaceName(api))
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
		url = strings.Replace(url, spec.UrlParamStr(param.Name.Source), ":"+param.Name.Source, -1)
	}
	return url
}

func generateKoaOperationRouting(w *gen.Writer, operation *spec.NamedOperation, validation string) {
	w.Line("router.%s('%s', async (ctx) => {", strings.ToLower(operation.Endpoint.Method), getKoaUrl(operation.Endpoint))
	w.Indent()
	apiCallParams := []string{}

	if operation.Body != nil || len(operation.Endpoint.UrlParams) > 0 || len(operation.HeaderParams) > 0 || len(operation.QueryParams) > 0 {
		if operation.Body != nil {
			w.Line("var body: %s", TsType(&operation.Body.Type.Definition))
		}
		if len(operation.Endpoint.UrlParams) > 0 {
			w.Line("var urlParams: %s", paramsTypeName(operation, "UrlParams"))
		}
		if len(operation.HeaderParams) > 0 {
			w.Line("var headerParams: %s", paramsTypeName(operation, "HeaderParams"))
		}
		if len(operation.QueryParams) > 0 {
			w.Line("var queryParams: %s", paramsTypeName(operation, "QueryParams"))
		}
		w.Line("try {")
		if operation.Body != nil {
			w.Line("  body = t.decode(%s.%s, ctx.request.body)", modelsPackage, runtimeType(validation, &operation.Body.Type.Definition))
			apiCallParams = append(apiCallParams, "body")
		}
		if len(operation.Endpoint.UrlParams) > 0 {
			w.Line("  urlParams = t.decode(%s, ctx.params)", paramsRuntimeTypeName(paramsTypeName(operation, "UrlParams")))
			apiCallParams = append(apiCallParams, "...urlParams")
		}
		if len(operation.HeaderParams) > 0 {
			w.Line("  headerParams = t.decode(%s, ctx.request.headers)", paramsRuntimeTypeName(paramsTypeName(operation, "HeaderParams")))
			apiCallParams = append(apiCallParams, "...headerParams")
		}
		if len(operation.QueryParams) > 0 {
			w.Line("  queryParams = t.decode(%s, ctx.request.query)", paramsRuntimeTypeName(paramsTypeName(operation, "QueryParams")))
			apiCallParams = append(apiCallParams, "...queryParams")
		}
		w.Line("} catch (error) {")
		w.Line("  ctx.throw(400, error)")
		w.Line("  return")
		w.Line("}")
	}
	w.Line("try {")
	w.Line("  let result = await service.%s({%s})", operation.Name.CamelCase(), strings.Join(apiCallParams, ", "))
	w.Line("  switch (result.status) {")
	for _, response := range operation.Responses {
		w.Line("    case '%s':", response.Name.FlatCase())
		w.Line("      ctx.status = %s", spec.HttpStatusCode(response.Name))
		if !response.Type.Definition.IsEmpty() {
			w.Line("ctx.body = t.encode(models.TMessage, result.data)")
		}
	}
	w.Line("  }")
	w.Line("} catch (error) {")
	w.Line("  ctx.throw(500)")
	w.Line("}")
	w.Unindent()
	w.Line("})")
}