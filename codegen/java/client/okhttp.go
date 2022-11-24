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
	w.Imports.Star(g.Packages.Json)
	w.Imports.Star(g.Packages.Utils)
	w.Imports.Star(g.Packages.Models(api.InHttp.InVersion))
	w.Imports.StaticStar(g.Packages.Utils.Subpackage(`ClientResponse`))
	w.Template(
		map[string]string{
			`JsonMapper`:     g.Models.JsonMapper()[0],
			`JsonMapperVar`:  g.Models.JsonMapper()[1],
			`JsonHelperInit`: g.Models.CreateJsonHelper(),
		}, `
public class [[.ClassName]] {
	private static final Logger logger = LoggerFactory.getLogger([[.ClassName]].class);

	private final String baseUrl;
	private final OkHttpClient client;
	private final Json json;

	public [[.ClassName]](String baseUrl, OkHttpClient client, [[.JsonMapper]] [[.JsonMapperVar]]) {
		this.baseUrl = baseUrl;
		this.client = client;
		this.json = new Json([[.JsonMapperVar]]);
	}

	public [[.ClassName]](String baseUrl, OkHttpClient client) {
		this(baseUrl, client, [[.JsonHelperInit]]);
	}

	public [[.ClassName]](String baseUrl) {
		this.baseUrl = baseUrl;
		this.json = new Json([[.JsonHelperInit]]);
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

func (g *OkHttpGenerator) generateClientMethod(w *writer.Writer, operation *spec.NamedOperation) {
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	requestBody := "null"
	w.Line(`public %s {`, operationSignature(g.Types, operation))
	if operation.BodyIs(spec.BodyString) {
		w.Line(`  var requestBody = RequestBody.create(body, MediaType.parse("text/plain"));`)
		requestBody = "requestBody"
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`  var bodyJson = json.%s;`, g.Models.JsonWrite("body", &operation.Body.Type.Definition))
		w.Line(`  var requestBody = RequestBody.create(bodyJson, MediaType.parse("application/json"));`)
		requestBody = "requestBody"
	}
	w.Line(`  var url = new UrlBuilder(baseUrl);`)
	if operation.InApi.InHttp.GetUrl() != "" {
		w.Line(`  url.addPathSegments("%s");`, strings.Trim(operation.InApi.InHttp.GetUrl(), "/"))
	}
	for _, urlPart := range operation.Endpoint.UrlParts {
		part := strings.Trim(urlPart.Part, "/")
		if urlPart.Param != nil {
			w.Line(`  url.addPathParameter(%s);`, urlPart.Param.Name.CamelCase())
		} else if len(part) > 0 {
			w.Line(`  url.addPathSegments("%s");`, part)
		}
	}
	for _, param := range operation.QueryParams {
		w.Line(`  url.addQueryParameter("%s", %s);`, param.Name.SnakeCase(), param.Name.CamelCase())
	}
	w.EmptyLine()
	w.Line(`  var request = new RequestBuilder("%s", url.build(), %s);`, methodName, requestBody)
	for _, param := range operation.HeaderParams {
		w.Line(`  request.addHeaderParameter("%s", %s);`, param.Name.Source, param.Name.CamelCase())
	}
	w.EmptyLine()
	w.Line(`  logger.info("Sending request, operationId: %s.%s, method: %s, url: %s");`, operation.InApi.Name.Source, operation.Name.Source, methodName, url)
	w.Line(`  var response = doRequest(client, request, logger);`)
	w.EmptyLine()
	for _, response := range operation.SuccessResponses() {
		w.Line(`  if (response.code() == %s) {`, spec.HttpStatusCode(response.Name))
		w.IndentWith(2)
		w.Line(`logger.info("Received response with status code {}", response.code());`)
		if response.BodyIs(spec.BodyEmpty) {
			w.Line(responseCreate(response, ""))
		}
		if response.BodyIs(spec.BodyString) {
			responseBodyString := "getResponseBodyString(response, logger)"
			w.Line(responseCreate(response, responseBodyString))
		}
		if response.BodyIs(spec.BodyJson) {
			w.Line(`var responseBodyString = getResponseBodyString(response, logger);`)
			responseBody := fmt.Sprintf(`json.%s`, g.Models.JsonRead("responseBodyString", &response.Type.Definition))
			w.Line(responseCreate(response, responseBody))
		}
		w.UnindentWith(2)
		w.Line(`  }`)
	}
	w.EmptyLine()
	generateThrowClientException(w.Indented(), `"Unexpected status code received: " + response.code()`, ``)
	w.Line(`}`)
}

func (g *OkHttpGenerator) Utils(responses *spec.Responses) []generator.CodeFile {
	files := []generator.CodeFile{}

	files = append(files, *g.generateRequestBuilder())
	files = append(files, *g.generateUrlBuilder())
	files = append(files, *g.generateStringify())
	files = append(files, *g.generateClientResponse())

	return files
}

func (g *OkHttpGenerator) generateRequestBuilder() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `RequestBuilder`)
	w.Lines(`
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
`)
	return w.ToCodeFile()
}

func (g *OkHttpGenerator) generateUrlBuilder() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `UrlBuilder`)
	w.Lines(`
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
`)
	return w.ToCodeFile()
}

func (g *OkHttpGenerator) generateStringify() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `Stringify`)
	w.Lines(`
public class Stringify {
    public static String paramToString(Object value) {
        if (value == null) {
            return null;
        }
        return String.valueOf(value);
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
import okhttp3.*;
import org.slf4j.*;
import [[.ErrorsPackage]].*;

import java.io.IOException;

public class ClientResponse {
	public static Response doRequest(OkHttpClient client, RequestBuilder request, Logger logger) {
		Response response;
		try {
			response = client.newCall(request.build()).execute();
		} catch (IOException e) {
			var errorMessage = "Failed to execute the request " + e.getMessage();
			logger.error(errorMessage);
			throw new ClientException(errorMessage, e);
		}
		return response;
	}

	public static String getResponseBodyString(Response response, Logger logger) {
		String responseBodyString;
		try {
			responseBodyString = response.body().string();
		} catch (IOException e) {
			var errorMessage = "Failed to convert response body to string " + e.getMessage();
			logger.error(errorMessage);
			throw new ClientException(errorMessage, e);
		}
		return responseBodyString;
	}
}
`)
	return w.ToCodeFile()
}

func (g *OkHttpGenerator) Exceptions(errors *spec.Responses) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *clientException(g.Packages.Errors))
	for _, errorResponse := range *errors {
		files = append(files, *inheritedClientException(g.Packages.Errors, g.Packages.ErrorsModels, g.Types, &errorResponse))
	}
	files = append(files, *g.errorsInterceptor(errors))
	return files
}

func (g *OkHttpGenerator) errorsInterceptor(errorsResponses *spec.Responses) *generator.CodeFile {
	w := writer.New(g.Packages.Errors, `ErrorsInterceptor`)
	w.Imports.Add(g.Models.ModelsUsageImports()...)
	w.Imports.Add(`java.io.IOException`)
	w.Imports.Add(`okhttp3.*`)
	w.Imports.Add(`org.jetbrains.annotations.NotNull`)
	w.Imports.Add(`org.slf4j.*`)
	w.Imports.Star(g.Packages.ErrorsModels)
	w.Imports.Star(g.Packages.Json)
	w.Imports.StaticStar(g.Packages.Utils.Subpackage(`ClientResponse`))
	w.Lines(`
public class [[.ClassName]] implements Interceptor {
	private static final Logger logger = LoggerFactory.getLogger([[.ClassName]].class);

	private final Json json;

	public [[.ClassName]](Json json) {
		this.json = json;
	}
	
	@NotNull
	@Override
	public Response intercept(@NotNull Chain chain) throws IOException {
		var request = chain.request();
		var response = chain.proceed(request);
`)
	for _, errorResponse := range *errorsResponses {
		w.Line(`    if (response.code() == %s) {`, spec.HttpStatusCode(errorResponse.Name))
		w.Line(`      var responseBodyString = getResponseBodyString(response, logger);`)
		w.Line(`      var responseBody = json.%s;`, g.Models.JsonRead("responseBodyString", &errorResponse.Type.Definition))
		w.Line(`      throw new %sException(responseBody);`, g.Types.Java(&errorResponse.Type.Definition))
		w.Line(`    }`)
	}
	w.Line(`		return response;`)
	w.Line(`  }`)
	w.Line(`}`)
	return w.ToCodeFile()
}
