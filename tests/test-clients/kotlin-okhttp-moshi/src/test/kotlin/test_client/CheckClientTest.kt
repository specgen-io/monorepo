package test_client

import org.junit.jupiter.api.Assertions.*
import org.junit.jupiter.api.Test
import test_client.clients.check.*
import test_client.errors.*
import test_client.models.*

class CheckClientTest {
    private val client = CheckClient("http://localhost:8081")

    private val intValue: Int = 123
    private val stringValue: String = "the value"
    private val bodyMessage = Message(intValue, stringValue)

    @Test
    fun checkEmpty_doesntThrowException() {
        assertDoesNotThrow(client::checkEmpty)
    }

    @Test
    fun checkEmptyResponse_doesntThrowException() {
        assertDoesNotThrow { client.checkEmptyResponse(bodyMessage) }
    }

    @Test
    fun checkForbidden_throwException() {
        assertThrows(ClientException::class.java) { client.checkForbidden() }
    }
}