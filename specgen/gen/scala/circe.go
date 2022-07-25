package scala

import (
	"fmt"
	"github.com/specgen-io/specgen/spec/v2"
	"github.com/specgen-io/specgen/v2/generator"
	"strings"
)

func GenerateCirceModels(specification *spec.Spec, packageName string, generatePath string) *generator.Sources {
	if packageName == "" {
		packageName = specification.Name.FlatCase()
	}
	mainPackage := NewPackage(generatePath, packageName, "")
	jsonPackage := mainPackage

	sources := generator.NewSources()
	jsonHelpers := generateJson(jsonPackage)
	taggedUnion := generateTaggedUnion(jsonPackage)
	sources.AddGenerated(taggedUnion, jsonHelpers)

	for _, version := range specification.Versions {
		versionClientPackage := mainPackage.Subpackage(version.Version.FlatCase())
		versionModelsPackage := versionClientPackage.Subpackage("models")
		versionFile := generateCirceModels(&version, versionModelsPackage, jsonPackage)
		sources.AddGenerated(versionFile)
	}

	return sources
}

func generateCirceModels(version *spec.Version, thepackage, taggedUnionPackage Package) *generator.CodeFile {
	w := NewScalaWriter()
	w.Line(`package %s`, thepackage.PackageName)
	w.EmptyLine()
	w.Line(`import enumeratum.values._`)
	w.Line(`import java.time._`)
	w.Line(`import java.time.format._`)
	w.Line(`import java.util.UUID`)
	w.Line(`import io.circe.Codec`)
	w.Line(`import io.circe.generic.extras.{Configuration, JsonKey}`)
	w.Line(`import io.circe.generic.extras.semiauto.{deriveConfiguredCodec, deriveUnwrappedCodec}`)
	w.Line(`import %s.taggedunion._`, taggedUnionPackage.PackageName)

	for _, model := range version.ResolvedModels {
		w.EmptyLine()
		if model.IsObject() {
			generateCirceObjectModel(w, model)
		} else if model.IsEnum() {
			generateCirceEnumModel(w, model)
		} else if model.IsOneOf() {
			generateCirceUnionModel(w, model)
		}
	}

	return &generator.CodeFile{
		Path:    thepackage.GetPath("Models.scala"),
		Content: w.String(),
	}
}

func generateCirceObjectModel(w *generator.Writer, model *spec.NamedModel) {
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

func generateCirceEnumModel(w *generator.Writer, model *spec.NamedModel) {
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

func generateCirceUnionModel(w *generator.Writer, model *spec.NamedModel) {
	traitName := model.Name.PascalCase()
	w.Line(`sealed trait %s`, traitName)
	w.EmptyLine()
	w.Line(`object %s {`, traitName)
	for _, item := range model.OneOf.Items {
		w.Line(`  case class %s(data: %s) extends %s`, item.Name.PascalCase(), ScalaType(&item.Type.Definition), traitName)
	}
	w.EmptyLine()
	for _, item := range model.OneOf.Items {
		caseName := item.Name.PascalCase()
		w.Line(`  implicit val codec%s: Codec[%s] = deriveUnwrappedCodec`, caseName, caseName)
		w.Line(`  implicit val tag%s: Tag[%s] = Tag("%s")`, caseName, caseName, item.Name.Source)
	}
	w.EmptyLine()
	unionConfig := "Config.wrapped"
	if model.OneOf.Discriminator != nil {
		unionConfig = fmt.Sprintf(`Config.discriminator("%s")`, *model.OneOf.Discriminator)
	}
	w.Line(`  implicit val config: Config[%s] = %s`, traitName, unionConfig)
	w.Line(`  implicit val codec: Codec[%s] = deriveUnionCodec`, traitName)
	w.Line(`}`)
}

func generateJson(thepackage Package) *generator.CodeFile {
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
	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thepackage.PackageName})
	return &generator.CodeFile{
		Path:    thepackage.GetPath("Json.scala"),
		Content: strings.TrimSpace(code),
	}
}

