package client

import (
	"fmt"
	"generator"
	"spec"
	"strings"
	"typescript/types"
	"typescript/validations"
	"typescript/writer"
)

type FetchGenerator struct {
	Modules    *Modules
	node       bool
	validation validations.Validation
	CommonGenerator
}

func (g *FetchGenerator) ApiClient(api *spec.Api) *generator.CodeFile {
	w := writer.New(g.Modules.Client(api))
	if g.node {
		w.Imports.LibNames(`url`, `URL`, `URLSearchParams`)
		w.Imports.Default(`node-fetch`, `fetch`)
	}
	w.Imports.Names(g.Modules.Params, `strParamsItems`, `stringify`)
	w.Imports.Names(g.Modules.Response, `checkRequiredErrors`)
	w.Imports.Star(g.Modules.Validation, `t`)
	w.Imports.Star(g.Modules.Models(api.InHttp.InVersion), types.ModelsPackage)
	w.Imports.Star(g.Modules.Errors, `errors`)
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
		if len(operation.Responses.Success()) > 1 {
			w.EmptyLine()
			generateOperationResponse(w, &operation)
		}
	}
	return w.ToCodeFile()
}

func (g *FetchGenerator) operation(w *writer.Writer, operation *spec.NamedOperation) {
	body := operation.Body
	hasQueryParams := len(operation.QueryParams) > 0
	hasHeaderParams := len(operation.HeaderParams) > 0
	w.Line(`%s => {`, operationSignature(operation))
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
		if operation.BodyIs(spec.RequestBodyString) {
			w.Line(`    "Content-Type": "text/plain"`)
		}
		if operation.BodyIs(spec.RequestBodyJson) {
			w.Line(`    "Content-Type": "application/json"`)
		}
		w.Line(`  })`)
		fetchConfigParts = append(fetchConfigParts, `headers: headers`)
	}
	w.Line("  const url = config.baseURL+`%s%s`", getUrl(operation.Endpoint), params)
	if operation.BodyIs(spec.RequestBodyString) {
		fetchConfigParts = append(fetchConfigParts, `body: parameters.body`)
	}
	if operation.BodyIs(spec.RequestBodyJson) {
		w.Line(`  const bodyJson = t.encode(%s, parameters.body)`, g.validation.RuntimeType(&body.Type.Definition))
		fetchConfigParts = append(fetchConfigParts, `body: JSON.stringify(bodyJson)`)
	}
	fetchConfig := strings.Join(fetchConfigParts, ", ")
	w.Line("  const response = await fetch(url, {%s})", fetchConfig)
	w.Line("  await checkRequiredErrors(response)")

	w.Line(`  switch (response.status) {`)
	for _, response := range operation.Responses.Success() {
		w.Line(`    case %s:`, spec.HttpStatusCode(response.Name))
		w.Line(`      %s`, g.operationReturn(response))
	}
	for _, response := range operation.Responses.NonRequiredErrors() {
		w.Line(`    case %s:`, spec.HttpStatusCode(response.Name))
		w.Line(`      throw %s`, g.errorResponse(&response.Response))
	}
	w.Line(`    default:`)
	w.Line("      throw new errors.ResponseException(`Unexpected status code ${ response.status }`)")
	w.Line(`  }`)

	w.Line(`},`)
}

func (g *FetchGenerator) operationReturn(response *spec.OperationResponse) string {
	body := g.responseBody(&response.Response)
	if len(response.Operation.Responses.Success()) == 1 {
		if body == `` {
			return `return`
		} else {
			return fmt.Sprintf(`return %s`, body)
		}
	} else {
		return fmt.Sprintf(`return %s`, newResponse(&response.Response, body))
	}
}

func (g *FetchGenerator) responseBody(response *spec.Response) string {
	if response.BodyIs(spec.ResponseBodyString) {
		return `await response.text()`
	}
	if response.BodyIs(spec.ResponseBodyJson) {
		data := fmt.Sprintf(`t.decode(%s, %s)`, g.validation.RuntimeType(&response.Type.Definition), `await response.json()`)
		return data
	}
	return ""
}

func (g *FetchGenerator) errorResponse(response *spec.Response) string {
	return fmt.Sprintf(`new errors.%s(%s)`, errorExceptionName(response), g.responseBody(response))
}

func (g *FetchGenerator) ErrorResponses(httpErrors *spec.HttpErrors) *generator.CodeFile {
	w := writer.New(g.Modules.Response)
	if g.node {
		w.Imports.LibNames("node-fetch", "Response")
	}
	w.Imports.Star(g.Modules.Validation, "t")
	w.Imports.Star(g.Modules.Errors, `errors`)
	w.Line(`export const checkRequiredErrors = async (response: Response) => {`)
	w.Line(`  switch (response.status) {`)
	for _, response := range httpErrors.Responses.Required() {
		w.Line(`    case %s:`, spec.HttpStatusCode(response.Name))
		w.Line(`      throw %s`, g.errorResponse(&response.Response))
	}
	w.Line(`  }`)
	w.Line(`}`)
	return w.ToCodeFile()
}
