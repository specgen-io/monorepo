package genscala

import (
	"github.com/ModaOperandi/spec"
	"github.com/vsapronov/gopoetry/scala"
	"path/filepath"
	"specgen/gen"
)

func buildReversedUnionMap(specification *spec.Spec) map[spec.Name][]spec.Name {
	result := map[spec.Name][]spec.Name {}
	for _, model := range specification.Models {
		if model.IsUnion() {
			for _, unionItem := range model.Union.Items {
				if _, ok := result[unionItem.Definition.Info.Model.Name]; !ok {
					result[unionItem.Definition.Info.Model.Name] = []spec.Name{}
				}
				result[unionItem.Definition.Info.Model.Name] = append(result[unionItem.Definition.Info.Model.Name], model.Name)
			}
		}
	}
	return result
}

func GenerateCirceModels(spec *spec.Spec, packageName string, outPath string) *gen.TextFile {
	unit := scala.Unit(packageName)
	unit.
		Import("enumeratum.values._").
		Import("java.time._").
		Import("java.time.format._").
		Import("java.util.UUID")

	unionsMap := buildReversedUnionMap(spec)

	for _, model := range spec.Models {
		if model.IsObject() {
			generateCirceObjectModel(model, unionsMap, unit)
		} else if model.IsEnum() {
			generateCirceEnumModel(model, unit)
		} else if model.IsUnion() {
			generateCirceUnionModel(model, unit)
		}
	}

	return &gen.TextFile{
		Path:    filepath.Join(outPath, "Models.scala"),
		Content: unit.Code(),
	}
}

func generateCirceObjectModel(model spec.NamedModel, unionMap map[spec.Name][]spec.Name, unit *scala.UnitDeclaration) {
	class := scala.Class(model.Name.PascalCase()).Case()

	if unions, ok := unionMap[model.Name]; ok {
		extends := class.Extends(unions[0].Source)
		for index := 1; index < len(unions); index ++ {
			extends.With(unions[index].Source)
		}
	}

	ctor := class.Contructor().ParamPerLine()
	for _, field := range model.Object.Fields {
		ctor.Param(field.Name.CamelCase(), ScalaType(&field.Type.Definition))
	}
	unit.AddDeclarations(class)
}

func generateCirceEnumModel(model spec.NamedModel, unit *scala.UnitDeclaration) {
	enumBase := scala.Class(model.Name.PascalCase()).Sealed().Abstract()
	enumBase.Extends("StringEnumEntry")
	enumBaseCtor := enumBase.Contructor()
	enumBaseCtor.Param("value", "String").Val()

	enumObject := scala.Object(model.Name.PascalCase()).Case()
	enumObject.Extends("StringEnum["+model.Name.PascalCase()+"]").With("StringCirceEnum["+model.Name.PascalCase()+"]")
	enumObject_ := enumObject.Define(true)
	for _, item := range model.Enum.Items {
		itemObject := scala.Object(item.Name.PascalCase()).Case()
		itemObject.Extends(model.Name.PascalCase() + `("` + item.Name.Source + `")`)
		enumObject_.AddCode(itemObject)
	}
	enumObject_.AddLn("val values = findValues")

	unit.AddDeclarations(enumBase)
	unit.AddDeclarations(enumObject)
}

func generateCirceUnionModel(model spec.NamedModel, unit *scala.UnitDeclaration) {
	trait := scala.Trait(model.Name.PascalCase()).Sealed()
	unit.AddDeclarations(trait)
}
