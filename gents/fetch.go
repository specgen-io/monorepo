package gents

import (
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func generateFetchApiClient(api spec.Api, node bool, validation string, validationModule module, modelsModule module, module module) *gen.TextFile {
	w := NewTsWriter()
	if node {
		w.Line(`import { URL, URLSearchParams } from 'url'`)
		w.Line(`import fetch from 'node-fetch'`)
	}
	w.Line(`import * as t from '%s'`, validationModule.GetImport(module))
	w.Line(`import * as %s from '%s'`, modelsPackage, modelsModule.GetImport(module))
	w.EmptyLine()
	w.Line(`export const client = (config: {baseURL: string}) => {`)
	w.Line(`  return {`)
	for i, operation := range api.Operations {
		if i > 0 {
			w.EmptyLine()
		}
		generateFetchOperation(w.IndentedWith(2), &operation, validation)
	}
	w.Line(`  }`)
	w.Line(`}`)
	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			w.EmptyLine()
			generateOperationResponse(w, &operation)
		}
	}
	return &gen.TextFile{module.GetPath(), w.String()}
}

func generateFetchOperation(w *gen.Writer, operation *spec.NamedOperation, validation string) {
	body := operation.Body
	hasQueryParams := len(operation.QueryParams) > 0
	hasHeaderParams := len(operation.HeaderParams) > 0
	w.Line(`%s: async (%s): Promise<%s> => {`, operation.Name.CamelCase(), createOperationParams(operation), responseType(operation, ""))
	params := ``
	if hasQueryParams {
		w.Line(`  const params = {`)
		for _, p := range operation.QueryParams {
			w.Line(`    "%s": String(parameters.%s),`, p.Name.Source, p.Name.CamelCase())
		}
		w.Line(`  }`)
		params = `?${new URLSearchParams(params)}`
	}
	headers := ``
	if hasHeaderParams {
		w.Line(`  const headers = {`)
		for _, p := range operation.HeaderParams {
			w.Line(`    "%s": String(parameters.%s),`, p.Name.Source, p.Name.CamelCase())
		}
		w.Line(`  }`)
		headers = `, headers: headers`
	}
	w.Line("  const url = config.baseURL+`%s%s`", getUrl(operation.Endpoint), params)
	fetchBody := ``
	if body != nil {
		w.Line(`  const bodyJson = t.encode(%s.%s, parameters.body)`, modelsPackage, runtimeType(validation, &body.Type.Definition))
		fetchBody = `, body: JSON.stringify(bodyJson)`
	}
	w.Line("  const response = await fetch(url, {method: '%s'%s%s})", strings.ToUpper(operation.Endpoint.Method), headers, fetchBody)
	w.Line(`  switch (response.status) {`)
	for _, response := range operation.Responses {
		w.Line(`    case %s:`, spec.HttpStatusCode(response.Name))
		w.Line(`      return Promise.resolve(%s)`, clientResponseResult(&response, validation, `await response.json()`))
	}
	w.Line(`    default:`)
	w.Line("      throw new Error(`Unexpected status code ${ response.status }`)")
	w.Line(`  }`)
	w.Line(`},`)
}
