package client

import (
	"github.com/specgen-io/specgen/v2/gents/modules"
	"github.com/specgen-io/specgen/v2/gents/responses"
	"github.com/specgen-io/specgen/v2/gents/types"
	"github.com/specgen-io/specgen/v2/gents/validations"
	"github.com/specgen-io/specgen/v2/gents/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

type axiosGenerator struct {
	validation validations.Validation
}

func (g *axiosGenerator) ApiClient(api spec.Api, validationModule, modelsModule, paramsModule, module modules.Module) *sources.CodeFile {
	w := writer.NewTsWriter()
	w.Line(`import { AxiosInstance, AxiosRequestConfig } from 'axios'`)
	w.Line(`import { strParamsItems, strParamsObject, stringify } from '%s'`, paramsModule.GetImport(module))
	w.Line(`import * as t from '%s'`, validationModule.GetImport(module))
	w.Line(`import * as %s from '%s'`, types.ModelsPackage, modelsModule.GetImport(module))
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
	return &sources.CodeFile{module.GetPath(), w.String()}
}

func (g *axiosGenerator) operation(w *sources.Writer, operation *spec.NamedOperation) {
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
	if hasHeaderParams {
		w.Line(`  const headers = strParamsObject({`)
		for _, p := range operation.HeaderParams {
			w.Line(`    "%s": parameters.%s,`, p.Name.Source, p.Name.CamelCase())
		}
		w.Line(`  })`)
		axiosConfigParts = append(axiosConfigParts, `headers: headers`)
	}
	axiosConfig := strings.Join(axiosConfigParts, `, `)
	if body != nil {
		if body.Type.Definition.Plain == spec.TypeString {
			w.Line("  const response = await axiosInstance.%s(`%s`, parameters.body, {%s})", strings.ToLower(operation.Endpoint.Method), getUrl(operation.Endpoint), axiosConfig)
		} else {
			w.Line(`  const bodyJson = t.encode(%s, parameters.body)`, g.validation.RuntimeTypeFromPackage(types.ModelsPackage, &body.Type.Definition))
			w.Line("  const response = await axiosInstance.%s(`%s`, bodyJson, {%s})", strings.ToLower(operation.Endpoint.Method), getUrl(operation.Endpoint), axiosConfig)
		}
	} else {
		w.Line("  const response = await axiosInstance.%s(`%s`, {%s})", strings.ToLower(operation.Endpoint.Method), getUrl(operation.Endpoint), axiosConfig)
	}
	w.Line(`  switch (response.status) {`)
	for _, response := range operation.Responses {
		w.Line(`    case %s:`, spec.HttpStatusCode(response.Name))
		w.Line(`      return Promise.resolve(%s)`, clientResponseResult(g.validation, &response, `response.data`, `response.data`))
	}
	w.Line(`    default:`)
	w.Line("      throw new Error(`Unexpected status code ${ response.status }`)")
	w.Line(`  }`)
	w.Line(`},`)
}
