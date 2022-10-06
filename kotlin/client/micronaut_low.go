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

var MicronautLow = "micronaut-low-level"

type MicronautLowGenerator struct {
	Types    *types.Types
	Models   models.Generator
	Packages *Packages
}

func NewMicronautLowGenerator(types *types.Types, models models.Generator, packages *Packages) *MicronautLowGenerator {
	return &MicronautLowGenerator{types, models, packages}
}

func (g *MicronautLowGenerator) Clients(version *spec.Version, clientVersionPackage packages.Package, modelsVersionPackage packages.Package, errorModelsPackage packages.Package, jsonPackage packages.Package, mainPackage packages.Package) []generator.CodeFile {
	utilsPackage := mainPackage.Subpackage("utils")

	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		apiPackage := clientVersionPackage.Subpackage(api.Name.SnakeCase())
		files = append(files, responses(&api, g.Types, apiPackage, modelsVersionPackage, errorModelsPackage)...)
		files = append(files, *g.client(&api, apiPackage, modelsVersionPackage, errorModelsPackage, jsonPackage, utilsPackage, mainPackage))
	}

	files = append(files, g.utils(utilsPackage)...)
	files = append(files, converters(mainPackage)...)
	files = append(files, staticConfigFiles(mainPackage)...)
	files = append(files, *clientException(mainPackage))

	return files
}

func (g *MicronautLowGenerator) client(api *spec.Api, apiPackage packages.Package, modelsVersionPackage packages.Package, errorModelsPackage packages.Package, jsonPackage packages.Package, utilsPackage packages.Package, mainPackage packages.Package) *generator.CodeFile {
	w := writer.NewKotlinWriter()
	w.Line(`package %s`, apiPackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(`io.micronaut.http.HttpHeaders.*`)
	imports.Add(`io.micronaut.http.HttpRequest.*`)
	imports.Add(`io.micronaut.http.client.HttpClient`)
	imports.Add(`io.micronaut.http.client.annotation.Client`)
	imports.Add(`jakarta.inject.Singleton`)
	imports.Add(`org.slf4j.*`)
	imports.Add(mainPackage.PackageStar)
	imports.Add(utilsPackage.PackageStar)
	imports.Add(modelsVersionPackage.PackageStar)
	imports.Add(g.Models.ModelsDefinitionsImports()...)
	imports.Add(g.Types.Imports()...)
	imports.Write(w)
	w.EmptyLine()
	className := clientName(api)
	w.Line(`@Singleton`)
	w.Line(`class %s(`, className)
	//TODO
	w.Line(`  @param:Client(ClientConfiguration.BASE_URL)`)
	w.Line(`  private val client: HttpClient,`)
	w.Line(`  private val objectMapper: ObjectMapper`)
	w.Line(`) {`)
	w.Line(`  private val logger: Logger = LoggerFactory.getLogger(%s::class.java)`, className)

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

func (g *MicronautLowGenerator) clientMethod(w *generator.Writer, operation *spec.NamedOperation) {
	methodName := operation.Endpoint.Method
	url := operation.FullUrl()

	w.Line(`fun %s {`, operationSignature(g.Types, operation))
	requestBody := "body"
	if operation.BodyIs(spec.BodyJson) {
		bodyJson, exception := g.Models.WriteJson("body", &operation.Body.Type.Definition)
		generateClientTryCatch(w.Indented(), "bodyJson",
			bodyJson,
			exception, `e`,
			`"Failed to serialize JSON " + e.message`)
		w.EmptyLine()
		requestBody = "bodyJson"
	}

	w.Line(`  val url = UrlBuilder("%s")`, getUrl(operation))

	for _, urlPart := range operation.Endpoint.UrlParts {
		if urlPart.Param != nil {
			w.Line(`  url.pathParam("%s", %s)`, urlPart.Param.Name.CamelCase(), urlPart.Param.Name.CamelCase())
		}
	}
	for _, param := range operation.QueryParams {
		w.Line(`  url.queryParam("%s", %s)`, param.Name.SnakeCase(), addBuilderParam(&param))
	}
	w.EmptyLine()
	w.Line(`  val request = RequestBuilder(%s)`, requestBuilderParams(methodName, requestBody, operation))
	if operation.BodyIs(spec.BodyJson) {
		w.Line(`  request.headerParam(CONTENT_TYPE, "application/json")`)
	}
	if operation.BodyIs(spec.BodyString) {
		w.Line(`  request.headerParam(CONTENT_TYPE, "text/plain")`)
	}
	for _, param := range operation.HeaderParams {
		w.Line(`  request.headerParam("%s", %s)`, param.Name.Source, addBuilderParam(&param))
	}
	w.EmptyLine()
	w.Line(`  logger.info("Sending request, operationId: %s.%s, method: %s, url: %s")`, operation.InApi.Name.Source, operation.Name.Source, methodName, url)
	w.EmptyLine()
	w.Line(`  val response = client.toBlocking().exchange(request.build(), String::class.java)`)
	w.EmptyLine()
	w.Line(`  return when (response.code()) {`)
	for _, response := range operation.Responses {
		statusCode := spec.HttpStatusCode(response.Name)
		if isSuccessfulStatusCode(statusCode) {
			w.Line(`    %s -> {`, statusCode)
			w.IndentWith(3)
			w.Line(`logger.info("Received response with status code {}", response.code())`)

			if response.BodyIs(spec.BodyEmpty) {
				responseCode := responseCreate(&response, ``)
				if responseCode != "" {
					w.Line(responseCode)
				}
			}
			if response.BodyIs(spec.BodyString) {
				generateClientTryCatch(w, `responseBody`,
					`response.body()!!.toString()`,
					`IOException`, `e`,
					`"Failed to convert response body to string " + e.message`)
				w.Line(responseCreate(&response, `responseBody`))
			}
			if response.BodyIs(spec.BodyJson) {
				responseBody, exception := g.Models.ReadJson(`response.body()!!.toString()`, &response.Type.Definition)
				generateClientTryCatch(w, `responseBody`,
					responseBody,
					exception, `e`,
					`"Failed to deserialize response body " + e.message`)
				w.Line(responseCreate(&response, `responseBody!!`))
			}
			w.UnindentWith(3)
			w.Line(`    }`)
		}
	}
	w.Line(`    else -> {`)
	generateThrowClientException(w.IndentedWith(3), `"Unexpected status code received: " + response.code()`, ``)
	w.Line(`    }`)
	w.Line(`  }`)
	w.Line(`}`)
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
	params := fmt.Sprintf(`%s, %s, ::%s`, urlParam, requestBody, methodName)
	if operation.BodyIs(spec.BodyEmpty) {
		params = fmt.Sprintf(`%s, ::%s`, urlParam, methodName)
	}

	return params
}

func (g *MicronautLowGenerator) utils(thePackage packages.Package) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *g.requestBuilder(thePackage))
	files = append(files, *g.urlBuilder(thePackage))
	return files
}

