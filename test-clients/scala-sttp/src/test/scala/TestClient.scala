package testservice

import org.scalatest.FlatSpec

import scala.concurrent.Await
import scala.concurrent.duration.Duration
import com.softwaremill.sttp.akkahttp.AkkaHttpBackend
import util.Util

import testservice.echo._
import testservice.models._

class EchoClientSpec extends FlatSpec {
  implicit val httpBackend = AkkaHttpBackend()

  "echoBodyString" should "return body with same text" in {
    val client = new EchoClient(Util.service_url)
    val expected = "some text"
    val responseFuture = client.echoBodyString(expected)
    val actual = Await.ready(responseFuture, Duration.Inf).value.get.get
    assert(actual == expected)
  }

  "echoBodyModel" should "return body with same members values" in {
    val client = new EchoClient(Util.service_url)
    val expected = Message(intField = 123, stringField = "some string")
    val responseFuture = client.echoBodyModel(expected)
    val actual = Await.ready(responseFuture, Duration.Inf).value.get.get
    assert(actual == expected)
  }

  "echoBodyArray" should "return body with same items" in {
    val client = new EchoClient(Util.service_url)
    val expected = List("str1", "str2")
    val responseFuture = client.echoBodyArray(expected)
    val actual = Await.ready(responseFuture, Duration.Inf).value.get.get
    assert(actual == expected)
  }

  "echoBodyMap" should "return body with same items" in {
    val client = new EchoClient(Util.service_url)
    val expected = Map("one" -> "the str1", "two" -> "the str2")
    val responseFuture = client.echoBodyMap(expected)
    val actual = Await.ready(responseFuture, Duration.Inf).value.get.get
    assert(actual == expected)
  }

  "echoQuery" should "return query params with same values" in {
    val client = new EchoClient(Util.service_url)
    val responseFuture = client.echoQuery(
      123,
      12345.toLong,
      1.23f,
      12.345,
      BigDecimal("12345"),
      true,
      "the string",
      Some("the string"),
      List("string 1", "string 2"),
      java.util.UUID.fromString("123e4567-e89b-12d3-a456-426655440000"),
      java.time.LocalDate.parse("2021-01-01"),
      List(java.time.LocalDate.parse("2021-01-02")),
      java.time.LocalDateTime.parse("2021-01-02T23:54"),
      Choice.SecondChoice,
      "value",
    )

    val actual = Await.ready(responseFuture, Duration.Inf).value.get.get
    val expected = Parameters(
      intField = 123,
      longField = 12345.toLong,
      floatField = 1.23f,
      doubleField = 12.345,
      decimalField = BigDecimal("12345"),
      boolField = true,
      stringField = "the string",
      stringOptField = Some("the string"),
      stringDefaultedField = "value",
      stringArrayField = List("string 1", "string 2"),
      uuidField = java.util.UUID.fromString("123e4567-e89b-12d3-a456-426655440000"),
      dateField = java.time.LocalDate.parse("2021-01-01"),
      dateArrayField = List(java.time.LocalDate.parse("2021-01-02")),
      datetimeField = java.time.LocalDateTime.parse("2021-01-02T23:54"),
      enumField = Choice.SecondChoice,
    )
    assert(actual == expected)
  }

  "echoHeader" should "return header params" in {
    val client = new EchoClient(Util.service_url)
    val responseFuture = client.echoHeader(
      123,
      12345.toLong,
      1.23f,
      12.345,
      BigDecimal("12345"),
      true,
      "the string",
      Some("the string"),
      List("string 1", "string 2"),
      java.util.UUID.fromString("123e4567-e89b-12d3-a456-426655440000"),
      java.time.LocalDate.parse("2021-01-01"),
      List(java.time.LocalDate.parse("2021-01-02")),
      java.time.LocalDateTime.parse("2021-01-02T23:54"),
      Choice.SecondChoice,
      "value",
    )
    val actual = Await.ready(responseFuture, Duration.Inf).value.get.get
    val expected = Parameters(
      intField = 123,
      longField = 12345.toLong,
      floatField = 1.23f,
      doubleField = 12.345,
      decimalField = BigDecimal("12345"),
      boolField = true,
      stringField = "the string",
      stringOptField = Some("the string"),
      stringDefaultedField = "value",
      stringArrayField = List("string 1", "string 2"),
      uuidField = java.util.UUID.fromString("123e4567-e89b-12d3-a456-426655440000"),
      dateField = java.time.LocalDate.parse("2021-01-01"),
      dateArrayField = List(java.time.LocalDate.parse("2021-01-02")),
      datetimeField = java.time.LocalDateTime.parse("2021-01-02T23:54"),
      enumField = Choice.SecondChoice,
    )
    assert(actual == expected)
  }

  "echoUrlParams" should "return url params" in {
    val client = new EchoClient(Util.service_url)
    val responseFuture = client.echoUrlParams(
      123,
      12345.toLong,
      1.23f,
      12.345,
      BigDecimal("12345"),
      true,
      "the string",
      java.util.UUID.fromString("123e4567-e89b-12d3-a456-426655440000"),
      java.time.LocalDate.parse("2021-01-01"),
      java.time.LocalDateTime.parse("2021-01-02T23:54"),
      Choice.SecondChoice,
    )
    val actual = Await.ready(responseFuture, Duration.Inf).value.get.get
    val expected = UrlParameters(
      intField = 123,
      longField = 12345.toLong,
      floatField = 1.23f,
      doubleField = 12.345,
      decimalField = BigDecimal("12345"),
      boolField = true,
      stringField = "the string",
      uuidField = java.util.UUID.fromString("123e4567-e89b-12d3-a456-426655440000"),
      dateField = java.time.LocalDate.parse("2021-01-01"),
      datetimeField = java.time.LocalDateTime.parse("2021-01-02T23:54"),
      enumField = Choice.SecondChoice,
    )
    assert(actual == expected)
  }
}
