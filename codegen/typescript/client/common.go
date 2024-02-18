package client

import (
	"strings"
	"typescript/writer"

	"generator"
	"spec"
	"typescript/validations"
)

func getUrl(endpoint spec.Endpoint) string {
	url := endpoint.Url
	for _, param := range endpoint.UrlParams {
		url = strings.Replace(url, spec.UrlParamStr(&param), "${stringify(parameters."+param.Name.CamelCase()+")}", -1)
	}
	return url
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
