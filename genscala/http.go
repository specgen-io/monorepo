package genscala

import (
	"github.com/specgen-io/specgen/v2/sources"
	"strings"
)

func generateStringParams(thepackage Package) *sources.CodeFile {
	code := `
package [[.PackageName]]

import java.time.format.DateTimeFormatter
import java.time.{LocalDate, LocalDateTime, LocalTime}
import java.util.UUID

import enumeratum.values.{StringEnum, StringEnumEntry}

import scala.collection.mutable.ListBuffer

import scala.collection.immutable.List

trait Codec[T] {
  def decode(s: String): T
  def encode(v: T): String
}

object ParamsTypesBindings {
  implicit val StringCodec: Codec[String] = new Codec[String] {
    def decode(s: String): String = s
    def encode(v: String): String = v
  }

  implicit val ByteCodec: Codec[Byte] = new Codec[Byte] {
    def decode(s: String): Byte = s.toByte
    def encode(v: Byte): String = v.toString
  }

  implicit val ShortCodec: Codec[Short] = new Codec[Short] {
    def decode(s: String): Short = s.toShort
    def encode(v: Short): String = v.toString
  }

  implicit val IntCodec: Codec[Int] = new Codec[Int] {
    def decode(s: String): Int = s.toInt
    def encode(v: Int): String = v.toString
  }

  implicit val LongCodec: Codec[Long] = new Codec[Long] {
    def decode(s: String): Long = s.toShort
    def encode(v: Long): String = v.toString
  }

  implicit val FloatCodec: Codec[Float] = new Codec[Float] {
    def decode(s: String): Float = s.toFloat
    def encode(v: Float): String = v.toString
  }

  implicit val DoubleCodec: Codec[Double] = new Codec[Double] {
    def decode(s: String): Double = s.toDouble
    def encode(v: Double): String = v.toString
  }

  implicit val BigDecimalCodec: Codec[BigDecimal] = new Codec[BigDecimal] {
    def decode(s: String): BigDecimal = BigDecimal(s)
    def encode(v: BigDecimal): String = v.toString
  }

  implicit val BooleanCodec: Codec[Boolean] = new Codec[Boolean] {
    def decode(s: String): Boolean = s.toBoolean
    def encode(v: Boolean): String = v.toString
  }

  implicit val CharCodec: Codec[Char] = new Codec[Char] {
    def decode(s: String): Char = if (s.length == 1) { s.charAt(0) } else { throw new Exception("Char query parameter supposed to have one symbol") }
    def encode(v: Char): String = v.toString
  }

  implicit val DateCodec: Codec[LocalDate] = new Codec[LocalDate] {
    def decode(s: String): LocalDate = LocalDate.parse(s, DateTimeFormatter.ISO_LOCAL_DATE)
    def encode(v: LocalDate): String = DateTimeFormatter.ISO_LOCAL_DATE.format(v)
  }

  implicit val DateTimeCodec: Codec[LocalDateTime] = new Codec[LocalDateTime] {
    def decode(s: String): LocalDateTime = LocalDateTime.parse(s, DateTimeFormatter.ISO_LOCAL_DATE_TIME)
    def encode(v: LocalDateTime): String = DateTimeFormatter.ISO_LOCAL_DATE_TIME.format(v)
  }

  implicit val TimeCodec: Codec[LocalTime] = new Codec[LocalTime] {
    def decode(s: String): LocalTime = LocalTime.parse(s, DateTimeFormatter.ISO_LOCAL_TIME)
    def encode(v: LocalTime): String = DateTimeFormatter.ISO_LOCAL_TIME.format(v)
  }

  implicit val UuidCodec: Codec[UUID] = new Codec[UUID] {
    def decode(s: String): UUID = UUID.fromString(s)
    def encode(v: UUID): String = v.toString
  }

  implicit def enumCodec[T <: StringEnumEntry](implicit E: StringEnum[T]): Codec[T] = new Codec[T] {
    def decode(s: String): T = E.withValue(s)
    def encode(v: T): String = v.value
  }

  class StringParamsReader(val values: Map[String, Seq[String]]) {
    def read[T](name: String)(implicit codec: Codec[T]): Option[T] =
      values.get(name).flatMap(_.headOption).map(codec.decode)
    def readList[T](name: String)(implicit codec: Codec[T]): Option[List[T]] =
      values.get(name).map(_.map(codec.decode).toList)
  }

  def stringify[T](value: T)(implicit codec: Codec[T]): String = codec.encode(value)

  class StringParamsWriter {
    val paramsList = ListBuffer[(String, String)]()

    def params = paramsList.toSeq

    def write[T](name: String, value: Option[T])(implicit codec: Codec[T]): Unit = {
      value match {
        case Some(value) => write(name, value)
        case None => ()
      }
    }

    def write[T](name: String, values: Seq[T])(implicit codec: Codec[T]): Unit = {
      values.foreach { value => write(name, value) }
    }

    def write[T](name: String, value: T)(implicit codec: Codec[T]): Unit = {
      paramsList += ((name, stringify(value)))
    }
  }
}`
	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thepackage.PackageName})
	return &sources.CodeFile{
		Path:    thepackage.GetPath("StringParams.scala"),
		Content: strings.TrimSpace(code)}
}
