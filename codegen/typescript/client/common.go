package client

import (
	"fmt"
	"strings"
	"typescript/writer"

	"generator"
	"spec"
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

func operationSignature(operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%s: async (%s): Promise<%s>`, operation.Name.CamelCase(), createOperationParams(operation), responseType(operation))
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
	if operation.Body.IsText() || operation.Body.IsJson() {
		operationParams = append(operationParams, "body: "+types.RequestBodyTsType(&operation.Body))
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

type CommonGenerator struct {
	Modules    *Modules
	validation validations.Validation
}

func (g *CommonGenerator) ParamsBuilder() *generator.CodeFile {
	w := writer.New(g.Modules.Params)
	w.Lines(`
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
}`)
	return w.ToCodeFile()
}
