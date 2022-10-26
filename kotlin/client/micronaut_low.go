package client

import (
	"fmt"
	"generator"
	"kotlin/imports"
	"kotlin/models"
	"kotlin/types"
	"kotlin/writer"
	"spec"
	"strings"
)

var MicronautLow = "micronaut-low-level"

type MicronautLowGenerator struct {
	Types    *types.Types
	Models   models.Generator
	Packages *Packages
}

func NewMicronautLowGenerator(types *types.Types, models models.Generator, packages *Packages) *MicronautLowGenerator {
	return &MicronautLowGenerator{types, models, packages}
}

func (g *MicronautLowGenerator) Clients(version *spec.Version) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, responses(&api, g.Types, g.Packages.Client(&api), g.Packages.Models(api.InHttp.InVersion), g.Packages.ErrorsModels)...)
		files = append(files, *g.client(&api))
	}
	files = append(files, converters(g.Packages.Converters)...)
	files = append(files, staticConfigFiles(g.Packages.Root, g.Packages.Json)...)
	return files
}

func (g *MicronautLowGenerator) client(api *spec.Api) *generator.CodeFile {
	w := writer.New(g.Packages.Client(api), clientName(api))
	imports := imports.New()
	imports.Add(g.Models.ModelsUsageImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Add(`io.micronaut.http.HttpHeaders.*`)
	imports.Add(`io.micronaut.http.HttpRequest.*`)
	imports.Add(`io.micronaut.http.client.HttpClient`)
	imports.Add(`io.micronaut.http.client.annotation.Client`)
	imports.Add(`jakarta.inject.Singleton`)
	imports.Add(`org.slf4j.*`)
	imports.Add(g.Packages.Root.PackageStar)
	imports.Add(g.Packages.Utils.PackageStar)
	imports.Add(g.Packages.Models(api.InHttp.InVersion).PackageStar)
	imports.Write(w)
	w.EmptyLine()
	w.Lines(`
@Singleton
class [[.ClassName]](
@param:Client(ClientConfiguration.BASE_URL)
private val client: HttpClient,
private val json: Json
) {
	private val logger: Logger = LoggerFactory.getLogger([[.ClassName]]::class.java)

    init {
`)
	w.IndentedWith(2).Lines(g.Models.CreateJsonHelper(`json`))
	w.Lines(`
		client = OkHttpClient()
	}
`)
	for _, operation := range api.Operations {
		w.EmptyLine()
		g.generateClientMethod(w.Indented(), &operation)
	}
	w.Line(`}`)
	return w.ToCodeFile()
}

func (g *MicronautLowGenerator) generateClientMethod(w generator.Writer, operation *spec.NamedOperation) {
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	w.Line(`fun %s {`, operationSignature(g.Types, operation))
	requestBody := "body"
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`  val bodyJson = json.%s`, g.Models.JsonWrite("body", &operation.Body.Type.Definition))
		requestBody = "bodyJson"
	}
	w.Line(`  val url = UrlBuilder("%s")`, getUrl(operation))
	for _, urlPart := range operation.Endpoint.UrlParts {
		if urlPart.Param != nil {
			w.Line(`  url.pathParam("%s", %s)`, urlPart.Param.Name.CamelCase(), urlPart.Param.Name.CamelCase())
		}
	}
	for _, param := range operation.QueryParams {
		w.Line(`  url.queryParam("%s", %s)`, param.Name.SnakeCase(), addBuilderParam(&param))
	}
	w.EmptyLine()
	w.Line(`  val request = RequestBuilder(%s)`, requestBuilderParams(methodName, requestBody, operation))
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`  request.headerParam(CONTENT_TYPE, "application/json")`)
	}
	if operation.BodyIs(spec.BodyString) {
		w.Line(`  request.headerParam(CONTENT_TYPE, "text/plain")`)
	}
	for _, param := range operation.HeaderParams {
		w.Line(`  request.headerParam("%s", %s)`, param.Name.Source, addBuilderParam(&param))
	}
	w.EmptyLine()
	w.Line(`  logger.info("Sending request, operationId: %s.%s, method: %s, url: %s")`, operation.InApi.Name.Source, operation.Name.Source, methodName, url)
	w.Line(`  val response = doRequest(client, request)`)
	w.EmptyLine()
	for _, response := range operation.Responses {
		statusCode := spec.HttpStatusCode(response.Name)
		if isSuccessfulStatusCode(statusCode) {
			w.Line(`  if (response.code == %s) {`, statusCode)
			w.IndentWith(2)
			w.Line(`logger.info("Received response with status code {}", response.code)`)
			if response.BodyIs(spec.BodyEmpty) {
				w.Line(responseCreate(&response, ""))
			}
			if response.BodyIs(spec.BodyString) {
				responseBodyString := "getResponseBodyString(response, logger)"
				w.Line(responseCreate(&response, responseBodyString))
			}
			if response.BodyIs(spec.BodyJson) {
				w.Line(`val responseBodyString = getResponseBodyString(response, logger)`)
				responseBody := fmt.Sprintf(`json.%s`, g.Models.JsonRead("responseBodyString", &response.Type.Definition))
				w.Line(responseCreate(&response, responseBody))
			}
			w.UnindentWith(2)
			w.Line(`  }`)
		}
	}
	w.Line(`  handleErrors(response, logger, json)`)
	w.EmptyLine()
	generateThrowClientException(w.Indented(), `"Unexpected status code received: " + response.code()`, ``)
	w.Line(`}`)
}

func getUrl(operation *spec.NamedOperation) string {
	url := strings.TrimRight(operation.Endpoint.UrlParts[0].Part, "/")
	if operation.InApi.InHttp.GetUrl() != "" {
		return strings.TrimRight(operation.InApi.InHttp.GetUrl(), "/") + url
	}
	return url
}

