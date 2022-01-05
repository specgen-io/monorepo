package gents

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func generateFetchApiClient(api spec.Api, node bool, validation string, validationModule, modelsModule, paramsModule, module module) *sources.CodeFile {
	w := NewTsWriter()
	if node {
		w.Line(`import { URL, URLSearchParams } from 'url'`)
		w.Line(`import fetch from 'node-fetch'`)
	}
	w.Line(`import { strParamsItems, stringify } from '%s'`, paramsModule.GetImport(module))
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
	return &sources.CodeFile{module.GetPath(), w.String()}
}

func generateFetchOperation(w *sources.Writer, operation *spec.NamedOperation, validation string) {
	body := operation.Body
	hasQueryParams := len(operation.QueryParams) > 0
	hasHeaderParams := len(operation.HeaderParams) > 0
	w.Line(`%s: async (%s): Promise<%s> => {`, operation.Name.CamelCase(), createOperationParams(operation), responseType(operation, ""))
	params := ``
	if hasQueryParams {
		w.Line(`  const query = strParamsItems({`)
		for _, p := range operation.QueryParams {
			w.Line(`    "%s": parameters.%s,`, p.Name.Source, p.Name.CamelCase())
		}
		w.Line(`  })`)
		params = `?${new URLSearchParams(query)}`
	}
	fetchConfig := fmt.Sprintf(`method: '%s'`, strings.ToUpper(operation.Endpoint.Method))
	if hasHeaderParams {
		w.Line(`  const headers = strParamsItems({`)
		for _, p := range operation.HeaderParams {
			w.Line(`    "%s": parameters.%s,`, p.Name.Source, p.Name.CamelCase())
		}
		w.Line(`  })`)
		fetchConfig += `, headers: headers`
	}
	w.Line("  const url = config.baseURL+`%s%s`", getUrl(operation.Endpoint), params)
	if body != nil {
		if body.Type.Definition.Plain == spec.TypeString {
			fetchConfig += `, body: parameters.body`
		} else {
			w.Line(`  const bodyJson = t.encode(%s, parameters.body)`, runtimeTypeFromPackage(validation, modelsPackage, &body.Type.Definition))
			fetchConfig += `, body: JSON.stringify(bodyJson)`
		}
	}
	w.Line("  const response = await fetch(url, {%s})", fetchConfig)
	w.Line(`  switch (response.status) {`)
	for _, response := range operation.Responses {
		w.Line(`    case %s:`, spec.HttpStatusCode(response.Name))
		w.Line(`      return Promise.resolve(%s)`, clientResponseResult(&response, validation, `await response.text()`, `await response.json()`))
	}
	w.Line(`    default:`)
	w.Line("      throw new Error(`Unexpected status code ${ response.status }`)")
	w.Line(`  }`)
	w.Line(`},`)
}
