package client

import (
	"strings"

	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/java/v2/packages"
)

func generateUtils(thePackage packages.Module) []generator.CodeFile {
	files := []generator.CodeFile{}

	files = append(files, *generateRequestBuilder(thePackage))
	files = append(files, *generateUrlBuilder(thePackage))
	files = append(files, *generateStringify(thePackage))

	return files
}

func generateRequestBuilder(thePackage packages.Module) *generator.CodeFile {
	code := `
package [[.PackageName]];

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
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("RequestBuilder.java"),
		Content: strings.TrimSpace(code),
	}
}

func generateUrlBuilder(thePackage packages.Module) *generator.CodeFile {
	code := `
package [[.PackageName]];

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
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("UrlBuilder.java"),
		Content: strings.TrimSpace(code),
	}
}

func generateStringify(thePackage packages.Module) *generator.CodeFile {
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

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("Stringify.java"),
		Content: strings.TrimSpace(code),
	}
}
