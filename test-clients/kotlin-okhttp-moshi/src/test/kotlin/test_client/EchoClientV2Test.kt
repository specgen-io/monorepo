package test_client

import org.junit.jupiter.api.Assertions.*
import org.junit.jupiter.api.Test
import test_client.v2.clients.echo.*
import test_client.v2.models.*

class EchoClientV2Test {
    private val client = EchoClient("http://localhost:8081")

    private val bodyMessage = Message(true, "the string")

    @Test
    fun echoBodyModel_responseIsEqualToRequest() {
        val response = client.echoBodyModel(bodyMessage)
        assertEquals(bodyMessage, response)
    }

    @Test
    fun echoBodyModel_doesntThrowException() {
        assertDoesNotThrow { client.echoBodyModel(bodyMessage) }
    }
}