package gents

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func GenerateClient(specification *spec.Spec, generatePath string, client string, validation string) error {
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
			if client == "axios" {
				sources = append(sources, *generateAxiosApiClient(api, validation, validationModule, modelsModule, apiModule))
			}
			if client == "node-fetch" {
				sources = append(sources, *generateFetchApiClient(api, true, validation, validationModule, modelsModule, apiModule))
			}
			if client == "browser-fetch" {
				sources = append(sources, *generateFetchApiClient(api, false, validation, validationModule, modelsModule, apiModule))
			}
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

func clientResponseResult(response *spec.NamedResponse, validation string, responseJson string) string {
	if response.Type.Definition.IsEmpty() {
		if len(response.Operation.Responses) == 1 {
			return ""
		}
		return fmt.Sprintf(`{ status: "%s" }`, response.Name.Source)
	} else {
		data := fmt.Sprintf(`t.decode(%s.%s, %s)`, modelsPackage, runtimeType(validation, &response.Type.Definition), responseJson)
		if len(response.Operation.Responses) == 1 {
			return data
		}
		return fmt.Sprintf(`{ status: "%s", data: %s }`, response.Name.Source, data)
	}
}
