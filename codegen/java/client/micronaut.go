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

var Micronaut = "micronaut"

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
	w.Imports.Add(`io.micronaut.http.*`)
	w.Imports.Add(`io.micronaut.http.client.*`)
	w.Imports.Add(`java.net.*`)
	w.Imports.Add(`org.slf4j.*`)
	w.Imports.Star(g.Packages.Errors)
	w.Imports.Star(g.Packages.ErrorsModels)
	w.Imports.Star(g.Packages.Json)
	w.Imports.Star(g.Packages.Utils)
	w.Imports.Star(g.Packages.Models(api.InHttp.InVersion))
	w.Imports.StaticStar(g.Packages.Utils.Subpackage(`Requestor`))
	w.Template(
		map[string]string{
			`JsonMapperType`: g.Models.JsonMapperType(),
			`JsonMapperInit`: g.Models.JsonMapperInit(),
		}, `
public class [[.ClassName]] {
	private static final Logger logger = LoggerFactory.getLogger([[.ClassName]].class);

	private final BlockingHttpClient client;
	private final Json json;
	private final ErrorsHandler errorsHandler;

	public [[.ClassName]](HttpClient client, [[.JsonMapperType]] mapper) {
		this.client = client.toBlocking();
		this.json = new Json(mapper);
		this.errorsHandler = new ErrorsHandler(json);
	}

	public [[.ClassName]](HttpClient client) {
		this(client, [[.JsonMapperInit]]);
	}

	public [[.ClassName]](String baseUrl) {
		this.json = new Json([[.JsonMapperInit]]);
		this.errorsHandler = new ErrorsHandler(json);
		try {
			this.client = HttpClient.create(new URL(baseUrl)).toBlocking();
		} catch (MalformedURLException e) {
			var errorMessage = "Failed to create URL object from string " + e.getMessage();
			logger.error(errorMessage);
			throw new ClientException(errorMessage, e);
		}
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
	w.Line(`var url = new UrlBuilder("%s");`, getUrl(operation))
	for _, urlPart := range operation.Endpoint.UrlParts {
		if urlPart.Param != nil {
			w.Line(`url.pathParam("%s", %s);`, urlPart.Param.Name.CamelCase(), urlPart.Param.Name.CamelCase())
		}
	}
	for _, param := range operation.QueryParams {
		w.Line(`url.queryParam("%s", %s);`, param.Name.SnakeCase(), param.Name.CamelCase())
	}
}

func (g *MicronautGenerator) createRequest(w *writer.Writer, operation *spec.NamedOperation) {
	var requestBody string
	if operation.Body.IsText() {
		requestBody = "body"
	}
	if operation.Body.IsJson() {
		w.Line(`var bodyJson = json.%s;`, g.Models.JsonWrite("body", &operation.Body.Type.Definition))
		requestBody = "bodyJson"
	}
	if operation.Body.IsBodyFormData() {
		w.EmptyLine()
		w.Line(`var body = new MultipartBodyBuilder();`)
		for _, param := range operation.Body.FormData {
			w.Line(`body.addPart("%s", %s);`, param.Name.SnakeCase(), param.Name.CamelCase())
		}
		requestBody = "body.build()"
	}
	if operation.Body.IsBodyFormUrlEncoded() {
		w.EmptyLine()
		w.Line(`var body = new UrlencodedFormBodyBuilder();`)
		for _, param := range operation.Body.FormUrlEncoded {
			w.Line(`body.addParameter("%s", %s);`, param.Name.SnakeCase(), param.Name.CamelCase())
		}
		requestBody = "body.build()"
	}
	w.EmptyLine()
	w.Line(`var request = new RequestBuilder(%s);`, requestBuilderParams(operation.Endpoint.Method, requestBody, operation))
	for _, param := range operation.HeaderParams {
		w.Line(`request.headerParam("%s", %s);`, param.Name.Source, param.Name.CamelCase())
	}
	if !operation.Body.IsEmpty() {
		w.Line(`request.contentType(%s);`, contentType(operation))
	}
}

func contentType(operation *spec.NamedOperation) string {
	if operation.Body.IsEmpty() {
		return ""
	} else if operation.Body.IsText() {
		return `MediaType.TEXT_PLAIN_TYPE`
	} else if operation.Body.IsJson() {
		return `MediaType.APPLICATION_JSON_TYPE`
	} else if operation.Body.IsBodyFormData() {
		return `MediaType.MULTIPART_FORM_DATA_TYPE`
	} else if operation.Body.IsBodyFormUrlEncoded() {
		return `MediaType.FORM`
	} else {
		panic(fmt.Sprintf("Unknown Contet Type"))
	}
}

func (g *MicronautGenerator) sendRequest(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line(`logger.info("Sending request, operationId: %s.%s, method: %s, url: %s");`, operation.InApi.Name.Source, operation.Name.Source, operation.Endpoint.Method, operation.FullUrl())
	w.Line(`var response = sendRequest(client, request);`)
	w.Line(`logger.info("Received response with status code ${response.code()}");`)
}

func (g *MicronautGenerator) processResponse(w *writer.Writer, operation *spec.NamedOperation) {
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

func (g *MicronautGenerator) generateClientMethod(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line(`public %s {`, operationSignature(g.Types, operation))
	w.Line(`  try {`)
	w.IndentWith(2)
	g.createUrl(w, operation)
	g.createRequest(w, operation)
	w.EmptyLine()
	g.sendRequest(w, operation)
	w.EmptyLine()
	w.Line(`errorsHandler.handle(response);`)
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

func (g *MicronautGenerator) successResponse(response *spec.OperationResponse) string {
	if response.Body.IsText() {
		return responseCreate(response, "response.body().toString()")
	}
	if response.Body.IsJson() {
		return responseCreate(response, fmt.Sprintf(`json.%s`, g.Models.JsonRead(`response.body().toString()`, &response.Body.Type.Definition)))
	}
	return responseCreate(response, "")
}

func (g *MicronautGenerator) errorResponse(response *spec.Response) string {
	var responseBody = ""
	if response.Body.IsText() {
		responseBody = "response.body().toString()"
	}
	if response.Body.IsJson() {
		responseBody = fmt.Sprintf(`json.%s`, g.Models.JsonRead(`response.body().toString()`, &response.Body.Type.Definition))
	}
	return fmt.Sprintf(`throw new %s(%s);`, errorExceptionClassName(response), responseBody)
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
	params := fmt.Sprintf(`%s, %s, HttpRequest::%s`, urlParam, requestBody, methodName)
	if operation.Body.IsEmpty() {
		params = fmt.Sprintf(`%s, HttpRequest::%s`, urlParam, methodName)
	}

	return params
}

func (g *MicronautGenerator) Utils(responses *spec.ErrorResponses) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *g.generateRequestBuilder())
	files = append(files, *g.generateUrlBuilder())
	files = append(files, *g.multipartBodyBuilder())
	files = append(files, *g.urlencodedFormBodyBuilder())
	files = append(files, *g.generateRequestor())
	return files
}

func (g *MicronautGenerator) generateRequestBuilder() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `RequestBuilder`)
	w.Lines(`
import io.micronaut.http.MediaType;
import io.micronaut.http.MutableHttpRequest;
import java.net.URI;
import java.util.List;
import java.util.function.*;

public class [[.ClassName]] {
	private final MutableHttpRequest<?> requestBuilder;

	public <T> [[.ClassName]](URI url, T body, BiFunction<URI, T, MutableHttpRequest<?>> method) {
		this.requestBuilder = method.apply(url, body);
	}

	public [[.ClassName]](URI url, Function<URI, MutableHttpRequest<?>> method) {
		this.requestBuilder = method.apply(url);
	}

	public [[.ClassName]] contentType(MediaType mediaType) {
		this.requestBuilder.contentType(mediaType);
		return this;
	}

	public [[.ClassName]] headerParam(String name, Object value) {
		if (value != null) {
			this.requestBuilder.header(name, String.valueOf(value));
		}
		return this;
	}

	public <T> [[.ClassName]] headerParam(String name, List<T> values) {
		for (T val : values) {
			this.headerParam(name, val);
		}
		return this;
	}

	public MutableHttpRequest<?> build() {
		return this.requestBuilder;
	}
}
`)
	return w.ToCodeFile()
}

func (g *MicronautGenerator) generateUrlBuilder() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `UrlBuilder`)
	w.Lines(`
import io.micronaut.http.uri.UriBuilder;

import java.net.URI;
import java.util.*;

public class [[.ClassName]] {
	private final UriBuilder uriBuilder;
	private final Map<String, Object> urlMap = new HashMap<>();

	public [[.ClassName]](String url) {
		this.uriBuilder = UriBuilder.of(url);
	}

	public [[.ClassName]] queryParam(String name, Object value) {
		if (value != null) {
			this.uriBuilder.queryParam(name, String.valueOf(value));
		}
		return this;
	}

	public <T> [[.ClassName]] queryParam(String name, List<T> values) {
		for (T val : values) {
			this.queryParam(name, val);
		}
		return this;
	}

	public [[.ClassName]] pathParam(String name, Object value) {
		this.uriBuilder.path("{" + name + "}");
		this.urlMap.put(name, value);
		return this;
	}

	public URI expand() {
		return this.uriBuilder.expand(Collections.checkedMap(this.urlMap, String.class, Object.class));
	}

	public URI build() {
		return this.uriBuilder.build();
	}
}
`)
	return w.ToCodeFile()
}

func (g *MicronautGenerator) multipartBodyBuilder() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `MultipartBodyBuilder`)
	w.Lines(`
import io.micronaut.http.client.multipart.MultipartBody;
import java.util.List;

public class [[.ClassName]] {
	private final MultipartBody.Builder multipartBodyBuilder;

	public [[.ClassName]]() {
		this.multipartBodyBuilder = MultipartBody.builder();
	}

	public [[.ClassName]] addPart(String name, Object value) {
		if (value != null) {
			this.multipartBodyBuilder.addPart(name, String.valueOf(value));
		}
		return this;
	}

	public <T> [[.ClassName]] addPart(String name, List<T> values) {
		for (T val : values) {
			this.multipartBodyBuilder.addPart(name, String.valueOf(val));
		}
		return this;
	}

	public MultipartBody build() {
		return this.multipartBodyBuilder.build();
	}
}
`)
	return w.ToCodeFile()
}

func (g *MicronautGenerator) urlencodedFormBodyBuilder() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `UrlencodedFormBodyBuilder`)
	w.Lines(`
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.util.*;
import java.util.stream.Collectors;

public class [[.ClassName]] {
	private final List<ParamValues> paramValues;

	public [[.ClassName]]() {
		this.paramValues = new ArrayList<>();
	}

	public [[.ClassName]] addParameter(String name, Object value) {
		if (value != null) {
			this.paramValues.add(new ParamValues(name, value));
		}
		return this;
	}

	public <T> [[.ClassName]] addParameter(String name, List<T> values) {
		for (T val : values) {
			this.paramValues.add(new ParamValues(name, val));
		}
		return this;
	}

	public String build() {
		return this.paramValues.stream()
			.map(entry -> encode(entry.name) + "=" + encode(entry.value.toString()))
			.collect(Collectors.joining("&"));
	}

	private static String encode(String string) {
		return URLEncoder.encode(string, StandardCharsets.UTF_8);
	}

	private static class ParamValues {
		private final String name;
		private final Object value;

		public ParamValues(String name, Object value) {
			this.name = name;
			this.value = value;
		}
	}
}
`)
	return w.ToCodeFile()
}

func (g *MicronautGenerator) generateRequestor() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `Requestor`)
	w.Lines(`
import io.micronaut.http.HttpResponse;
import io.micronaut.http.client.BlockingHttpClient;
import io.micronaut.http.client.exceptions.HttpClientResponseException;

public class [[.ClassName]] {
	public static HttpResponse<?> sendRequest(BlockingHttpClient client, RequestBuilder request) {
		try {
			return client.exchange(request.build(), String.class);
		} catch (HttpClientResponseException e) {
			return e.getResponse();
		}
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
	w.Imports.Add(g.Models.ModelsUsageImports()...)
	w.Imports.Star(g.Packages.Json)
	w.Imports.Star(g.Packages.ErrorsModels)
	w.Lines(`
import io.micronaut.http.HttpResponse;

public class [[.ClassName]] {
	private final Json json;

	public [[.ClassName]](Json json) {
		this.json = json;
	}

	public void handle(HttpResponse<?> response) {
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
	}
}
`)
	return w.ToCodeFile()
}
