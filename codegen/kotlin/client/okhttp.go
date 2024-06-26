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
			`JsonMapperType`: g.Models.JsonMapperType(),
			`JsonMapperInit`: g.Models.JsonMapperInit(),
		}, `
class [[.ClassName]] {
    private val logger: Logger = LoggerFactory.getLogger([[.ClassName]]::class.java)

    private var baseUrl: String
    private var client: OkHttpClient
    private var json: Json

    constructor(baseUrl: String, client: OkHttpClient, mapper: [[.JsonMapperType]]) {
        this.baseUrl = baseUrl
        this.client = client
        this.json = Json(mapper)
    }

    constructor(baseUrl: String, client: OkHttpClient) : this(baseUrl, client, [[.JsonMapperInit]])

    constructor(baseUrl: String) {
        this.baseUrl = baseUrl
        this.json = Json([[.JsonMapperInit]])
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

func (g *OkHttpGenerator) createUrl(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line(`val url = UrlBuilder(baseUrl)`)
	if operation.InApi.InHttp.GetUrl() != "" {
		w.Line(`url.addPathSegments("%s")`, trimSlash(operation.InApi.InHttp.GetUrl()))
	}
	for _, urlPart := range operation.Endpoint.UrlParts {
		part := trimSlash(urlPart.Part)
		if urlPart.Param != nil {
			w.Line(`url.addPathParameter(%s)`, urlPart.Param.Name.CamelCase())
		} else if len(part) > 0 {
			w.Line(`url.addPathSegments("%s")`, part)
		}
	}
	for _, param := range operation.QueryParams {
		w.Line(`url.addQueryParameter("%s", %s)`, param.Name.SnakeCase(), addBuilderParam(&param))
	}
}

func (g *OkHttpGenerator) requestContentType(operation *spec.NamedOperation) string {
	switch operation.Body.Kind() {
	case spec.BodyEmpty:
		return ""
	case spec.BodyText:
		return `text/plain`
	case spec.BodyJson:
		return `application/json`
	case spec.BodyBinary:
		return `application/octet-stream`
	case spec.BodyFormData:
		return `multipart/form-data`
	case spec.BodyFormUrlEncoded:
		return ""
	default:
		panic(fmt.Sprintf("Unknown Content Type"))
	}
}

func (g *OkHttpGenerator) createRequest(w *writer.Writer, operation *spec.NamedOperation) {
	requestBody := "null"
	if operation.Body.IsText() || operation.Body.IsBinary() {
		w.Line(`val requestBody = body.toRequestBody("%s".toMediaTypeOrNull())`, g.requestContentType(operation))
		requestBody = "requestBody"
	}
	if operation.Body.IsJson() {
		w.Line(`val requestBody = json.%s.toRequestBody("%s".toMediaTypeOrNull())`, g.Models.WriteJson("body", &operation.Body.Type.Definition), g.requestContentType(operation))
		requestBody = "requestBody"
	}
	if operation.Body.IsBodyFormData() {
		w.Line(`val body = MultipartBodyBuilder(MultipartBody.FORM)`)
		for _, param := range operation.Body.FormData {
			if param.Type.Definition.String() == spec.TypeFile {
				w.Line(`body.addFormDataPart("%s", fileName, %s)`, param.Name.Source, addBuilderParam(&param))
			} else {
				w.Line(`body.addFormDataPart("%s", %s)`, param.Name.SnakeCase(), addBuilderParam(&param))
			}
		}
		requestBody = "body.build()"
	}
	if operation.Body.IsBodyFormUrlEncoded() {
		w.Line(`val body = UrlencodedFormBodyBuilder()`)
		for _, param := range operation.Body.FormUrlEncoded {
			w.Line(`body.add("%s", %s)`, param.Name.SnakeCase(), addBuilderParam(&param))
		}
		requestBody = "body.build()"
	}
	w.Line(`val request = RequestBuilder("%s", url.build(), %s)`, operation.Endpoint.Method, requestBody)
	for _, param := range operation.HeaderParams {
		w.Line(`request.addHeaderParameter("%s", %s)`, param.Name.Source, addBuilderParam(&param))
	}
}

func (g *OkHttpGenerator) sendRequest(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line(`logger.info("Sending request, operationId: %s.%s, method: %s, url: %s")`, operation.InApi.Name.Source, operation.Name.Source, operation.Endpoint.Method, operation.FullUrl())
	w.Line(`val response = client.newCall(request.build()).execute()`)
	w.Line(`logger.info("Received response with status code ${response.code}")`)
}

func (g *OkHttpGenerator) processResponse(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line(`when (response.code) {`)
	for _, response := range operation.Responses.Success() {
		w.Line(`    %s -> %s`, spec.HttpStatusCode(response.Name), g.successResponse(response))
	}
	for _, response := range operation.Responses.NonRequiredErrors() {
		w.Line(`    %s -> %s`, spec.HttpStatusCode(response.Name), g.errorResponse(&response.Response))
	}
	w.Line(`    else -> throw ResponseException("Unexpected status code received: ${response.code}")`)
	w.Line(`}`)
}

func (g *OkHttpGenerator) generateClientMethod(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line(`fun %s {`, operationSignature(g.Types, operation))
	w.Line(`   try {`)
	w.IndentWith(2)
	g.createUrl(w, operation)
	w.EmptyLine()
	g.createRequest(w, operation)
	w.EmptyLine()
	g.sendRequest(w, operation)
	w.EmptyLine()
	g.processResponse(w, operation)
	w.UnindentWith(2)
	w.Lines(`
    } catch (ex: Throwable) {
        logger.error(ex.message)
        throw ClientException(ex)
    }
}
`)
}

func (g *OkHttpGenerator) successResponse(response *spec.OperationResponse) string {
	if response.Body.IsText() {
		return responseCreate(response, "response.body!!.string()")
	}
	if response.Body.IsJson() {
		return responseCreate(response, fmt.Sprintf(`json.%s`, g.Models.ReadJson(`response.body!!.charStream()`, &response.Body.Type.Definition)))
	}
	if response.Body.IsBinary() || response.Body.IsFile() {
		return responseCreate(response, "response.body!!.charStream()")
	}
	return responseCreate(response, "")
}

func (g *OkHttpGenerator) errorResponse(response *spec.Response) string {
	var responseBody = ""
	if response.Body.IsText() {
		responseBody = "response.body!!.string()"
	}
	if response.Body.IsJson() {
		responseBody = fmt.Sprintf(`json.%s`, g.Models.ReadJson(`response.body!!.charStream()`, &response.Body.Type.Definition))
	}
	if response.Body.IsBinary() || response.Body.IsFile() {
		responseBody = "response.body!!.charStream()"
	}
	return fmt.Sprintf(`throw %s(%s)`, errorExceptionClassName(response), responseBody)
}

func (g *OkHttpGenerator) Utils() []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *g.generateRequestBuilder())
	files = append(files, *g.generateUrlBuilder())
	files = append(files, *g.generateMultipartBodyBuilder())
	files = append(files, *g.generateUrlencodedFormBodyBuilder())
	return files
}

func (g *OkHttpGenerator) generateRequestBuilder() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `RequestBuilder`)
	w.Lines(`
