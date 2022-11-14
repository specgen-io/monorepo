package client

import (
	"generator"
	"spec"
	"strings"
	"typescript/responses"
	"typescript/types"
	"typescript/validations"
	"typescript/writer"
)

type axiosGenerator struct {
	Modules    *Modules
	validation validations.Validation
}

func (g *axiosGenerator) ApiClient(api *spec.Api) *generator.CodeFile {
	w := writer.New(g.Modules.Client(api))
	w.Imports.LibNames(`axios`, `AxiosInstance`, `AxiosRequestConfig`)
	w.Imports.Names(g.Modules.Params, `strParamsItems`, `strParamsObject`, `stringify`)
	w.Imports.Star(g.Modules.Validation, `t`)
	w.Imports.Star(g.Modules.Models(api.InHttp.InVersion), types.ModelsPackage)
	w.EmptyLine()
	w.Line(`export const client = (axiosInstance: AxiosInstance) => {`)
	w.Line(`  return {`)
	w.Line(`    axiosInstance,`)
	for _, operation := range api.Operations {
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

func (g *axiosGenerator) operation(w *writer.Writer, operation *spec.NamedOperation) {
	body := operation.Body
	hasQueryParams := len(operation.QueryParams) > 0
	hasHeaderParams := len(operation.HeaderParams) > 0
	w.EmptyLine()
	w.Line(`%s: async (%s): Promise<%s> => {`, operation.Name.CamelCase(), createOperationParams(operation), responses.ResponseType(operation, ""))
	axiosConfigParts := []string{}
	if hasQueryParams {
		w.Line(`  const query = strParamsItems({`)
		for _, p := range operation.QueryParams {
			w.Line(`    "%s": parameters.%s,`, p.Name.Source, p.Name.CamelCase())
		}
		w.Line(`  })`)
		axiosConfigParts = append(axiosConfigParts, `params: new URLSearchParams(query)`)
	}
	if hasHeaderParams || body != nil {
		w.Line(`  const headers = strParamsObject({`)
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
		axiosConfigParts = append(axiosConfigParts, `headers: headers`)
	}
	axiosConfig := strings.Join(axiosConfigParts, `, `)
	if operation.BodyIs(spec.BodyEmpty) {
		w.Line("  const response = await axiosInstance.%s(`%s`, {%s})", strings.ToLower(operation.Endpoint.Method), getUrl(operation.Endpoint), axiosConfig)
	}
	if operation.BodyIs(spec.BodyString) {
		w.Line("  const response = await axiosInstance.%s(`%s`, parameters.body, {%s})", strings.ToLower(operation.Endpoint.Method), getUrl(operation.Endpoint), axiosConfig)
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`  const bodyJson = t.encode(%s, parameters.body)`, g.validation.RuntimeType(&body.Type.Definition))
		w.Line("  const response = await axiosInstance.%s(`%s`, bodyJson, {%s})", strings.ToLower(operation.Endpoint.Method), getUrl(operation.Endpoint), axiosConfig)
	}
	w.Line(`  switch (response.status) {`)
	if len(operation.Responses) == 1 {
		response := operation.Responses[0]
		body := clientResponseBody(g.validation, &response.Response, `response.data`, `response.data`)
		w.Line(`    case %s:`, spec.HttpStatusCode(response.Name))
		w.Line(`      return Promise.resolve(%s)`, body)
	} else {
		for _, response := range operation.Responses {
			body := clientResponseBody(g.validation, &response.Response, `response.data`, `response.data`)
			w.Line(`    case %s:`, spec.HttpStatusCode(response.Name))
			w.Line(`      return Promise.resolve(%s)`, responses.New(&response.Response, body))
		}
	}
	w.Line(`    default:`)
	w.Line("      throw new Error(`Unexpected status code ${ response.status }`)")
	w.Line(`  }`)

	w.Line(`},`)
}