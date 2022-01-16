package client

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gents/modules"
	"github.com/specgen-io/specgen/v2/gents/responses"
	"github.com/specgen-io/specgen/v2/gents/types"
	"github.com/specgen-io/specgen/v2/gents/validation"
	"github.com/specgen-io/specgen/v2/gents/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

type fetchGenerator struct {
	node       bool
	validation validation.Validation
}

func (g *fetchGenerator) generateApiClient(api spec.Api, validationModule, modelsModule, paramsModule, module modules.Module) *sources.CodeFile {
	w := writer.NewTsWriter()
	if g.node {
		w.Line(`import { URL, URLSearchParams } from 'url'`)
		w.Line(`import fetch from 'node-fetch'`)
	}
	w.Line(`import { strParamsItems, stringify } from '%s'`, paramsModule.GetImport(module))
	w.Line(`import * as t from '%s'`, validationModule.GetImport(module))
	w.Line(`import * as %s from '%s'`, types.ModelsPackage, modelsModule.GetImport(module))
	w.EmptyLine()
	w.Line(`export const client = (config: {baseURL: string}) => {`)
	w.Line(`  return {`)
	for i, operation := range api.Operations {
		if i > 0 {
			w.EmptyLine()
		}
		g.generateFetchOperation(w.IndentedWith(2), &operation)
	}
	w.Line(`  }`)
	w.Line(`}`)
	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			w.EmptyLine()
			responses.GenerateOperationResponse(w, &operation)
		}
	}
	return &sources.CodeFile{module.GetPath(), w.String()}
}

func (g *fetchGenerator) generateFetchOperation(w *sources.Writer, operation *spec.NamedOperation) {
	body := operation.Body
	hasQueryParams := len(operation.QueryParams) > 0
	hasHeaderParams := len(operation.HeaderParams) > 0
	w.Line(`%s: async (%s): Promise<%s> => {`, operation.Name.CamelCase(), createOperationParams(operation), responses.ResponseType(operation, ""))
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
			w.Line(`  const bodyJson = t.encode(%s, parameters.body)`, g.validation.RuntimeTypeFromPackage(types.ModelsPackage, &body.Type.Definition))
			fetchConfig += `, body: JSON.stringify(bodyJson)`
		}
	}
	w.Line("  const response = await fetch(url, {%s})", fetchConfig)
	w.Line(`  switch (response.status) {`)
	for _, response := range operation.Responses {
		w.Line(`    case %s:`, spec.HttpStatusCode(response.Name))
		w.Line(`      return Promise.resolve(%s)`, clientResponseResult(g.validation, &response, `await response.text()`, `await response.json()`))
	}
	w.Line(`    default:`)
	w.Line("      throw new Error(`Unexpected status code ${ response.status }`)")
	w.Line(`  }`)
	w.Line(`},`)
}
