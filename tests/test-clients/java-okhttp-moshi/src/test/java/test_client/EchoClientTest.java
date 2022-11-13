package test_client;

import org.junit.jupiter.api.Test;
import test_client.clients.echo.EchoClient;
import test_client.models.Everything;
import test_client.models.Parameters;
import test_client.models.UrlParameters;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static test_client.Constants.*;

public class EchoClientTest {
	private final EchoClient client = new EchoClient(BASE_URL);

	@Test
	public void echoBodyString_responseIsEqualToRequest() {
		String response = client.echoBodyString(BODY_STRING);

		assertEquals(BODY_STRING, response);
	}

	@Test
	public void echoBodyString_doesntThrowException() {
		assertDoesNotThrow(() -> client.echoBodyString(BODY_STRING));
	}

	@Test
	public void echoBodyModel_responseIsEqualToRequest() {
		var response = client.echoBodyModel(MESSAGE);

		assertEquals(MESSAGE, response);
	}

	@Test
	public void echoBodyModel_doesntThrowException() {
		assertDoesNotThrow(() -> client.echoBodyModel(MESSAGE));
	}

	@Test
	public void echoBodyArray_responseIsEqualToRequest() {
		var response = client.echoBodyArray(STRING_ARRAY_VALUE);

		assertEquals(STRING_ARRAY_VALUE, response);
	}

	@Test
	public void echoBodyArray_doesntThrowException() {
		assertDoesNotThrow(() -> client.echoBodyArray(STRING_ARRAY_VALUE));
	}

	@Test
	public void echoBodyMap_responseIsEqualToRequest() {
		var response = client.echoBodyMap(STRING_HASH_MAP);

		assertEquals(STRING_HASH_MAP, response);
	}

	@Test
	public void echoBodyMap_doesntThrowException() {
		assertDoesNotThrow(() -> client.echoBodyMap(STRING_HASH_MAP));
	}

	@Test
	public void echoQuery_responseIsEqualToRequest() {
		var request = new Parameters(
			INT_VALUE,
			LONG_VALUE,
			FLOAT_VALUE,
			DOUBLE_VALUE,
			DECIMAL_VALUE,
			BOOL_VALUE,
			STRING_VALUE,
			STRING_OPT_VALUE,
			STRING_DEFAULTED_VALUE,
			STRING_ARRAY_VALUE,
			UUID_VALUE,
			DATE_VALUE,
			DATE_ARRAY_VALUE,
			DATETIME_VALUE,
			ENUM_VALUE
		);
		var response = client.echoQuery(
			INT_VALUE,
			LONG_VALUE,
			FLOAT_VALUE,
			DOUBLE_VALUE,
			DECIMAL_VALUE,
			BOOL_VALUE,
			STRING_VALUE,
			STRING_OPT_VALUE,
			STRING_DEFAULTED_VALUE,
			STRING_ARRAY_VALUE,
			UUID_VALUE,
			DATE_VALUE,
			DATE_ARRAY_VALUE,
			DATETIME_VALUE,
			ENUM_VALUE
		);

		assertEquals(request, response);
	}

	@Test
	public void echoQuery_doesntThrowException() {
		assertDoesNotThrow(() -> client.echoQuery(
			INT_VALUE,
			LONG_VALUE,
			FLOAT_VALUE,
			DOUBLE_VALUE,
			DECIMAL_VALUE,
			BOOL_VALUE,
			STRING_VALUE,
			STRING_OPT_VALUE,
			STRING_DEFAULTED_VALUE,
			STRING_ARRAY_VALUE,
			UUID_VALUE,
			DATE_VALUE,
			DATE_ARRAY_VALUE,
			DATETIME_VALUE,
			ENUM_VALUE
		));
	}


