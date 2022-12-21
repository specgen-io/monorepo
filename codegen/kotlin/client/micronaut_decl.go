package client

import (
	"fmt"
	"generator"
	"github.com/pinzolo/casee"
	"kotlin/models"
	"kotlin/types"
	"kotlin/writer"
	"spec"
	"strings"
)

var MicronautDecl = "micronaut-declarative"

type MicronautDeclGenerator struct {
	Types    *types.Types
	Models   models.Generator
	Packages *Packages
}

func NewMicronautDeclGenerator(types *types.Types, models models.Generator, packages *Packages) *MicronautDeclGenerator {
	return &MicronautDeclGenerator{types, models, packages}
}

func (g *MicronautDeclGenerator) Clients(version *spec.Version) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, g.responses(&api)...)
		files = append(files, *g.client(&api))
	}
	files = append(files, converters(g.Packages.Converters)...)
	files = append(files, staticConfigFiles(g.Packages.Root, g.Packages.Json)...)
	return files
}

func (g *MicronautDeclGenerator) client(api *spec.Api) *generator.CodeFile {
	w := writer.New(g.Packages.Client(api), clientName(api))
	w.Imports.Add(`io.micronaut.http.*`)
	w.Imports.Add(`io.micronaut.http.annotation.*`)
	w.Imports.Add(`io.micronaut.http.client.annotation.Client`)
	w.Imports.PackageStar(g.Packages.Root)
	w.Imports.PackageStar(g.Packages.Models(api.InHttp.InVersion))
	w.Imports.Add(g.Types.Imports()...)
	w.EmptyLine()
	w.Line(`@Client(ClientConfiguration.BASE_URL)`)
	w.Line(`interface [[.ClassName]] {`)
	for _, operation := range api.Operations {
		w.EmptyLine()
		g.clientMethod(w.Indented(), &operation)
	}
	w.Line(`}`)
	return w.ToCodeFile()
}

func (g *MicronautDeclGenerator) clientMethod(w *writer.Writer, operation *spec.NamedOperation) {
	methodName := casee.ToPascalCase(operation.Endpoint.Method)
	url := operation.FullUrl()

	if operation.BodyIs(spec.BodyString) {
		w.Line(`@%s(value = "%s", processes = [MediaType.TEXT_PLAIN])`, methodName, url)
	} else if operation.BodyIs(spec.BodyJson) {
		w.Line(`@%s(value = "%s", processes = [MediaType.APPLICATION_JSON])`, methodName, url)
	} else {
		w.Line(`@%s(value = "%s")`, methodName, url)
	}
	w.Line(`fun %s`, g.operationSignature(operation))
}

func (g *MicronautDeclGenerator) operationSignature(operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%s(%s): %s`,
		operation.Name.CamelCase(),
		strings.Join(g.operationParameters(operation), ", "),
		g.operationReturnType(operation),
	)
}

func (g *MicronautDeclGenerator) operationReturnType(operation *spec.NamedOperation) string {
	if len(operation.Responses.Success()) == 1 {
		return g.Types.Kotlin(&operation.Responses.Success()[0].Type.Definition)
	}
	return "HttpResponse<String>"
}

func (g *MicronautDeclGenerator) operationParameters(operation *spec.NamedOperation) []string {
	params := []string{}
	if operation.Body != nil {
		params = append(params, fmt.Sprintf("@Body body: %s", g.Types.Kotlin(&operation.Body.Type.Definition)))
	}
	for _, param := range operation.QueryParams {
		params = append(params, fmt.Sprintf(`@QueryValue(value = "%s") %s: %s`, param.Name.Source, param.Name.CamelCase(), g.Types.Kotlin(&param.Type.Definition)))
	}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf(`@Header(value = "%s") %s: %s`, param.Name.Source, param.Name.CamelCase(), g.Types.Kotlin(&param.Type.Definition)))
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, fmt.Sprintf(`@PathVariable(value = "%s") %s: %s`, param.Name.Source, param.Name.CamelCase(), g.Types.Kotlin(&param.Type.Definition)))
	}
	return params
}

func (g *MicronautDeclGenerator) responses(api *spec.Api) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, operation := range api.Operations {
		if len(operation.Responses.Success()) > 1 {
			files = append(files, *g.response(&operation))
		}
	}
	return files
}

func (g *MicronautDeclGenerator) response(operation *spec.NamedOperation) *generator.CodeFile {
	w := writer.New(g.Packages.Client(operation.InApi), responseName(operation))
	w.Imports.Add(`com.fasterxml.jackson.core.type.TypeReference`)
	w.Imports.Add(`io.micronaut.http.HttpResponse`)
	w.Imports.PackageStar(g.Packages.Json)
	w.Imports.PackageStar(g.Packages.Models(operation.InApi.InHttp.InVersion))
	w.Imports.PackageStar(g.Packages.Utils)
	w.Imports.PackageStar(g.Packages.Errors)
	w.EmptyLine()
	w.Line(`open class [[.ClassName]] {`)
	for index, response := range operation.Responses {
		if index > 0 {
			w.EmptyLine()
		}
		g.implementations(w.Indented(), &response)
	}
	w.EmptyLine()
	g.createObjectMethod(w.Indented(), operation)
	w.Line(`}`)
	return w.ToCodeFile()
}

func (g *MicronautDeclGenerator) implementations(w *writer.Writer, response *spec.OperationResponse) {
	responseImplementationName := response.Name.PascalCase()
	if !response.Type.Definition.IsEmpty() {
		w.Line(`class %s(val body: %s) : %s()`, responseImplementationName, g.Types.Kotlin(&response.Type.Definition), responseName(response.Operation))
	} else {
		w.Line(`class %s : %s()`, responseImplementationName, responseName(response.Operation))
	}
}

func (g *MicronautDeclGenerator) createObjectMethod(w *writer.Writer, operation *spec.NamedOperation) {
	w.Lines(`
companion object {
	fun create(json: Json, response: HttpResponse<String>): EchoSuccessResponse {
		val responseBodyString = getResponseBodyString(response)
		return when(response.code()) {
`)
	for _, response := range operation.Responses {
		if !response.BodyIs(spec.BodyEmpty) {
			w.Line(`      %s -> %s(json.%s)`, spec.HttpStatusCode(response.Name), response.Name.PascalCase(), g.Models.ReadJson("responseBodyString", &response.Type.Definition))
		}
	}
	w.Lines(`
			else -> throw ClientException("Unknown status code ${response.code()}")
		}
	}
}
`)
}

func responseName(operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%sResponse`, operation.Name.PascalCase())
}

func (g *MicronautDeclGenerator) Utils() []generator.CodeFile {
	return []generator.CodeFile{*g.generateClientResponse()}
}

func (g *MicronautDeclGenerator) generateClientResponse() *generator.CodeFile {
	w := writer.New(g.Packages.Utils, `ClientResponse`)
	w.Template(
		map[string]string{
			`ErrorsPackage`: g.Packages.Errors.PackageName,
		}, `
import io.micronaut.http.HttpResponse
import [[.ErrorsPackage]].*
import java.io.IOException

fun <T> getResponseBodyString(response: HttpResponse<T>): String {
	return try {
		response.body()!!.toString()
	} catch (e: IOException) {
		throw ClientException("Failed to convert response body to string " + e.message, e)
	}
}
`)
	return w.ToCodeFile()
}

func (g *MicronautDeclGenerator) Exceptions(errors *spec.ErrorResponses) []generator.CodeFile {
	return []generator.CodeFile{*clientException(g.Packages.Errors)}
}
