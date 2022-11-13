package test_client;

import org.junit.jupiter.api.Test;
import test_client.v2.clients.echo.EchoClient;
import test_client.v2.models.Message;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static test_client.Constants.BASE_URL;

public class EchoClientV2Test {
	private final EchoClient client = new EchoClient(BASE_URL);

	private final Message message = new Message(true, "the string");

	@Test
	public void echoBodyModel_responseIsEqualToRequest() {
		Message response = client.echoBodyModel(message);
		assertEquals(message, response);
	}

	@Test
	public void echoBodyModel_doesntThrowException() {
		assertDoesNotThrow(() -> client.echoBodyModel(message));
	}
}
