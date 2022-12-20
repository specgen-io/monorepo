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
	w.Imports.Add(`io.micronaut.http.client.*`)
	w.Imports.Add(`java.net.URL`)
	w.Imports.Add(`org.slf4j.*`)
	w.Imports.PackageStar(g.Packages.Errors)
	w.Imports.PackageStar(g.Packages.Json)
	w.Imports.PackageStar(g.Packages.Utils)
	w.Imports.PackageStar(g.Packages.Models(api.InHttp.InVersion))
	w.Template(
		map[string]string{
			`JsonMapperType`: g.Models.JsonMapperType(),
			`JsonMapperInit`: g.Models.JsonMapperInit(),
		}, `
class [[.ClassName]] {
    private val logger: Logger = LoggerFactory.getLogger([[.ClassName]]::class.java)

    private var client: BlockingHttpClient
    private var json: Json
    private val errorsHandler: ErrorsHandler

    constructor(client: HttpClient, mapper: [[.JsonMapperType]]) {
        this.client = client.toBlocking()
        this.json = Json(mapper)
        this.errorsHandler = ErrorsHandler(json)
    }

    constructor(client: HttpClient) : this(client, [[.JsonMapperInit]])

    constructor(baseUrl: String) {
        this.json = Json([[.JsonMapperInit]])
        this.client = HttpClient.create(URL(baseUrl)).toBlocking()
        this.errorsHandler = ErrorsHandler(json)
    }
`)
	for _, operation := range api.Operations {
		w.EmptyLine()
		g.generateClientMethod(w.Indented(), &operation)
	}
	w.Line(`}`)
	return w.ToCodeFile()
}

func (g *MicronautGenerator) createUrl(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line(`val url = UrlBuilder("%s")`, getUrl(operation))
	for _, urlPart := range operation.Endpoint.UrlParts {
		if urlPart.Param != nil {
			w.Line(`url.pathParam("%s", %s)`, urlPart.Param.Name.CamelCase(), urlPart.Param.Name.CamelCase())
		}
	}
	for _, param := range operation.QueryParams {
		w.Line(`url.queryParam("%s", %s)`, param.Name.SnakeCase(), addBuilderParam(&param))
	}
}

func (g *MicronautGenerator) createRequest(w *writer.Writer, operation *spec.NamedOperation) {
	requestBody := "body"
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`val bodyJson = json.%s`, g.Models.WriteJson("body", &operation.Body.Type.Definition))
		requestBody = "bodyJson"
	}
	w.EmptyLine()
	w.Line(`val request = RequestBuilder(%s)`, requestBuilderParams(operation.Endpoint.Method, requestBody, operation))
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`request.headerParam(CONTENT_TYPE, "application/json")`)
	}
	if operation.BodyIs(spec.BodyString) {
		w.Line(`request.headerParam(CONTENT_TYPE, "text/plain")`)
	}
	for _, param := range operation.HeaderParams {
		w.Line(`request.headerParam("%s", %s)`, param.Name.Source, addBuilderParam(&param))
	}
}

func (g *MicronautGenerator) sendRequest(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line(`logger.info("Sending request, operationId: %s.%s, method: %s, url: %s")`, operation.InApi.Name.Source, operation.Name.Source, operation.Endpoint.Method, operation.FullUrl())
	w.Line(`val response = client.sendRequest(request)`)
	w.Line(`logger.info("Received response with status code ${response.code()}")`)
}

func (g *MicronautGenerator) processResponse(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line(`when (response.code()) {`)
	for _, response := range operation.Responses.Success() {
		w.Line(`    %s -> %s`, spec.HttpStatusCode(response.Name), g.successResponse(response))
	}
	for _, response := range operation.Responses.NonRequiredErrors() {
		w.Line(`    %s -> %s`, spec.HttpStatusCode(response.Name), g.errorResponse(&response.Response))
	}
	w.Line(`    else -> throw ResponseException("Unexpected status code received: ${response.code()}")`)
	w.Line(`}`)
}

func (g *MicronautGenerator) generateClientMethod(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line(`fun %s {`, operationSignature(g.Types, operation))
	w.Line(`    try {`)
	w.IndentWith(2)
	g.createUrl(w, operation)
	g.createRequest(w, operation)
	w.EmptyLine()
	g.sendRequest(w, operation)
	w.EmptyLine()
	w.Line(`errorsHandler.handle(response)`)
	w.EmptyLine()
	g.processResponse(w, operation)
	w.UnindentWith(2)
	w.Lines(`
    } catch (ex: Throwable) {
        logger.error(ex.message)
        throw ClientException(ex)
    }
`)
	w.Line(`}`)
}

func (g *MicronautGenerator) successResponse(response *spec.OperationResponse) string {
	if response.BodyIs(spec.BodyString) {
		return responseCreate(response, "response.body()!!.toString()")
	}
	if response.BodyIs(spec.BodyJson) {
		return responseCreate(response, fmt.Sprintf(`json.%s`, g.Models.ReadJson(`response.body()!!.toString()`, &response.Type.Definition)))
	}
	return responseCreate(response, "")
}

func (g *MicronautGenerator) errorResponse(response *spec.Response) string {
	var responseBody = ""
	if response.BodyIs(spec.BodyString) {
		responseBody = "response.body()!!.string()"
	}
	if response.BodyIs(spec.BodyJson) {
		responseBody = fmt.Sprintf(`json.%s`, g.Models.ReadJson(`response.body()!!.toString()`, &response.Type.Definition))
	}
	return fmt.Sprintf(`throw %s(%s)`, errorExceptionClassName(response), responseBody)
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

func (g *MicronautGenerator) Utils() []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *g.generateRequestBuilder())
	files = append(files, *g.generateUrlBuilder())
	files = append(files, *g.generateRequestor())
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

func (g *MicronautGenerator) generateRequestor() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `Requestor`)
	w.Lines(`
import io.micronaut.http.HttpResponse
import io.micronaut.http.client.BlockingHttpClient
import io.micronaut.http.client.exceptions.HttpClientResponseException

fun BlockingHttpClient.sendRequest(request: RequestBuilder): HttpResponse<*> {
    return try {
        this.exchange(request.build(), String::class.java)
    } catch (ex: HttpClientResponseException) {
        ex.response
    }
}
`)
	return w.ToCodeFile()
}

func (g *MicronautGenerator) Exceptions(errors *spec.ErrorResponses) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *clientException(g.Packages.Errors))
	files = append(files, *responseException(g.Packages.Errors))
	for _, errorResponse := range *errors {
		files = append(files, *errorResponseException(g.Packages.Errors, g.Packages.ErrorsModels, &errorResponse.Response))
	}
	files = append(files, *g.errorsHandler(errors))
	return files
}

func (g *MicronautGenerator) errorsHandler(errorsResponses *spec.ErrorResponses) *generator.CodeFile {
	w := writer.New(g.Packages.Errors, `ErrorsHandler`)
	w.Imports.Add(`io.micronaut.http.HttpResponse`)
	w.Imports.PackageStar(g.Packages.Json)
	w.Lines(`
class [[.ClassName]](private val json: Json) {
	fun handle(response: HttpResponse<*>) {
`)
	w.IndentWith(2)
	w.Line(`when (response.code()) {`)
	for _, errorResponse := range errorsResponses.Required() {
		w.Line(`    %s -> %s`, spec.HttpStatusCode(errorResponse.Name), g.errorResponse(&errorResponse.Response))
	}
	w.Line(`}`)
	w.UnindentWith(2)
	w.Lines(`
	}
}
`)
	return w.ToCodeFile()
}
