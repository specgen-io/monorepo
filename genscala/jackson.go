package genscala

import (
	"fmt"
	"github.com/ModaOperandi/spec"
	"github.com/vsapronov/gopoetry/scala"
	"path/filepath"
	"specgen/gen"
)

func GenerateModels(spec *spec.Spec, packageName string, outPath string) *gen.TextFile {
	unit := scala.Unit(packageName)
	unit.
		Import("com.fasterxml.jackson.databind.JsonNode").
		Import("com.fasterxml.jackson.annotation.JsonProperty").
		Import("spec.enum._")

	for _, model := range spec.Models {
		generateModel(model, unit)
	}

	return &gen.TextFile{
		Path:    filepath.Join(outPath, "Models.scala"),
		Content: unit.Code(),
	}
}

func generateModel(model spec.NamedModel, unit *scala.UnitDeclaration) {
	if model.IsObject() {
		generateObjectModel(model, unit)
	} else {
		generateEnumModel(model, unit)
	}
}

func generateObjectModel(model spec.NamedModel, unit *scala.UnitDeclaration) {
	class := scala.Class(model.Name.PascalCase()).Case()
	ctor := class.Contructor().ParamPerLine()
	for _, field := range model.Object.Fields {
		attribute := fmt.Sprintf("JsonProperty(\"%s\")", field.Name.Source)
		ctor.Param(field.Name.CamelCase(), ScalaType(&field.Type)).Attribute(attribute)
	}
	unit.AddDeclarations(class)
}

func generateEnumModel(model spec.NamedModel, unit *scala.UnitDeclaration) {
	enumTrait := scala.Trait(model.Name.PascalCase()).Sealed().Attribute("StringEnum")

	enumObject := scala.Object(model.Name.PascalCase())
	enumObject_ := enumObject.Define(true)
	for _, item := range model.Enum.Items {
		itemObject := scala.Object(item.Name.PascalCase()).Case().Extends(model.Name.PascalCase())
		itemObject.Define(false).Def("toString").NoParams().Override().Define().Add("\"" + item.Name.Source + "\"")
		enumObject_.AddCode(itemObject)
	}

	unit.AddDeclarations(enumTrait)
	unit.AddDeclarations(enumObject)
}
