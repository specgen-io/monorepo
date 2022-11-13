package testservice.v2.models

import java.time.{LocalDate, LocalDateTime, LocalTime}
import java.time.format.DateTimeFormatter
import java.util.UUID
import io.circe.{Decoder, Encoder, Json}

import testservice.Jsoner

import org.scalatest.FlatSpec

class TestV2Models extends FlatSpec {
  def check[T](data: T, jsonStr: String) (implicit encoder: Encoder[T], decoder: Decoder[T]): Unit = {
    val actualJson = Jsoner.write(data)
    assert(jsonStr == actualJson)
    val actualData = Jsoner.read[T](jsonStr)
    assert(Right(data) == actualData)
  }

  "object fields" should "be serializable" in {
    check(Message(field = "some string"), """{"field":"some string"}""")
  }
}