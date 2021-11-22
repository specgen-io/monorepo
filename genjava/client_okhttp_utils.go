package genjava

import (
	"github.com/specgen-io/specgen/v2/gen"
	"strings"
)

func generateUtils(thePackage Module) []gen.TextFile {
	files := []gen.TextFile{}

	files = append(files, *generateRequestBuilder(thePackage))
	files = append(files, *generateUrlBuilder(thePackage))
	files = append(files, *generateStringify(thePackage))

	return files
}

func generateRequestBuilder(thePackage Module) *gen.TextFile {
	code := `
package [[.PackageName]];

import okhttp3.*;

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

	public <T> RequestBuilder addHeaderParameter(String name, T[] values) {
		for (T val : values) {
			this.addHeaderParameter(name, val);
		}
		return this;
	}

	public Request build() {
		return this.requestBuilder.build();
	}
}
`

	code, _ = gen.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &gen.TextFile{
		Path:    thePackage.GetPath("RequestBuilder.java"),
		Content: strings.TrimSpace(code),
	}
}

func generateUrlBuilder(thePackage Module) *gen.TextFile {
	code := `
package [[.PackageName]];

import okhttp3.HttpUrl;

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

    public <T> UrlBuilder addQueryParameter(String name, T[] values) {
        for (T val : values) {
            this.addQueryParameter(name, val);
        }
        return this;
    }

    public UrlBuilder addPathSegment(Object value) {
        var valueStr = Stringify.paramToString(value);
        this.urlBuilder.addPathSegment(valueStr);
        return this;
    }

    public HttpUrl build() {
        return this.urlBuilder.build();
    }
}
`

	code, _ = gen.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &gen.TextFile{
		Path:    thePackage.GetPath("UrlBuilder.java"),
		Content: strings.TrimSpace(code),
	}
}

func generateStringify(thePackage Module) *gen.TextFile {
	code := `
package [[.PackageName]];

public class Stringify {
    public static String paramToString(Object value) {
        if (value == null) {
            return null;
        }
        return String.valueOf(value);
    }
}
`

	code, _ = gen.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &gen.TextFile{
		Path:    thePackage.GetPath("Stringify.java"),
		Content: strings.TrimSpace(code),
	}
}
