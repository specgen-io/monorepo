package gents

import (
	"fmt"
	"github.com/specgen-io/spec"
	"specgen/gen"
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

	apiCallParamsObject := generateParametersParsing(
		w, operation, validation,
		"ctx.request.body",
		"ctx.request.headers",
		"ctx.params",
		"ctx.request.query",
		"ctx.throw(400, error)",
	)

	w.Line("try {")
	w.Line("  let result = await service.%s(%s)", operation.Name.CamelCase(), apiCallParamsObject)
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