func (g *MicronautLowGenerator) requestBuilder(thePackage packages.Package) *generator.CodeFile {
	code := `
package [[.PackageName]]

import io.micronaut.http.MutableHttpRequest
import java.net.URI

class RequestBuilder {
    private var requestBuilder: MutableHttpRequest<Any>

    constructor(url: URI, body: Any?, method: (URI, Any?) -> MutableHttpRequest<Any>) {
        this.requestBuilder = method(url, body)
    }

    constructor(url: URI, method: (URI) -> MutableHttpRequest<Any>) {
        this.requestBuilder = method(url)
    }

    fun headerParam(name: String, value: Any): RequestBuilder {
        val valueStr = value.toString()
        this.requestBuilder.header(name, valueStr)
        return this
    }

    fun <T> headerParam(name: String, values: List<T>): RequestBuilder {
        for (value in values) {
            this.headerParam(name, value!!)
        }
        return this
    }

    fun build(): MutableHttpRequest<Any> {
        return this.requestBuilder
    }
}
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("RequestBuilder.kt"),
		Content: strings.TrimSpace(code),
	}
}

func (g *MicronautLowGenerator) urlBuilder(thePackage packages.Package) *generator.CodeFile {
	code := `
package [[.PackageName]]

import io.micronaut.http.uri.UriBuilder
import java.net.URI
import java.util.*

class UrlBuilder(url: String) {
    private val uriBuilder: UriBuilder = UriBuilder.of(url)
    private val urlMap: MutableMap<String, Any> = mutableMapOf()

    fun queryParam(name: String, value: Any): UrlBuilder {
        val valueStr = value.toString()
        uriBuilder.queryParam(name, valueStr)
        return this
    }

    fun <T> queryParam(name: String, values: List<T>): UrlBuilder {
        for (value in values) {
            this.queryParam(name, value!!)
        }
        return this
    }

    fun pathParam(name: String, value: Any): UrlBuilder {
        this.uriBuilder.path("{$name}")
        this.urlMap += mapOf(name to value)

        return this
    }

    fun expand(): URI {
        return this.uriBuilder.expand(
            Collections.checkedMap(
                this.urlMap, String::class.java, Any::class.java
            )
        )
    }

    fun build(): URI {
        return this.uriBuilder.build()
    }
}
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("UrlBuilder.kt"),
		Content: strings.TrimSpace(code),
	}
}