import okhttp3.*

class [[.ClassName]](method: String, url: HttpUrl, body: RequestBody?) {
	private val requestBuilder: Request.Builder = Request.Builder().url(url).method(method, body)

	fun addHeaderParameter(name: String, value: Any) {
		val valueStr = value.toString()
		this.requestBuilder.addHeader(name, valueStr)
	}

	fun <T> addHeaderParameter(name: String, values: List<T>) {
		for (value in values) {
			this.addHeaderParameter(name, value!!)
		}
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

class [[.ClassName]](baseUrl: String) {
	private val urlBuilder: HttpUrl.Builder = baseUrl.toHttpUrl().newBuilder()

	fun addQueryParameter(name: String, value: Any) {
		val valueStr = value.toString()
		urlBuilder.addQueryParameter(name, valueStr)
	}

	fun <T> addQueryParameter(name: String, values: List<T>) {
		for (value in values) {
			this.addQueryParameter(name, value!!)
		}
	}

	fun addPathSegments(value: String) {
		this.urlBuilder.addPathSegments(value)
	}

	fun addPathParameter(value: Any) {
		val valueStr = value.toString()
		this.urlBuilder.addPathSegment(valueStr)
	}

	fun build(): HttpUrl {
		return this.urlBuilder.build()
	}
}
`)
	return w.ToCodeFile()
}

func (g *OkHttpGenerator) generateMultipartBodyBuilder() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `MultipartBodyBuilder`)
	w.Lines(`
import okhttp3.*
import okhttp3.MediaType.Companion.toMediaTypeOrNull
import okhttp3.RequestBody.Companion.asRequestBody
import okhttp3.RequestBody.Companion.toRequestBody
import java.io.File
import java.net.URLConnection

class [[.ClassName]](type: MediaType) {
	private val multipartBodyBuilder: MultipartBody.Builder = MultipartBody.Builder().setType(type)

	fun addFormDataPart(name: String, value: Any) {
		multipartBodyBuilder.addFormDataPart(name, value.toString())
	}

	fun <T> addFormDataPart(name: String, values: List<T>) {
		for (value in values) {
			multipartBodyBuilder.addFormDataPart(name, value.toString())
		}
	}

	fun addFormDataPart(fieldName: String, file: File) {
		multipartBodyBuilder.addFormDataPart(fieldName, file.name, file.asRequestBody(getFileContentType(file.name)))
	}

	fun addFormDataPart(fieldName: String, fileName: String, file: ByteArray) {
		multipartBodyBuilder.addFormDataPart(fieldName, fileName, file.toRequestBody(getFileContentType(fileName)))
	}

	private fun getFileContentType(fileName: String): MediaType {
		return URLConnection.getFileNameMap().getContentTypeFor(fileName).toMediaTypeOrNull()!!
	}

	fun build(): MultipartBody {
		return multipartBodyBuilder.build()
	}
}
`)
	return w.ToCodeFile()
}

func (g *OkHttpGenerator) generateUrlencodedFormBodyBuilder() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `UrlencodedFormBodyBuilder`)
	w.Lines(`
import okhttp3.FormBody

class [[.ClassName]] {
	private val formBodyBuilder: FormBody.Builder = FormBody.Builder()

	fun add(name: String, value: Any) {
		formBodyBuilder.add(name, value.toString())
	}

	fun <T> add(name: String, values: List<T>) {
		for (value in values) {
			formBodyBuilder.add(name, value.toString())
		}
	}

	fun build(): FormBody {
		return formBodyBuilder.build()
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
		files = append(files, *errorResponseException(g.Packages.Errors, g.Packages.ErrorsModels, &errorResponse.Response))
	}
	files = append(files, *g.errorsInterceptor(errors))
	return files
}

func (g *OkHttpGenerator) errorsInterceptor(errorsResponses *spec.ErrorResponses) *generator.CodeFile {
	w := writer.New(g.Packages.Errors, `ErrorsInterceptor`)
	w.Imports.Add(`okhttp3.*`)
	w.Imports.PackageStar(g.Packages.ErrorsModels)
	w.Imports.PackageStar(g.Packages.Json)
	w.Lines(`
class [[.ClassName]](private var json: Json) : Interceptor {
	override fun intercept(chain: Interceptor.Chain): Response {
`)
	w.IndentWith(2)
	w.Line(`val request: Request = chain.request()`)
	w.Line(`val response: Response = chain.proceed(request)`)
	w.Line(`when (response.code) {`)
	for _, errorResponse := range errorsResponses.Required() {
		w.Line(`    %s -> %s`, spec.HttpStatusCode(errorResponse.Name), g.errorResponse(&errorResponse.Response))
	}
	w.Line(`}`)
	w.Line(`return response`)
	w.UnindentWith(2)
	w.Lines(`
	}
}
`)
	return w.ToCodeFile()
}
