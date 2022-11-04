package client

import (
	"fmt"
	"strings"

	"generator"
	"spec"
	"typescript/module"
	"typescript/types"
	"typescript/validations"
)

func getUrl(endpoint spec.Endpoint) string {
	url := endpoint.Url
	for _, param := range endpoint.UrlParams {
		url = strings.Replace(url, spec.UrlParamStr(&param), "${stringify(parameters."+param.Name.CamelCase()+")}", -1)
	}
	return url
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
			tsParams = append(tsParams, param.Name.CamelCase()+requiredSign+": "+types.TsType(paramType))
		}
	}
	return tsParams
}

func createOperationParams(operation *spec.NamedOperation) string {
	operationParams := []string{}
	operationParams = append(operationParams, createParams(operation.HeaderParams, true)...)
	if operation.Body != nil {
		operationParams = append(operationParams, "body: "+types.TsType(&operation.Body.Type.Definition))
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

func clientResponseBody(validation validations.Validation, response *spec.Response, textResponseData, jsonResponseData string) string {
	if response.BodyIs(spec.BodyEmpty) {
		return ""
	}
	if response.BodyIs(spec.BodyString) {
		return textResponseData
	} else {
		data := fmt.Sprintf(`t.decode(%s, %s)`, validation.RuntimeType(&response.Type.Definition), jsonResponseData)
		return data
	}
}

func generateParamsBuilder(module module.Module) *generator.CodeFile {
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

	return &generator.CodeFile{module.GetPath(), strings.TrimSpace(code)}
}
