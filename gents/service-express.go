package gents

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
	"strings"
)

var express = "express"

func generateExpressSpecRouter(specification *spec.Spec, generatePath string) *gen.TextFile {
	w := NewTsWriter()
	w.Line("import {Router} from 'express'")
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
	w.Line("  let router = Router()")
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			w.Line("  router.use('%s', %s.%s(%s))", expressVersionUrl(&version), versionModule(&version, "routing"), apiRouterName(&api), apiServiceParamName(&api))
		}
	}
	w.Line("  return router")
	w.Line("}")

	return &gen.TextFile{filepath.Join(generatePath, "spec_router.ts"), w.String()}
}

func expressVersionUrl(version *spec.Version) string {
	url := version.Http.GetUrl()
	if url == "" {
		return "/"
	}
	return url
}

func generateExpressVersionRouting(version *spec.Version, validation string, generatePath string) *gen.TextFile {
	w := NewTsWriter()

	w.Line("import {Router} from 'express'")
	w.Line("import {Request, Response} from 'express'")
	w.Line(importEncoding(validation))
	w.Line("import * as %s from './%s'", modelsPackage, versionFilename(version, "models", ""))
	w.Line("import * as services from './%s'", versionFilename(version, "services", ""))

	for _, api := range version.Http.Apis {
		for _, operation := range api.Operations {
			generateParams(w, paramsTypeName(&operation, "HeaderParams"), true, operation.HeaderParams, validation)
			generateParams(w, paramsTypeName(&operation, "UrlParams"), false, operation.Endpoint.UrlParams, validation)
			generateParams(w, paramsTypeName(&operation, "QueryParams"), false, operation.QueryParams, validation)
		}

		generateExpressApiRouting(w, &api, validation)
	}

	filename := versionFilename(version, "routing", "ts")
	return &gen.TextFile{filepath.Join(generatePath, filename), w.String()}
}

func generateExpressApiRouting(w *gen.Writer, api *spec.Api, validation string) {
	w.EmptyLine()
	w.Line("export let %s = (service: services.%s) => {", apiRouterName(api), serviceInterfaceName(api))
	w.Line("  let router = Router()")
	for _, operation := range api.Operations {
		w.EmptyLine()
		generateExpressOperationRouting(w.Indented(), &operation, validation)
	}
	w.EmptyLine()
	w.Line("  return router")
	w.Line("}")
}

func getExpressUrl(endpoint spec.Endpoint) string {
	url := endpoint.Url
	for _, param := range endpoint.UrlParams {
		url = strings.Replace(url, spec.UrlParamStr(param.Name.Source), ":"+param.Name.Source, -1)
	}
	return url
}

func generateExpressOperationRouting(w *gen.Writer, operation *spec.NamedOperation, validation string) {
	w.Line("router.%s('%s', async (request: Request, response: Response) => {", strings.ToLower(operation.Endpoint.Method), getExpressUrl(operation.Endpoint))
	w.Indent()

	apiCallParamsObject := generateParametersParsing(
		w, operation, validation,
		"request.body",
		"request.headers",
		"request.params",
		"request.query",
		"response.status(400).send()",
	)

	w.Line("try {")
	w.Line("  let result = await service.%s(%s)", operation.Name.CamelCase(), apiCallParamsObject)
	w.Line("  switch (result.status) {")
	for _, response := range operation.Responses {
		w.Line("    case '%s':", response.Name.FlatCase())
		responseBody := ".send()"
		if !response.Type.Definition.IsEmpty() {
			responseBody = fmt.Sprintf(".type('json').send(JSON.stringify(t.encode(%s.%s, result.data)))", modelsPackage, runtimeType(validation, &response.Type.Definition))
		}
		w.Line("      response.status(%s)%s", spec.HttpStatusCode(response.Name), responseBody)
	}
	w.Line("  }")
	w.Line("} catch (error) {")
	w.Line("  response.status(500).send()")
	w.Line("}")
	w.Unindent()
	w.Line("})")
}