	@Test
	public void echoHeader_responseIsEqualToRequest() {
		var request = new Parameters(
			INT_VALUE,
			LONG_VALUE,
			FLOAT_VALUE,
			DOUBLE_VALUE,
			DECIMAL_VALUE,
			BOOL_VALUE,
			STRING_VALUE,
			STRING_OPT_VALUE,
			STRING_DEFAULTED_VALUE,
			STRING_ARRAY_VALUE,
			UUID_VALUE,
			DATE_VALUE,
			DATE_ARRAY_VALUE,
			DATETIME_VALUE,
			ENUM_VALUE
		);
		var response = client.echoHeader(
			INT_VALUE,
			LONG_VALUE,
			FLOAT_VALUE,
			DOUBLE_VALUE,
			DECIMAL_VALUE,
			BOOL_VALUE,
			STRING_VALUE,
			STRING_OPT_VALUE,
			STRING_DEFAULTED_VALUE,
			STRING_ARRAY_VALUE,
			UUID_VALUE,
			DATE_VALUE,
			DATE_ARRAY_VALUE,
			DATETIME_VALUE,
			ENUM_VALUE
		);

		assertEquals(request, response);
	}

	@Test
	public void echoHeader_doesntThrowException() {
		assertDoesNotThrow(() -> client.echoHeader(
			INT_VALUE,
			LONG_VALUE,
			FLOAT_VALUE,
			DOUBLE_VALUE,
			DECIMAL_VALUE,
			BOOL_VALUE,
			STRING_VALUE,
			STRING_OPT_VALUE,
			STRING_DEFAULTED_VALUE,
			STRING_ARRAY_VALUE,
			UUID_VALUE,
			DATE_VALUE,
			DATE_ARRAY_VALUE,
			DATETIME_VALUE,
			ENUM_VALUE
		));
	}

	@Test
	public void echoUrlParams_responseIsEqualToRequest() {
		var request = new UrlParameters(
			INT_VALUE,
			LONG_VALUE,
			FLOAT_VALUE,
			DOUBLE_VALUE,
			DECIMAL_VALUE,
			BOOL_VALUE,
			STRING_VALUE,
			UUID_VALUE,
			DATE_VALUE,
			DATETIME_VALUE,
			ENUM_VALUE
		);
		var response = client.echoUrlParams(
			INT_VALUE,
			LONG_VALUE,
			FLOAT_VALUE,
			DOUBLE_VALUE,
			DECIMAL_VALUE,
			BOOL_VALUE,
			STRING_VALUE,
			UUID_VALUE,
			DATE_VALUE,
			DATETIME_VALUE,
			ENUM_VALUE
		);

		assertEquals(request, response);
	}

	@Test
	public void echoUrlParams_doesntThrowException() {
		assertDoesNotThrow(() -> client.echoUrlParams(
			INT_VALUE,
			LONG_VALUE,
			FLOAT_VALUE,
			DOUBLE_VALUE,
			DECIMAL_VALUE,
			BOOL_VALUE,
			STRING_VALUE,
			UUID_VALUE,
			DATE_VALUE,
			DATETIME_VALUE,
			ENUM_VALUE
		));
	}

	@Test
	public void echoEverything_responseIsEqualToRequest() {
		var request = new Everything(
			MESSAGE,
			FLOAT_VALUE,
			BOOL_VALUE,
			UUID_VALUE,
			DATETIME_VALUE,
			DATE_VALUE,
			DECIMAL_VALUE
		);
		var response = client.echoEverything(
			MESSAGE,
			FLOAT_VALUE,
			BOOL_VALUE,
			UUID_VALUE,
			DATETIME_VALUE,
			DATE_VALUE,
			DECIMAL_VALUE
		);

		assertThat(request).usingRecursiveComparison().isEqualTo(response);
	}

	@Test
	public void echoEverything_doesntThrowException() {
		assertDoesNotThrow(() -> client.echoEverything(
			MESSAGE,
			FLOAT_VALUE,
			BOOL_VALUE,
			UUID_VALUE,
			DATETIME_VALUE,
			DATE_VALUE,
			DECIMAL_VALUE
		));
	}
}
