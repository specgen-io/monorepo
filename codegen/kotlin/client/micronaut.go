package client

import (
	"fmt"
	"generator"
	"kotlin/models"
	"kotlin/types"
	"kotlin/writer"
	"spec"
	"strings"
)

var Micronaut = "micronaut-low-level"

type MicronautGenerator struct {
	Types    *types.Types
	Models   models.Generator
	Packages *Packages
}

func NewMicronautGenerator(types *types.Types, models models.Generator, packages *Packages) *MicronautGenerator {
	return &MicronautGenerator{types, models, packages}
}

func (g *MicronautGenerator) Clients(version *spec.Version) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, responses(&api, g.Types, g.Packages.Client(&api), g.Packages.Models(api.InHttp.InVersion), g.Packages.ErrorsModels)...)
		files = append(files, *g.client(&api))
	}
	files = append(files, converters(g.Packages.Converters)...)
	return files
}

func (g *MicronautGenerator) client(api *spec.Api) *generator.CodeFile {
	w := writer.New(g.Packages.Client(api), clientName(api))
	w.Imports.Add(g.Models.ModelsUsageImports()...)
	w.Imports.Add(g.Types.Imports()...)
	w.Imports.Add(`io.micronaut.http.HttpHeaders.*`)
	w.Imports.Add(`io.micronaut.http.HttpRequest.*`)
	w.Imports.Add(`io.micronaut.http.client.HttpClient`)
	w.Imports.Add(`java.net.URL`)
	w.Imports.Add(`org.slf4j.*`)
	w.Imports.PackageStar(g.Packages.Errors)
	w.Imports.PackageStar(g.Packages.Json)
	w.Imports.PackageStar(g.Packages.Utils)
	w.Imports.PackageStar(g.Packages.Models(api.InHttp.InVersion))
	w.Imports.Package(g.Packages.Utils.Subpackage(`ClientResponse`).Subpackage(`doRequest`))
	w.Imports.Package(g.Packages.Utils.Subpackage(`ClientResponse`).Subpackage(`getResponseBodyString`))
	w.EmptyLine()
	w.Lines(`
class [[.ClassName]](private val baseUrl: String) {
	private val logger: Logger = LoggerFactory.getLogger([[.ClassName]]::class.java)

    private var client: HttpClient
    private val json: Json

	init {
`)
	w.IndentedWith(2).Lines(g.Models.CreateJsonHelper(`json`))
	w.Lines(`
        client = HttpClient.create(URL(baseUrl))
	}
`)
	for _, operation := range api.Operations {
		w.EmptyLine()
		g.generateClientMethod(w.Indented(), &operation)
	}
	w.Line(`}`)
	return w.ToCodeFile()
}

func (g *MicronautGenerator) generateClientMethod(w *writer.Writer, operation *spec.NamedOperation) {
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
	w.Line(`  val response = doRequest(client, request, logger)`)
	w.EmptyLine()
	for _, response := range operation.SuccessResponses() {
		w.Line(`  if (response.code() == %s) {`, spec.HttpStatusCode(response.Name))
		w.IndentWith(2)
		w.Line(`logger.info("Received response with status code {}", response.code())`)
		if response.BodyIs(spec.BodyEmpty) {
			w.Line(responseCreate(response, ""))
		}
		if response.BodyIs(spec.BodyString) {
			responseBodyString := "getResponseBodyString(response, logger)"
			w.Line(responseCreate(response, responseBodyString))
		}
		if response.BodyIs(spec.BodyJson) {
			w.Line(`val responseBodyString = getResponseBodyString(response, logger)`)
			responseBody := fmt.Sprintf(`json.%s`, g.Models.JsonRead("responseBodyString", &response.Type.Definition))
			w.Line(responseCreate(response, responseBody))
		}
		w.UnindentWith(2)
		w.Line(`  }`)
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

func (g *MicronautGenerator) Utils(responses *spec.ErrorResponses) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *g.generateRequestBuilder())
	files = append(files, *g.generateUrlBuilder())
	files = append(files, *g.generateClientResponse())
	files = append(files, *g.generateErrorsHandler(responses))
	return files
}

func (g *MicronautGenerator) generateRequestBuilder() *generator.CodeFile {
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

func (g *MicronautGenerator) generateUrlBuilder() *generator.CodeFile {
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

func (g *MicronautGenerator) generateClientResponse() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `ClientResponse`)
	w.Template(
		map[string]string{
			`ErrorsPackage`: g.Packages.Errors.PackageName,
		}, `
import io.micronaut.http.HttpResponse
import io.micronaut.http.client.HttpClient
import io.micronaut.http.client.exceptions.HttpClientResponseException
import org.slf4j.Logger
import [[.ErrorsPackage]].*
import java.io.IOException

object ClientResponse {
	fun doRequest(client: HttpClient, request: RequestBuilder, logger: Logger): HttpResponse<String> {
		return try {
			client.toBlocking().exchange(request.build(), String::class.java)
		} catch (e: HttpClientResponseException) {
			val errorMessage = "Failed to execute the request " + e.message
			logger.error(errorMessage)
			throw ClientException(errorMessage, e)
		}
	}

	fun <T> getResponseBodyString(response: HttpResponse<T>, logger: Logger): String {
		return try {
			response.body()!!.toString()
		} catch (e: IOException) {
			val errorMessage = "Failed to convert response body to string " + e.message
			logger.error(errorMessage)
			throw ClientException(errorMessage, e)
		}
	}
}
`)
	return w.ToCodeFile()
}

func (g *MicronautGenerator) generateErrorsHandler(errorsResponses *spec.ErrorResponses) *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `ErrorsHandler`)
	w.Imports.Add(g.Models.ModelsUsageImports()...)
	w.Imports.Add(`io.micronaut.http.*`)
	w.Imports.Add(`org.slf4j.*`)
	w.Imports.PackageStar(g.Packages.Errors)
	w.Imports.PackageStar(g.Packages.ErrorsModels)
	w.Imports.PackageStar(g.Packages.Json)
	w.Imports.Package(g.Packages.Utils.Subpackage(`ClientResponse`).Subpackage(`getResponseBodyString`))
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

func (g *MicronautGenerator) Exceptions(errors *spec.ErrorResponses) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *clientException(g.Packages.Errors))
	for _, errorResponse := range *errors {
		if errorResponse.Required {
			files = append(files, *inheritedClientException(g.Packages.Errors, g.Packages.ErrorsModels, g.Types, &errorResponse.Response))
		}
	}
	return files
}
