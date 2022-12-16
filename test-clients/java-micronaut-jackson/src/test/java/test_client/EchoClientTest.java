package test_client;

import org.junit.jupiter.api.Test;
import test_client.clients.echo.*;
import test_client.models.*;

import java.math.BigDecimal;
import java.time.*;
import java.util.*;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.*;

public class EchoClientTest {
	private final String baseUrl = "http://localhost:8081";
	private final EchoClient client = new EchoClient(baseUrl);

	private final int intValue = 123;
	private final long longValue = 12345;
	private final float floatValue = 1.23f;
	private final double doubleValue = 12.345;
	private final BigDecimal decimalValue = new BigDecimal("12345");
	private final boolean boolValue = true;
	private final String stringValue = "the value";
	private final String stringOptValue = "the value";
	private final String stringDefaultedValue = "value";
	private final List<String> stringArrayValue = Arrays.asList("the str1", "the str2");
	private final UUID uuidValue = UUID.fromString("123e4567-e89b-12d3-a456-426655440000");
	private final LocalDate dateValue = LocalDate.parse("2020-01-01");
	private final List<LocalDate> dateArrayValue = Arrays.asList(LocalDate.parse("2020-01-01"), LocalDate.parse("2020-01-02"));
	private final LocalDateTime datetimeValue = LocalDateTime.parse("2019-11-30T17:45:55");
	private final Choice enumValue = Choice.SECOND_CHOICE;

	private final String bodyString = "TWFueSBoYW5kcyBtYWtlIGxpZ2h0IHdvcmsu";
	private final Message message = new Message(intValue, stringValue);
	private final HashMap<String, String> hashMap = new HashMap<>() {{
		put("string_field", "the value");
		put("string_field_2", "the value_2");
	}};

	@Test
	public void echoBodyString_responseIsEqualToRequest() {
		String response = client.echoBodyString(bodyString);
		assertEquals(bodyString, response);
	}

	@Test
	public void echoBodyString_doesntThrowException() {
		assertDoesNotThrow(() -> client.echoBodyString(bodyString));
	}

	@Test
	public void echoBodyModel_responseIsEqualToRequest() {
		var response = client.echoBodyModel(message);
		assertEquals(message, response);
	}

	@Test
	public void echoBodyModel_doesntThrowException() {
		assertDoesNotThrow(() -> client.echoBodyModel(message));
	}

	@Test
	public void echoBodyArray_responseIsEqualToRequest() {
		var response = client.echoBodyArray(stringArrayValue);
		assertEquals(stringArrayValue, response);
	}

	@Test
	public void echoBodyArray_doesntThrowException() {
		assertDoesNotThrow(() -> client.echoBodyArray(stringArrayValue));
	}

	@Test
	public void echoBodyMap_responseIsEqualToRequest() {
		var response = client.echoBodyMap(hashMap);
		assertEquals(hashMap, response);
	}

	@Test
	public void echoBodyMap_doesntThrowException() {
		assertDoesNotThrow(() -> client.echoBodyMap(hashMap));
	}

	@Test
	public void echoQuery_responseIsEqualToRequest() {
		var request = new Parameters(
			intValue,
			longValue,
			floatValue,
			doubleValue,
			decimalValue,
			boolValue,
			stringValue,
			stringOptValue,
			stringDefaultedValue,
			stringArrayValue,
			uuidValue,
			dateValue,
			dateArrayValue,
			datetimeValue,
			enumValue
		);
		var response = client.echoQuery(
			intValue,
			longValue,
			floatValue,
			doubleValue,
			decimalValue,
			boolValue,
			stringValue,
			stringOptValue,
			stringDefaultedValue,
			stringArrayValue,
			uuidValue,
			dateValue,
			dateArrayValue,
			datetimeValue,
			enumValue
		);
		assertEquals(request, response);
	}

