package responses

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/genkotlin/modules"
	"github.com/specgen-io/specgen/v2/genkotlin/types"
	"github.com/specgen-io/specgen/v2/genkotlin/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func GenerateResponsesSignatures(operation *spec.NamedOperation) string {
	if len(operation.Responses) == 1 {
		for _, response := range operation.Responses {
			if !response.Type.Definition.IsEmpty() {
				return fmt.Sprintf(`%s(%s): %s`, operation.Name.CamelCase(), joinDelimParams(addOperationResponseParams(operation)), types.KotlinType(&response.Type.Definition))
			} else {
				return fmt.Sprintf(`%s(%s)`, operation.Name.CamelCase(), joinDelimParams(addOperationResponseParams(operation)))
			}
		}
	}
	if len(operation.Responses) > 1 {
		return fmt.Sprintf(`%s(%s): %s`, operation.Name.CamelCase(), joinDelimParams(addOperationResponseParams(operation)), InterfaceName(operation))
	}
	return ""
}

func addOperationResponseParams(operation *spec.NamedOperation) []string {
	params := []string{}
	if operation.Body != nil {
		params = append(params, fmt.Sprintf("body: %s", types.KotlinType(&operation.Body.Type.Definition)))
	}
	for _, param := range operation.QueryParams {
		params = append(params, fmt.Sprintf("%s: %s", param.Name.CamelCase(), types.KotlinType(&param.Type.Definition)))
	}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf("%s: %s", param.Name.CamelCase(), types.KotlinType(&param.Type.Definition)))
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, fmt.Sprintf("%s: %s", param.Name.CamelCase(), types.KotlinType(&param.Type.Definition)))
	}
	return params
}

func GenerateResponseInterface(operation *spec.NamedOperation, apiPackage modules.Module, modelsVersionPackage modules.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	w := writer.NewKotlinWriter()
	w.Line(`package %s`, apiPackage.PackageName)
	w.EmptyLine()
	w.Line(`import %s`, modelsVersionPackage.PackageStar)
	w.EmptyLine()
	w.Line(`interface %s {`, InterfaceName(operation))
	for index, response := range operation.Responses {
		if index > 0 {
			w.EmptyLine()
		}
		generateResponsesImplementations(w.Indented(), &response)
	}
	w.Line(`}`)

	files = append(files, sources.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.kt", InterfaceName(operation))),
		Content: w.String(),
	})
	return files
}

func generateResponsesImplementations(w *sources.Writer, response *spec.NamedResponse) {
	serviceResponseImplementationName := response.Name.PascalCase()
	w.Line(`class %s : %s {`, serviceResponseImplementationName, InterfaceName(response.Operation))
	if !response.Type.Definition.IsEmpty() {
		w.Line(`  private lateinit var body: %s`, types.KotlinType(&response.Type.Definition))
		w.EmptyLine()
		w.Line(`  constructor()`)
		w.EmptyLine()
		w.Line(`  constructor(body: %s) {`, types.KotlinType(&response.Type.Definition))
		w.Line(`    this.body = body`)
		w.Line(`  }`)
	}
	w.Line(`}`)
}

func InterfaceName(operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%sResponse`, operation.Name.PascalCase())
}

func joinDelimParams(params []string) string {
	return strings.Join(params, ", ")
}
