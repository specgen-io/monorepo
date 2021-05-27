package gents

import (
	"fmt"
	spec "github.com/specgen-io/spec"
	"gopkg.in/specgen-io/specgen.v2/gen"
	"path/filepath"
	"strings"
)

func GenerateAxiosClient(serviceFile string, generatePath string, models string) error {
	spec, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}

	files := generateAxiosClients(spec, models == "io-ts", generatePath)

	if models == "io-ts" {
		iots := generateIoTs(filepath.Join(generatePath, "io-ts.ts"))
		files = append(files, *iots)
	} else {
		superstruct := generateSuperstruct(filepath.Join(generatePath, "superstruct.ts"))
		files = append(files, *superstruct)
	}
	err = gen.WriteFiles(files, true)
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

func generateClientApiClass(w *gen.Writer, api spec.Api) {
	w.EmptyLine()
	w.Line(`export const %sClient = (axiosInstance: AxiosInstance) => {`, api.Name.CamelCase())
	w.Line(`  return {`)
	w.Line(`    axiosInstance,`)
	for _, operation := range api.Operations {
		generateOperation(w.IndentedWith(2), &operation)
	}
	w.Line(`  }`)
	w.Line(`}`)
}

func generateOperation(w *gen.Writer, operation *spec.NamedOperation) {
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
		w.Line(`  const bodyJson = t.encode(%s, parameters.body)`, IoTsType(&body.Type.Definition))
		w.Line("  const response = await axiosInstance.%s(`%s`, bodyJson, config)", strings.ToLower(operation.Endpoint.Method), getUrl(operation.Endpoint))
	} else {
		w.Line("  const response = await axiosInstance.%s(`%s`, config)", strings.ToLower(operation.Endpoint.Method), getUrl(operation.Endpoint))
	}
	w.Line(`  switch (response.status) {`)
	for _, response := range operation.Responses {
		dataParam := ``
		if !response.Type.Definition.IsEmpty() {
			dataParam = fmt.Sprintf(`, data: t.decode(%s, response.data)`, IoTsType(&response.Type.Definition))
		}
		w.Line(`    case %s:`, spec.HttpStatusCode(response.Name))
		w.Line(`      return Promise.resolve({ status: "%s"%s })`, response.Name.Source, dataParam)
	}
	w.Line(`    default:`)
	w.Line("      throw new Error(`Unexpected status code ${ response.status }`)")
	w.Line(`  }`)
	w.Line(`},`)
}

func generateAxiosClient(w *gen.Writer, version *spec.Version) {
	w.Line(`import { AxiosInstance, AxiosRequestConfig } from 'axios'`)
	for _, api := range version.Http.Apis {
		generateClientApiClass(w, api)
		for _, operation := range api.Operations {
			generateIoTsResponse(w, &operation)
		}
	}
}

func generateAxiosClients(spec *spec.Spec, ioTsModels bool, path string) []gen.TextFile {
	files := []gen.TextFile{}
	for _, version := range spec.Versions {
		w := NewTsWriter()
		if ioTsModels {
			generateIoTsModels(w, &version)
		} else {
			generateSuperstructModels(w, &version)
		}
		w.EmptyLine()
		generateAxiosClient(w, &version)
		filename := "index.ts"
		if version.Version.Source != "" {
			filename = version.Version.FlatCase() + ".ts"
		}
		file := gen.TextFile{filepath.Join(path, filename), w.String()}
		files = append(files, file)
	}
	return files
}
