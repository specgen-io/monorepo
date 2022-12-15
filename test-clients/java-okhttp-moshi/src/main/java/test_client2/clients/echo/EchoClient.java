package test_client2.clients.echo;

import com.squareup.moshi.Moshi;
import com.squareup.moshi.Types;
import java.lang.reflect.ParameterizedType;
import java.math.BigDecimal;
import java.time.*;
import java.util.*;
import java.io.*;
import okhttp3.*;
import org.slf4j.*;
import test_client2.errors.*;
import test_client2.json.*;
import test_client2.utils.*;
import test_client2.models.*;
import static test_client2.utils.ClientResponse.*;

public class EchoClient {
	private static final Logger logger = LoggerFactory.getLogger(EchoClient.class);

	private final String baseUrl;
	private final OkHttpClient client;
	private final Json json;

	public EchoClient(String baseUrl, OkHttpClient client, Moshi moshi) {
		this.baseUrl = baseUrl;
		this.client = client;
		this.json = new Json(moshi);
	}

	public EchoClient(String baseUrl, OkHttpClient client) {
		this(baseUrl, client, CustomMoshiAdapters.setup(new Moshi.Builder()).build());
	}

	public EchoClient(String baseUrl) {
		this.baseUrl = baseUrl;
		this.json = new Json(CustomMoshiAdapters.setup(new Moshi.Builder()).build());
		this.client = new OkHttpClient().newBuilder().addInterceptor(new ErrorsInterceptor(json)).build();
	}

	public String echoBodyString(String body) {
		var requestBody = RequestBody.create(body, MediaType.parse("text/plain"));
		var url = new UrlBuilder(baseUrl);
		url.addPathSegments("echo/body_string");

		var request = new RequestBuilder("POST", url.build(), requestBody);

		logger.info("Sending request, operationId: echo.echo_body_string, method: POST, url: /echo/body_string");
		var response = doRequest(client, request, logger);

		if (response.code() == 200) {
			logger.info("Received response with status code {}", response.code());
			return getResponseBodyString(response, logger);
		}

		var errorMessage = "Unexpected status code received: " + response.code();
		logger.error(errorMessage);
		throw new ClientException(errorMessage);
	}

	public Message echoBodyModel(Message body) {
		var bodyJson = json.write(Message.class, body);
		var requestBody = RequestBody.create(bodyJson, MediaType.parse("application/json"));
		var url = new UrlBuilder(baseUrl);
		url.addPathSegments("echo/body_model");

		var request = new RequestBuilder("POST", url.build(), requestBody);

		logger.info("Sending request, operationId: echo.echo_body_model, method: POST, url: /echo/body_model");
		var response = doRequest(client, request, logger);

		if (response.code() == 200) {
			logger.info("Received response with status code {}", response.code());
			var responseBodyString = getResponseBodyString(response, logger);
			return json.read(responseBodyString, Message.class);
		}

		var errorMessage = "Unexpected status code received: " + response.code();
		logger.error(errorMessage);
		throw new ClientException(errorMessage);
	}

	public List<String> echoBodyArray(List<String> body) {
		var bodyJson = json.write(Types.newParameterizedType(List.class, String.class), body);
		var requestBody = RequestBody.create(bodyJson, MediaType.parse("application/json"));
		var url = new UrlBuilder(baseUrl);
		url.addPathSegments("echo/body_array");

		var request = new RequestBuilder("POST", url.build(), requestBody);

		logger.info("Sending request, operationId: echo.echo_body_array, method: POST, url: /echo/body_array");
		var response = doRequest(client, request, logger);

		if (response.code() == 200) {
			logger.info("Received response with status code {}", response.code());
			var responseBodyString = getResponseBodyString(response, logger);
			return json.read(responseBodyString, Types.newParameterizedType(List.class, String.class));
		}

		var errorMessage = "Unexpected status code received: " + response.code();
		logger.error(errorMessage);
		throw new ClientException(errorMessage);
	}

	public Map<String, String> echoBodyMap(Map<String, String> body) {
		var bodyJson = json.write(Types.newParameterizedType(Map.class, String.class, String.class), body);
		var requestBody = RequestBody.create(bodyJson, MediaType.parse("application/json"));
		var url = new UrlBuilder(baseUrl);
		url.addPathSegments("echo/body_map");

		var request = new RequestBuilder("POST", url.build(), requestBody);

		logger.info("Sending request, operationId: echo.echo_body_map, method: POST, url: /echo/body_map");
		var response = doRequest(client, request, logger);

		if (response.code() == 200) {
			logger.info("Received response with status code {}", response.code());
			var responseBodyString = getResponseBodyString(response, logger);
			return json.read(responseBodyString, Types.newParameterizedType(Map.class, String.class, String.class));
		}

		var errorMessage = "Unexpected status code received: " + response.code();
		logger.error(errorMessage);
		throw new ClientException(errorMessage);
	}