func requestBuilderParams(methodName, requestBody string, operation *spec.NamedOperation) string {
	urlParam := "url.build()"
	if &operation.Endpoint.UrlParams != nil {
		urlParam = "url.expand()"
	}
	params := fmt.Sprintf(`%s, %s, ::%s`, urlParam, requestBody, methodName)
	if operation.BodyIs(spec.BodyEmpty) {
		params = fmt.Sprintf(`%s, ::%s`, urlParam, methodName)
	}

	return params
}

func (g *MicronautLowGenerator) Utils(responses *spec.Responses) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *g.generateRequestBuilder())
	files = append(files, *g.generateUrlBuilder())
	files = append(files, *g.generateClientResponse())
	files = append(files, *g.generateErrorsHandler(responses))
	return files
}

func (g *MicronautLowGenerator) generateRequestBuilder() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `RequestBuilder`)
	w.Lines(`
import io.micronaut.http.MutableHttpRequest
import java.net.URI

class RequestBuilder {
    private var generateRequestBuilder: MutableHttpRequest<Any>

    constructor(url: URI, body: Any?, method: (URI, Any?) -> MutableHttpRequest<Any>) {
        this.generateRequestBuilder = method(url, body)
    }

    constructor(url: URI, method: (URI) -> MutableHttpRequest<Any>) {
        this.generateRequestBuilder = method(url)
    }

    fun headerParam(name: String, value: Any): RequestBuilder {
        val valueStr = value.toString()
        this.generateRequestBuilder.header(name, valueStr)
        return this
    }

    fun <T> headerParam(name: String, values: List<T>): RequestBuilder {
        for (value in values) {
            this.headerParam(name, value!!)
        }
        return this
    }

    fun build(): MutableHttpRequest<Any> {
        return this.generateRequestBuilder
    }
}
`)
	return w.ToCodeFile()
}

func (g *MicronautLowGenerator) generateUrlBuilder() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `UrlBuilder`)
	w.Lines(`
import io.micronaut.http.uri.UriBuilder
import java.net.URI
import java.util.*

class UrlBuilder(url: String) {
    private val uriBuilder: UriBuilder = UriBuilder.of(url)
    private val urlMap: MutableMap<String, Any> = mutableMapOf()

    fun queryParam(name: String, value: Any): UrlBuilder {
        val valueStr = value.toString()
        uriBuilder.queryParam(name, valueStr)
        return this
    }

    fun <T> queryParam(name: String, values: List<T>): UrlBuilder {
        for (value in values) {
            this.queryParam(name, value!!)
        }
        return this
    }

    fun pathParam(name: String, value: Any): UrlBuilder {
        this.uriBuilder.path("{$name}")
        this.urlMap += mapOf(name to value)

        return this
    }

    fun expand(): URI {
        return this.uriBuilder.expand(
            Collections.checkedMap(
                this.urlMap, String::class.java, Any::class.java
            )
        )
    }

    fun build(): URI {
        return this.uriBuilder.build()
    }
}
`)
	return w.ToCodeFile()
}

func (g *MicronautLowGenerator) generateClientResponse() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `ClientResponse`)
	w.Template(
		map[string]string{
			`ErrorsPackage`: g.Packages.Errors.PackageName,
		}, `
import io.micronaut.http.HttpResponse
import io.micronaut.http.client.HttpClient
import org.slf4j.Logger
import [[.ErrorsPackage]].*
import java.io.IOException

object ClientResponse {
    fun doRequest(client: HttpClient, request: RequestBuilder): HttpResponse<String> =
        client.toBlocking().exchange(request.build(), String::class.java)

    fun <T> getResponseBodyString(response: HttpResponse<T>, logger: Logger): String {
		return try {
            response.body()!!.toString()
		} catch (e: IOException) {
            val errorMessage = "Failed to convert response body to string " + e.message
            logger.error(errorMessage)
            throw ClientException(errorMessage, e)
        }
        return responseBodyString
    }
}
`)
	return w.ToCodeFile()
}

func (g *MicronautLowGenerator) generateErrorsHandler(errorsResponses *spec.Responses) *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `ErrorsHandler`)
	imports := imports.New()
	imports.Add(g.Models.ModelsUsageImports()...)
	imports.Add(`io.micronaut.http.*`)
	imports.Add(`org.slf4j.*`)
	imports.Add(g.Packages.ErrorsModels.PackageStar)
	imports.Add(g.Packages.Json.PackageStar)
	imports.Add(g.Packages.Utils.Subpackage(`ClientResponse`).PackageStar)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`fun <T> handleErrors(response: HttpResponse<T>, logger: Logger, json: Json) {`)
	for _, errorResponse := range *errorsResponses {
		w.Line(`  if (response.code() == %s) {`, spec.HttpStatusCode(errorResponse.Name))
		w.Line(`    val responseBodyString = getResponseBodyString(response, logger)`)
		w.Line(`    val responseBody = json.%s`, g.Models.JsonRead("responseBodyString", &errorResponse.Type.Definition))
		w.Line(`    throw %sException(responseBody)`, g.Types.Kotlin(&errorResponse.Type.Definition))
		w.Line(`  }`)
	}
	w.Line(`}`)

	return w.ToCodeFile()
}

func (g *MicronautLowGenerator) Exceptions(errors *spec.Responses) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *clientException(g.Packages.Errors))
	for _, errorResponse := range *errors {
		files = append(files, *inheritedClientException(g.Packages.Errors, g.Packages.ErrorsModels, g.Types, &errorResponse))
	}
	return files
}
