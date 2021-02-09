package gents

import (
	"bytes"
	"strings"
	"github.com/ModaOperandi/spec"
	"path/filepath"
	"specgen/gen"
)

func GenerateAxiosClient(serviceFile string, generatePath string) error {
	spec, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}

	models := GenerateIoTsModels(spec, filepath.Join(generatePath, "models.ts"))
	client := generateAxiosClient(spec, filepath.Join(generatePath, "client.ts"))

	files := []gen.TextFile{*models, *client}
	err = gen.WriteFiles(files, true)
	if err != nil {
		return err
	}

	return nil
}

func generateAxiosClient(spec *spec.Spec, outPath string) *gen.TextFile {
	w := new(bytes.Buffer)
	generateClients(spec, w)
	return &gen.TextFile{Path: outPath, Content: w.String()}
}

func getUrl(endpoint spec.Endpoint) string {
	url := endpoint.Url
	for _, param := range endpoint.UrlParams {
		url = strings.Replace(url, spec.UrlParamStr(param.Name.Source), "${"+param.Name.CamelCase()+"}", -1)
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
			tsParams = append(tsParams, param.Name.CamelCase()+requiredSign+": "+TsType(&param.Type.Definition))
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
	return strings.Join(operationParams, ", ")
}