package testservice.models

import java.time.{LocalDate, LocalDateTime, LocalTime}
import java.time.format.DateTimeFormatter
import java.util.UUID

import io.circe.{Decoder, Encoder, Json}
import testservice.Jsoner

import org.scalatest.FlatSpec

class ModelsSpec extends FlatSpec {
  def check[T](data: T, jsonStr: String) (implicit encoder: Encoder[T], decoder: Decoder[T]): Unit = {
    val actualJson = Jsoner.write(data)
    assert(jsonStr == actualJson)
    val actualData = Jsoner.read[T](jsonStr)
    assert(Right(data) == actualData)
  }

  "object fields" should "be serializable" in {
    check(Message(123), """{"field":123}""")
  }

  "object fields of different cases" should "be serializable" in {
    val data = MessageCases("snake_case value", "camelCase value")
    val jsonStr = """{"snake_case":"snake_case value","camelCase":"camelCase value"}"""
    check(data, jsonStr)
  }

  "numeric fields" should "be serializable" in {
    val data = NumericFields(
      intField = 0,
      longField = 0,
      floatField = 1.23f,
      doubleField = 1.23,
      decimalField = BigDecimal("1.23"),
    )
    val jsonStr = """{"int_field":0,"long_field":0,"float_field":1.23,"double_field":1.23,"decimal_field":1.23}"""
    check(data, jsonStr)
  }

  "non numeric fields" should "be serializable" in {
    val data = NonNumericFields(
      booleanField = true,
      stringField = "the string",
      uuidField = UUID.fromString("123e4567-e89b-12d3-a456-426655440000"),
      dateField = LocalDate.parse("2019-11-30", DateTimeFormatter.ISO_LOCAL_DATE),
      datetimeField = LocalDateTime.parse("2019-11-30T17:45:55", DateTimeFormatter.ISO_LOCAL_DATE_TIME),
    )
    val jsonStr = """{"boolean_field":true,"string_field":"the string","uuid_field":"123e4567-e89b-12d3-a456-426655440000","date_field":"2019-11-30","datetime_field":"2019-11-30T17:45:55"}"""
    check(data, jsonStr)
  }

  "array fields" should "be serializable" in {
    val data = ArrayFields(
      intArrayField = List(1, 2, 3),
      stringArrayField = List("one", "two", "three"),
    )
    val jsonStr = """{"int_array_field":[1,2,3],"string_array_field":["one","two","three"]}"""
    check(data, jsonStr)
  }

  "map fields" should "be serializable" in {
    val data = MapFields(
      intMapField = Map("one" -> 1, "two" -> 2),
      stringMapField = Map("one" -> "first", "two" -> "second"),
    )
    val jsonStr = """{"int_map_field":{"one":1,"two":2},"string_map_field":{"one":"first","two":"second"}}"""
    check(data, jsonStr)
  }

  "option fields" should "be serializable" in {
    val data = OptionalFields(
      intOptionField = Some(123),
      stringOptionField = Some("the string"),
    )
    val jsonStr = """{"int_option_field":123,"string_option_field":"the string"}"""
    check(data, jsonStr)
  }

  "option null fields" should "be serializable" in {
    val data = OptionalFields(
      intOptionField = None,
      stringOptionField = None,
    )
    val jsonStr = """{}"""
    check(data, jsonStr)
  }

  "enum field" should "be serializable" in {
    val data = EnumFields(enumField = Choice.SecondChoice)
    val jsonStr = """{"enum_field":"Two"}"""
    check(data, jsonStr)
  }

  "json field" should "be serializable" in {
    val data = RawJsonField(
      jsonField = Json.fromFields(Map(
        "the_array" -> Json.fromValues(List(Json.fromInt(1), Json.fromString("some string"))),
        "the_object" -> Json.fromFields(Map(
          "the_bool" -> Json.fromBoolean(true),
          "the_string" -> Json.fromString("some value"),
        )),
        "the_scalar" -> Json.fromInt(123),
      )),
    )
    val jsonStr = """{"json_field":{"the_array":[1,"some string"],"the_object":{"the_bool":true,"the_string":"some value"},"the_scalar":123}}"""
    check(data, jsonStr)
  }

  "oneOf" should "be serializable" in {
    val data: OrderEventWrapper = OrderEventWrapper.Changed(OrderChanged(UUID.fromString("58d5e212-165b-4ca0-909b-c86b9cee0111"), 3))
    val jsonStr = """{"changed":{"id":"58d5e212-165b-4ca0-909b-c86b9cee0111","quantity":3}}"""
    check(data, jsonStr)
  }

  "oneOf discriminator" should "be serializable" in {
    val data: OrderEventDiscriminator = OrderEventDiscriminator.Changed(OrderChanged(UUID.fromString("58d5e212-165b-4ca0-909b-c86b9cee0111"), 3))
    val jsonStr = """{"_type":"changed","id":"58d5e212-165b-4ca0-909b-c86b9cee0111","quantity":3}"""
    check(data, jsonStr)
  }
}