package responses

import (
	"fmt"

	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/spec/v2"
	"github.com/specgen-io/specgen/typescript/v2/types"
)

func GenerateOperationResponse(w *generator.Writer, operation *spec.NamedOperation) {
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

func New(response *spec.Response, body string) string {
	if body == `` {
		return fmt.Sprintf(`{ status: "%s" }`, response.Name.Source)
	} else {
		return fmt.Sprintf(`{ status: "%s", data: %s }`, response.Name.Source, body)
	}
}

func responseTypeName(operation *spec.NamedOperation) string {
	return operation.Name.PascalCase() + "Response"
}
