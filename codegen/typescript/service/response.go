package service

import (
	"typescript/writer"

	"spec"
	"typescript/types"
)

func GenerateOperationResponse(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line("export type %s =", responseTypeName(operation))
	for _, response := range operation.Responses {
		if !response.Body.IsEmpty() {
			w.Line(`  | { status: "%s", data: %s }`, response.Name.Source, types.ResponseBodyTsType(&response.Body))
		} else {
			w.Line(`  | { status: "%s" }`, response.Name.Source)
		}
	}
}

func ResponseType(operation *spec.NamedOperation, servicePackage string) string {
	if len(operation.Responses) == 1 {
		response := operation.Responses[0]
		if response.Body.IsEmpty() {
			return "void"
		}
		return types.ResponseBodyTsType(&response.Body)
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
