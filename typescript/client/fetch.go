package client

import (
	"fmt"
	"strings"
	"typescript/imports"

	"generator"
	"spec"
	"typescript/responses"
	"typescript/types"
	"typescript/validations"
	"typescript/writer"
)

type fetchGenerator struct {
	Modules    *Modules
	node       bool
	validation validations.Validation
}

func (g *fetchGenerator) ApiClient(api *spec.Api) *generator.CodeFile {
	w := writer.New(g.Modules.Client(api))
	imports := imports.New(g.Modules.Client(api))
	if g.node {
		imports.LibNames(`url`, `URL`, `URLSearchParams`)
		imports.Default(`node-fetch`, `fetch`)
	}
	imports.Names(g.Modules.Params, `strParamsItems`, `stringify`)
	imports.Star(g.Modules.Validation, `t`)
	imports.Star(g.Modules.Models(api.InHttp.InVersion), types.ModelsPackage)
	imports.Write(w)
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
	return w.ToCodeFile()
}

func (g *fetchGenerator) operation(w generator.Writer, operation *spec.NamedOperation) {
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
		w.Line(`  const bodyJson = t.encode(%s, parameters.body)`, g.validation.RuntimeType(&body.Type.Definition))
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
