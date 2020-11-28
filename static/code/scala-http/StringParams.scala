package [[.PackageName]]

import java.time.format.DateTimeFormatter
import java.time.{LocalDate, LocalDateTime, LocalTime}
import java.util.UUID

import enumeratum.values.{StringEnum, StringEnumEntry}

import scala.reflect.runtime.universe._
import scala.collection.mutable

class StringParamsReader(val reader: String => Option[String]) {

  def read[T](name: String)(implicit tag: TypeTag[T]): Option[T] =
    reader(name).map(StringParams.readValue[T])

  def readEnum[T <: StringEnumEntry: StringEnum](name: String)(implicit tag: TypeTag[T]): Option[T] =
    reader(name).map(StringParams.readEnumValue[T])
}

class StringParamsWriter {
  private val paramsMap = new mutable.HashMap[String, String]()
  def params(): Map[String, String] = paramsMap.toMap

  def write[T](name: String, value: Option[T])(implicit tag: TypeTag[T]) = {
    value match {
      case Some(value) => paramsMap(name) = StringParams.writeValue(value)
      case None => ()
    }
  }

  def write[T](name: String, value: T)(implicit tag: TypeTag[T]) = {
    paramsMap(name) = StringParams.writeValue(value)
  }
}

object StringParams {
  def readEnumValue[T <: StringEnumEntry: StringEnum](value: String)(implicit tag: TypeTag[T]): T = {
    implicitly[StringEnum[T]].withValue(value)
  }

  def readValue[T](value: String)(implicit tag: TypeTag[T]): T = {
    val typeName = tag.tpe.toString
    (typeName match {
      case "Byte" => value.toByte
      case "Short" => value.toShort
      case "Int" => value.toInt
      case "Long" => value.toLong
      case "Float" => value.toFloat
      case "Double" => value.toDouble
      case "BigDecimal" => BigDecimal(value)
      case "Boolean" => value.toBoolean
      case "Char" => value.charAt(0)
      case "String" => value
      case "java.time.LocalDate" => LocalDate.parse(value, DateTimeFormatter.ISO_LOCAL_DATE)
      case "java.time.LocalDateTime" => LocalDateTime.parse(value, DateTimeFormatter.ISO_LOCAL_DATE_TIME)
      case "java.time.LocalTime" => LocalTime.parse(value, DateTimeFormatter.ISO_LOCAL_TIME)
      case "java.util.UUID" => UUID.fromString(value)
      case _ => throw new Exception(s"Unsupported query parameter type: $typeName")
    }).asInstanceOf[T]
  }

  def writeValue[T](value: T)(implicit tag: TypeTag[T]): String = {
    value match {
      case value: StringEnumEntry => value.value
      case value: LocalDate => DateTimeFormatter.ISO_LOCAL_DATE.format(value)
      case value: LocalDateTime => DateTimeFormatter.ISO_LOCAL_DATE_TIME.format(value)
      case value: LocalTime => DateTimeFormatter.ISO_LOCAL_TIME.format(value)
      case value => value.toString
    }
  }

}