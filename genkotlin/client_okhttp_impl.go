package genkotlin

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func generateClientsImplementations(version *spec.Version, thePackage Module, modelsVersionPackage Module, jsonPackage Module, utilsPackage Module, mainPackage Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		apiPackage := thePackage.Subpackage(api.Name.SnakeCase())
		files = append(files, generateClient(&api, apiPackage, modelsVersionPackage, jsonPackage, utilsPackage, mainPackage)...)
	}
	return files
}

func generateClient(api *spec.Api, apiPackage Module, modelsVersionPackage Module, jsonPackage Module, utilsPackage Module, mainPackage Module) []sources.CodeFile {
	files := []sources.CodeFile{}

	w := NewKotlinWriter()
	w.Line(`package %s`, apiPackage.PackageName)
	w.EmptyLine()
	w.Line(`import com.fasterxml.jackson.core.JsonProcessingException`)
	w.Line(`import com.fasterxml.jackson.core.type.TypeReference`)
	w.Line(`import com.fasterxml.jackson.databind.ObjectMapper`)
	w.Line(`import com.fasterxml.jackson.module.kotlin.jacksonObjectMapper`)
	w.Line(`import okhttp3.*`)
	w.Line(`import okhttp3.MediaType.Companion.toMediaTypeOrNull`)
	w.Line(`import okhttp3.RequestBody.Companion.toRequestBody`)
	w.Line(`import org.slf4j.*`)
	w.Line(`import java.io.*`)
	w.Line(`import java.math.BigDecimal`)
	w.Line(`import java.time.*`)
	w.Line(`import java.util.*`)
	w.EmptyLine()
	w.Line(`import %s`, mainPackage.PackageStar)
	w.Line(`import %s`, utilsPackage.PackageStar)
	w.Line(`import %s`, modelsVersionPackage.PackageStar)
	w.Line(`import %s.setupObjectMapper`, jsonPackage.PackageName)
	w.EmptyLine()
	className := clientName(api)
	w.Line(`class %s(`, className)
	w.Line(`  private val baseUrl: String,`)
	w.Line(`  private val client: OkHttpClient = OkHttpClient(),`)
	w.Line(`  private val objectMapper: ObjectMapper = setupObjectMapper(jacksonObjectMapper())`)
	w.Line(`) {`)
	w.Line(`  private val logger: Logger = LoggerFactory.getLogger(%s::class.java)`, className)
	for _, operation := range api.Operations {
		w.EmptyLine()
		generateClientMethod(w.Indented(), &operation)
	}
	w.Line(`}`)

	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			files = append(files, generateResponseInterface(&operation, apiPackage, modelsVersionPackage)...)
		}
	}

	files = append(files, sources.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.kt", className)),
		Content: w.String(),
	})

	return files
}

func generateClientMethod(w *sources.Writer, operation *spec.NamedOperation) {
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()
	requestBody := "null"

	w.Line(`fun %s {`, generateResponsesSignatures(operation))
	bodyDataVar := "bodyJson"
	mediaType := "application/json"
	if operation.Body != nil {
		if operation.Body.Type.Definition.Plain == spec.TypeString {
			bodyDataVar = "body"
			mediaType = "text/plain"
		} else {
			generateClientTryCatch(w.Indented(), bodyDataVar,
				`objectMapper.writeValueAsString(body)`,
				`JsonProcessingException`, `e`,
				`"Failed to serialize JSON " + e.message`)
			w.EmptyLine()
		}
	}
	w.Line(`  val url = UrlBuilder(baseUrl)`)
	if operation.Api.Apis.GetUrl() != "" {
		w.Line(`  url.addPathSegment("%s")`, TrimSlash(operation.Api.Apis.GetUrl()))
	}
	for _, urlPart := range operation.Endpoint.UrlParts {
		part := TrimSlash(urlPart.Part)
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
			responseJavaType := KotlinType(&response.Type.Definition)
			responseBody := fmt.Sprintf(` objectMapper.readValue(response.body!!.string(), object: TypeReference<%s>(){})`, responseJavaType)
			if response.Type.Definition.Plain == spec.TypeString {
				responseBody = `response.body!!.string()`
			}
			generateClientTryCatch(w, "responseBody",
				responseBody,
				`IOException`, `e`,
				`"Failed to deserialize response body " + e.message`)
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
