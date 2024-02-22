package client

import (
	"fmt"
	"generator"
	"java/models"
	"java/types"
	"java/writer"
	"spec"
	"strings"
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
	w.Imports.Add(`org.slf4j.*`)
	w.Imports.Star(g.Packages.Errors)
	w.Imports.Star(g.Packages.ErrorsModels)
	w.Imports.Star(g.Packages.Json)
	w.Imports.Star(g.Packages.Utils)
	w.Imports.Star(g.Packages.Models(api.InHttp.InVersion))
	w.Template(
		map[string]string{
			`JsonMapperType`: g.Models.JsonMapperType(),
			`JsonMapperInit`: g.Models.JsonMapperInit(),
		}, `
public class [[.ClassName]] {
	private static final Logger logger = LoggerFactory.getLogger([[.ClassName]].class);

	private final String baseUrl;
	private final OkHttpClient client;
	private final Json json;

	public [[.ClassName]](String baseUrl, OkHttpClient client, [[.JsonMapperType]] mapper) {
		this.baseUrl = baseUrl;
		this.client = client;
		this.json = new Json(mapper);
	}

	public [[.ClassName]](String baseUrl, OkHttpClient client) {
		this(baseUrl, client, [[.JsonMapperInit]]);
	}

	public [[.ClassName]](String baseUrl) {
		this.baseUrl = baseUrl;
		this.json = new Json([[.JsonMapperInit]]);
		this.client = new OkHttpClient().newBuilder().addInterceptor(new ErrorsInterceptor(json)).build();
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
	w.Line(`var url = new UrlBuilder(baseUrl);`)
	if operation.InApi.InHttp.GetUrl() != "" {
		w.Line(`url.addPathSegments("%s");`, strings.Trim(operation.InApi.InHttp.GetUrl(), "/"))
	}
	for _, urlPart := range operation.Endpoint.UrlParts {
		part := strings.Trim(urlPart.Part, "/")
		if urlPart.Param != nil {
			w.Line(`url.addPathParameter(%s);`, urlPart.Param.Name.CamelCase())
		} else if len(part) > 0 {
			w.Line(`url.addPathSegments("%s");`, part)
		}
	}
	for _, param := range operation.QueryParams {
		w.Line(`url.addQueryParameter("%s", %s);`, param.Name.SnakeCase(), param.Name.CamelCase())
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
	w.EmptyLine()
	if operation.Body.IsText() || operation.Body.IsBinary() {
		w.Line(`var requestBody = RequestBody.create(body, MediaType.parse("%s"));`, g.requestContentType(operation))
		requestBody = "requestBody"
		w.EmptyLine()
	}
	if operation.Body.IsJson() {
		w.Line(`var bodyJson = json.%s;`, g.Models.JsonWrite("body", &operation.Body.Type.Definition))
		w.Line(`var requestBody = RequestBody.create(bodyJson, MediaType.parse("%s"));`, g.requestContentType(operation))
		requestBody = "requestBody"
		w.EmptyLine()
	}
	if operation.Body.IsBodyFormData() {
		w.Line(`var body = new MultipartBodyBuilder(MultipartBody.FORM);`)
		for _, param := range operation.Body.FormData {
			w.Line(`body.addFormDataPart("%s", %s);`, param.Name.SnakeCase(), param.Name.CamelCase())
		}
		requestBody = "body.build()"
		w.EmptyLine()
	}
	if operation.Body.IsBodyFormUrlEncoded() {
		w.Line(`var body = new UrlencodedFormBodyBuilder();`)
		for _, param := range operation.Body.FormUrlEncoded {
			w.Line(`body.add("%s", %s);`, param.Name.SnakeCase(), param.Name.CamelCase())
		}
		requestBody = "body.build()"
		w.EmptyLine()
	}

	w.Line(`var request = new RequestBuilder("%s", url.build(), %s);`, operation.Endpoint.Method, requestBody)
	for _, param := range operation.HeaderParams {
		w.Line(`request.addHeaderParameter("%s", %s);`, param.Name.Source, param.Name.CamelCase())
	}
}

func (g *OkHttpGenerator) sendRequest(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line(`logger.info("Sending request, operationId: %s.%s, method: %s, url: %s");`, operation.InApi.Name.Source, operation.Name.Source, operation.Endpoint.Method, operation.FullUrl())
	w.Line(`var response = client.newCall(request.build()).execute();`)
	w.Line(`logger.info("Received response with status code {}", response.code());`)
}

func (g *OkHttpGenerator) processResponse(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line(`switch (response.code()) {`)
	for _, response := range operation.Responses.Success() {
		w.Line(`  case %s:`, spec.HttpStatusCode(response.Name))
		w.Line(`    %s`, g.successResponse(response))
	}
	for _, response := range operation.Responses.NonRequiredErrors() {
		w.Line(`  case %s:`, spec.HttpStatusCode(response.Name))
		w.Line(`    %s`, g.errorResponse(&response.Response))
	}
	w.Line(`  default:`)
	w.Line(`    throw new ResponseException(String.format("Unexpected status code received: {}", response.code()));`)
	w.Line(`}`)
}

func (g *OkHttpGenerator) generateClientMethod(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line(`public %s {`, operationSignature(g.Types, operation))
	w.Line(`  try {`)
	w.IndentWith(2)
	g.createUrl(w, operation)
	g.createRequest(w, operation)
	w.EmptyLine()
	g.sendRequest(w, operation)
	w.EmptyLine()
	g.processResponse(w, operation)
	w.UnindentWith(2)
	w.Lines(`
	} catch (Throwable ex) {
		logger.error(ex.getMessage());
		throw new ClientException((ex));
	}
}
`)
}

func (g *OkHttpGenerator) successResponse(response *spec.OperationResponse) string {
	if response.Body.IsText() {
		return responseCreate(response, "response.body().string()")
	}
	if response.Body.IsJson() {
		return responseCreate(response, fmt.Sprintf(`json.%s`, g.Models.JsonRead(`response.body().charStream()`, &response.Body.Type.Definition)))
	}
	if response.Body.IsBinary() {
		return responseCreate(response, "response.body().charStream()")
	}
	return responseCreate(response, "")
}

func (g *OkHttpGenerator) errorResponse(response *spec.Response) string {
	var responseBody = ""
	if response.Body.IsText() {
		responseBody = "response.body().string()"
	}
	if response.Body.IsJson() {
		responseBody = fmt.Sprintf(`json.%s`, g.Models.JsonRead(`response.body().charStream()`, &response.Body.Type.Definition))
	}
	if response.Body.IsBinary() {
		responseBody = "response.body().charStream()"
	}
	return fmt.Sprintf(`throw new %s(%s);`, errorExceptionClassName(response), responseBody)
}

func (g *OkHttpGenerator) Utils(responses *spec.ErrorResponses) []generator.CodeFile {
	files := []generator.CodeFile{}

	files = append(files, *g.generateRequestBuilder())
	files = append(files, *g.generateUrlBuilder())
	files = append(files, *g.multipartBodyBuilder())
	files = append(files, *g.urlencodedFormBodyBuilder())

	return files
}

func (g *OkHttpGenerator) generateRequestBuilder() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `RequestBuilder`)
	w.Lines(`
import java.util.List;
import okhttp3.*;

public class [[.ClassName]] {
	private final Request.Builder requestBuilder;

	public [[.ClassName]](String method, HttpUrl url, RequestBody body) {
		this.requestBuilder = new Request.Builder().url(url).method(method, body);
	}

	public void addHeaderParameter(String name, Object value) {
		if (value != null) {
			this.requestBuilder.addHeader(name, String.valueOf(value));
		}
	}

	public <T> void addHeaderParameter(String name, List<T> values) {
		for (T val : values) {
			this.addHeaderParameter(name, val);
		}
	}

	public Request build() {
		return this.requestBuilder.build();
	}
}
`)
	return w.ToCodeFile()
}

func (g *OkHttpGenerator) generateUrlBuilder() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `UrlBuilder`)
	w.Lines(`
import java.util.List;
import okhttp3.HttpUrl;

public class [[.ClassName]] {
	private final HttpUrl.Builder urlBuilder;

	public [[.ClassName]](String baseUrl) {
		this.urlBuilder = HttpUrl.get(baseUrl).newBuilder();
	}

	public void addQueryParameter(String name, Object value) {
		if (value != null) {
			this.urlBuilder.addQueryParameter(name, String.valueOf(value));
		}
	}

	public <T> void addQueryParameter(String name, List<T> values) {
		for (T val : values) {
			this.addQueryParameter(name, val);
		}
	}

	public void addPathSegments(String value) {
		this.urlBuilder.addPathSegments(value);
	}

	public void addPathParameter(Object value) {
		this.urlBuilder.addPathSegment(String.valueOf(value));
	}

	public HttpUrl build() {
		return this.urlBuilder.build();
	}
}
`)
	return w.ToCodeFile()
}

func (g *OkHttpGenerator) multipartBodyBuilder() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `MultipartBodyBuilder`)
	w.Lines(`
import java.io.File;
import java.util.List;
import okhttp3.*;

public class [[.ClassName]] {
	private final MediaType contentType;
	private final MultipartBody.Builder multipartBodyBuilder;

	public MultipartBodyBuilder(MediaType contentType) {
		this.contentType = contentType;
		this.multipartBodyBuilder = new MultipartBody.Builder().setType(contentType);
	}

	public void addFormDataPart(String name, Object value) {
		if (value != null) {
			this.multipartBodyBuilder.addFormDataPart(name, String.valueOf(value));
		}
	}

	public <T> void addFormDataPart(String name, List<T> values) {
		for (T val : values) {
			this.multipartBodyBuilder.addFormDataPart(name, String.valueOf(val));
		}
	}

	public void addFormDataPart(String fieldName, File file) {
		this.multipartBodyBuilder.addFormDataPart(fieldName, file.getName(), RequestBody.create(file, this.contentType));
	}

	public void addFormDataPart(String fieldName, String fileName, byte[] file) {
		this.multipartBodyBuilder.addFormDataPart(fieldName, fileName, RequestBody.create(file, this.contentType));
	}

	public MultipartBody build() {
		return this.multipartBodyBuilder.build();
	}
}
`)
	return w.ToCodeFile()
}

func (g *OkHttpGenerator) urlencodedFormBodyBuilder() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `UrlencodedFormBodyBuilder`)
	w.Lines(`
import java.util.List;
import okhttp3.FormBody;

public class [[.ClassName]] {
	private final FormBody.Builder formBodyBuilder;

	public [[.ClassName]]() {
		this.formBodyBuilder = new FormBody.Builder();
	}

	public void add(String name, Object value) {
		if (value != null) {
			this.formBodyBuilder.add(name, String.valueOf(value));
		}
	}

	public <T> void add(String name, List<T> values) {
		for (T val : values) {
			this.formBodyBuilder.add(name, String.valueOf(val));
		}
	}

	public FormBody build() {
		return this.formBodyBuilder.build();
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
	w.Imports.Add(g.Models.ModelsUsageImports()...)
	w.Imports.Add(`java.io.IOException`)
	w.Imports.Add(`okhttp3.*`)
	w.Imports.Add(`org.jetbrains.annotations.NotNull`)
	w.Imports.Star(g.Packages.ErrorsModels)
	w.Imports.Star(g.Packages.Json)
	w.Lines(`
public class [[.ClassName]] implements Interceptor {
	private final Json json;

	public [[.ClassName]](Json json) {
		this.json = json;
	}
	
	@NotNull
	@Override
	public Response intercept(@NotNull Chain chain) throws IOException {
		var request = chain.request();
		var response = chain.proceed(request);

		switch (response.code()) {
`)
	w.IndentWith(2)
	for _, errorResponse := range errorsResponses.Required() {
		w.Line(`  case %s:`, spec.HttpStatusCode(errorResponse.Name))
		w.Line(`    %s`, g.errorResponse(&errorResponse.Response))
	}
	w.UnindentWith(2)
	w.Lines(`
		}
		return response;
	}
}
`)
	return w.ToCodeFile()
}
