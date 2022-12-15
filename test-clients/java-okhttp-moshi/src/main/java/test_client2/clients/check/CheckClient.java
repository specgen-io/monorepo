package test_client2.clients.check;

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

public class CheckClient {
	private static final Logger logger = LoggerFactory.getLogger(CheckClient.class);

	private final String baseUrl;
	private final OkHttpClient client;
	private final Json json;

	public CheckClient(String baseUrl, OkHttpClient client, Moshi moshi) {
		this.baseUrl = baseUrl;
		this.client = client;
		this.json = new Json(moshi);
	}

	public CheckClient(String baseUrl, OkHttpClient client) {
		this(baseUrl, client, CustomMoshiAdapters.setup(new Moshi.Builder()).build());
	}

	public CheckClient(String baseUrl) {
		this.baseUrl = baseUrl;
		this.json = new Json(CustomMoshiAdapters.setup(new Moshi.Builder()).build());
		this.client = new OkHttpClient().newBuilder().addInterceptor(new ErrorsInterceptor(json)).build();
	}

	public void checkEmpty() {
		var url = new UrlBuilder(baseUrl);
		url.addPathSegments("check/empty");

		var request = new RequestBuilder("GET", url.build(), null);

		logger.info("Sending request, operationId: check.check_empty, method: GET, url: /check/empty");
		var response = doRequest(client, request, logger);

		if (response.code() == 200) {
			logger.info("Received response with status code {}", response.code());
			return;
		}

		var errorMessage = "Unexpected status code received: " + response.code();
		logger.error(errorMessage);
		throw new ClientException(errorMessage);
	}

	public void checkEmptyResponse(Message body) {
		var bodyJson = json.write(Message.class, body);
		var requestBody = RequestBody.create(bodyJson, MediaType.parse("application/json"));
		var url = new UrlBuilder(baseUrl);
		url.addPathSegments("check/empty_response");

		var request = new RequestBuilder("POST", url.build(), requestBody);

		logger.info("Sending request, operationId: check.check_empty_response, method: POST, url: /check/empty_response");
		var response = doRequest(client, request, logger);

		if (response.code() == 200) {
			logger.info("Received response with status code {}", response.code());
			return;
		}

		var errorMessage = "Unexpected status code received: " + response.code();
		logger.error(errorMessage);
		throw new ClientException(errorMessage);
	}

	public Message checkForbidden() {
		var url = new UrlBuilder(baseUrl);
		url.addPathSegments("check/forbidden");

		var request = new RequestBuilder("GET", url.build(), null);

		logger.info("Sending request, operationId: check.check_forbidden, method: GET, url: /check/forbidden");
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

	public void checkConflict() {
		var url = new UrlBuilder(baseUrl);
		url.addPathSegments("check/conflict");

		var request = new RequestBuilder("GET", url.build(), null);

		logger.info("Sending request, operationId: check.check_conflict, method: GET, url: /check/conflict");
		var response = doRequest(client, request, logger);

		if (response.code() == 200) {
			logger.info("Received response with status code {}", response.code());
			return;
		}

		var errorMessage = "Unexpected status code received: " + response.code();
		logger.error(errorMessage);
		throw new ClientException(errorMessage);
	}

	public void checkBadRequest() {
		var url = new UrlBuilder(baseUrl);
		url.addPathSegments("check/bad_request");

		var request = new RequestBuilder("GET", url.build(), null);

		logger.info("Sending request, operationId: check.check_bad_request, method: GET, url: /check/bad_request");
		var response = doRequest(client, request, logger);

		if (response.code() == 200) {
			logger.info("Received response with status code {}", response.code());
			return;
		}

		var errorMessage = "Unexpected status code received: " + response.code();
		logger.error(errorMessage);
		throw new ClientException(errorMessage);
	}

	public void sameOperationName() {
		var url = new UrlBuilder(baseUrl);
		url.addPathSegments("check/same_operation_name");

		var request = new RequestBuilder("GET", url.build(), null);

		logger.info("Sending request, operationId: check.same_operation_name, method: GET, url: /check/same_operation_name");
		var response = doRequest(client, request, logger);

		if (response.code() == 200) {
			logger.info("Received response with status code {}", response.code());
			return;
		}

		var errorMessage = "Unexpected status code received: " + response.code();
		logger.error(errorMessage);
		throw new ClientException(errorMessage);
	}
}