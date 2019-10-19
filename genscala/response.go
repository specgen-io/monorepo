package genscala

import (
	"github.com/ModaOperandi/spec"
	"github.com/vsapronov/gopoetry/scala"
)

func responseType(operation spec.NamedOperation) string {
	return operation.Name.PascalCase() + "Response"
}

func generateResponse(responseTypeName string, responses spec.Responses) (*scala.ClassDeclaration, *scala.ClassDeclaration) {
	class := scala.Class(responseTypeName).Sealed()
	ctor := class.Contructor().Private()
	ctor.Param("status", "Int")
	ctor.Param("body", "Option[Any]")
	class_ := class.Define(true)
	class_.AddLn("def toResult()(implicit json: JsonMapper) = OperationResult(status, body.map(json.writeValueAsString(_)))")

	object := scala.Object(responseTypeName)
	object_ := object.Define(true)
	for _, response := range responses {
		responseClass := scala.Class(response.Name.PascalCase()).Case()
		responseClassCtor := responseClass.Contructor()
		statusValue := spec.HttpStatusCode(response.Name)
		bodyValue := "None"
		if !response.Type.IsEmpty() {
			responseClassCtor.Param("body", ScalaType(&response.Type))
			bodyValue = "Some(body)"
		}
		responseClass.Extends(responseTypeName + "(" + statusValue + ", " + bodyValue + ")")
		object_.AddCode(responseClass)
	}

	create := object_.Def("fromResult")
	create.Param("result", "OperationResult")
	create.ImplicitParam("json", "JsonMapper")
	match := create.Define().Add("result.status match ").Block(true)
	for _, response := range responses {
		responseParam := ""
		if !response.Type.IsEmpty() {
			responseParam = "json.readValue(result.body.get, classOf[" + ScalaType(&response.Type) + "])"
		}
		match.AddLn("case " + spec.HttpStatusCode(response.Name) + " => " + response.Name.PascalCase() + "(" + responseParam + ")")
	}
	return class, object
}

func generateApiInterfaceResponse(api spec.Api, apiTraitName string) *scala.ClassDeclaration {
	apiObject := scala.Object(apiTraitName)
	apiObject_ := apiObject.Define(true)

	for _, operation := range api.Operations {
		responseTypeName := responseType(operation)
		class, object := generateResponse(responseTypeName, operation.Responses)
		apiObject_.AddCode(class)
		apiObject_.AddCode(object)
	}
	return apiObject
}
