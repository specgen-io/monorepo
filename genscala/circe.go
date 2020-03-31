package genscala

import (
	"github.com/ModaOperandi/spec"
	"github.com/vsapronov/gopoetry/scala"
	"fmt"
	"path/filepath"
	"specgen/gen"
)

func GenerateCirceModels(spec *spec.Spec, packageName string, outPath string) *gen.TextFile {
	unit := scala.Unit(packageName)
	unit.
		Import("enumeratum.values._").
		Import("java.time._").
		Import("java.time.format._").
		Import("java.util.UUID")

	for _, model := range spec.Models {
		if model.IsObject() {
			generateCirceObjectModel(model, unit)
		} else if model.IsEnum() {
			generateCirceEnumModel(model, unit)
		} else if model.IsUnion() {
			generateCirceUnionModel(model, unit, packageName)
		}
	}

	generateCirceUnionItemsCodecs(unit, spec.Models, )

	return &gen.TextFile{
		Path:    filepath.Join(outPath, "Models.scala"),
		Content: unit.Code(),
	}
}

func generateCirceObjectModel(model spec.NamedModel, unit *scala.UnitDeclaration) {
	class := scala.Class(model.Name.PascalCase()).Case()

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

func generateCirceUnionModel(model spec.NamedModel, unit *scala.UnitDeclaration, packageName string) {
	trait := scala.Trait(model.Name.PascalCase()).Sealed()

	object := scala.Object(model.Name.PascalCase())
	objectDefinition := object.Define(true)
	for _, item := range model.Union.Items {
		itemClass := scala.Class(item.Name.PascalCase()).Case()
		itemClass.Extends(model.Name.PascalCase())
		itemClass.Contructor().Param("data", packageName+"."+ScalaType(&item.Type.Definition))
		objectDefinition.AddCode(itemClass)
	}
	unit.AddDeclarations(trait)
	unit.AddDeclarations(object)
}

func generateCirceUnionItemsCodecs(unit *scala.UnitDeclaration, models spec.Models) {
	unit.
		Import("io.circe.generic.extras.{AutoDerivation, Configuration}").
		Import("io.circe.{Decoder, Encoder}").
		Import("io.circe.generic.extras.semiauto.{deriveUnwrappedDecoder, deriveUnwrappedEncoder}")

	object := scala.Object("json")
	object.Extends("AutoDerivation")
	objectDefinition := object.Define(true)
	objectDefinition.AddLn("implicit val auto = Configuration.default.withSnakeCaseMemberNames.withSnakeCaseConstructorNames.withDefaults")
	for _, model := range models {
		if model.IsUnion() {
			for _, item := range model.Union.Items {
				itemTypeName := model.Name.PascalCase()+"."+item.Name.PascalCase()
				itemCodecName := model.Name.PascalCase()+item.Name.PascalCase()
				objectDefinition.AddLn(fmt.Sprintf("implicit val encoder%s: Encoder[%s] = deriveUnwrappedEncoder", itemCodecName, itemTypeName))
				objectDefinition.AddLn(fmt.Sprintf("implicit val decoder%s: Decoder[%s] = deriveUnwrappedDecoder", itemCodecName, itemTypeName))
			}
		}
	}
	unit.AddDeclarations(object)
}