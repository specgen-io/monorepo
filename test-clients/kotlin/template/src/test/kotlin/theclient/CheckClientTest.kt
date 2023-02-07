package theclient

import org.junit.jupiter.api.*
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.assertThrows
import theclient.clients.check.CheckClient
import theclient.errors.*
import theclient.errors.models.*
import theclient.models.*

class CheckClientTest {
    private val client = CheckClient("http://localhost:8081")

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
        Assertions.assertTrue(exception.cause is ForbiddenException)
    }

    @Test
    fun checkConflict_checkExceptionBody() {
        val exception = assertThrows<ClientException> { client.checkConflict() }.cause as ConflictException
        val expectedBody = ConflictMessage("Conflict with the current state of the target resource")
        Assertions.assertEquals(exception.body, expectedBody)
    }

    @Test
    fun checkConflict_throwException() {
        val exception = assertThrows<ClientException> { client.checkConflict() }
        Assertions.assertTrue(exception.cause is ConflictException)
    }

    @Test
    fun checkBadRequest_checkExceptionBody() {
        val exception = assertThrows<ClientException> { client.checkBadRequest() }.cause as BadRequestException
        val expectedBody = BadRequestError("Failed to execute request", ErrorLocation.UNKNOWN, null)
        Assertions.assertEquals(exception.body, expectedBody)
    }

    @Test
    fun checkBadRequest_throwException() {
        val exception = assertThrows<ClientException> { client.checkBadRequest() }
        Assertions.assertTrue(exception.cause is BadRequestException)
    }
}
