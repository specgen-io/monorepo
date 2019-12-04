package genscala

import (
	"github.com/ModaOperandi/spec"
	"github.com/vsapronov/gopoetry/scala"
)

func responseType(operation spec.NamedOperation) string {
	return operation.Name.PascalCase() + "Response"
}

func generateResponse(responseTypeName string, responses spec.Responses) (*scala.TraitDeclaration, *scala.ClassDeclaration) {
	trait := scala.Trait(responseTypeName).Sealed()
	trait.Define(false).Add("def toResult(): OperationResult")

	object := scala.Object(responseTypeName)
	object_ := object.Define(true)
	for _, response := range responses {
		responseClass := scala.Class(response.Name.PascalCase()).Case()
		responseClassCtor := responseClass.Contructor()
		statusValue := spec.HttpStatusCode(response.Name)
		bodyValue := "None"
		if !response.Type.Definition.IsEmpty() {
			responseClassCtor.Param("body", ScalaType(&response.Type.Definition))
			bodyValue = "Some(Jsoner.write(body))"
		}
		responseClass.Extends(responseTypeName + ` { def toResult = OperationResult(` + statusValue + `, ` + bodyValue + `)}`)
		object_.AddCode(responseClass)
	}

	create := object_.Def("fromResult")
	create.Param("result", "OperationResult")
	match := create.Define().Add("result.status match ").Block(true)
	for _, response := range responses {
		responseParam := ""
		if !response.Type.Definition.IsEmpty() {
			responseParam = `Jsoner.read[` + ScalaType(&response.Type.Definition) + `](result.body.get)`
		}
		match.AddLn("case " + spec.HttpStatusCode(response.Name) + " => " + response.Name.PascalCase() + "(" + responseParam + ")")
	}
	return trait, object
}

func generateApiInterfaceResponse(api spec.Api, apiTraitName string) *scala.ClassDeclaration {
	apiObject := scala.Object(apiTraitName)
	apiObject_ := apiObject.Define(true)

	apiObject_.AddLn("import spec.circe.json._")
	apiObject_.AddLn("implicit val jsonerConfig = Jsoner.config")

	for _, operation := range api.Operations {
		responseTypeName := responseType(operation)
		trait, object := generateResponse(responseTypeName, operation.Responses)
		apiObject_.AddCode(trait)
		apiObject_.AddCode(object)
	}
	return apiObject
}
