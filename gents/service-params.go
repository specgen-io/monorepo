package gents

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func paramsRuntimeTypeName(typeName string) string {
	return fmt.Sprintf("T%s", typeName)
}

func paramsTypeName(operation *spec.NamedOperation, namePostfix string) string {
	return fmt.Sprintf("%s%s", operation.Name.PascalCase(), namePostfix)
}

func generateParams(w *gen.Writer, typeName string, params []spec.NamedParam, validation string) {
	if validation == Superstruct {
		generateSuperstructParams(w, typeName, params)
		return
	}
	if validation == IoTs {
		generateIoTsParams(w, typeName, params)
		return
	}
	panic(fmt.Sprintf("Unknown validation: %s", validation))
}

func generateParamsStaticCode(module module) *gen.TextFile {
	code := `
export function zipHeaders(headers: string[]): Record<string, string | string[]> {
  let result: Record<string, string | string[]> = {}

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

	return &gen.TextFile{module.GetPath(), strings.TrimSpace(code)}
}
