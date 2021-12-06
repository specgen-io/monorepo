package gents

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
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
		modelsModule := versionModule.Submodule("models")
		sources = append(sources, *generateVersionModels(&version, validation, validationModule, modelsModule))
		for _, api := range version.Http.Apis {
			apiModule := versionModule.Submodule(api.Name.SnakeCase())
			sources = append(sources, *generateApiClient(api, validation, validationModule, modelsModule, apiModule))
		}
	}

	err := gen.WriteFiles(sources, true)

	if err != nil {
		return err
	}

	return nil
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

func generateApiClient(api spec.Api, validation string, validationModule module, modelsModule module, module module) *gen.TextFile {
	w := NewTsWriter()
	w.Line(`import { AxiosInstance, AxiosRequestConfig } from 'axios'`)
	w.Line(`import * as t from '%s'`, validationModule.GetImport(module))
	w.Line(`import * as %s from '%s'`, modelsPackage, modelsModule.GetImport(module))
	w.EmptyLine()
	w.Line(`export const client = (axiosInstance: AxiosInstance) => {`)
	w.Line(`  return {`)
	w.Line(`    axiosInstance,`)
	for _, operation := range api.Operations {
		generateOperation(w.IndentedWith(2), &operation, validation)
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

func generateOperation(w *gen.Writer, operation *spec.NamedOperation, validation string) {
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
		w.Line(`      return Promise.resolve(%s)`, clientResponseResult(&response, validation))
	}
	w.Line(`    default:`)
	w.Line("      throw new Error(`Unexpected status code ${ response.status }`)")
	w.Line(`  }`)
	w.Line(`},`)
}

func clientResponseResult(response *spec.NamedResponse, validation string) string {
	if response.Type.Definition.IsEmpty() {
		if len(response.Operation.Responses) == 1 {
			return ""
		}
		return fmt.Sprintf(`{ status: "%s" }`, response.Name.Source)
	} else {
		data := `response.data`
		if response.Type.Definition.Plain != spec.TypeString {
			data = fmt.Sprintf(`t.decode(%s.%s, response.data)`, modelsPackage, runtimeType(validation, &response.Type.Definition))
		}
		if len(response.Operation.Responses) == 1 {
			return data
		}
		return fmt.Sprintf(`{ status: "%s", data: %s }`, response.Name.Source, data)
	}
}
