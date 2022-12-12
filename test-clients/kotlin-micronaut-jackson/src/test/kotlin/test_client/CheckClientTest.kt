package test_client

import org.junit.jupiter.api.Assertions.*
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.assertThrows
import test_client.clients.check.CheckClient
import test_client.errors.*
import test_client.errors.models.*
import test_client.models.*

class CheckClientTest {
    private val baseUrl = "http://localhost:8081"
    private val client = CheckClient(baseUrl)

    private val bodyMessage = Message(123, "the string")

    @Test
    fun checkEmpty_doesntThrowException() {
        assertDoesNotThrow { client.checkEmpty() }
    }

    @Test
    fun checkEmptyResponse_doesntThrowException() {
        assertDoesNotThrow { client.checkEmptyResponse(bodyMessage) }
    }

    @Test
    fun checkForbidden_throwException() {
        val exception = assertThrows<ClientException> { client.checkForbidden() }
        assertTrue(exception.cause is ForbiddenException)
    }

    @Test
    fun checkConflict_checkExceptionBody() {
        val exception = assertThrows<ClientException> { client.checkConflict() }.cause as ConflictException
        val expectedBody = ConflictMessage("Conflict with the current state of the target resource")
        assertEquals(exception.body, expectedBody)
    }

    @Test
    fun checkConflict_throwException() {
        val exception = assertThrows<ClientException> { client.checkConflict() }
        assertTrue(exception.cause is ConflictException)
    }

    @Test
    fun checkBadRequest_checkExceptionBody() {
        val exception = assertThrows<ClientException> { client.checkBadRequest() }.cause as BadRequestErrorException
        val expectedBody = BadRequestError("Failed to execute request", ErrorLocation.UNKNOWN, null)
        assertEquals(exception.body, expectedBody)
    }

    @Test
    fun checkBadRequest_throwException() {
        val exception = assertThrows<ClientException> { client.checkBadRequest() }
        assertTrue(exception.cause is BadRequestErrorException)
    }
}