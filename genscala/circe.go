package genscala

import (
	"github.com/ModaOperandi/spec"
	"github.com/vsapronov/gopoetry/scala"
	"path/filepath"
	"specgen/gen"
)

func GenerateCirceModels(spec *spec.Spec, packageName string, outPath string) *gen.TextFile {
	modelsMap := buildModelsMap(spec.Models)

	unit := scala.Unit(packageName)
	unit.
		Import("enumeratum.values._").
		Import("java.time._").
		Import("java.time.format._").
		Import("java.util.UUID")

	for _, model := range spec.Models {
		if model.IsObject() {
			generateCirceObjectModel(model, modelsMap, unit)
		} else {
			generateCirceEnumModel(model, unit)
		}
	}

	return &gen.TextFile{
		Path:    filepath.Join(outPath, "Models.scala"),
		Content: unit.Code(),
	}
}

func generateCirceObjectModel(model spec.NamedModel, modelsMap ModelsMap, unit *scala.UnitDeclaration) {
	class := scala.Class(model.Name.PascalCase()).Case()
	ctor := class.Contructor().ParamPerLine()
	for _, field := range model.Object.Fields {
		param := ctor.Param(field.Name.CamelCase(), ScalaType(&field.Type))
		if field.Default != nil {
			param.Init(scala.Code(DefaultValue(&field.Type, *field.Default, modelsMap)))
		}
	}
	unit.AddDeclarations(class)
}

func generateCirceEnumModel(model spec.NamedModel, unit *scala.UnitDeclaration) {
	enumBase := scala.Class(model.Name.PascalCase()).Sealed().Abstract().Extends("StringEnumEntry")
	enumBaseCtor := enumBase.Contructor()
	enumBaseCtor.Param("value", "String").Val()

	enumObject := scala.Object(model.Name.PascalCase()).Case().Extends("StringEnum["+model.Name.PascalCase()+"]", "StringCirceEnum["+model.Name.PascalCase()+"]")
	enumObject_ := enumObject.Define(true)
	for _, item := range model.Enum.Items {
		itemObject := scala.Object(item.Name.PascalCase()).Case().Extends(model.Name.PascalCase() + `("` + item.Name.Source + `")`)
		enumObject_.AddCode(itemObject)
	}
	enumObject_.AddLn("val values = findValues")

	unit.AddDeclarations(enumBase)
	unit.AddDeclarations(enumObject)
}
