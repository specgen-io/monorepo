package client

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/genjava/packages"
	"github.com/specgen-io/specgen/v2/genjava/responses"
	"github.com/specgen-io/specgen/v2/genjava/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func (g *Generator) Clients(version *spec.Version, thePackage packages.Module, modelsVersionPackage packages.Module, jsonPackage packages.Module, utilsPackage packages.Module, mainPackage packages.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		apiPackage := thePackage.Subpackage(api.Name.SnakeCase())
		files = append(files, g.client(&api, apiPackage, modelsVersionPackage, jsonPackage, utilsPackage, mainPackage)...)
	}
	return files
}

func (g *Generator) client(api *spec.Api, apiPackage packages.Module, modelsVersionPackage packages.Module, jsonPackage packages.Module, utilsPackage packages.Module, mainPackage packages.Module) []sources.CodeFile {
	files := []sources.CodeFile{}

	w := writer.NewJavaWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	w.Line(`import com.fasterxml.jackson.databind.ObjectMapper;`)
	w.Line(`import okhttp3.*;`)
	w.Line(`import org.slf4j.*;`)
	w.Line(`import java.io.*;`)
	w.Line(`import java.math.BigDecimal;`)
	w.Line(`import java.time.*;`)
	w.Line(`import java.util.*;`)
	w.EmptyLine()
	w.Line(`import %s;`, mainPackage.PackageStar)
	w.Line(`import %s;`, utilsPackage.PackageStar)
	w.Line(`import %s;`, modelsVersionPackage.PackageStar)
	w.Line(`import %s.Json;`, jsonPackage.PackageName)
	w.EmptyLine()
	className := clientName(api)
	w.Line(`public class %s {`, className)
	w.Line(`  private static final Logger logger = LoggerFactory.getLogger(%s.class);`, className)
	w.EmptyLine()
	w.Line(`  private final String baseUrl;`)
	w.Line(`  private final ObjectMapper objectMapper;`)
	w.Line(`  private final OkHttpClient client;`)
	w.EmptyLine()
	w.Line(`  public %s(String baseUrl) {`, className)
	w.Line(`    this.baseUrl = baseUrl;`)
	w.Line(`    this.objectMapper = new ObjectMapper();`)
	w.Line(`    Json.setupObjectMapper(objectMapper);`)
	w.Line(`    this.client = new OkHttpClient();`)
	w.Line(`  }`)
	for _, operation := range api.Operations {
		w.EmptyLine()
		g.generateClientMethod(w.Indented(), &operation)
	}
	w.Line(`}`)

	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			files = append(files, responses.Interfaces(g.Types, &operation, apiPackage, modelsVersionPackage)...)
		}
	}

	files = append(files, sources.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", className)),
		Content: w.String(),
	})

	return files
}

func (g *Generator) generateClientMethod(w *sources.Writer, operation *spec.NamedOperation) {
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	requestBody := "null"

	w.Line(`public %s {`, responses.Signature(g.Types, operation))
	if operation.Body != nil {
		bodyDataVar := "bodyJson"
		mediaType := "application/json"
		if operation.Body.Type.Definition.Plain == spec.TypeString {
			bodyDataVar = "body"
			mediaType = "text/plain"
		} else {
			w.Line(`  String bodyJson;`)
			generateClientTryCatch(w.Indented(),
				fmt.Sprintf(`bodyJson = %s;`, g.Models.WriteJson("body")),
				`IOException`, `e`,
				`"Failed to serialize JSON " + e.getMessage()`)
			w.EmptyLine()
		}
		w.Line(`  var requestBody = RequestBody.create(%s, MediaType.parse("%s"));`, bodyDataVar, mediaType)
		requestBody = "requestBody"
	}
	w.Line(`  var url = new UrlBuilder(baseUrl);`)
	if operation.Api.Apis.GetUrl() != "" {
		w.Line(`  url.addPathSegment("%s");`, strings.Trim(operation.Api.Apis.GetUrl(), "/"))
	}
	for _, urlPart := range operation.Endpoint.UrlParts {
		part := strings.Trim(urlPart.Part, "/")
		if urlPart.Param != nil {
			w.Line(`  url.addPathSegment(%s);`, urlPart.Param.Name.CamelCase())
		} else if len(part) > 0 {
			w.Line(`  url.addPathSegment("%s");`, part)
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
	w.Line(`  logger.info("Sending request, operationId: %s.%s, method: %s, url: %s");`, operation.Api.Name.Source, operation.Name.Source, methodName, url)
	w.Line(`  Response response;`)
	generateClientTryCatch(w.Indented(),
		`response = client.newCall(request.build()).execute();`,
		`IOException`, `e`,
		`"Failed to execute the request " + e.getMessage()`)
	w.EmptyLine()
	w.Line(`  switch (response.code()) {`)
	for _, response := range operation.Responses {
		w.Line(`    case %s:`, spec.HttpStatusCode(response.Name))
		w.IndentWith(3)
		w.Line(`logger.info("Received response with status code {}", response.code());`)
		if !response.Type.Definition.IsEmpty() {
			responseJavaType := g.Types.Java(&response.Type.Definition)
			w.Line(`%s responseBody;`, responseJavaType)
			valueTypeName := fmt.Sprintf("%s.class", responseJavaType)
			responseBody := g.Models.ReadJson("response.body().string()", valueTypeName)
			if response.Type.Definition.Plain == spec.TypeString {
				responseBody = `response.body().string()`
			}
			generateClientTryCatch(w,
				fmt.Sprintf(`responseBody = %s;`, responseBody),
				`IOException`, `e`,
				`"Failed to deserialize response body " + e.getMessage()`)
			if len(operation.Responses) > 1 {
				w.Line(`return new %s.%s(responseBody);`, responses.InterfaceName(operation), response.Name.PascalCase())
			} else {
				w.Line(`return responseBody;`)
			}
		} else {
			if len(operation.Responses) > 1 {
				w.Line(`return new %s.%s();`, responses.InterfaceName(operation), response.Name.PascalCase())
			} else {
				w.Line(`return;`)
			}
		}
		w.UnindentWith(3)
	}
	w.Line(`    default:`)
	generateThrowClientException(w.IndentedWith(3), `"Unexpected status code received: " + response.code()`, ``)
	w.Line(`  }`)
	w.Line(`}`)
}

func generateTryCatch(w *sources.Writer, exceptionObject string, codeBlock func(w *sources.Writer), exceptionHandler func(w *sources.Writer)) {
	w.Line(`try {`)
	codeBlock(w.Indented())
	w.Line(`} catch (%s) {`, exceptionObject)
	exceptionHandler(w.Indented())
	w.Line(`}`)
}

func generateClientTryCatch(w *sources.Writer, statement string, exceptionType, exceptionVar, errorMessage string) {
	generateTryCatch(w, exceptionType+` `+exceptionVar,
		func(w *sources.Writer) {
			w.Line(statement)
		},
		func(w *sources.Writer) {
			generateThrowClientException(w, errorMessage, exceptionVar)
		})
}

func generateThrowClientException(w *sources.Writer, errorMessage string, wrapException string) {
	w.Line(`var errorMessage = %s;`, errorMessage)
	w.Line(`logger.error(errorMessage);`)
	params := "errorMessage"
	if wrapException != "" {
		params += ", " + wrapException
	}
	w.Line(`throw new ClientException(%s);`, params)
}
