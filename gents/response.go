package gents

import (
	spec "github.com/specgen-io/spec"
	"specgen/gen"
)

func generateIoTsResponse(w *gen.Writer, operation *spec.NamedOperation) {
	w.Line("")
	w.Line("export type %s =", responseTypeName(operation))
	for _, response := range operation.Responses {
		if !response.Type.Definition.IsEmpty() {
			w.Line(`  | { status: "%s", data: %s }`, response.Name.Source, TsType(&response.Type.Definition))
		} else {
			w.Line(`  | { status: "%s" }`, response.Name.Source)
		}
	}
}