package genscala

import (
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
)

func responseType(operation *spec.NamedOperation) string {
	if len(operation.Responses) == 1 {
		response := operation.Responses[0]
		return ScalaType(&response.Type.Definition)
	} else {
		return responseTypeName(operation)
	}
}

func responseTypeName(operation *spec.NamedOperation) string {
	return operation.Name.PascalCase() + "Response"
}

func generateResponse(w *gen.Writer, operation *spec.NamedOperation) {
	if len(operation.Responses) > 1 {
		w.Line(`sealed trait %s`, responseTypeName(operation))
		w.Line(`object %s {`, responseTypeName(operation))
		for _, response := range operation.Responses {
			var bodyParam = ""
			if !response.Type.Definition.IsEmpty() {
				bodyParam = `body: ` + ScalaType(&response.Type.Definition)
			}
			w.Line(`  case class %s(%s) extends %s`, response.Name.PascalCase(), bodyParam, responseTypeName(operation))
		}
		w.Line(`}`)
	}
}
