package [[.PackageName]]

import java.time.format.DateTimeFormatter
import java.time.{LocalDate, LocalDateTime, LocalTime}

import enumeratum.values.{StringEnum, StringEnumEntry}
import play.api.mvc.QueryStringBindable

trait Codec[T] {
  def decode(s: String): T
  def encode(v: T): String
}

object QueryParamsBindings {
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

  implicit val ByteCodec: Codec[Byte] = new Codec[Byte] {
    def decode(s: String): Byte = s.toByte
    def encode(v: Byte): String = v.toString
  }

  implicit val ShortCodec: Codec[Short] = new Codec[Short] {
    def decode(s: String): Short = s.toShort
    def encode(v: Short): String = v.toString
  }

  implicit val BigDecimalCodec: Codec[BigDecimal] = new Codec[BigDecimal] {
    def decode(s: String): BigDecimal = BigDecimal(s)
    def encode(v: BigDecimal): String = v.toString
  }

  implicit val CharCodec: Codec[Char] = new Codec[Char] {
    def decode(s: String): Char = if (s.length == 1) { s.charAt(0) } else { throw new Exception("Char query parameter supposed to have one symbol") }
    def encode(v: Char): String = v.toString
  }

  implicit def enumCodec[T <: StringEnumEntry](implicit E: StringEnum[T]): Codec[T] = new Codec[T] {
    def decode(s: String): T = E.withValue(s)
    def encode(v: T): String = v.value
  }

  implicit def bindableParser[T](implicit stringBinder: QueryStringBindable[String], codec: Codec[T]): QueryStringBindable[T] = new QueryStringBindable[T] {
    override def bind(key: String, params: Map[String, Seq[String]]): Option[Either[String, T]] =
      for {
        dateStr <- stringBinder.bind(key, params)
      } yield {
        dateStr match {
          case Right(value) =>
            try {
              Right(codec.decode(value))
            } catch {
              case t: Throwable => Left(s"Unable to bind from key: $key, error: ${t.getMessage}")
            }
          case _ => Left(s"Unable to bind from key: $key")
        }
      }

    override def unbind(key: String, value: T): String = stringBinder.unbind(key, codec.encode(value))
  }
}