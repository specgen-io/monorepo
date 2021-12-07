package gents

import (
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func generateAxiosApiClient(api spec.Api, validation string, validationModule, modelsModule, paramsModule, module module) *gen.TextFile {
	w := NewTsWriter()
	w.Line(`import { AxiosInstance, AxiosRequestConfig } from 'axios'`)
	w.Line(`import { params, stringify } from '%s'`, paramsModule.GetImport(module))
	w.Line(`import * as t from '%s'`, validationModule.GetImport(module))
	w.Line(`import * as %s from '%s'`, modelsPackage, modelsModule.GetImport(module))
	w.EmptyLine()
	w.Line(`export const client = (axiosInstance: AxiosInstance) => {`)
	w.Line(`  return {`)
	w.Line(`    axiosInstance,`)
	for _, operation := range api.Operations {
		generateAxiosOperation(w.IndentedWith(2), &operation, validation)
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

func generateAxiosOperation(w *gen.Writer, operation *spec.NamedOperation, validation string) {
	body := operation.Body
	hasQueryParams := len(operation.QueryParams) > 0
	hasHeaderParams := len(operation.HeaderParams) > 0
	w.EmptyLine()
	w.Line(`%s: async (%s): Promise<%s> => {`, operation.Name.CamelCase(), createOperationParams(operation), responseType(operation, ""))
	if hasQueryParams {
		w.Line(`  const params = {`)
		for _, p := range operation.QueryParams {
			w.Line(`    "%s": parameters.%s,`, p.Name.Source, p.Name.CamelCase())
		}
		w.Line(`  }`)
	}
	if hasHeaderParams {
		w.Line(`  const headers = {`)
		for _, p := range operation.HeaderParams {
			w.Line(`    "%s": parameters.%s,`, p.Name.Source, p.Name.CamelCase())
		}
		w.Line(`  }`)
	}
	params := ``
	if hasQueryParams {
		params = `params: params,`
	}
	headers := ``
	if hasHeaderParams {
		headers = `headers: headers,`
	}
	w.Line(`  const config: AxiosRequestConfig = {%s%s}`, params, headers)
	if body != nil {
		if body.Type.Definition.Plain == spec.TypeString {
			w.Line("  const response = await axiosInstance.%s(`%s`, parameters.body, config)", strings.ToLower(operation.Endpoint.Method), getUrl(operation.Endpoint))
		} else {
			w.Line(`  const bodyJson = t.encode(%s.%s, parameters.body)`, modelsPackage, runtimeType(validation, &body.Type.Definition))
			w.Line("  const response = await axiosInstance.%s(`%s`, bodyJson, config)", strings.ToLower(operation.Endpoint.Method), getUrl(operation.Endpoint))
		}
	} else {
		w.Line("  const response = await axiosInstance.%s(`%s`, config)", strings.ToLower(operation.Endpoint.Method), getUrl(operation.Endpoint))
	}
	w.Line(`  switch (response.status) {`)
	for _, response := range operation.Responses {
		w.Line(`    case %s:`, spec.HttpStatusCode(response.Name))
		w.Line(`      return Promise.resolve(%s)`, clientResponseResult(&response, validation, `response.data`))
	}
	w.Line(`    default:`)
	w.Line("      throw new Error(`Unexpected status code ${ response.status }`)")
	w.Line(`  }`)
	w.Line(`},`)
}