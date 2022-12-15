package test_client2.utils;

import okhttp3.*;
import org.slf4j.*;
import test_client2.errors.*;

import java.io.IOException;

public class ClientResponse {
	public static Response doRequest(OkHttpClient client, RequestBuilder request, Logger logger) {
		Response response;
		try {
			response = client.newCall(request.build()).execute();
		} catch (IOException e) {
			var errorMessage = "Failed to execute the request " + e.getMessage();
			logger.error(errorMessage);
			throw new ClientException(errorMessage, e);
		}
		return response;
	}

	public static String getResponseBodyString(Response response, Logger logger) {
		String responseBodyString;
		try {
			responseBodyString = response.body().string();
		} catch (IOException e) {
			var errorMessage = "Failed to convert response body to string " + e.getMessage();
			logger.error(errorMessage);
			throw new ClientException(errorMessage, e);
		}
		return responseBodyString;
	}
}