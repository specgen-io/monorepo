package client

import (
	"fmt"
	"github.com/specgen-io/specgen/spec/v2"
	"github.com/specgen-io/specgen/specgen/v2/gen/ts/modules"
	"github.com/specgen-io/specgen/specgen/v2/gen/ts/responses"
	"github.com/specgen-io/specgen/specgen/v2/gen/ts/types"
	"github.com/specgen-io/specgen/specgen/v2/gen/ts/validations"
	"github.com/specgen-io/specgen/specgen/v2/gen/ts/writer"
	"github.com/specgen-io/specgen/specgen/v2/generator"
	"strings"
)

type fetchGenerator struct {
	node       bool
	validation validations.Validation
}

func (g *fetchGenerator) ApiClient(api spec.Api, validationModule, modelsModule, paramsModule, module modules.Module) *generator.CodeFile {
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
		g.operation(w.IndentedWith(2), &operation)
	}
	w.Line(`  }`)
	w.Line(`}`)
	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			w.EmptyLine()
			responses.GenerateOperationResponse(w, &operation)
		}
	}
	return &generator.CodeFile{module.GetPath(), w.String()}
}

func (g *fetchGenerator) operation(w *generator.Writer, operation *spec.NamedOperation) {
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
	fetchConfigParts := []string{fmt.Sprintf(`method: '%s'`, strings.ToUpper(operation.Endpoint.Method))}
	if hasHeaderParams || body != nil {
		w.Line(`  const headers = strParamsItems({`)
		for _, p := range operation.HeaderParams {
			w.Line(`    "%s": parameters.%s,`, p.Name.Source, p.Name.CamelCase())
		}
		if operation.BodyIs(spec.BodyString) {
			w.Line(`    "Content-Type": "text/plain"`)
		}
		if operation.BodyIs(spec.BodyJson) {
			w.Line(`    "Content-Type": "application/json"`)
		}
		w.Line(`  })`)
		fetchConfigParts = append(fetchConfigParts, `headers: headers`)
	}
	w.Line("  const url = config.baseURL+`%s%s`", getUrl(operation.Endpoint), params)
	if operation.BodyIs(spec.BodyString) {
		fetchConfigParts = append(fetchConfigParts, `body: parameters.body`)
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`  const bodyJson = t.encode(%s, parameters.body)`, g.validation.RuntimeTypeFromPackage(types.ModelsPackage, &body.Type.Definition))
		fetchConfigParts = append(fetchConfigParts, `body: JSON.stringify(bodyJson)`)
	}
	fetchConfig := strings.Join(fetchConfigParts, ", ")
	w.Line("  const response = await fetch(url, {%s})", fetchConfig)

	w.Line(`  switch (response.status) {`)
	if len(operation.Responses) == 1 {
		response := operation.Responses[0]
		body := clientResponseBody(g.validation, &response.Response, `await response.text()`, `await response.json()`)
		w.Line(`    case %s:`, spec.HttpStatusCode(response.Name))
		w.Line(`      return Promise.resolve(%s)`, body)
	} else {
		for _, response := range operation.Responses {
			body := clientResponseBody(g.validation, &response.Response, `await response.text()`, `await response.json()`)
			w.Line(`    case %s:`, spec.HttpStatusCode(response.Name))
			w.Line(`      return Promise.resolve(%s)`, responses.New(&response.Response, body))
		}
	}
	w.Line(`    default:`)
	w.Line("      throw new Error(`Unexpected status code ${ response.status }`)")
	w.Line(`  }`)

	w.Line(`},`)
}