	public Parameters echoQuery(int intQuery, long longQuery, float floatQuery, double doubleQuery, BigDecimal decimalQuery, boolean boolQuery, String stringQuery, String stringOptQuery, String stringDefaultedQuery, List<String> stringArrayQuery, UUID uuidQuery, LocalDate dateQuery, List<LocalDate> dateArrayQuery, LocalDateTime datetimeQuery, Choice enumQuery) {
		var url = new UrlBuilder(baseUrl);
		url.addPathSegments("echo/query");
		url.addQueryParameter("int_query", intQuery);
		url.addQueryParameter("long_query", longQuery);
		url.addQueryParameter("float_query", floatQuery);
		url.addQueryParameter("double_query", doubleQuery);
		url.addQueryParameter("decimal_query", decimalQuery);
		url.addQueryParameter("bool_query", boolQuery);
		url.addQueryParameter("string_query", stringQuery);
		url.addQueryParameter("string_opt_query", stringOptQuery);
		url.addQueryParameter("string_defaulted_query", stringDefaultedQuery);
		url.addQueryParameter("string_array_query", stringArrayQuery);
		url.addQueryParameter("uuid_query", uuidQuery);
		url.addQueryParameter("date_query", dateQuery);
		url.addQueryParameter("date_array_query", dateArrayQuery);
		url.addQueryParameter("datetime_query", datetimeQuery);
		url.addQueryParameter("enum_query", enumQuery);

		var request = new RequestBuilder("GET", url.build(), null);

		logger.info("Sending request, operationId: echo.echo_query, method: GET, url: /echo/query");
		var response = doRequest(client, request, logger);

		if (response.code() == 200) {
			logger.info("Received response with status code {}", response.code());
			var responseBodyString = getResponseBodyString(response, logger);
			return json.read(responseBodyString, Parameters.class);
		}

		var errorMessage = "Unexpected status code received: " + response.code();
		logger.error(errorMessage);
		throw new ClientException(errorMessage);
	}

	public Parameters echoHeader(int intHeader, long longHeader, float floatHeader, double doubleHeader, BigDecimal decimalHeader, boolean boolHeader, String stringHeader, String stringOptHeader, String stringDefaultedHeader, List<String> stringArrayHeader, UUID uuidHeader, LocalDate dateHeader, List<LocalDate> dateArrayHeader, LocalDateTime datetimeHeader, Choice enumHeader) {
		var url = new UrlBuilder(baseUrl);
		url.addPathSegments("echo/header");

		var request = new RequestBuilder("GET", url.build(), null);
		request.addHeaderParameter("Int-Header", intHeader);
		request.addHeaderParameter("Long-Header", longHeader);
		request.addHeaderParameter("Float-Header", floatHeader);
		request.addHeaderParameter("Double-Header", doubleHeader);
		request.addHeaderParameter("Decimal-Header", decimalHeader);
		request.addHeaderParameter("Bool-Header", boolHeader);
		request.addHeaderParameter("String-Header", stringHeader);
		request.addHeaderParameter("String-Opt-Header", stringOptHeader);
		request.addHeaderParameter("String-Defaulted-Header", stringDefaultedHeader);
		request.addHeaderParameter("String-Array-Header", stringArrayHeader);
		request.addHeaderParameter("Uuid-Header", uuidHeader);
		request.addHeaderParameter("Date-Header", dateHeader);
		request.addHeaderParameter("Date-Array-Header", dateArrayHeader);
		request.addHeaderParameter("Datetime-Header", datetimeHeader);
		request.addHeaderParameter("Enum-Header", enumHeader);

		logger.info("Sending request, operationId: echo.echo_header, method: GET, url: /echo/header");
		var response = doRequest(client, request, logger);

		if (response.code() == 200) {
			logger.info("Received response with status code {}", response.code());
			var responseBodyString = getResponseBodyString(response, logger);
			return json.read(responseBodyString, Parameters.class);
		}

		var errorMessage = "Unexpected status code received: " + response.code();
		logger.error(errorMessage);
		throw new ClientException(errorMessage);
	}

