package client

import (
	"fmt"
	"generator"
	"kotlin/models"
	"kotlin/types"
	"kotlin/writer"
	"spec"
)

var OkHttp = "okhttp"

type OkHttpGenerator struct {
	Types    *types.Types
	Models   models.Generator
	Packages *Packages
}

func NewOkHttpGenerator(types *types.Types, models models.Generator, packages *Packages) *OkHttpGenerator {
	return &OkHttpGenerator{types, models, packages}
}

func (g *OkHttpGenerator) Clients(version *spec.Version) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, responses(&api, g.Types, g.Packages.Client(&api), g.Packages.Models(api.InHttp.InVersion), g.Packages.ErrorsModels)...)
		files = append(files, *g.client(&api))
	}
	return files
}

func (g *OkHttpGenerator) client(api *spec.Api) *generator.CodeFile {
	w := writer.New(g.Packages.Client(api), clientName(api))
	w.Imports.Add(g.Models.ModelsUsageImports()...)
	w.Imports.Add(g.Types.Imports()...)
	w.Imports.Add(`okhttp3.*`)
	w.Imports.Add(`okhttp3.MediaType.Companion.toMediaTypeOrNull`)
	w.Imports.Add(`okhttp3.RequestBody.Companion.toRequestBody`)
	w.Imports.Add(`org.slf4j.*`)
	w.Imports.PackageStar(g.Packages.Errors)
	w.Imports.PackageStar(g.Packages.Json)
	w.Imports.PackageStar(g.Packages.Utils)
	w.Imports.PackageStar(g.Packages.Models(api.InHttp.InVersion))
	w.Template(
		map[string]string{
			`JsonMapper`:     g.Models.JsonMapper()[0],
			`JsonMapperVar`:  g.Models.JsonMapper()[1],
			`InitJsonHelper`: g.Models.InitJsonHelper(),
		}, `
class [[.ClassName]] {
	private val logger: Logger = LoggerFactory.getLogger([[.ClassName]]::class.java)

	private var baseUrl: String
	private var client: OkHttpClient
	private var json: Json

	constructor(baseUrl: String, client: OkHttpClient, [[.JsonMapperVar]]: [[.JsonMapper]]) {
		this.baseUrl = baseUrl
		this.client = client
		this.json = Json([[.JsonMapperVar]])
	}

	constructor(baseUrl: String, client: OkHttpClient) : this(baseUrl, client, [[.InitJsonHelper]])

	constructor(baseUrl: String) {
		this.baseUrl = baseUrl
		this.json = Json([[.InitJsonHelper]])
		this.client = OkHttpClient().newBuilder().addInterceptor(ErrorsInterceptor(json)).build()
	}
`)
	for _, operation := range api.Operations {
		w.EmptyLine()
		g.generateClientMethod(w.Indented(), &operation)
	}
	w.Line(`}`)
	return w.ToCodeFile()
}

func (g *OkHttpGenerator) generateClientMethod(w *writer.Writer, operation *spec.NamedOperation) {
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	w.Line(`fun %s {`, operationSignature(g.Types, operation))
	w.Line(`  try {`)
	requestBody := "null"
	if operation.BodyIs(spec.BodyString) {
		w.Line(`    val requestBody = body.toRequestBody("text/plain".toMediaTypeOrNull())`)
		requestBody = "requestBody"
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`    val requestBody = json.%s.toRequestBody("application/json".toMediaTypeOrNull())`, g.Models.WriteJson("body", &operation.Body.Type.Definition))
		requestBody = "requestBody"
	}
	w.Line(`    val url = UrlBuilder(baseUrl)`)
	if operation.InApi.InHttp.GetUrl() != "" {
		w.Line(`    url.addPathSegments("%s")`, trimSlash(operation.InApi.InHttp.GetUrl()))
	}
	for _, urlPart := range operation.Endpoint.UrlParts {
		part := trimSlash(urlPart.Part)
		if urlPart.Param != nil {
			w.Line(`    url.addPathParameter(%s)`, urlPart.Param.Name.CamelCase())
		} else if len(part) > 0 {
			w.Line(`    url.addPathSegments("%s")`, part)
		}
	}
	for _, param := range operation.QueryParams {
		w.Line(`    url.addQueryParameter("%s", %s)`, param.Name.SnakeCase(), addBuilderParam(&param))
	}
	w.EmptyLine()
	w.Line(`    val request = RequestBuilder("%s", url.build(), %s)`, methodName, requestBody)
	for _, param := range operation.HeaderParams {
		w.Line(`    request.addHeaderParameter("%s", %s)`, param.Name.Source, addBuilderParam(&param))
	}
	w.EmptyLine()
	w.Line(`    logger.info("Sending request, operationId: %s.%s, method: %s, url: %s")`, operation.InApi.Name.Source, operation.Name.Source, methodName, url)
	w.Line(`    val response = client.newCall(request.build()).execute()`)
	w.Line(`    logger.info("Received response with status code ${response.code}")`)
	w.EmptyLine()
	w.IndentWith(2)
	w.Line(`when (response.code) {`)
	for _, response := range operation.Responses.Success() {
		var responseBody string
		if response.BodyIs(spec.BodyEmpty) {
			responseBody = responseCreate(response, "")
		}
		if response.BodyIs(spec.BodyString) {
			responseBody = responseCreate(response, "response.body!!.string()")
		}
		if response.BodyIs(spec.BodyJson) {
			responseBody = responseCreate(response, fmt.Sprintf(`json.%s`, g.Models.ReadJson(`response.body!!.charStream()`, &response.Type.Definition)))
		}
		w.Line(`  %s -> %s`, spec.HttpStatusCode(response.Name), responseBody)
	}
	for _, errorResponse := range operation.Responses.NonRequiredErrors() {
		w.Line(`  %s -> throw %sException()`, spec.HttpStatusCode(errorResponse.Name), errorResponse.Name.PascalCase())
	}
	w.Line(`  else -> throw ResponseException("Unexpected status code received: ${response.code}")`)
	w.Line(`}`)
	w.UnindentWith(2)
	w.Lines(`
	} catch (ex: Throwable) {
		logger.error(ex.message)
		throw ClientException(ex)
	}
`)
	w.Line(`}`)
}

