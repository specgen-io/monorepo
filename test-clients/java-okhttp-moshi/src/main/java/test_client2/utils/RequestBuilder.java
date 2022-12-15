package test_client2.utils;

import okhttp3.*;
import java.util.List;

public class RequestBuilder {
	private final Request.Builder requestBuilder;

	public RequestBuilder(String method, HttpUrl url, RequestBody body) {
		this.requestBuilder = new Request.Builder().url(url).method(method, body);
	}

	public RequestBuilder addHeaderParameter(String name, Object value) {
		var valueStr = Stringify.paramToString(value);
		if (valueStr != null) {
			this.requestBuilder.addHeader(name, valueStr);
		}
		return this;
	}

	public <T> RequestBuilder addHeaderParameter(String name, List<T> values) {
		for (T val : values) {
			this.addHeaderParameter(name, val);
		}
		return this;
	}

	public Request build() {
		return this.requestBuilder.build();
	}
}