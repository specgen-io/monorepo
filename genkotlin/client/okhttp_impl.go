package client

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/genkotlin/imports"
	"github.com/specgen-io/specgen/v2/genkotlin/modules"
	"github.com/specgen-io/specgen/v2/genkotlin/responses"
	"github.com/specgen-io/specgen/v2/genkotlin/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func (g *Generator) Clients(version *spec.Version, thePackage modules.Module, modelsVersionPackage modules.Module, jsonPackage modules.Module, utilsPackage modules.Module, mainPackage modules.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		apiPackage := thePackage.Subpackage(api.Name.SnakeCase())
		files = append(files, g.client(&api, apiPackage, modelsVersionPackage, jsonPackage, utilsPackage, mainPackage)...)
	}
	return files
}

func (g *Generator) client(api *spec.Api, apiPackage modules.Module, modelsVersionPackage modules.Module, jsonPackage modules.Module, utilsPackage modules.Module, mainPackage modules.Module) []sources.CodeFile {
	files := []sources.CodeFile{}

	w := writer.NewKotlinWriter()
	w.Line(`package %s`, apiPackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.Models.JsonImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Add(`okhttp3.*`)
	imports.Add(`okhttp3.MediaType.Companion.toMediaTypeOrNull`)
	imports.Add(`okhttp3.RequestBody.Companion.toRequestBody`)
	imports.Add(`org.slf4j.*`)
	imports.Add(mainPackage.PackageStar)
	imports.Add(utilsPackage.PackageStar)
	imports.Add(fmt.Sprintf(`%s.setupObjectMapper`, jsonPackage.PackageName))
	imports.Write(w)
	w.EmptyLine()
	className := clientName(api)
	w.Line(`class %s(`, className)
	w.Line(`  private val baseUrl: String,`)
	w.Line(`  private val client: OkHttpClient = OkHttpClient(),`)
	g.Models.InitJsonMapper(w.Indented())
	w.Line(`) {`)
	w.Line(`  private val logger: Logger = LoggerFactory.getLogger(%s::class.java)`, className)
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
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.kt", className)),
		Content: w.String(),
	})

	return files
}

func (g *Generator) generateClientMethod(w *sources.Writer, operation *spec.NamedOperation) {
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	requestBody := "null"

	w.Line(`fun %s {`, responses.Signature(g.Types, operation))
	bodyDataVar := "bodyJson"
	mediaType := "application/json"
	if operation.Body != nil {
		if operation.Body.Type.Definition.Plain == spec.TypeString {
			bodyDataVar = "body"
			mediaType = "text/plain"
		} else {
			bodyJson, exception := g.Models.WriteJson("body", &operation.Body.Type.Definition)
			generateClientTryCatch(w.Indented(), bodyDataVar,
				bodyJson,
				exception, `e`,
				`"Failed to serialize JSON " + e.message`)
			w.EmptyLine()
		}
	}
	w.Line(`  val url = UrlBuilder(baseUrl)`)
	if operation.Api.Apis.GetUrl() != "" {
		w.Line(`  url.addPathSegment("%s")`, trimSlash(operation.Api.Apis.GetUrl()))
	}
	for _, urlPart := range operation.Endpoint.UrlParts {
		part := trimSlash(urlPart.Part)
		if urlPart.Param != nil {
			w.Line(`  url.addPathSegment(%s)`, urlPart.Param.Name.CamelCase())
		} else if len(part) > 0 {
			w.Line(`  url.addPathSegment("%s")`, part)
		}
	}
	for _, param := range operation.QueryParams {
		w.Line(`  url.addQueryParameter("%s", %s)`, param.Name.SnakeCase(), addBuilderParam(&param))
	}
	w.EmptyLine()
	if operation.Body != nil {
		w.Line(`  val requestBody = %s.toRequestBody("%s".toMediaTypeOrNull())`, bodyDataVar, mediaType)
		requestBody = "requestBody"
	}
	w.Line(`  val request = RequestBuilder("%s", url.build(), %s)`, methodName, requestBody)
	for _, param := range operation.HeaderParams {
		w.Line(`  request.addHeaderParameter("%s", %s)`, param.Name.Source, addBuilderParam(&param))
	}
	w.EmptyLine()
	w.Line(`  logger.info("Sending request, operationId: %s.%s, method: %s, url: %s")`, operation.Api.Name.Source, operation.Name.Source, methodName, url)
	generateClientTryCatch(w.Indented(), "response",
		`client.newCall(request.build()).execute()`,
		`IOException`, `e`,
		`"Failed to execute the request " + e.message`)
	w.EmptyLine()
	w.Line(`  return when (response.code) {`)
	for _, response := range operation.Responses {
		w.Line(`    %s -> {`, spec.HttpStatusCode(response.Name))
		w.IndentWith(3)
		w.Line(`logger.info("Received response with status code {}", response.code)`)
		if !response.Type.Definition.IsEmpty() {
			if response.Type.Definition.Plain == spec.TypeString {
				generateClientTryCatch(w, `responseBody`,
					`response.body()!!.string()`,
					`IOException`, `e`,
					`"Failed to convert response body to string " + e.message`)
			} else {
				responseBody, exception := g.Models.ReadJson(`response.body()!!.string()`, &response.Type.Definition)
				generateClientTryCatch(w, `responseBody`,
					responseBody,
					exception, `e`,
					`"Failed to deserialize response body " + e.message`)
			}
			if len(operation.Responses) > 1 {
				w.Line(`%s.%s(responseBody)`, serviceResponseInterfaceName(operation), response.Name.PascalCase())
			} else {
				w.Line(`responseBody`)
			}
		} else {
			if len(operation.Responses) > 1 {
				w.Line(`%s.%s()`, serviceResponseInterfaceName(operation), response.Name.PascalCase())
			}
		}
		w.UnindentWith(3)
		w.Line(`    }`)
	}
	w.Line(`    else -> {`)
	generateThrowClientException(w.IndentedWith(3), `"Unexpected status code received: " + response.code`, ``)
	w.Line(`    }`)
	w.Line(`  }`)
	w.Line(`}`)
}

func generateTryCatch(w *sources.Writer, valName string, exceptionObject string, codeBlock func(w *sources.Writer), exceptionHandler func(w *sources.Writer)) {
	w.Line(`val %s = try {`, valName)
	codeBlock(w.Indented())
	w.Line(`} catch (%s) {`, exceptionObject)
	exceptionHandler(w.Indented())
	w.Line(`}`)
}

func generateClientTryCatch(w *sources.Writer, valName string, statement string, exceptionType, exceptionVar, errorMessage string) {
	generateTryCatch(w, valName, exceptionVar+`: `+exceptionType,
		func(w *sources.Writer) {
			w.Line(statement)
		},
		func(w *sources.Writer) {
			generateThrowClientException(w, errorMessage, exceptionVar)
		})
}

func generateThrowClientException(w *sources.Writer, errorMessage string, wrapException string) {
	w.Line(`val errorMessage = %s`, errorMessage)
	w.Line(`logger.error(errorMessage)`)
	params := "errorMessage"
	if wrapException != "" {
		params += ", " + wrapException
	}
	w.Line(`throw ClientException(%s)`, params)
}

func trimSlash(param string) string {
	return strings.Trim(param, "/")
}
