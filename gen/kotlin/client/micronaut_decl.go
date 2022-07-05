package client

import (
	"fmt"
	"github.com/pinzolo/casee"
	"github.com/specgen-io/specgen/v2/gen/kotlin/imports"
	"github.com/specgen-io/specgen/v2/gen/kotlin/models"
	"github.com/specgen-io/specgen/v2/gen/kotlin/modules"
	"github.com/specgen-io/specgen/v2/gen/kotlin/types"
	"github.com/specgen-io/specgen/v2/gen/kotlin/writer"
	"github.com/specgen-io/specgen/v2/generator"
	"github.com/specgen-io/specgen/v2/spec"
	"strconv"
)

var MicronautDecl = "micronaut-declarative"

type MicronautDeclGenerator struct {
	Types  *types.Types
	Models models.Generator
}

func NewMicronautDeclGenerator(types *types.Types, models models.Generator) *MicronautDeclGenerator {
	return &MicronautDeclGenerator{types, models}
}

func (g *MicronautDeclGenerator) ClientImplementation(version *spec.Version, thePackage modules.Module, modelsVersionPackage modules.Module, jsonPackage modules.Module, mainPackage modules.Module) []generator.CodeFile {
	files := []generator.CodeFile{}

	for _, api := range version.Http.Apis {
		apiPackage := thePackage.Subpackage(api.Name.SnakeCase())
		files = append(files, g.client(&api, apiPackage, modelsVersionPackage, jsonPackage, mainPackage)...)
	}

	return files
}

func (g *MicronautDeclGenerator) client(api *spec.Api, apiPackage modules.Module, modelsVersionPackage modules.Module, jsonPackage modules.Module, mainPackage modules.Module) []generator.CodeFile {
	files := []generator.CodeFile{}

	w := writer.NewKotlinWriter()
	w.Line(`package %s`, apiPackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(`io.micronaut.http.*`)
	imports.Add(`io.micronaut.http.annotation.*`)
	imports.Add(`io.micronaut.http.client.annotation.Client`)
	imports.Add(mainPackage.PackageName + `.ServerConfiguration`)
	imports.Add(modelsVersionPackage.PackageStar)
	imports.Add(g.Types.Imports()...)
	imports.Write(w)
	w.EmptyLine()
	interfaceName := clientName(api)
	w.Line(`@Client(ServerConfiguration.BASE_URL)`)
	w.Line(`interface %s {`, interfaceName)
	for _, operation := range api.Operations {
		w.EmptyLine()
		g.clientMethod(w.Indented(), &operation)
	}
	w.Line(`}`)

	for _, operation := range api.Operations {
		if responsesNumber(&operation) > 1 {
			files = append(files, response(g.Types, &operation, apiPackage, modelsVersionPackage)...)
		}
	}

	files = append(files, generator.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.kt", interfaceName)),
		Content: w.String(),
	})

	return files
}

func (g *MicronautDeclGenerator) clientMethod(w *generator.Writer, operation *spec.NamedOperation) {
	methodName := casee.ToPascalCase(operation.Endpoint.Method)
	url := operation.FullUrl()

	if operation.BodyIs(spec.BodyString) {
		w.Line(`@%s(value = "%s", processes = [MediaType.TEXT_PLAIN])`, methodName, url)
	} else if operation.BodyIs(spec.BodyJson) {
		w.Line(`@%s(value = "%s", processes = [MediaType.APPLICATION_JSON])`, methodName, url)
	} else {
		w.Line(`@%s(value = "%s")`, methodName, url)
	}
	w.Line(`fun %s`, declSignature(g.Types, operation))
}

