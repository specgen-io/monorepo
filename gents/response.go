package gents

import (
	"github.com/specgen-io/specgen/v2/spec"
	"github.com/specgen-io/specgen/v2/gen"
)

func generateOperationResponse(w *gen.Writer, operation *spec.NamedOperation) {
	w.Line("export type %s =", responseTypeName(operation))
	for _, response := range operation.Responses {
		if !response.Type.Definition.IsEmpty() {
			w.Line(`  | { status: "%s", data: %s }`, response.Name.Source, TsType(&response.Type.Definition))
		} else {
			w.Line(`  | { status: "%s" }`, response.Name.Source)
		}
	}
}
