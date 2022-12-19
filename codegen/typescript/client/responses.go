package client

import (
	"fmt"
	"typescript/writer"

	"spec"
	"typescript/types"
)

func generateOperationResponse(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line("export type %s =", responseTypeName(operation))
	for _, response := range operation.Responses {
		if !response.Type.Definition.IsEmpty() {
			w.Line(`  | { status: "%s", data: %s }`, response.Name.Source, types.TsType(&response.Type.Definition))
		} else {
			w.Line(`  | { status: "%s" }`, response.Name.Source)
		}
	}
}

func responseType(operation *spec.NamedOperation) string {
	successResponses := operation.Responses.Success()
	if len(successResponses) == 1 {
		if successResponses[0].Definition.Type.Definition.IsEmpty() {
			return "void"
		} else {
			return types.TsType(&successResponses[0].Definition.Type.Definition)
		}
	} else {
		return responseTypeName(operation)
	}
}

func newResponse(response *spec.Response, body string) string {
	if body == `` {
		return fmt.Sprintf(`{ status: "%s" }`, response.Name.Source)
	} else {
		return fmt.Sprintf(`{ status: "%s", data: %s }`, response.Name.Source, body)
	}
}

func responseTypeName(operation *spec.NamedOperation) string {
	return operation.Name.PascalCase() + "Response"
}
