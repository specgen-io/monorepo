package genscala

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/vsapronov/gopoetry/scala"
	"path/filepath"
	"strings"
)

func GenerateCirceModels(spec *spec.Spec, packageName string, outPath string) []gen.TextFile {
	files := []gen.TextFile{}
	for _, version := range spec.Versions {
		versionFile := generateCirceModels(&version, packageName, outPath)
		files = append(files, *versionFile)
	}
	return files
}

func generateCirceModels(version *spec.Version, packageName string, outPath string) *gen.TextFile {
	if version.Version.Source != "" {
		packageName = fmt.Sprintf("%s.%s", packageName, version.Version.FlatCase())
	}
	unit := Unit(packageName)
	unit.
		Import("enumeratum.values._").
		Import("java.time._").
		Import("java.time.format._").
		Import("java.util.UUID").
		Import("io.circe.Codec").
		Import("io.circe.generic.extras.{Configuration, JsonKey}").
		Import("io.circe.generic.extras.semiauto.{deriveConfiguredCodec, deriveUnwrappedCodec}")

	for _, model := range version.ResolvedModels {
		if model.IsObject() {
			class, object := generateCirceObjectModel(model)
			unit.AddDeclarations(class, object)
		} else if model.IsEnum() {
			enumClass, enumObject := generateCirceEnumModel(model)
			unit.AddDeclarations(enumClass, enumObject)
		} else if model.IsOneOf() {
			trait, object := generateCirceUnionModel(model, packageName)
			unit.AddDeclarations(trait)
			unit.AddDeclarations(object)
		}
	}

	return &gen.TextFile{
		Path:    filepath.Join(outPath, fmt.Sprintf("%sModels.scala", version.Version.PascalCase())),
		Content: unit.Code(),
	}
}

func generateCirceObjectModel(model *spec.NamedModel) (scala.Writable, scala.Writable) {
	ctor := Constructor().ParamPerLine()
	for _, field := range model.Object.Fields {
		jsonKey := fmt.Sprintf(`JsonKey("%s")`, field.Name.Source)
		fieldParam := Param(scalaCamelCase(field.Name), ScalaType(&field.Type.Definition)).Attribute(jsonKey)
		ctor.AddParams(fieldParam)
	}
	modelClass := CaseClass(model.Name.PascalCase()).Constructor(ctor)

	modelObject := Object(model.Name.PascalCase())
	modelObject.Add(
		Line("implicit val config = Configuration.default"),
		Line("implicit val codec: Codec[%s] = deriveConfiguredCodec", model.Name.PascalCase()),
	)
	return modelClass, modelObject
}

func generateCirceEnumModel(model *spec.NamedModel) (scala.Writable, scala.Writable) {
	className := model.Name.PascalCase()

	enumClass :=
		Class(className).Sealed().Abstract().Extends("StringEnumEntry").
			Constructor(Constructor().
				Val("value", "String"),
			)

	enumObject :=
		CaseObject(className).Extends("StringEnum[" + className + "]").With("StringCirceEnum[" + className + "]")

	for _, item := range model.Enum.Items {
		enumObject.Add(
			CaseObject(item.Name.PascalCase()).Extends(className + `("` + item.Value + `")`),
		)
	}
	enumObject.Add(Line("val values = findValues"))

	return enumClass, enumObject
}

func generateCirceUnionModel(model *spec.NamedModel, packageName string) (scala.Writable, scala.Writable) {
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

	object.Add(
		Line("implicit val config = Configuration.default.withSnakeCaseConstructorNames"),
		Line("implicit val codec: Codec[%s] = deriveConfiguredCodec", model.Name.PascalCase()),
	)

	for _, item := range model.OneOf.Items {
		itemTypeName := model.Name.PascalCase() + "." + item.Name.PascalCase()
		object.Add(
			Line("implicit val codec%s: Codec[%s] = deriveUnwrappedCodec", item.Name.PascalCase(), itemTypeName),
		)
	}

	return trait, object
}

func generateJson(packageName string, path string) *gen.TextFile {
	code := `
package [[.PackageName]]

object Jsoner {

  import io.circe._
  import io.circe.syntax._
  import io.circe.parser._

  def readThrowing[T](jsonStr: String)(implicit decoder: Decoder[T]): T = {
    parse(jsonStr) match {
      case Right(parsed) => {
        parsed.as[T] match {
          case Right(decoded) => decoded
          case Left(failure) => throw failure
        }
      }
      case Left(failure) => throw failure
    }
  }

  def read[T](jsonStr: String)(implicit decoder: Decoder[T]): Either[Error, T] =
    parse(jsonStr).flatMap(json => json.as[T])

  def write[T](value: T, formatted: Boolean = false)(implicit encoder: Encoder[T]): String = {
    if (formatted) {
      Printer.spaces2.copy(dropNullValues = true).print(value.asJson)
    } else {
      Printer.noSpaces.copy(dropNullValues = true).print(value.asJson)
    }
  }
}`
	code, _ = gen.ExecuteTemplate(code, struct{ PackageName string }{packageName})
	return &gen.TextFile{path, strings.TrimSpace(code)}
}
