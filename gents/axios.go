package gents

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"strings"
)

func GenerateAxiosClient(specification *spec.Spec, generatePath string, validation string) error {
	sources := []gen.TextFile{}
	module := Module(generatePath)
	validationModule := module.Submodule(validation)
	validationFile := generateValidation(validation, validationModule)
	sources = append(sources, *validationFile)
	for _, version := range specification.Versions {
		versionModule := module.Submodule(version.Version.FlatCase())
		clientModule := versionModule.Submodule("index")
		sources = append(sources, *generateAxiosClient(&version, validation, validationModule, clientModule))
	}
	modelsFiles := generateModels(specification, validation, validationModule, module)
	sources = append(sources, modelsFiles...)

	err := gen.WriteFiles(sources, true)

	if err != nil {
		return err
	}

	return nil
}

func generateAxiosClient(version *spec.Version, validation string, validationModule module, module module) *gen.TextFile {
	w := NewTsWriter()
	w.Line(`import { AxiosInstance, AxiosRequestConfig } from 'axios'`)
	w.Line(`import * as t from '%s'`, validationModule.GetImport(module))
	w.Line(`import * as %s from './models'`, modelsPackage)
	for _, api := range version.Http.Apis {
		generateClientApiClass(w, api, validation)
		for _, operation := range api.Operations {
			w.EmptyLine()
			generateOperationResponse(w, &operation)
		}
	}
	return &gen.TextFile{module.GetPath(), w.String()}
}

func getUrl(endpoint spec.Endpoint) string {
	url := endpoint.Url
	for _, param := range endpoint.UrlParams {
		url = strings.Replace(url, spec.UrlParamStr(param.Name.Source), "${parameters."+param.Name.CamelCase()+"}", -1)
	}
	return url
}

func responseTypeName(operation *spec.NamedOperation) string {
	return operation.Name.PascalCase() + "Response"
}

func createParams(params []spec.NamedParam, required bool) []string {
	tsParams := []string{}
	for _, param := range params {
		isRequired := param.Default == nil && !param.Type.Definition.IsNullable()
		if isRequired == required {
			requiredSign := ""
			if !isRequired {
				requiredSign = "?"
			}
			paramType := &param.Type.Definition
			if !isRequired && !paramType.IsNullable() {
				paramType = spec.Nullable(paramType)
			}
			tsParams = append(tsParams, param.Name.CamelCase()+requiredSign+": "+TsType(paramType))
		}
	}
	return tsParams
}

func createOperationParams(operation *spec.NamedOperation) string {
	operationParams := []string{}
	operationParams = append(operationParams, createParams(operation.HeaderParams, true)...)
	if operation.Body != nil {
		operationParams = append(operationParams, "body: "+TsType(&operation.Body.Type.Definition))
	}
	operationParams = append(operationParams, createParams(operation.Endpoint.UrlParams, true)...)
	operationParams = append(operationParams, createParams(operation.QueryParams, true)...)
	operationParams = append(operationParams, createParams(operation.HeaderParams, false)...)
	operationParams = append(operationParams, createParams(operation.QueryParams, false)...)
	if len(operationParams) == 0 {
		return ""
	}
	return fmt.Sprintf("parameters: {%s}", strings.Join(operationParams, ", "))
}

func generateClientApiClass(w *gen.Writer, api spec.Api, validation string) {
	w.EmptyLine()
	w.Line(`export const %sClient = (axiosInstance: AxiosInstance) => {`, api.Name.CamelCase())
	w.Line(`  return {`)
	w.Line(`    axiosInstance,`)
	for _, operation := range api.Operations {
		generateOperation(w.IndentedWith(2), &operation, validation)
	}
	w.Line(`  }`)
	w.Line(`}`)
}

func generateOperation(w *gen.Writer, operation *spec.NamedOperation, validation string) {
	body := operation.Body
	hasQueryParams := len(operation.QueryParams) > 0
	hasHeaderParams := len(operation.HeaderParams) > 0
	w.EmptyLine()
	w.Line(`%s: async (%s): Promise<%s> => {`, operation.Name.CamelCase(), createOperationParams(operation), responseTypeName(operation))
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
		w.Line(`  const bodyJson = t.encode(%s.%s, parameters.body)`, modelsPackage, runtimeType(validation, &body.Type.Definition))
		w.Line("  const response = await axiosInstance.%s(`%s`, bodyJson, config)", strings.ToLower(operation.Endpoint.Method), getUrl(operation.Endpoint))
	} else {
		w.Line("  const response = await axiosInstance.%s(`%s`, config)", strings.ToLower(operation.Endpoint.Method), getUrl(operation.Endpoint))
	}
	w.Line(`  switch (response.status) {`)
	for _, response := range operation.Responses {
		dataParam := ``
		if !response.Type.Definition.IsEmpty() {
			dataParam = fmt.Sprintf(`, data: t.decode(%s.%s, response.data)`, modelsPackage, runtimeType(validation, &response.Type.Definition))
		}
		w.Line(`    case %s:`, spec.HttpStatusCode(response.Name))
		w.Line(`      return Promise.resolve({ status: "%s"%s })`, response.Name.Source, dataParam)
	}
	w.Line(`    default:`)
	w.Line("      throw new Error(`Unexpected status code ${ response.status }`)")
	w.Line(`  }`)
	w.Line(`},`)
}
