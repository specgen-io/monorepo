package test_client;

import org.junit.jupiter.api.Test;
import test_client.clients.check.CheckClient;
import test_client.errors.*;
import test_client.errors.models.*;
import test_client.models.*;

import static org.junit.jupiter.api.Assertions.*;

public class CheckClientTest {
	private final String baseUrl = "http://localhost:8081";
	private final CheckClient client = new CheckClient(baseUrl);

	private final Message message = new Message(123, "the value");

	@Test
	public void checkEmpty_doesntThrowException() {
		assertDoesNotThrow(client::checkEmpty);
	}

	@Test
	public void checkEmptyResponse_doesntThrowException() {
		assertDoesNotThrow(() -> client.checkEmptyResponse(message));
	}

	@Test
	public void checkForbidden_doesntThrowException() {
		var exception = assertThrows(ClientException.class, client::checkForbidden);
		assertTrue(exception.getCause() instanceof ForbiddenException);
	}

	@Test
	public void checkConflict_checkExceptionBody() {
		var exception = (ConflictException) assertThrows(ClientException.class, client::checkConflict).getCause();
		var expectedBody = new ConflictMessage("Conflict with the current state of the target resource");
		assertEquals(exception.getBody(), expectedBody);
	}

	@Test
	public void checkConflict_throwException() {
		var exception = assertThrows(ClientException.class, client::checkConflict);
		assertTrue(exception.getCause() instanceof ConflictException);
	}

	@Test
	public void checkBadRequest_checkExceptionBody() {
		var exception = (BadRequestException) assertThrows(ClientException.class, client::checkBadRequest).getCause();
		var expectedBody = new BadRequestError("Failed to execute request", ErrorLocation.UNKNOWN, null);
		assertEquals(exception.getBody(), expectedBody);
	}

	@Test
	public void checkBadRequest_throwException() {
		var exception = assertThrows(ClientException.class, client::checkBadRequest);
		assertTrue(exception.getCause() instanceof BadRequestException);
	}
}
