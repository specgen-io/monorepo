package gents

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
)

func paramsRuntimeTypeName(typeName string) string {
	return fmt.Sprintf("T%s", typeName)
}

func paramsTypeName(operation *spec.NamedOperation, namePostfix string) string {
	return fmt.Sprintf("%s%s", operation.Name.PascalCase(), namePostfix)
}

func generateParams(w *gen.Writer, typeName string, isHeader bool, params []spec.NamedParam, validation string) {
	if validation == Superstruct {
		generateSuperstructParams(w, typeName, isHeader, params)
		return
	}
	if validation == IoTs {
		generateIoTsParams(w, typeName, isHeader, params)
		return
	}
	panic(fmt.Sprintf("Unknown validation: %s", validation))
}
