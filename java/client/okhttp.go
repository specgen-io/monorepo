package client

import (
	"fmt"
	"strings"

	"generator"
	"java/imports"
	"java/packages"
	"java/writer"
	"spec"
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
	imports.Add(g.Packages.Root.PackageStar)
	imports.Add(g.Packages.Json.PackageStar)
	imports.Add(g.Packages.Utils.PackageStar)
	imports.Add(g.Packages.Models(api.InHttp.InVersion).PackageStar)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`public class [[.ClassName]] {`)
	w.Line(`  private static final Logger logger = LoggerFactory.getLogger([[.ClassName]].class);`)
	w.EmptyLine()
	w.Line(`  private String baseUrl;`)
	w.Line(`  private OkHttpClient client;`)
	g.CreateJsonMapperField(w.Indented(), "")
	w.EmptyLine()
	w.Line(`  public [[.ClassName]](String baseUrl) {`)
	w.Line(`    this.baseUrl = baseUrl;`)
	g.InitJsonMapper(w.IndentedWith(2))
	w.Line(`    this.client = new OkHttpClient();`)
	w.Line(`  }`)
	for _, operation := range api.Operations {
		w.EmptyLine()
		g.generateClientMethod(w.Indented(), &operation)
	}
	w.Line(`}`)
	return w.ToCodeFile()
}

func (g *Generator) generateClientMethod(w *generator.Writer, operation *spec.NamedOperation) {
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	requestBody := "null"

	w.Line(`public %s {`, operationSignature(g.Types, operation))
	if operation.BodyIs(spec.BodyString) {
		w.Line(`  var requestBody = RequestBody.create(body, MediaType.parse("text/plain"));`)
		requestBody = "requestBody"
	}
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`  String bodyJson;`)
		bodyJson, exception := g.WriteJson("body", &operation.Body.Type.Definition)
		generateClientTryCatch(w.Indented(),
			fmt.Sprintf(`bodyJson = %s;`, bodyJson),
			exception, `e`,
			`"Failed to serialize JSON " + e.getMessage()`)
		w.EmptyLine()
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
	w.Line(`  Response response;`)
	generateClientTryCatch(w.Indented(),
		`response = client.newCall(request.build()).execute();`,
		`IOException`, `e`,
		`"Failed to execute the request " + e.getMessage()`)
	w.EmptyLine()
	w.Line(`  switch (response.code()) {`)
	for _, response := range operation.Responses {
		w.Line(`    case %s: {`, spec.HttpStatusCode(response.Name))
		w.IndentWith(3)
		w.Line(`logger.info("Received response with status code {}", response.code());`)
		if response.BodyIs(spec.BodyEmpty) {
			w.Line(responseCreate(&response, ""))
		}
		if response.BodyIs(spec.BodyString) {
			w.Line(`%s responseBody;`, g.Types.Java(&response.Type.Definition))
			generateClientTryCatch(w,
				fmt.Sprintf(`responseBody = response.body().string();`),
				`IOException`, `e`,
				`"Failed to convert response body to string " + e.getMessage()`)
			w.Line(responseCreate(&response, `responseBody`))
		}
		if response.BodyIs(spec.BodyJson) {
			w.Line(`%s responseBody;`, g.Types.Java(&response.Type.Definition))
			responseBody, exception := g.ReadJson("response.body().string()", &response.Type.Definition)
			generateClientTryCatch(w,
				fmt.Sprintf(`responseBody = %s;`, responseBody),
				exception, `e`,
				`"Failed to deserialize response body " + e.getMessage()`)
			w.Line(responseCreate(&response, `responseBody`))
		}
		w.UnindentWith(3)
		w.Line(`    }`)
	}
	w.Line(`    default:`)
	generateThrowClientException(w.IndentedWith(3), `"Unexpected status code received: " + response.code()`, ``)
	w.Line(`  }`)
	w.Line(`}`)
}

func (g *Generator) responses(version *spec.Version) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		for _, operation := range api.Operations {
			if len(operation.Responses) > 1 {
				files = append(files, *g.responseInterface(g.Types, &operation))
			}
		}
	}
	return files
}

func generateTryCatch(w *generator.Writer, exceptionObject string, codeBlock func(w *generator.Writer), exceptionHandler func(w *generator.Writer)) {
	w.Line(`try {`)
	codeBlock(w.Indented())
	w.Line(`} catch (%s) {`, exceptionObject)
	exceptionHandler(w.Indented())
	w.Line(`}`)
}

func generateClientTryCatch(w *generator.Writer, statement string, exceptionType, exceptionVar, errorMessage string) {
	generateTryCatch(w, exceptionType+` `+exceptionVar,
		func(w *generator.Writer) {
			w.Line(statement)
		},
		func(w *generator.Writer) {
			generateThrowClientException(w, errorMessage, exceptionVar)
		})
}

func generateThrowClientException(w *generator.Writer, errorMessage string, wrapException string) {
	w.Line(`var errorMessage = %s;`, errorMessage)
	w.Line(`logger.error(errorMessage);`)
	params := "errorMessage"
	if wrapException != "" {
		params += ", " + wrapException
	}
	w.Line(`throw new ClientException(%s);`, params)
}

func (g *Generator) Utils() []generator.CodeFile {
	files := []generator.CodeFile{}

	files = append(files, *generateRequestBuilder(g.Packages.Utils))
	files = append(files, *generateUrlBuilder(g.Packages.Utils))
	files = append(files, *generateStringify(g.Packages.Utils))

	return files
}

func generateRequestBuilder(thePackage packages.Package) *generator.CodeFile {
	code := `
package [[.PackageName]];

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
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("RequestBuilder.java"),
		Content: strings.TrimSpace(code),
	}
}

func generateUrlBuilder(thePackage packages.Package) *generator.CodeFile {
	code := `
package [[.PackageName]];

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
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("UrlBuilder.java"),
		Content: strings.TrimSpace(code),
	}
}

func generateStringify(thePackage packages.Package) *generator.CodeFile {
	code := `
package [[.PackageName]];

public class Stringify {
    public static String paramToString(Object value) {
        if (value == null) {
            return null;
        }
        return String.valueOf(value);
    }
}
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("Stringify.java"),
		Content: strings.TrimSpace(code),
	}
}