func generateTaggedUnion(thepackage Package) *generator.CodeFile {
	code := `
package [[.PackageName]]

object taggedunion {
  import io.circe._
  import shapeless.{:+:, HNil, CNil, Coproduct, Generic, Inl, Inr, Lazy}

  implicit val cnilEncoder: Encoder[CNil] =
    Encoder.instance(cnil => throw new Exception("Inconceivable!"))

  implicit val cnilDecoder: Decoder[CNil] =
    Decoder.instance(a => Left(DecodingFailure("Discriminator field or wrapped object not recognized", a.history)))

  implicit def coproductEncoder[H, T <: Coproduct](
    implicit hEncoder: Lazy[Encoder[H]],
    tag: Lazy[Tag[H]],
    config: Lazy[Config[H]],
    tEncoder: Encoder[T]
  ): Encoder[H :+: T] = Encoder.instance{
    case Inl(head) =>
      encodeWithConfig(config.value, tag.value, hEncoder.value)(head)
    case Inr(tail) => tEncoder(tail)
  }

  implicit def coproductDecoder[H, T <: Coproduct](
    implicit hDecoder: Lazy[Decoder[H]],
    tag: Lazy[Tag[H]],
    config: Lazy[Config[H]],
    tDecoder: Decoder[T]
  ): Decoder[H :+: T] = Decoder.instance{ cursor =>
    decodeWithConfig(cursor, tag.value, config.value, hDecoder) match {
      case Left(NotMatch)   => tDecoder.tryDecode(cursor).map(Inr(_))
      case Right(value)     => Right(Inl[H, T](value))
      case Left(InternalParsingError(circeError))           => Left[DecodingFailure, H :+: T](circeError)
      case Left(DiscriminatorFieldCorrupted(circeError))    => Left[DecodingFailure, H :+: T](circeError)
    }
  }

  implicit val hnilEncoder: Encoder[HNil] =
    Encoder.instance(hnil => throw new Exception("Inconceivable!"))

  implicit val hnilDecoder: Decoder[HNil] =
    Decoder.instance(a => Left(DecodingFailure("Inconceivable!", a.history)))

  def deriveUnionEncoder[A, R](implicit gen: Generic.Aux[A, R], enc: Encoder[R]): Encoder[A] = {
    Encoder.instance(a => enc(gen.to(a)))
  }

  def deriveUnionDecoder[A, R](implicit gen: Generic.Aux[A, R], dec: Decoder[R]): Decoder[A] = {
    Decoder.instance(a => dec(a).map(gen.from))
  }

  def deriveUnionCodec[A, R](
    implicit
    gen: Generic.Aux[A, R],
    enc: Encoder[R],
    dec: Decoder[R]
  ): Codec[A] = new Codec[A] {
    override def apply(a: A): Json = enc(gen.to(a))
    override def apply(c: HCursor): Decoder.Result[A] = dec(c).map(gen.from)
  }

  trait Tag[A]{
    def name: String
  }

  object Tag{
    def apply[A](string: String): Tag[A] =
      new Tag[A] {
        override def name: String = string
      }
  }

  sealed trait CoproductConfig[-A]
  case class Discriminator[A](fieldName: String) extends CoproductConfig[A]
  case class WrappedObject[A]() extends CoproductConfig[A]

  trait Config[-A]{
    def coproduct: CoproductConfig[A]
  }

  object Config{
    def discriminator[A](fieldName: String) = new Config[A]{
      def coproduct: CoproductConfig[A] = Discriminator(fieldName)
    }

    def wrapped[A] = new Config[A]{
      def coproduct: CoproductConfig[A] = WrappedObject()
    }
  }

  def encodeWithConfig[H](config: Config[H], tag: Tag[H], encoder: Encoder[H]) = Encoder.instance[H]{ value =>
    config.coproduct match {
      case discriminator: Discriminator[H] =>
        encoder(value)
          .mapObject(js =>
            js.+:(discriminator.fieldName, Json.fromString(tag.name))
          )
      case _: WrappedObject[H] =>
        Json.obj((tag.name, encoder(value)))
    }
  }

  sealed trait DecodingError
  case object NotMatch extends DecodingError
  case class InternalParsingError(circeError: DecodingFailure) extends DecodingError
  case class DiscriminatorFieldCorrupted(circeError: DecodingFailure) extends DecodingError

  def decodeWithConfig[H](cursor: HCursor, tag: Tag[H], config: Config[H], hDecoder: Lazy[Decoder[H]]): Either[DecodingError, H] = {
    config.coproduct match {
      case Discriminator(fieldName) =>
        cursor.downField(fieldName).as[String] match {
          case Left(value)  =>
            Left(DiscriminatorFieldCorrupted(value))

          case Right(discriminator) =>
            if(discriminator == tag.name){
              hDecoder.value.tryDecode(cursor) match {
                case Left(value)  => Left(InternalParsingError(value))
                case Right(value) => Right(value)
              }
            } else {
              Left(NotMatch)
            }
        }

      case WrappedObject() => {
        hDecoder.value.tryDecode(cursor.downField(tag.name)) match {
          case Left(value)  =>
            Left(NotMatch)

          case Right(discriminator) =>
            Right(discriminator)
        }
      }
    }
  }
}`
	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thepackage.PackageName})
	return &generator.CodeFile{
		Path:    thepackage.GetPath("TaggedUnion.scala"),
		Content: strings.TrimSpace(code),
	}
}