func (g *OkHttpGenerator) Utils(responses *spec.ErrorResponses) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *g.generateRequestBuilder())
	files = append(files, *g.generateUrlBuilder())
	return files
}

func (g *OkHttpGenerator) generateRequestBuilder() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `RequestBuilder`)
	w.Lines(`
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

	fun <T> addHeaderParameter(name: String, values: List<T>): RequestBuilder {
		for (value in values) {
			this.addHeaderParameter(name, value!!)
		}
		return this
	}

	fun build(): Request {
		return this.requestBuilder.build()
	}
}
`)
	return w.ToCodeFile()
}

func (g *OkHttpGenerator) generateUrlBuilder() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `UrlBuilder`)
	w.Lines(`
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

	fun <T> addQueryParameter(name: String, values: List<T>): UrlBuilder {
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
`)
	return w.ToCodeFile()
}

func (g *OkHttpGenerator) Exceptions(errors *spec.ErrorResponses) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *clientException(g.Packages.Errors))
	files = append(files, *responseException(g.Packages.Errors))
	for _, errorResponse := range *errors {
		files = append(files, *inheritedClientException(g.Packages.Errors, g.Packages.ErrorsModels, g.Types, &errorResponse.Response))
	}
	files = append(files, *g.errorsInterceptor(errors))
	return files
}

func (g *OkHttpGenerator) errorsInterceptor(errorsResponses *spec.ErrorResponses) *generator.CodeFile {
	w := writer.New(g.Packages.Errors, `ErrorsInterceptor`)
	w.Imports.Add(`okhttp3.*`)
	w.Imports.Add(`org.slf4j.*`)
	w.Imports.PackageStar(g.Packages.ErrorsModels)
	w.Imports.PackageStar(g.Packages.Json)
	w.Lines(`
class [[.ClassName]](private var json: Json) : Interceptor {
	private val logger: Logger = LoggerFactory.getLogger([[.ClassName]]::class.java)

	override fun intercept(chain: Interceptor.Chain): Response {
		val request: Request = chain.request()
		val response: Response = chain.proceed(request)

		when (response.code) {
`)
	w.IndentWith(3)
	for _, errorResponse := range errorsResponses.Required() {
		w.Line(`%s -> {`, spec.HttpStatusCode(errorResponse.Name))
		responseBody := "responseBody"
		if errorResponse.BodyIs(spec.BodyEmpty) {
			responseBody = ""
		}
		if errorResponse.BodyIs(spec.BodyString) {
			w.Line(`  val %s = response.body!!.string()`, responseBody)
			w.Line(`  logger.error(%s)`, responseBody)
		}
		if errorResponse.BodyIs(spec.BodyJson) {
			w.Line(`  val %s: %s = json.%s`, responseBody, g.Types.Kotlin(&errorResponse.Type.Definition), g.Models.ReadJson(`response.body!!.charStream()`, &errorResponse.Type.Definition))
			w.Line(`  logger.error(%s.message)`, responseBody)
		}
		w.Line(`  throw %sException(%s)`, errorResponse.Name.PascalCase(), responseBody)
		w.Line(`}`)
	}
	w.UnindentWith(3)
	w.Lines(`
		}
		return response
	}
}
`)
	return w.ToCodeFile()
}
