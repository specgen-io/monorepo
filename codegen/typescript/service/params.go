package service

import (
	"fmt"
	"typescript/writer"

	"generator"
	"spec"
)

func paramsTypeName(operation *spec.NamedOperation, namePostfix string) string {
	return fmt.Sprintf("%s%s", operation.Name.PascalCase(), namePostfix)
}

func urlParamsType(operation *spec.NamedOperation) string {
	return paramsTypeName(operation, "UrlParams")
}

func headersType(operation *spec.NamedOperation) string {
	return paramsTypeName(operation, "HeaderParams")
}

func queryType(operation *spec.NamedOperation) string {
	return paramsTypeName(operation, "QueryParams")
}

func (g *Generator) ParamsStaticCode() *generator.CodeFile {
	w := writer.New(g.Modules.Params)
	w.Lines(`
export function zipHeaders(headers: string[]): Record<string, string | string[]> {
  const result: Record<string, string | string[]> = {}

  for (let i = 0; i < headers.length / 2; i++) {
      const key: string = headers[i*2]
      const value: string = headers[i*2+1]

      if (key in result) {
          const existingValue = result[key]
          if (Array.isArray(existingValue)) {
              existingValue.push(value)
          } else {
              result[key] = [existingValue, value]                
          }
      } else {
          result[key] = value
      }
  }
  return result
}`)
	return w.ToCodeFile()
}
