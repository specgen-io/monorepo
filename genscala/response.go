package genscala

import (
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
)

func responseType(operation spec.NamedOperation) string {
	return operation.Name.PascalCase() + "Response"
}

func generateResponse(w *gen.Writer, responseTypeName string, responses spec.Responses) {
	w.Line(`sealed trait %s`, responseTypeName)
	w.Line(`object %s {`, responseTypeName)
	for _, response := range responses {
		var bodyParam = ""
		if !response.Type.Definition.IsEmpty() {
			bodyParam = `body: ` + ScalaType(&response.Type.Definition)
		}
		w.Line(`  case class %s(%s) extends %s`, response.Name.PascalCase(), bodyParam, responseTypeName)
	}
	w.Line(`}`)
}

func generateApiInterfaceResponse(w *gen.Writer, api *spec.Api, apiTraitName string) {
	w.Line(`object %s {`, apiTraitName)
	for _, operation := range api.Operations {
		responseTypeName := responseType(operation)
		generateResponse(w.Indented(), responseTypeName, operation.Responses)
	}
	w.Line(`}`)
}
