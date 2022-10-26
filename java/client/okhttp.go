package client

import (
	"fmt"
	"generator"
	"java/imports"
	"java/packages"
	"java/writer"
	"spec"
	"strings"
)

func (g *Generator) Clients(version *spec.Version) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *g.client(&api))
	}
	files = append(files, g.responses(version)...)
	return files
}

func (g *Generator) client(api *spec.Api) *generator.CodeFile {
	w := writer.New(g.Packages.Client(api), clientName(api))
	imports := imports.New()
	imports.Add(g.ModelsUsageImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Add(`okhttp3.*`)
	imports.Add(`org.slf4j.*`)
	imports.Add(g.Packages.Errors.PackageStar)
	imports.Add(g.Packages.Json.PackageStar)
	imports.Add(g.Packages.Utils.PackageStar)
	imports.Add(g.Packages.Models(api.InHttp.InVersion).PackageStar)
	imports.AddStatic(g.Packages.Utils.Subpackage(`ClientResponse`).PackageStar)
	imports.AddStatic(g.Packages.Utils.Subpackage(`ErrorsHandler`).PackageStar)
	imports.Write(w)
	w.EmptyLine()
	w.Lines(`
public class [[.ClassName]] {
	private static final Logger logger = LoggerFactory.getLogger([[.ClassName]].class);

	private final String baseUrl;
	private final OkHttpClient client;
	private final Json json;

	public [[.ClassName]](String baseUrl) {
		this.baseUrl = baseUrl;
`)
	w.IndentedWith(2).Lines(g.CreateJsonHelper(`this.json`))
	w.Lines(`
		this.client = new OkHttpClient();
	}
`)
	for _, operation := range api.Operations {
		w.EmptyLine()
		g.generateClientMethod(w.Indented(), &operation)
	}
	w.Line(`}`)
	return w.ToCodeFile()
}

func (g *Generator) generateClientMethod(w generator.Writer, operation *spec.NamedOperation) {
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	requestBody := "null"
	w.Line(`public %s {`, operationSignature(g.Types, operation))
	if operation.BodyIs(spec.BodyString) {
		w.Line(`  var requestBody = RequestBody.create(body, MediaType.parse("text/plain"));`)
		requestBody = "requestBody"
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`  var bodyJson = json.%s;`, g.JsonWrite("body", &operation.Body.Type.Definition))
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
	for _, response := range operation.Responses {
		if isSuccessfulStatusCode(spec.HttpStatusCode(response.Name)) {
			w.Line(`  if (response.code() == %s) {`, spec.HttpStatusCode(response.Name))
			w.IndentWith(2)
			w.Line(`logger.info("Received response with status code {}", response.code());`)
			if response.BodyIs(spec.BodyEmpty) {
				w.Line(responseCreate(&response, ""))
			}
			if response.BodyIs(spec.BodyString) {
				responseBodyString := "getResponseBodyString(response, logger)"
				w.Line(responseCreate(&response, responseBodyString))
			}
			if response.BodyIs(spec.BodyJson) {
				w.Line(`var responseBodyString = getResponseBodyString(response, logger);`)
				responseBody := fmt.Sprintf(`json.%s`, g.JsonRead("responseBodyString", &response.Type.Definition))
				w.Line(responseCreate(&response, responseBody))
			}
			w.UnindentWith(2)
			w.Line(`  }`)
		}
	}
	w.Line(`  handleErrors(response, logger, json);`)
	w.EmptyLine()
	generateThrowClientException(w.Indented(), `"Unexpected status code received: " + response.code()`, ``)
	w.Line(`}`)
}

func (g *Generator) responses(version *spec.Version) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		for _, operation := range api.Operations {
			if responsesNumber(&operation) > 1 {
				files = append(files, *g.responseInterface(g.Types, &operation))
			}
		}
	}
	return files
}

func generateThrowClientException(w generator.Writer, errorMessage string, wrapException string) {
	w.Line(`var errorMessage = %s;`, errorMessage)
	w.Line(`logger.error(errorMessage);`)
	params := "errorMessage"
	if wrapException != "" {
		params += ", " + wrapException
	}
	w.Line(`throw new ClientException(%s);`, params)
}

func (g *Generator) Utils(responses *spec.Responses) []generator.CodeFile {
	files := []generator.CodeFile{}

	files = append(files, *generateRequestBuilder(g.Packages.Utils))
	files = append(files, *generateUrlBuilder(g.Packages.Utils))
	files = append(files, *generateStringify(g.Packages.Utils))
	files = append(files, *generateClientResponse(g.Packages.Utils, g.Packages.Errors))
	files = append(files, *g.generateErrorsHandler(g.Packages.Utils, responses))

	return files
}

func generateRequestBuilder(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `RequestBuilder`)
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

func generateUrlBuilder(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `UrlBuilder`)
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

func generateStringify(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `Stringify`)
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

func generateClientResponse(thePackage packages.Package, errorsPackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `ClientResponse`)
	w.Template(
		map[string]string{
			`ErrorsPackage`: errorsPackage.PackageName,
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

func (g *Generator) generateErrorsHandler(thePackage packages.Package, errorsResponses *spec.Responses) *generator.CodeFile {
	w := writer.New(thePackage, `ErrorsHandler`)
	imports := imports.New()
	imports.Add(g.ModelsUsageImports()...)
	imports.Add(`okhttp3.*`)
	imports.Add(`org.slf4j.*`)
	imports.Add(g.Packages.Errors.PackageStar)
	imports.Add(g.Packages.ErrorsModels.PackageStar)
	imports.Add(g.Packages.Json.PackageStar)
	imports.AddStatic(g.Packages.Utils.Subpackage(`ClientResponse`).PackageStar)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`public class [[.ClassName]] {`)
	w.Line(`  public static void handleErrors(Response response, Logger logger, Json json) {`)
	for _, errorResponse := range *errorsResponses {
		w.Line(`    if (response.code() == %s) {`, spec.HttpStatusCode(errorResponse.Name))
		w.Line(`      var responseBodyString = getResponseBodyString(response, logger);`)
		w.Line(`      var responseBody = json.%s;`, g.JsonRead("responseBodyString", &errorResponse.Type.Definition))
		w.Line(`      throw new %sException(responseBody);`, g.Types.Java(&errorResponse.Type.Definition))
		w.Line(`    }`)
	}
	w.Line(`  }`)
	w.Line(`}`)

	return w.ToCodeFile()
}
