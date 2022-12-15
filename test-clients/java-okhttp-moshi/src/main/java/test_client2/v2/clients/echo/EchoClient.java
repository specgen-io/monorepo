package test_client2.v2.clients.echo;

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
import test_client2.v2.models.*;
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

	public Message echoBodyModel(Message body) {
		var bodyJson = json.write(Message.class, body);
		var requestBody = RequestBody.create(bodyJson, MediaType.parse("application/json"));
		var url = new UrlBuilder(baseUrl);
		url.addPathSegments("v2");
		url.addPathSegments("echo/body_model");

		var request = new RequestBuilder("POST", url.build(), requestBody);

		logger.info("Sending request, operationId: echo.echo_body_model, method: POST, url: /v2/echo/body_model");
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
}