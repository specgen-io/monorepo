package client

import (
	"fmt"
	"strings"

	"generator"
	"kotlin/imports"
	"kotlin/models"
	"kotlin/packages"
	"kotlin/types"
	"kotlin/writer"
	"spec"
)

var OkHttp = "okhttp"

type OkHttpGenerator struct {
	Types  *types.Types
	Models models.Generator
}

func NewOkHttpGenerator(types *types.Types, models models.Generator) *OkHttpGenerator {
	return &OkHttpGenerator{types, models}
}

func (g *OkHttpGenerator) Clients(version *spec.Version, thePackage packages.Package, modelsVersionPackage packages.Package, errorModelsPackage packages.Package, jsonPackage packages.Package, mainPackage packages.Package) []generator.CodeFile {
	files := []generator.CodeFile{}

	utilsPackage := thePackage.Subpackage("utils")
	for _, api := range version.Http.Apis {
		apiPackage := thePackage.Subpackage(api.Name.SnakeCase())
		files = append(files, responses(&api, g.Types, apiPackage, modelsVersionPackage, errorModelsPackage)...)
		files = append(files, *g.client(&api, apiPackage, modelsVersionPackage, errorModelsPackage, jsonPackage, utilsPackage, mainPackage))
	}
	files = append(files, g.utils(utilsPackage)...)
	files = append(files, *clientException(mainPackage))

	return files
}

func (g *OkHttpGenerator) client(api *spec.Api, apiPackage packages.Package, modelsVersionPackage packages.Package, errorModelsPackage packages.Package, jsonPackage packages.Package, utilsPackage packages.Package, mainPackage packages.Package) *generator.CodeFile {
	w := writer.NewKotlinWriter()
	w.Line(`package %s`, apiPackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.Models.ModelsDefinitionsImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Add(`okhttp3.*`)
	imports.Add(`okhttp3.MediaType.Companion.toMediaTypeOrNull`)
	imports.Add(`okhttp3.RequestBody.Companion.toRequestBody`)
	imports.Add(`org.slf4j.*`)
	imports.Add(mainPackage.PackageStar)
	imports.Add(utilsPackage.PackageStar)
	imports.Add(modelsVersionPackage.PackageStar)
	imports.Add(g.Models.SetupImport(jsonPackage))
	imports.Write(w)
	w.EmptyLine()
	className := clientName(api)
	w.Line(`class %s(private val baseUrl: String) {`, className)
	w.Line(`  %s`, g.Models.CreateJsonMapperField(""))
	w.Line(`  private val client: OkHttpClient`)
	w.EmptyLine()
	w.Line(`  private val logger: Logger = LoggerFactory.getLogger(%s::class.java)`, className)
	w.EmptyLine()
	w.Line(`  init {`)
	g.Models.InitJsonMapper(w.IndentedWith(2))
	w.Line(`    client = OkHttpClient()`)
	w.Line(`  }`)

	for _, operation := range api.Operations {
		w.EmptyLine()
		g.clientMethod(w.Indented(), &operation)
	}
	w.Line(`}`)

	return &generator.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.kt", className)),
		Content: w.String(),
	}
}

