package responses

import (
	"github.com/specgen-io/specgen/v2/gents/types"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateOperationResponse(w *sources.Writer, operation *spec.NamedOperation) {
	w.Line("export type %s =", responseTypeName(operation))
	for _, response := range operation.Responses {
		if !response.Type.Definition.IsEmpty() {
			w.Line(`  | { status: "%s", data: %s }`, response.Name.Source, types.TsType(&response.Type.Definition))
		} else {
			w.Line(`  | { status: "%s" }`, response.Name.Source)
		}
	}
}

func ResponseType(operation *spec.NamedOperation, servicePackage string) string {
	if len(operation.Responses) == 1 {
		response := operation.Responses[0]
		if response.Definition.Type.Definition.IsEmpty() {
			return "void"
		}
		return types.TsType(&response.Definition.Type.Definition)
	}
	result := responseTypeName(operation)
	if servicePackage != "" {
		result = "service." + result
	}
	return result
}

func responseTypeName(operation *spec.NamedOperation) string {
	return operation.Name.PascalCase() + "Response"
}
