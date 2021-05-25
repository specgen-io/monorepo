package genscala

import (
	spec "github.com/specgen-io/spec"
	"github.com/vsapronov/gopoetry/scala"
)

func responseType(operation spec.NamedOperation) string {
	return operation.Name.PascalCase() + "Response"
}

func generateResponse(responseTypeName string, responses spec.Responses) (*scala.TraitDeclaration, *scala.ClassDeclaration) {
	trait := scala.Trait(responseTypeName).Sealed()

	object := scala.Object(responseTypeName)

	for _, response := range responses {
		var bodyParam scala.Writable = nil
		if !response.Type.Definition.IsEmpty() {
			bodyParam = Param(`body`, ScalaType(&response.Type.Definition))
		}
		responseClass :=
			CaseClass(response.Name.PascalCase()).Extends(responseTypeName).Constructor(Constructor().AddParams(bodyParam))

		object.Add(responseClass)
	}

	return trait, object
}

func generateApiInterfaceResponse(api spec.Api, apiTraitName string) *scala.ClassDeclaration {
	apiObject := Object(apiTraitName)
	for _, operation := range api.Operations {
		responseTypeName := responseType(operation)
		trait, object := generateResponse(responseTypeName, operation.Responses)
		apiObject.Add(trait, object)
	}
	return apiObject
}
