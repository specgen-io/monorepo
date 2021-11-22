package genjava

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func generateClientsImplementations(version *spec.Version, thePackage Module, modelsVersionPackage Module, utilsPackage Module) []gen.TextFile {
	files := []gen.TextFile{}
	for _, api := range version.Http.Apis {
		apiPackage := thePackage.Subpackage(api.Name.SnakeCase())
		files = append(files, generateClient(&api, apiPackage, modelsVersionPackage, utilsPackage)...)
	}
	return files
}

func generateClient(api *spec.Api, apiPackage Module, modelsVersionPackage Module, utilsPackage Module) []gen.TextFile {
	files := []gen.TextFile{}

	w := NewJavaWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	w.Line(`import com.fasterxml.jackson.core.*;`)
	w.Line(`import com.fasterxml.jackson.databind.*;`)
	w.Line(`import okhttp3.*;`)
	w.Line(`import org.slf4j.*;`)
	w.Line(`import java.io.*;`)
	w.Line(`import java.math.BigDecimal;`)
	w.Line(`import java.time.*;`)
	w.Line(`import java.util.*;`)
	w.EmptyLine()
	w.Line(`import %s;`, utilsPackage.PackageStar)
	w.Line(`import %s;`, modelsVersionPackage.PackageStar)
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
		generateClientMethod(w.Indented(), operation)
	}
	w.Line(`}`)

	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			files = append(files, generateResponseInterface(operation, apiPackage, modelsVersionPackage)...)
		}
	}

	files = append(files, gen.TextFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", className)),
		Content: w.String(),
	})

	return files
}

func generateClientMethod(w *gen.Writer, operation spec.NamedOperation) {
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	requestBody := "null"

	w.Line(`public %s {`, generateResponsesSignatures(operation))
	if operation.Body != nil {
		w.Line(`  String bodyJson;`)
		w.Line(`  try {`)
		w.Line(`    bodyJson = objectMapper.writeValueAsString(body);`)
		w.Line(`  } catch (JsonProcessingException e) {`)
		generateErrorMessage(w, `"Failed to serialize JSON "`, ` + e.getMessage(), e`)
		w.Line(`  }`)
		w.EmptyLine()
		w.Line(`  var requestBody = RequestBody.create(bodyJson, MediaType.parse("application/json"));`)
		requestBody = "requestBody"
	}
	w.Line(`  var url = new UrlBuilder(baseUrl);`)
	w.Line(`  url.addPathSegment("%s");`, addRequestUrlParams(operation))
	generateUrlBuilding(w, operation)
	w.EmptyLine()
	w.Line(`  var request = new RequestBuilder("%s", url.build(), %s);`, methodName, requestBody)
	generateRequestBuilding(w, operation)
	w.EmptyLine()
	w.Line(`  logger.info("Sending request, operationId: %s.%s, method: %s, url: %s");`, operation.Api.Name.Source, operation.Name.Source, methodName, url)
	w.Line(`  Response response;`)
	w.Line(`  try {`)
	w.Line(`    response = client.newCall(request.build()).execute();`)
	w.Line(`  } catch (IOException e) {`)
	generateErrorMessage(w, `"Failed to execute the request "`, ` + e.getMessage(), e`)
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  switch (response.code()) {`)
	w.Indent()
	for _, response := range operation.Responses {
		w.Line(`  case %s:`, spec.HttpStatusCode(response.Name))
		if !response.Type.Definition.IsEmpty() {
			w.Indent()
			w.Line(`  try {`)
			w.Line(`    logger.info("Received response with status code {}", response.code());`)
			if len(operation.Responses) > 1 {
				w.Line(`    return new %s(objectMapper.readValue(response.body().string(), %s.class));`, serviceResponseImplName(operation, response), JavaType(&response.Type.Definition))
			} else {
				w.Line(`    return objectMapper.readValue(response.body().string(), %s.class);`, JavaType(&response.Type.Definition))
			}
			w.Line(`  } catch (IOException e) {`)
			generateErrorMessage(w, `"Failed to deserialize response body "`, ` + e.getMessage(), e`)
			w.Line(`  }`)
		} else {
			w.Line(`  logger.info("Received response with status code {}", response.code());`)
			if len(operation.Responses) > 1 {
				w.Line(`  return new %s();`, serviceResponseImplName(operation, response))
			} else {
				w.Line(`  return;`)
			}
		}
	}
	w.Unindent()
	w.Line(`  default:`)
	generateErrorMessage(w, `"Unexpected status code received: " + response.code()`)
	w.Line(`  }`)
	w.Line(`}`)
}

func generateErrorMessage(w *gen.Writer, messages ...string) {
	w.Indent()
	w.Line(`  var errorMessage = %s;`, messages[0])
	w.Line(`  logger.error(errorMessage);`)
	w.Line(`  throw new ClientException(%s);`, JoinParams(messages))
	w.Unindent()
}

func addRequestUrlParams(operation spec.NamedOperation) string {
	if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
		return JoinParams(getUrl(operation))
	} else {
		return strings.TrimPrefix(operation.Endpoint.Url, "/")
	}
}

func getUrl(operation spec.NamedOperation) []string {
	reminder := strings.TrimPrefix(operation.Endpoint.Url, "/")
	urlParams := []string{}
	if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
		for _, param := range operation.Endpoint.UrlParams {
			parts := strings.Split(reminder, spec.UrlParamStr(param.Name.Source))
			urlParams = append(urlParams, strings.TrimSuffix(parts[0], "/"))
			reminder = parts[1]
		}
	}
	return urlParams
}

func generateUrlBuilding(w *gen.Writer, operation spec.NamedOperation) {
	w.Indent()
	for _, param := range operation.QueryParams {
		w.Line(`url.addQueryParameter("%s", %s);`, param.Name.SnakeCase(), param.Name.CamelCase())
	}
	for _, param := range operation.Endpoint.UrlParams {
		w.Line(`url.addPathSegment(%s);`, param.Name.CamelCase())
	}
	w.Unindent()
}

func generateRequestBuilding(w *gen.Writer, operation spec.NamedOperation) {
	w.Indent()
	for _, param := range operation.HeaderParams {
		w.Line(`request.addHeaderParameter("%s", %s);`, param.Name.SnakeCase(), param.Name.CamelCase())
	}
	w.Unindent()
}
