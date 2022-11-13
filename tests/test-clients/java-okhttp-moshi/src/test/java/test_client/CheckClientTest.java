package test_client;

import org.junit.jupiter.api.Test;
import test_client.clients.check.CheckClient;
import test_client.errors.ClientException;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static test_client.Constants.BASE_URL;
import static test_client.Constants.MESSAGE;

public class CheckClientTest {
	private final CheckClient client = new CheckClient(BASE_URL);

	@Test
	public void checkEmpty_doesntThrowException() {
		assertDoesNotThrow(client::checkEmpty);
	}

	@Test
	public void checkEmptyResponse_doesntThrowException() {
		assertDoesNotThrow(() -> client.checkEmptyResponse(MESSAGE));
	}

	@Test
	public void checkForbidden_doesntThrowException() {
		assertThrows(ClientException.class, client::checkForbidden);
	}
}
