package genjava

import (
	"github.com/specgen-io/specgen/v2/gen"
	"strings"
)

func generateRequestBuilder(thePackage Module) *gen.TextFile {
	code := `
package [[.PackageName]];

import okhttp3.*;

public class RequestBuilder {
	private final HttpUrl.Builder urlBuilder;
	private final Request.Builder requestBuilder;

	public RequestBuilder(String url) {
		this.requestBuilder = new Request.Builder();
		this.urlBuilder = HttpUrl.get(url).newBuilder();
	}

	public void post(RequestBody requestBody) {
		this.requestBuilder.post(requestBody);
	}

	public static String paramToString(Object value) {
		if (value == null) {
			return null;
		}
		return String.valueOf(value);
	}

	public void addQueryParameter(String name, Object value) {
		var valueStr = paramToString(value);
		if (valueStr != null) {
			this.urlBuilder.addQueryParameter(name, valueStr);
		}
	}

	public <T> void addQueryParameter(String name, T[] values) {
		for (T val : values) {
			this.addQueryParameter(name, val);
		}
	}

	public void addHeaderParameter(String name, Object value) {
		var valueStr = paramToString(value);
		if (valueStr != null) {
			this.requestBuilder.addHeader(name, valueStr);
		}
	}

	public <T> void addHeaderParameter(String name, T[] values) {
		for (T val : values) {
			this.addHeaderParameter(name, val);
		}
	}

	public void addPathSegment(Object value) {
		var valueStr = paramToString(value);
		this.urlBuilder.addPathSegment(valueStr);
	}

	public Request build() {
		var url = this.urlBuilder.build();
		return this.requestBuilder.url(url).build();
	}
}
`

	code, _ = gen.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &gen.TextFile{
		Path:    thePackage.GetPath("RequestBuilder.java"),
		Content: strings.TrimSpace(code),
	}
}

func generateClientException(thePackage Module) *gen.TextFile {
	code := `
package [[.PackageName]];

public class ClientException extends RuntimeException {
	public ClientException() {
		super();
	}

	public ClientException(String message) {
		super(message);
	}

	public ClientException(String message, Throwable cause) {
		super(message, cause);
	}

	public ClientException(Throwable cause) {
		super(cause);
	}
}
`

	code, _ = gen.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &gen.TextFile{
		Path:    thePackage.GetPath("ClientException.java"),
		Content: strings.TrimSpace(code),
	}
}
