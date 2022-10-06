package client

import (
	"fmt"
	"strings"

	"generator"
	"github.com/pinzolo/casee"
	"kotlin/imports"
	"kotlin/models"
	"kotlin/packages"
	"kotlin/types"
	"kotlin/writer"
	"spec"
)

var MicronautDecl = "micronaut-declarative"

type MicronautDeclGenerator struct {
	Types  *types.Types
	Models models.Generator
}

func NewMicronautDeclGenerator(types *types.Types, models models.Generator) *MicronautDeclGenerator {
	return &MicronautDeclGenerator{types, models}
}

func (g *MicronautDeclGenerator) Clients(version *spec.Version, thePackage packages.Package, modelsVersionPackage packages.Package, errorModelsPackage packages.Package, jsonPackage packages.Package, mainPackage packages.Package) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		apiPackage := thePackage.Subpackage(api.Name.SnakeCase())
		files = append(files, g.responses(&api, apiPackage, modelsVersionPackage)...)
		files = append(files, *g.client(&api, apiPackage, modelsVersionPackage, mainPackage))
	}
	files = append(files, converters(mainPackage)...)
	files = append(files, staticConfigFiles(mainPackage)...)

	return files
}

func (g *MicronautDeclGenerator) client(api *spec.Api, apiPackage packages.Package, modelsVersionPackage packages.Package, mainPackage packages.Package) *generator.CodeFile {
	w := writer.NewKotlinWriter()
	w.Line(`package %s`, apiPackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(`io.micronaut.http.*`)
	imports.Add(`io.micronaut.http.annotation.*`)
	imports.Add(`io.micronaut.http.client.annotation.Client`)
	imports.Add(mainPackage.PackageStar)
	imports.Add(modelsVersionPackage.PackageStar)
	imports.Add(g.Types.Imports()...)
	imports.Write(w)
	w.EmptyLine()
	interfaceName := clientName(api)
	//TODO
	w.Line(`@Client(ClientConfiguration.BASE_URL)`)
	w.Line(`interface %s {`, interfaceName)
	for _, operation := range api.Operations {
		w.EmptyLine()
		g.clientMethod(w.Indented(), &operation)
	}
	w.Line(`}`)

	return &generator.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.kt", interfaceName)),
		Content: w.String(),
	}
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
	w.Line(`fun %s`, g.operationSignature(operation))
}

func (g *MicronautDeclGenerator) operationSignature(operation *spec.NamedOperation) string {
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

	if len(successResponses(operation)) == 1 {
		for _, response := range operation.Responses {
			if !response.Type.Definition.IsEmpty() {
				return fmt.Sprintf(`%s(%s): %s`, operation.Name.CamelCase(), strings.Join(params, ", "), g.Types.Kotlin(&response.Type.Definition))
			} else {
				return fmt.Sprintf(`%s(%s)`, operation.Name.CamelCase(), strings.Join(params, ", "))
			}
		}
	}
	if len(successResponses(operation)) > 1 {
		return fmt.Sprintf(`%s(%s): HttpResponse<String>`, operation.Name.CamelCase(), strings.Join(params, ", "))
	}
	return ""
}

func (g *MicronautDeclGenerator) responses(api *spec.Api, apiPackage packages.Package, modelsVersionPackage packages.Package) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, operation := range api.Operations {
		if len(successResponses(&operation)) > 1 {
			files = append(files, *g.response(g.Types, &operation, apiPackage, modelsVersionPackage))
		}
	}
	return files
}

func (g *MicronautDeclGenerator) response(types *types.Types, operation *spec.NamedOperation, apiPackage packages.Package, modelsVersionPackage packages.Package) *generator.CodeFile {
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

	return &generator.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.kt", responseName(operation))),
		Content: w.String(),
	}
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
