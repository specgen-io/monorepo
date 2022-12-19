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

type AxiosGenerator struct {
	Modules    *Modules
	validation validations.Validation
	CommonGenerator
}

func (g *AxiosGenerator) ApiClient(api *spec.Api) *generator.CodeFile {
	w := writer.New(g.Modules.Client(api))
	w.Imports.Default(`axios`, `axios`)
	w.Imports.LibNames(`axios`, `AxiosInstance`, `AxiosRequestConfig`)
	w.Imports.Names(g.Modules.Params, `strParamsItems`, `strParamsObject`, `stringify`)
	w.Imports.Star(g.Modules.Validation, `t`)
	w.Imports.Star(g.Modules.Models(api.InHttp.InVersion), types.ModelsPackage)
	w.Imports.Star(g.Modules.Errors, `errors`)
	w.EmptyLine()
	w.Line(`export const client = (config?: AxiosRequestConfig) => {`)
	w.Line(`  const axiosInstance = axios.create({...config, validateStatus: () => true})`)
	w.Line(`  axiosInstance.interceptors.response.use(errors.checkRequiredErrors)`)
	w.EmptyLine()
	w.Line(`  return {`)
	w.Line(`    axiosInstance,`)
	for _, operation := range api.Operations {
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

func (g *AxiosGenerator) operation(w *writer.Writer, operation *spec.NamedOperation) {
	body := operation.Body
	hasQueryParams := len(operation.QueryParams) > 0
	hasHeaderParams := len(operation.HeaderParams) > 0
	w.EmptyLine()
	w.Line(`%s: async (%s): Promise<%s> => {`, operation.Name.CamelCase(), createOperationParams(operation), responseType(operation))
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
	for _, response := range operation.Responses.Success() {
		w.Line(`    case %s: %s`, spec.HttpStatusCode(response.Name), g.operationReturn(response))
	}
	for _, response := range operation.Responses.NonRequiredErrors() {
		w.Line(`    case %s: throw new errors.%s(%s)`, spec.HttpStatusCode(response.Name), errorExceptionName(&response.Response), g.responseBody(&response.Response))
	}
	w.Line("    default: throw new errors.ResponseException(`Unexpected status code ${ response.status }`)")
	w.Line(`  }`)

	w.Line(`},`)
}

func (g *AxiosGenerator) operationReturn(response *spec.OperationResponse) string {
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
