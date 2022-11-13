package testservice.v2

import org.scalatest.FlatSpec

import scala.concurrent.Await
import scala.concurrent.duration.Duration
import com.softwaremill.sttp.akkahttp.AkkaHttpBackend
import util.Util

import testservice.v2.echo._
import testservice.v2.models._

class EchoClientSpec extends FlatSpec {
  implicit val httpBackend = AkkaHttpBackend()

  "echBodyModel" should "return body with same members values" in {
    val client = new EchoClient(Util.service_url)
    val expected = Message(boolField = true, stringField = "some string")
    val responseFuture = client.echoBodyModel(expected)
    val actual = Await.ready(responseFuture, Duration.Inf).value.get.get
    assert(actual == expected)
  }
}