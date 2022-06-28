package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen/ts/modules"
	"github.com/specgen-io/specgen/v2/gen/ts/validations/common"
	"github.com/specgen-io/specgen/v2/generator"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func paramsTypeName(operation *spec.NamedOperation, namePostfix string) string {
	return fmt.Sprintf("%s%s", operation.Name.PascalCase(), namePostfix)
}

func urlParamsRuntimeType(operation *spec.NamedOperation) string {
	return common.ParamsRuntimeTypeName(paramsTypeName(operation, "UrlParams"))
}

func headersRuntimeType(operation *spec.NamedOperation) string {
	return common.ParamsRuntimeTypeName(paramsTypeName(operation, "HeaderParams"))
}

func queryRuntimeType(operation *spec.NamedOperation) string {
	return common.ParamsRuntimeTypeName(paramsTypeName(operation, "QueryParams"))
}

func generateParamsStaticCode(module modules.Module) *generator.CodeFile {
	code := `
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
}`

	return &generator.CodeFile{module.GetPath(), strings.TrimSpace(code)}
}
