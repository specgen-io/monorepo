package gents

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"strings"
)

func stringParamsRuntimeTypeName(operationName *spec.Name, namePostfix string) string {
	return fmt.Sprintf("T%s%s", operationName.PascalCase(), namePostfix)
}

func stringParamsTypeName(operationName *spec.Name, namePostfix string) string {
	return fmt.Sprintf("%s%s", operationName.PascalCase(), namePostfix)
}

func generateParams(w *gen.Writer, operationName *spec.Name, namePostfix string, isHeader bool, params []spec.NamedParam) {
	if len(params) > 0 {
		w.EmptyLine()
		w.Line("const %s = t.type({", stringParamsRuntimeTypeName(operationName, namePostfix))
		for _, param := range params {
			paramName := param.Name.Source
			if isHeader {
				paramName = isTsIdentifier(strings.ToLower(param.Name.Source))
			}
			w.Line("  %s: %s,", paramName, StringSsTypeDefaulted(&param))
		}
		w.Line("})")
		w.EmptyLine()
		w.Line("type %s = t.Infer<typeof %s>", stringParamsTypeName(operationName, namePostfix), stringParamsRuntimeTypeName(operationName, namePostfix))
	}
}
