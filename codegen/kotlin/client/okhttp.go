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
	w.Imports.Package(g.Packages.Utils.Subpackage(`ClientResponse`).Subpackage(`doRequest`))
	w.Imports.Package(g.Packages.Utils.Subpackage(`ClientResponse`).Subpackage(`getResponseBodyString`))
	w.EmptyLine()
	w.Lines(`
class [[.ClassName]](private val baseUrl: String) {
	private val logger: Logger = LoggerFactory.getLogger([[.ClassName]]::class.java)

	private val client: OkHttpClient
	private val json: Json

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

func (g *OkHttpGenerator) generateClientMethod(w *writer.Writer, operation *spec.NamedOperation) {
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	w.Line(`fun %s {`, operationSignature(g.Types, operation))
	requestBody := "null"
	if operation.BodyIs(spec.BodyString) {
		w.Line(`  val requestBody = body.toRequestBody("text/plain".toMediaTypeOrNull())`)
		requestBody = "requestBody"
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`  val bodyJson = json.%s`, g.Models.JsonWrite("body", &operation.Body.Type.Definition))
		w.Line(`  val requestBody = bodyJson.toRequestBody("application/json".toMediaTypeOrNull())`)
		requestBody = "requestBody"
	}
	w.Line(`  val url = UrlBuilder(baseUrl)`)
	if operation.InApi.InHttp.GetUrl() != "" {
		w.Line(`  url.addPathSegments("%s")`, trimSlash(operation.InApi.InHttp.GetUrl()))
	}
	for _, urlPart := range operation.Endpoint.UrlParts {
		part := trimSlash(urlPart.Part)
		if urlPart.Param != nil {
			w.Line(`  url.addPathParameter(%s)`, urlPart.Param.Name.CamelCase())
		} else if len(part) > 0 {
			w.Line(`  url.addPathSegments("%s")`, part)
		}
	}
	for _, param := range operation.QueryParams {
		w.Line(`  url.addQueryParameter("%s", %s)`, param.Name.SnakeCase(), addBuilderParam(&param))
	}
	w.EmptyLine()
	w.Line(`  val request = RequestBuilder("%s", url.build(), %s)`, methodName, requestBody)
	for _, param := range operation.HeaderParams {
		w.Line(`  request.addHeaderParameter("%s", %s)`, param.Name.Source, addBuilderParam(&param))
	}
	w.EmptyLine()
	w.Line(`  logger.info("Sending request, operationId: %s.%s, method: %s, url: %s")`, operation.InApi.Name.Source, operation.Name.Source, methodName, url)
	w.Line(`  val response = doRequest(client, request, logger)`)
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
	generateThrowClientException(w.Indented(), `"Unexpected status code received: " + response.code`, ``)
	w.Line(`}`)
}

func (g *OkHttpGenerator) Utils(responses *spec.Responses) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *g.generateRequestBuilder())
	files = append(files, *g.generateUrlBuilder())
	files = append(files, *g.generateClientResponse())
	files = append(files, *g.generateErrorsHandler(responses))
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

func (g *OkHttpGenerator) generateClientResponse() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `ClientResponse`)
	w.Template(
		map[string]string{
			`ErrorsPackage`: g.Packages.Errors.PackageName,
		}, `
import okhttp3.*
import org.slf4j.Logger
import [[.ErrorsPackage]].*
import java.io.IOException

object ClientResponse {
	fun doRequest(client: OkHttpClient, request: RequestBuilder, logger: Logger): Response {
		return try {
			client.newCall(request.build()).execute()
		} catch (e: IOException) {
			val errorMessage = "Failed to execute the request " + e.message
			logger.error(errorMessage)
			throw ClientException(errorMessage, e)
		}
	}

	fun getResponseBodyString(response: Response, logger: Logger): String {
		return try {
			response.body!!.string()
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

func (g *OkHttpGenerator) generateErrorsHandler(errorsResponses *spec.Responses) *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `ErrorsHandler`)
	w.Imports.Add(g.Models.ModelsUsageImports()...)
	w.Imports.Add(`okhttp3.*`)
	w.Imports.Add(`org.slf4j.*`)
	w.Imports.PackageStar(g.Packages.Errors)
	w.Imports.PackageStar(g.Packages.ErrorsModels)
	w.Imports.PackageStar(g.Packages.Json)
	w.Imports.Package(g.Packages.Utils.Subpackage(`ClientResponse`).Subpackage(`getResponseBodyString`))
	w.EmptyLine()
	w.Line(`fun handleErrors(response: Response, logger: Logger, json: Json) {`)
	for _, errorResponse := range *errorsResponses {
		w.Line(`  if (response.code == %s) {`, spec.HttpStatusCode(errorResponse.Name))
		w.Line(`    val responseBodyString = getResponseBodyString(response, logger)`)
		w.Line(`    val responseBody = json.%s`, g.Models.JsonRead("responseBodyString", &errorResponse.Type.Definition))
		w.Line(`    throw %sException(responseBody)`, g.Types.Kotlin(&errorResponse.Type.Definition))
		w.Line(`  }`)
	}
	w.Line(`}`)

	return w.ToCodeFile()
}

func (g *OkHttpGenerator) Exceptions(errors *spec.Responses) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *clientException(g.Packages.Errors))
	for _, errorResponse := range *errors {
		files = append(files, *inheritedClientException(g.Packages.Errors, g.Packages.ErrorsModels, g.Types, &errorResponse))
	}
	return files
}
