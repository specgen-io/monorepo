package gents

import (
	spec "github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
)

func generateOperationResponse(w *gen.Writer, operation *spec.NamedOperation, modelsPackage *string) {
	w.Line("export type %s =", responseTypeName(operation))
	for _, response := range operation.Responses {
		if !response.Type.Definition.IsEmpty() {
			w.Line(`  | { status: "%s", data: %s }`, response.Name.Source, PackagedTsType(&response.Type.Definition, modelsPackage))
		} else {
			w.Line(`  | { status: "%s" }`, response.Name.Source)
		}
	}
}