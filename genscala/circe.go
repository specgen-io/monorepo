package genscala

import (
	"github.com/ModaOperandi/spec"
	"github.com/vsapronov/gopoetry/scala"
	"path/filepath"
	"specgen/gen"
)

func GenerateCirceModels(spec *spec.Spec, packageName string, outPath string) *gen.TextFile {
	unit := Unit(packageName)
	unit.
		Import("enumeratum.values._").
		Import("java.time._").
		Import("java.time.format._").
		Import("java.util.UUID").
		Import("io.circe.generic.extras.{AutoDerivation, Configuration}").
		Import("io.circe.{Decoder, Encoder}").
		Import("io.circe.generic.extras.semiauto.{deriveUnwrappedDecoder, deriveUnwrappedEncoder}")

	for _, model := range spec.Models {
		if model.IsObject() {
			class := generateCirceObjectModel(model)
			unit.AddDeclarations(class)
		} else if model.IsEnum() {
			enumClass, enumObject := generateCirceEnumModel(model)
			unit.AddDeclarations(enumClass, enumObject)
		} else if model.IsOneOf() {
			trait, object := generateCirceUnionModel(model, packageName)
			unit.AddDeclarations(trait)
			unit.AddDeclarations(object)
		}
	}

	unit.AddDeclarations(generateCirceUnionItemsCodecs(spec.Models))

	return &gen.TextFile{
		Path:    filepath.Join(outPath, "Models.scala"),
		Content: unit.Code(),
	}
}

func generateCirceObjectModel(model spec.NamedModel) *scala.ClassDeclaration {
	ctor := Constructor().ParamPerLine()
	for _, field := range model.Object.Fields {
		ctor.Param(field.Name.CamelCase(), ScalaType(&field.Type.Definition))
	}
	modelClass := CaseClass(model.Name.PascalCase()).Constructor(ctor)
	return modelClass
}

func generateCirceEnumModel(model spec.NamedModel) (scala.Writable, scala.Writable) {
	className := model.Name.PascalCase()

	enumClass :=
		Class(className).Sealed().Abstract().Extends("StringEnumEntry").
			Constructor(Constructor().
				Val("value", "String"),
			)

	enumObject :=
		CaseObject(className).Extends("StringEnum["+className+"]").With("StringCirceEnum["+className+"]")

	for _, item := range model.Enum.Items {
		enumObject.Add(
			CaseObject(item.Name.PascalCase()).Extends(className + `("` + item.Value + `")`),
		)
	}
	enumObject.Add(Line("val values = findValues"))

	return enumClass, enumObject
}

func generateCirceUnionModel(model spec.NamedModel, packageName string) (scala.Writable, scala.Writable) {
	trait := Trait(model.Name.PascalCase()).Sealed()

	object := Object(model.Name.PascalCase())
	for _, item := range model.OneOf.Items {
		itemClass :=
			CaseClass(item.Name.PascalCase()).Extends(model.Name.PascalCase()).
				Constructor(Constructor().
					Param("data", packageName+"."+ScalaType(&item.Type.Definition)),
				)
		object.Add(itemClass)
	}

	return trait, object
}

func generateCirceUnionItemsCodecs(models spec.Models) *scala.ClassDeclaration {
	object :=
		Object("json").Extends("AutoDerivation").
			Add(Line("implicit val auto = Configuration.default.withSnakeCaseMemberNames.withSnakeCaseConstructorNames.withDefaults"))
	for _, model := range models {
		if model.IsOneOf() {
			for _, item := range model.OneOf.Items {
				itemTypeName := model.Name.PascalCase()+"."+item.Name.PascalCase()
				itemCodecName := model.Name.PascalCase()+item.Name.PascalCase()
				object.Add(
					Line("implicit val encoder%s: Encoder[%s] = deriveUnwrappedEncoder", itemCodecName, itemTypeName),
					Line("implicit val decoder%s: Decoder[%s] = deriveUnwrappedDecoder", itemCodecName, itemTypeName),
				)
			}
		}
	}
	return object
}