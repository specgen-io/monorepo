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

func (g *Generator) Clients(version *spec.Version, thePackage packages.Package, modelsVersionPackage packages.Package, errorModelsPackage packages.Package, jsonPackage packages.Package, utilsPackage packages.Package, mainPackage packages.Package) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		apiPackage := thePackage.Subpackage(api.Name.SnakeCase())
		files = append(files, g.client(&api, apiPackage, modelsVersionPackage, errorModelsPackage, jsonPackage, utilsPackage, mainPackage)...)
	}
	return files
}

func (g *Generator) client(api *spec.Api, apiPackage packages.Package, modelsVersionPackage packages.Package, errorModelsPackage packages.Package, jsonPackage packages.Package, utilsPackage packages.Package, mainPackage packages.Package) []generator.CodeFile {
	files := []generator.CodeFile{}

	w := writer.NewJavaWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.Models.ModelsDefinitionsImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Add(`okhttp3.*`)
	imports.Add(`org.slf4j.*`)
	imports.Add(mainPackage.PackageStar)
	imports.Add(jsonPackage.PackageStar)
	imports.Add(utilsPackage.PackageStar)
	imports.Add(modelsVersionPackage.PackageStar)
	imports.Write(w)
	w.EmptyLine()
	className := clientName(api)
	w.Line(`public class %s {`, className)
	w.Line(`  private static final Logger logger = LoggerFactory.getLogger(%s.class);`, className)
	w.EmptyLine()
	w.Line(`  private String baseUrl;`)
	w.Line(`  private OkHttpClient client;`)
	g.Models.CreateJsonMapperField(w.Indented(), "")
	w.EmptyLine()
	w.Line(`  public %s(String baseUrl) {`, className)
	w.Line(`    this.baseUrl = baseUrl;`)
	g.Models.InitJsonMapper(w.IndentedWith(2))
	w.Line(`    this.client = new OkHttpClient();`)
	w.Line(`  }`)
	for _, operation := range api.Operations {
		w.EmptyLine()
		g.generateClientMethod(w.Indented(), &operation)
	}
	w.Line(`}`)

	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			files = append(files, reponseInterface(g.Types, &operation, apiPackage, modelsVersionPackage, errorModelsPackage)...)
		}
	}

	files = append(files, generator.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", className)),
		Content: w.String(),
	})

	return files
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
		bodyJson, exception := g.Models.WriteJson("body", &operation.Body.Type.Definition)
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
			responseBody, exception := g.Models.ReadJson("response.body().string()", &response.Type.Definition)
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