func declSignature(types *types.Types, operation *spec.NamedOperation) string {
	if responsesNumber(operation) == 1 {
		for _, response := range operation.Responses {
			if !response.Type.Definition.IsEmpty() {
				return fmt.Sprintf(`%s(%s): %s`, operation.Name.CamelCase(), joinParams(parameters(operation, types)), types.Kotlin(&response.Type.Definition))
			} else {
				return fmt.Sprintf(`%s(%s)`, operation.Name.CamelCase(), joinParams(parameters(operation, types)))
			}
		}
	}
	if responsesNumber(operation) > 1 {
		return fmt.Sprintf(`%s(%s): HttpResponse<String>`, operation.Name.CamelCase(), joinParams(parameters(operation, types)))
	}
	return ""
}

func parameters(operation *spec.NamedOperation, types *types.Types) []string {
	params := []string{}
	if operation.Body != nil {
		params = append(params, fmt.Sprintf("@Body body: %s", types.Kotlin(&operation.Body.Type.Definition)))
	}
	for _, param := range operation.QueryParams {
		params = append(params, fmt.Sprintf("@QueryValue %s: %s", param.Name.CamelCase(), types.Kotlin(&param.Type.Definition)))
	}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf("@Header %s: %s", param.Name.CamelCase(), types.Kotlin(&param.Type.Definition)))
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, fmt.Sprintf("@PathVariable %s: %s", param.Name.CamelCase(), types.Kotlin(&param.Type.Definition)))
	}
	return params
}

func response(types *types.Types, operation *spec.NamedOperation, apiPackage modules.Module, modelsVersionPackage modules.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	w := writer.NewKotlinWriter()
	w.Line(`package %s`, apiPackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(`com.fasterxml.jackson.databind.ObjectMapper`)
	imports.Add(`io.micronaut.http.HttpResponse`)
	imports.Add(modelsVersionPackage.PackageStar)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`open class %s {`, responseName(operation))
	for index, response := range operation.Responses {
		if index > 0 {
			w.EmptyLine()
		}
		implementations(w.Indented(), types, &response)
	}
	w.EmptyLine()
	createObjectMethod(w.Indented(), types, operation)
	w.Line(`}`)

	files = append(files, generator.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.kt", responseName(operation))),
		Content: w.String(),
	})
	return files
}

func implementations(w *generator.Writer, types *types.Types, response *spec.OperationResponse) {
	responseImplementationName := response.Name.PascalCase()
	if !response.Type.Definition.IsEmpty() {
		w.Line(`class %s(val body: %s) : %s()`, responseImplementationName, types.Kotlin(&response.Type.Definition), responseName(response.Operation))
	} else {
		w.Line(`class %s : %s()`, responseImplementationName, responseName(response.Operation))
	}
}

func createObjectMethod(w *generator.Writer, types *types.Types, operation *spec.NamedOperation) {
	w.Line(`companion object {`)
	w.Line(`  fun create(json: ObjectMapper, response: HttpResponse<String>): %s {`, responseName(operation.Responses[0].Operation))
	w.Line(`    return when(response.code()) {`)
	for _, response := range operation.Responses {
		if !response.BodyIs(spec.BodyEmpty) {
			w.Line(`      %s -> %s(json.readValue(response.body(), %s::class.java))`, spec.HttpStatusCode(response.Name), response.Name.PascalCase(), types.Kotlin(&response.Type.Definition))
		}
	}
	w.Line(`      else -> throw Exception("Unknown status code ${response.code()}")`)
	w.Line(`    }`)
	w.Line(`  }`)
	w.Line(`}`)
}

func responseName(operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%sResponse`, operation.Name.PascalCase())
}

func isSuccessfulStatusCode(statusCodeStr string) bool {
	statusCode, _ := strconv.Atoi(statusCodeStr)
	if statusCode >= 200 && statusCode <= 299 {
		return true
	}
	return false
}

func responsesNumber(operation *spec.NamedOperation) int {
	count := 0
	for _, response := range operation.Responses {
		if isSuccessfulStatusCode(spec.HttpStatusCode(response.Name)) {
			count++
		}
	}
	return count
}