	@Test
	public void echoQuery_doesntThrowException() {
		assertDoesNotThrow(() -> client.echoQuery(
			intValue,
			longValue,
			floatValue,
			doubleValue,
			decimalValue,
			boolValue,
			stringValue,
			stringOptValue,
			stringDefaultedValue,
			stringArrayValue,
			uuidValue,
			dateValue,
			dateArrayValue,
			datetimeValue,
			enumValue
		));
	}

	@Test
	public void echoHeader_responseIsEqualToRequest() {
		var request = new Parameters(
			intValue,
			longValue,
			floatValue,
			doubleValue,
			decimalValue,
			boolValue,
			stringValue,
			stringOptValue,
			stringDefaultedValue,
			stringArrayValue,
			uuidValue,
			dateValue,
			dateArrayValue,
			datetimeValue,
			enumValue
		);
		var response = client.echoHeader(
			intValue,
			longValue,
			floatValue,
			doubleValue,
			decimalValue,
			boolValue,
			stringValue,
			stringOptValue,
			stringDefaultedValue,
			stringArrayValue,
			uuidValue,
			dateValue,
			dateArrayValue,
			datetimeValue,
			enumValue
		);
		assertEquals(request, response);
	}

	@Test
	public void echoHeader_doesntThrowException() {
		assertDoesNotThrow(() -> client.echoHeader(
			intValue,
			longValue,
			floatValue,
			doubleValue,
			decimalValue,
			boolValue,
			stringValue,
			stringOptValue,
			stringDefaultedValue,
			stringArrayValue,
			uuidValue,
			dateValue,
			dateArrayValue,
			datetimeValue,
			enumValue
		));
	}

	@Test
	public void echoUrlParams_responseIsEqualToRequest() {
		var request = new UrlParameters(
			intValue,
			longValue,
			floatValue,
			doubleValue,
			decimalValue,
			boolValue,
			stringValue,
			uuidValue,
			dateValue,
			datetimeValue,
			enumValue
		);
		var response = client.echoUrlParams(
			intValue,
			longValue,
			floatValue,
			doubleValue,
			decimalValue,
			boolValue,
			stringValue,
			uuidValue,
			dateValue,
			datetimeValue,
			enumValue
		);
		assertEquals(request, response);
	}

	@Test
	public void echoUrlParams_doesntThrowException() {
		assertDoesNotThrow(() -> client.echoUrlParams(
			intValue,
			longValue,
			floatValue,
			doubleValue,
			decimalValue,
			boolValue,
			stringValue,
			uuidValue,
			dateValue,
			datetimeValue,
			enumValue
		));
	}

	@Test
	public void echoEverything_responseIsEqualToRequest() {
		var request = new Everything(
			message,
			floatValue,
			boolValue,
			uuidValue,
			datetimeValue,
			dateValue,
			decimalValue
		);
		var response = client.echoEverything(
			message,
			floatValue,
			boolValue,
			uuidValue,
			datetimeValue,
			dateValue,
			decimalValue
		);
		assertEquals(request, response);
	}

	@Test
	public void echoEverything_doesntThrowException() {
		assertDoesNotThrow(() -> client.echoEverything(
			message,
			floatValue,
			boolValue,
			uuidValue,
			datetimeValue,
			dateValue,
			decimalValue
		));
	}

	@Test
	public void echoSuccess_ok() {
		var expected = new EchoSuccessResponse.Ok(new OkResult("ok"));
		var response = client.echoSuccess("ok");

		assertThat(expected).usingRecursiveComparison().isEqualTo(response);
	}

	@Test
	public void echoSuccess_created() {
		var expected = new EchoSuccessResponse.Created(new CreatedResult("created"));
		var response = client.echoSuccess("created");

		assertThat(expected).usingRecursiveComparison().isEqualTo(response);
	}

	@Test
	public void echoSuccess_accepted() {
		var expected = new EchoSuccessResponse.Accepted(new AcceptedResult("accepted"));
		var response = client.echoSuccess("accepted");

		assertThat(expected).usingRecursiveComparison().isEqualTo(response);
	}
}
