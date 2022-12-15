package test_client;

import org.junit.jupiter.api.Test;
import test_client.v2.clients.echo.EchoClient;
import test_client.v2.models.Message;

import static org.junit.jupiter.api.Assertions.*;

public class EchoClientV2Test {
	private final String baseUrl = "http://localhost:8081";
	private final EchoClient client = new EchoClient(baseUrl);

	private final Message message = new Message(true, "the string");

	@Test
	public void echoBodyModel_responseIsEqualToRequest() {
		var response = client.echoBodyModel(message);
		assertEquals(message, response);
	}

	@Test
	public void echoBodyModel_doesntThrowException() {
		assertDoesNotThrow(() -> client.echoBodyModel(message));
	}
}
