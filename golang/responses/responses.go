package responses

import (
	"fmt"
	"generator"
	"golang/types"
	"golang/writer"
	"spec"
)

func NewErrorResponse(response spec.Response, body string) string {
	return fmt.Sprintf(`&httperrors.%s{Body: %s}`, response.Name.PascalCase(), body)
}

func NewResponse(response *spec.OperationResponse, body string) string {
	return fmt.Sprintf(`%s{%s: &%s}`, ResponseTypeName(response.Operation), response.Name.PascalCase(), body)
}

func GenerateOperationResponseStruct(w *generator.Writer, operation *spec.NamedOperation) {
	w.Line(`type %s struct {`, ResponseTypeName(operation))
	responses := [][]string{}
	for _, response := range operation.Responses {
		responses = append(responses, []string{
			response.Name.PascalCase(),
			types.GoType(spec.Nullable(&response.Type.Definition)),
		})
	}
	writer.WriteAlignedLines(w.Indented(), responses)
	w.Line(`}`)
}

func ResponseTypeName(operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%sResponse`, operation.Name.PascalCase())
}
