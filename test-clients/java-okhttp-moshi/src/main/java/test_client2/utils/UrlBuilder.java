package test_client2.utils;

import okhttp3.HttpUrl;
import java.util.List;

public class UrlBuilder {
		private final HttpUrl.Builder urlBuilder;

		public UrlBuilder(String baseUrl) {
				this.urlBuilder = HttpUrl.get(baseUrl).newBuilder();
		}

		public UrlBuilder addQueryParameter(String name, Object value) {
				var valueStr = Stringify.paramToString(value);
				if (valueStr != null) {
						this.urlBuilder.addQueryParameter(name, valueStr);
				}
				return this;
		}

		public <T> UrlBuilder addQueryParameter(String name, List<T> values) {
				for (T val : values) {
						this.addQueryParameter(name, val);
				}
				return this;
		}

		public UrlBuilder addPathSegments(String value) {
				this.urlBuilder.addPathSegments(value);
				return this;
		}

		public UrlBuilder addPathParameter(Object value) {
				var valueStr = Stringify.paramToString(value);
				this.urlBuilder.addPathSegment(valueStr);
				return this;
		}

		public HttpUrl build() {
				return this.urlBuilder.build();
		}
}