	public UrlParameters echoUrlParams(int intUrl, long longUrl, float floatUrl, double doubleUrl, BigDecimal decimalUrl, boolean boolUrl, String stringUrl, UUID uuidUrl, LocalDate dateUrl, LocalDateTime datetimeUrl, Choice enumUrl) {
		var url = new UrlBuilder(baseUrl);
		url.addPathSegments("echo/url_params");
		url.addPathParameter(intUrl);
		url.addPathParameter(longUrl);
		url.addPathParameter(floatUrl);
		url.addPathParameter(doubleUrl);
		url.addPathParameter(decimalUrl);
		url.addPathParameter(boolUrl);
		url.addPathParameter(stringUrl);
		url.addPathParameter(uuidUrl);
		url.addPathParameter(dateUrl);
		url.addPathParameter(datetimeUrl);
		url.addPathParameter(enumUrl);

		var request = new RequestBuilder("GET", url.build(), null);

		logger.info("Sending request, operationId: echo.echo_url_params, method: GET, url: /echo/url_params/{int_url}/{long_url}/{float_url}/{double_url}/{decimal_url}/{bool_url}/{string_url}/{uuid_url}/{date_url}/{datetime_url}/{enum_url}");
		var response = doRequest(client, request, logger);

		if (response.code() == 200) {
			logger.info("Received response with status code {}", response.code());
			var responseBodyString = getResponseBodyString(response, logger);
			return json.read(responseBodyString, UrlParameters.class);
		}

		var errorMessage = "Unexpected status code received: " + response.code();
		logger.error(errorMessage);
		throw new ClientException(errorMessage);
	}

	public Everything echoEverything(Message body, float floatQuery, boolean boolQuery, UUID uuidHeader, LocalDateTime datetimeHeader, LocalDate dateUrl, BigDecimal decimalUrl) {
		var bodyJson = json.write(Message.class, body);
		var requestBody = RequestBody.create(bodyJson, MediaType.parse("application/json"));
		var url = new UrlBuilder(baseUrl);
		url.addPathSegments("echo/everything");
		url.addPathParameter(dateUrl);
		url.addPathParameter(decimalUrl);
		url.addQueryParameter("float_query", floatQuery);
		url.addQueryParameter("bool_query", boolQuery);

		var request = new RequestBuilder("POST", url.build(), requestBody);
		request.addHeaderParameter("Uuid-Header", uuidHeader);
		request.addHeaderParameter("Datetime-Header", datetimeHeader);

		logger.info("Sending request, operationId: echo.echo_everything, method: POST, url: /echo/everything/{date_url}/{decimal_url}");
		var response = doRequest(client, request, logger);

		if (response.code() == 200) {
			logger.info("Received response with status code {}", response.code());
			var responseBodyString = getResponseBodyString(response, logger);
			return json.read(responseBodyString, Everything.class);
		}

		var errorMessage = "Unexpected status code received: " + response.code();
		logger.error(errorMessage);
		throw new ClientException(errorMessage);
	}

	public void sameOperationName() {
		var url = new UrlBuilder(baseUrl);
		url.addPathSegments("echo/same_operation_name");

		var request = new RequestBuilder("GET", url.build(), null);

		logger.info("Sending request, operationId: echo.same_operation_name, method: GET, url: /echo/same_operation_name");
		var response = doRequest(client, request, logger);

		if (response.code() == 200) {
			logger.info("Received response with status code {}", response.code());
			return;
		}

		var errorMessage = "Unexpected status code received: " + response.code();
		logger.error(errorMessage);
		throw new ClientException(errorMessage);
	}

	public EchoSuccessResponse echoSuccess(String resultStatus) {
		var url = new UrlBuilder(baseUrl);
		url.addPathSegments("echo/success");
		url.addQueryParameter("result_status", resultStatus);

		var request = new RequestBuilder("GET", url.build(), null);

		logger.info("Sending request, operationId: echo.echo_success, method: GET, url: /echo/success");
		var response = doRequest(client, request, logger);

		if (response.code() == 200) {
			logger.info("Received response with status code {}", response.code());
			var responseBodyString = getResponseBodyString(response, logger);
			return new EchoSuccessResponse.Ok(json.read(responseBodyString, OkResult.class));
		}
		if (response.code() == 201) {
			logger.info("Received response with status code {}", response.code());
			var responseBodyString = getResponseBodyString(response, logger);
			return new EchoSuccessResponse.Created(json.read(responseBodyString, CreatedResult.class));
		}
		if (response.code() == 202) {
			logger.info("Received response with status code {}", response.code());
			var responseBodyString = getResponseBodyString(response, logger);
			return new EchoSuccessResponse.Accepted(json.read(responseBodyString, AcceptedResult.class));
		}

		var errorMessage = "Unexpected status code received: " + response.code();
		logger.error(errorMessage);
		throw new ClientException(errorMessage);
	}
}