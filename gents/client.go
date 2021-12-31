package gents

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func GenerateClient(specification *spec.Spec, generatePath string, client string, validation string) *sources.Sources {
	sources := sources.NewSources()
	module := Module(generatePath)

	validationModule := module.Submodule(validation)
	sources.AddGenerated(generateValidation(validation, validationModule))
	paramsModule := module.Submodule("params")
	sources.AddGenerated(generateParamsBuilder(paramsModule))
	for _, version := range specification.Versions {
		versionModule := module.Submodule(version.Version.FlatCase())
		modelsModule := versionModule.Submodule("models")
		sources.AddGenerated(generateVersionModels(&version, validation, validationModule, modelsModule))
		for _, api := range version.Http.Apis {
			apiModule := versionModule.Submodule(api.Name.SnakeCase())
			if client == "axios" {
				sources.AddGenerated(generateAxiosApiClient(api, validation, validationModule, modelsModule, paramsModule, apiModule))
			}
			if client == "node-fetch" {
				sources.AddGenerated(generateFetchApiClient(api, true, validation, validationModule, modelsModule, paramsModule, apiModule))
			}
			if client == "browser-fetch" {
				sources.AddGenerated(generateFetchApiClient(api, false, validation, validationModule, modelsModule, paramsModule, apiModule))
			}
		}
	}

	return sources
}

func getUrl(endpoint spec.Endpoint) string {
	url := endpoint.Url
	for _, param := range endpoint.UrlParams {
		url = strings.Replace(url, spec.UrlParamStr(&param), "${stringify(parameters."+param.Name.CamelCase()+")}", -1)
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

func clientResponseResult(response *spec.NamedResponse, validation string, textResposneData, jsonResponseData string) string {
	if response.Type.Definition.IsEmpty() {
		if len(response.Operation.Responses) == 1 {
			return ""
		}
		return fmt.Sprintf(`{ status: "%s" }`, response.Name.Source)
	} else {
		data := textResposneData
		if response.Type.Definition.Plain != spec.TypeString {
			data = fmt.Sprintf(`t.decode(%s.%s, %s)`, modelsPackage, runtimeType(validation, &response.Type.Definition), jsonResponseData)
		}
		if len(response.Operation.Responses) == 1 {
			return data
		}
		return fmt.Sprintf(`{ status: "%s", data: %s }`, response.Name.Source, data)
	}
}

func generateParamsBuilder(module module) *sources.CodeFile {
	code := `
export function stringify(value: ScalarParam): string {
  if (value instanceof Date) {
      return value.toISOString()
  }
  return String(value)
}

export function stringifyX(value: ScalarParam | ScalarParam[]): string | string[] {
  if (Array.isArray(value)) {
    return value.map(stringify)
  } else {
    return stringify(value)
  }
}

type ScalarParam = string | number | boolean | Date
type ParamType = undefined | ScalarParam | ScalarParam[]

type ParamItem = [string, string]

export function strParamsItems(params: Record<string, ParamType>): ParamItem[] {
  return Object.entries(params)
      .filter(([key, value]) => value != undefined)
      .map(([key, value]): ParamItem[] => {
        if (Array.isArray(value)) {
          return value.map(item => [key, stringify(item)])
        } else {
          return [[key, stringify(value!)]]
        }
      })
      .flat()
}

export function strParamsObject(params: Record<string, ParamType>): Record<string, string | string[]> {
  return Object.keys(params)
      .filter(paramName => params[paramName] != undefined)
      .reduce((obj, paramName) => ({...obj, [paramName]: stringifyX(params[paramName]!)}), {} as Record<string, string | string[]>)
}`

	return &sources.CodeFile{module.GetPath(), strings.TrimSpace(code)}
}