func (g *OkHttpGenerator) clientMethod(w *generator.Writer, operation *spec.NamedOperation) {
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()

	w.Line(`fun %s {`, operationSignature(g.Types, operation))
	requestBody := "null"
	if operation.BodyIs(spec.BodyString) {
		w.Line(`  val requestBody = body.toRequestBody("text/plain".toMediaTypeOrNull())`)
		requestBody = "requestBody"
	}
	if operation.BodyIs(spec.BodyJson) {
		bodyJson, exception := g.Models.WriteJson("body", &operation.Body.Type.Definition)
		generateClientTryCatch(w.Indented(), "bodyJson",
			bodyJson,
			exception, `e`,
			`"Failed to serialize JSON " + e.message`)
		w.EmptyLine()
		w.Line(`  val requestBody = bodyJson.toRequestBody("application/json".toMediaTypeOrNull())`)
		requestBody = "requestBody"
	}

	w.Line(`  val url = UrlBuilder(baseUrl)`)
	if operation.InApi.InHttp.GetUrl() != "" {
		w.Line(`  url.addPathSegments("%s")`, trimSlash(operation.InApi.InHttp.GetUrl()))
	}
	for _, urlPart := range operation.Endpoint.UrlParts {
		part := trimSlash(urlPart.Part)
		if urlPart.Param != nil {
			w.Line(`  url.addPathParameter(%s)`, urlPart.Param.Name.CamelCase())
		} else if len(part) > 0 {
			w.Line(`  url.addPathSegments("%s")`, part)
		}
	}
	for _, param := range operation.QueryParams {
		w.Line(`  url.addQueryParameter("%s", %s)`, param.Name.SnakeCase(), addBuilderParam(&param))
	}
	w.EmptyLine()
	w.Line(`  val request = RequestBuilder("%s", url.build(), %s)`, methodName, requestBody)
	for _, param := range operation.HeaderParams {
		w.Line(`  request.addHeaderParameter("%s", %s)`, param.Name.Source, addBuilderParam(&param))
	}
	w.EmptyLine()
	w.Line(`  logger.info("Sending request, operationId: %s.%s, method: %s, url: %s")`, operation.InApi.Name.Source, operation.Name.Source, methodName, url)
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

		if response.BodyIs(spec.BodyEmpty) {
			responseCode := responseCreate(&response, ``)
			if responseCode != "" {
				w.Line(responseCode)
			}
		}
		if response.BodyIs(spec.BodyString) {
			generateClientTryCatch(w, `responseBody`,
				`response.body!!.string()`,
				`IOException`, `e`,
				`"Failed to convert response body to string " + e.message`)
			w.Line(responseCreate(&response, `responseBody`))
		}
		if response.BodyIs(spec.BodyJson) {
			responseBody, exception := g.Models.ReadJson(`response.body!!.string()`, &response.Type.Definition)
			generateClientTryCatch(w, `responseBody`,
				responseBody,
				exception, `e`,
				`"Failed to deserialize response body " + e.message`)
			w.Line(responseCreate(&response, `responseBody!!`))
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

func (g *OkHttpGenerator) utils(thePackage packages.Package) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *g.requestBuilder(thePackage))
	files = append(files, *g.urlBuilder(thePackage))
	return files
}

func (g *OkHttpGenerator) requestBuilder(thePackage packages.Package) *generator.CodeFile {
	code := `
package [[.PackageName]]

import okhttp3.*

class RequestBuilder(method: String, url: HttpUrl, body: RequestBody?) {
    private val requestBuilder: Request.Builder

    init {
        requestBuilder = Request.Builder().url(url).method(method, body)
    }

    fun addHeaderParameter(name: String, value: Any): RequestBuilder {
        val valueStr = value.toString()
        this.requestBuilder.addHeader(name, valueStr)
        return this
    }

    fun <T> addHeaderParameter(name: String, values: List<T>): RequestBuilder {
        for (value in values) {
            this.addHeaderParameter(name, value!!)
        }
        return this
    }

    fun build(): Request {
        return this.requestBuilder.build()
    }
}
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("RequestBuilder.kt"),
		Content: strings.TrimSpace(code),
	}
}

func (g *OkHttpGenerator) urlBuilder(thePackage packages.Package) *generator.CodeFile {
	code := `
package [[.PackageName]]

import okhttp3.HttpUrl
import okhttp3.HttpUrl.Companion.toHttpUrl

class UrlBuilder(baseUrl: String) {
    private val urlBuilder: HttpUrl.Builder

    init {
        this.urlBuilder = baseUrl.toHttpUrl().newBuilder()
    }

    fun addQueryParameter(name: String, value: Any): UrlBuilder {
        val valueStr = value.toString()
        urlBuilder.addQueryParameter(name, valueStr)
        return this
    }

    fun <T> addQueryParameter(name: String, values: List<T>): UrlBuilder {
        for (value in values) {
            this.addQueryParameter(name, value!!)
        }
        return this
    }

    fun addPathSegments(value: String): UrlBuilder {
        this.urlBuilder.addPathSegments(value)
        return this
    }

    fun addPathParameter(value: Any): UrlBuilder {
        val valueStr = value.toString()
        this.urlBuilder.addPathSegment(valueStr)
        return this
    }

    fun build(): HttpUrl {
        return this.urlBuilder.build()
    }
}
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("UrlBuilder.kt"),
		Content: strings.TrimSpace(code),
	}
}
