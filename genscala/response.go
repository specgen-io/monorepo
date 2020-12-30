package genscala

import (
	"fmt"
	"github.com/ModaOperandi/spec"
	"github.com/vsapronov/gopoetry/scala"
)

func responseType(operation spec.NamedOperation) string {
	return operation.Name.PascalCase() + "Response"
}

func generateResponseCases(responses spec.Responses) *scala.StatementsDeclaration {
	cases := scala.Statements()
	for _, response := range responses {
		responseParam := ``
		if !response.Type.Definition.IsEmpty() {
			responseParam = fmt.Sprintf(`Jsoner.read[%s](result.body.get)`, ScalaType(&response.Type.Definition))
		}
		cases.Add(Line(`case %s => %s(%s)`, spec.HttpStatusCode(response.Name), response.Name.PascalCase(), responseParam))
	}
	return cases
}

func generateResponse(responseTypeName string, responses spec.Responses) (*scala.TraitDeclaration, *scala.ClassDeclaration) {
	trait := scala.Trait(responseTypeName).Sealed().MembersInline(scala.Code("def toResult(): OperationResult"))

	object := scala.Object(responseTypeName)

	for _, response := range responses {
		statusValue := spec.HttpStatusCode(response.Name)
		bodyValue := `None`
		if !response.Type.Definition.IsEmpty() {
			bodyValue = `Some(Jsoner.write(body))`
		}
		baseType := fmt.Sprintf(`%s { def toResult = OperationResult(%s, %s)}`, responseTypeName, statusValue, bodyValue)

		var bodyParam scala.Writable = nil
		if !response.Type.Definition.IsEmpty() {
			bodyParam = Param(`body`, ScalaType(&response.Type.Definition))
		}
		responseClass :=
			CaseClass(response.Name.PascalCase()).Extends(baseType).Constructor(Constructor().AddParams(bodyParam))

		object.Add(responseClass)
	}

	create :=
		scala.Def(`fromResult`).Param(`result`, `OperationResult`).
			BodyInline(
				Code(`result.status match `),
				Scope(generateResponseCases(responses)),
			)

	object.Add(create)
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
