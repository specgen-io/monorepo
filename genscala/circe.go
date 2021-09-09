package genscala

import (
	"fmt"
	"github.com/specgen-io/spec"
	"specgen/gen"
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
		packageName = fmt.Sprintf(`%s.%s`, packageName, version.Version.FlatCase())
	}

	w := NewScalaWriter()
	w.Line(`package %s`, packageName)
	w.EmptyLine()
	w.Line(`import enumeratum.values._`)
	w.Line(`import java.time._`)
	w.Line(`import java.time.format._`)
	w.Line(`import java.util.UUID`)
	w.Line(`import io.circe.Codec`)
	w.Line(`import io.circe.generic.extras.{Configuration, JsonKey}`)
	w.Line(`import io.circe.generic.extras.semiauto.{deriveConfiguredCodec, deriveUnwrappedCodec}`)

	for _, model := range version.ResolvedModels {
		w.EmptyLine()
		if model.IsObject() {
			generateCirceObjectModel(w, model)
		} else if model.IsEnum() {
			generateCirceEnumModel(w, model)
		} else if model.IsOneOf() {
			generateCirceUnionModel(w, model, packageName)
		}
	}

	return &gen.TextFile{
		Path:    filepath.Join(outPath, fmt.Sprintf("%sModels.scala", version.Version.PascalCase())),
		Content: w.String(),
	}
}

func generateCirceObjectModel(w *gen.Writer, model *spec.NamedModel) {
	w.Line(`case class %s(`, model.Name.PascalCase())
	for index, field := range model.Object.Fields {
		comma := ","
		if index == len(model.Object.Fields)-1 {
			comma = ""
		}
		w.Line(`  @JsonKey("%s") %s: %s%s`, field.Name.Source, scalaCamelCase(field.Name), ScalaType(&field.Type.Definition), comma)
	}
	w.Line(`)`)
	w.EmptyLine()
	w.Line(`object %s {`, model.Name.PascalCase())
	w.Line(`  implicit val config = Configuration.default`)
	w.Line(`  implicit val codec: Codec[%s] = deriveConfiguredCodec`, model.Name.PascalCase())
	w.Line(`}`)
}

func generateCirceEnumModel(w *gen.Writer, model *spec.NamedModel) {
	className := model.Name.PascalCase()
	w.Line(`sealed abstract class %s(val value: String) extends StringEnumEntry`, className)
	w.EmptyLine()
	w.Line(`case object %s extends StringEnum[%s] with StringCirceEnum[%s] {`, className, className, className)
	for _, item := range model.Enum.Items {
		w.Line(`  case object %s extends %s("%s")`, item.Name.PascalCase(), className, item.Value)
	}
	w.Line(`  val values = findValues`)
	w.Line(`}`)
}

func generateCirceUnionModel(w *gen.Writer, model *spec.NamedModel, packageName string) {
	traitName := model.Name.PascalCase()
	w.Line(`sealed trait %s`, traitName)
	w.EmptyLine()
	w.Line(`object %s {`, traitName)
	for _, item := range model.OneOf.Items {
		w.Line(`  case class %s(data: %s.%s) extends %s`, item.Name.PascalCase(), packageName, ScalaType(&item.Type.Definition), traitName)
	}
	w.Line(`  implicit val config = Configuration.default.withSnakeCaseConstructorNames`)
	w.Line(`  implicit val codec: Codec[%s] = deriveConfiguredCodec`, traitName)

	for _, item := range model.OneOf.Items {
		itemTypeName := model.Name.PascalCase() + "." + item.Name.PascalCase()
		w.Line(`  implicit val codec%s: Codec[%s] = deriveUnwrappedCodec`, item.Name.PascalCase(), itemTypeName)
	}
	w.Line(`}`)
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
