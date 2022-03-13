package client

import (
	"github.com/specgen-io/specgen/v2/gen/kotlin/modules"
	"github.com/specgen-io/specgen/v2/sources"
	"strings"
)

func generateUtils(thePackage modules.Module) []sources.CodeFile {
	files := []sources.CodeFile{}

	files = append(files, *generateRequestBuilder(thePackage))
	files = append(files, *generateUrlBuilder(thePackage))

	return files
}

func generateRequestBuilder(thePackage modules.Module) *sources.CodeFile {
	code := `
package [[.PackageName]]

import okhttp3.*

class RequestBuilder(method: String, url: HttpUrl, body: RequestBody?) {
    private val requestBuilder: Request.Builder

    init {
        requestBuilder = Request.Builder().url(url).method(method, body)
    }

    fun addHeaderParameter(name: String, value: Any): RequestBuilder {
        val valueStr = value.toString()
        this.requestBuilder.addHeader(name, valueStr)
        return this
    }

    fun <T> addHeaderParameter(name: String, values: Array<T>): RequestBuilder {
        for (value in values) {
            this.addHeaderParameter(name, value!!)
        }
        return this
    }

    fun build(): Request {
        return this.requestBuilder.build()
    }
}
`

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("RequestBuilder.kt"),
		Content: strings.TrimSpace(code),
	}
}

func generateUrlBuilder(thePackage modules.Module) *sources.CodeFile {
	code := `
package [[.PackageName]]

import okhttp3.HttpUrl
import okhttp3.HttpUrl.Companion.toHttpUrl

class UrlBuilder(baseUrl: String) {
    private val urlBuilder: HttpUrl.Builder

    init {
        this.urlBuilder = baseUrl.toHttpUrl().newBuilder()
    }

    fun addQueryParameter(name: String, value: Any): UrlBuilder {
        val valueStr = value.toString()
        urlBuilder.addQueryParameter(name, valueStr)
        return this
    }

    fun <T> addQueryParameter(name: String, values: Array<T>): UrlBuilder {
        for (value in values) {
            this.addQueryParameter(name, value!!)
        }
        return this
    }

    fun addPathSegments(value: String): UrlBuilder {
        this.urlBuilder.addPathSegments(value)
        return this
    }

    fun addPathParameter(value: Any): UrlBuilder {
        val valueStr = value.toString()
        this.urlBuilder.addPathSegment(valueStr)
        return this
    }

    fun build(): HttpUrl {
        return this.urlBuilder.build()
    }
}
`

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("UrlBuilder.kt"),
		Content: strings.TrimSpace(code),
	}
}
