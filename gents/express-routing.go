package gents

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
	"strings"
)

func generateExpressVersionRouting(version *spec.Version, generatePath string) *gen.TextFile {
	w := NewTsWriter()

	w.Line("import express from 'express'")
	w.Line("import {Request, Response} from 'express'")
	w.Line("import * as t from './superstruct'")
	w.Line("import * as models from './%s'", versionFilename(version,"models", ""))
	w.Line("import * as services from './%s'", versionFilename(version,"services", ""))

	for _, api := range version.Http.Apis {
		for _, operation := range api.Operations {
			generateParams(w, &operation.Name, "HeaderParams", true, operation.HeaderParams)
			generateParams(w, &operation.Name, "UrlParams", false, operation.Endpoint.UrlParams)
			generateParams(w, &operation.Name, "QueryParams", false, operation.QueryParams)
		}

		generateExpressApiRouting(w, &api)
	}

	filename := versionFilename(version, "routing", "ts")
	return &gen.TextFile{filepath.Join(generatePath, filename), w.String()}
}

func generateExpressApiRouting(w *gen.Writer, api *spec.Api) {
	w.EmptyLine()
	w.Line("export let %sRouter = (service: services.%s): express.Router => {", api.Name.CamelCase(), serviceInterfaceName(api))
	w.Line("  let router = express.Router()")
	for _, operation := range api.Operations {
		w.EmptyLine()
		generateExpressOperationRouting(w.Indented(), &operation)
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

func generateExpressOperationRouting(w *gen.Writer, operation *spec.NamedOperation) {
	w.Line("router.%s('%s', async (request: Request, response: Response) => {", strings.ToLower(operation.Endpoint.Method), getExpressUrl(operation.Endpoint))
	w.Line("  try {")
	generateExpressOperationRoutingCode(w.IndentedWith(2), operation)
	w.Line("  } catch (error) {")
	w.Line("    response.status(500).send()")
	w.Line("  }")
	w.Line("})")
}

func generateExpressOperationRoutingCode(w *gen.Writer, operation *spec.NamedOperation) {
	apiCallParams := []string{}
	if operation.Body != nil {
		w.Line("let body = t.decode(models.%s, request.body)", SsType(&operation.Body.Type.Definition))
		apiCallParams = append(apiCallParams, "body")
	}
	if len(operation.Endpoint.UrlParams) > 0 {
		w.Line("let urlParams = t.decode(%s, request.params)", stringParamsRuntimeTypeName(&operation.Name, "UrlParams"))
		apiCallParams = append(apiCallParams, "...urlParams")
	}
	if len(operation.HeaderParams) > 0 {
		w.Line("let header = t.decode(%s, request.headers)", stringParamsRuntimeTypeName(&operation.Name, "HeaderParams"))
		apiCallParams = append(apiCallParams, "...header")
	}
	if len(operation.QueryParams) > 0 {
		w.Line("let query = t.decode(%s, request.query)", stringParamsRuntimeTypeName(&operation.Name, "QueryParams"))
		apiCallParams = append(apiCallParams, "...query")
	}
	w.Line("let result = await service.%s({%s})", operation.Name.CamelCase(), strings.Join(apiCallParams, ", "))
	w.Line("switch (result.status) {")
	for _, response := range operation.Responses {
		w.Line("  case '%s':", response.Name.FlatCase())
		responseBody := ".send()"
		if !response.Type.Definition.IsEmpty() {
			responseBody = fmt.Sprintf(".type('json').send(JSON.stringify(t.encode(models.%s, result.data)))", SsType(&response.Type.Definition))
		}
		w.Line("    response.status(%s)%s", spec.HttpStatusCode(response.Name), responseBody)
	}
	w.Line("}")
}