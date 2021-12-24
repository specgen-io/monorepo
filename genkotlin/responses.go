package genkotlin

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
)

func generateResponsesSignatures(operation *spec.NamedOperation) string {
	if len(operation.Responses) == 1 {
		for _, response := range operation.Responses {
			if !response.Type.Definition.IsEmpty() {
				return fmt.Sprintf(`%s(%s): %s`, operation.Name.CamelCase(), JoinDelimParams(addOperationResponseParams(operation)), KotlinType(&response.Type.Definition))
			} else {
				return fmt.Sprintf(`%s(%s)`, operation.Name.CamelCase(), JoinDelimParams(addOperationResponseParams(operation)))
			}
		}
	}
	if len(operation.Responses) > 1 {
		return fmt.Sprintf(`%s(%s): %s`, operation.Name.CamelCase(), JoinDelimParams(addOperationResponseParams(operation)), serviceResponseInterfaceName(operation))
	}
	return ""
}

func addOperationResponseParams(operation *spec.NamedOperation) []string {
	params := []string{}
	if operation.Body != nil {
		params = append(params, fmt.Sprintf("body: %s", KotlinType(&operation.Body.Type.Definition)))
	}
	for _, param := range operation.QueryParams {
		params = append(params, fmt.Sprintf("%s: %s", param.Name.CamelCase(), KotlinType(&param.Type.Definition)))
	}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf("%s: %s", param.Name.CamelCase(), KotlinType(&param.Type.Definition)))
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, fmt.Sprintf("%s: %s", param.Name.CamelCase(), KotlinType(&param.Type.Definition)))
	}
	return params
}

func generateResponseInterface(operation *spec.NamedOperation, apiPackage Module, modelsVersionPackage Module) []gen.TextFile {
	files := []gen.TextFile{}
	w := NewKotlinWriter()
	w.Line(`package %s`, apiPackage.PackageName)
	w.EmptyLine()
	w.Line(`import %s`, modelsVersionPackage.PackageStar)
	w.EmptyLine()
	w.Line(`interface %s {`, serviceResponseInterfaceName(operation))
	for index, response := range operation.Responses {
		if index > 0 {
			w.EmptyLine()
		}
		generateResponsesImplementations(w.Indented(), &response)
	}
	w.Line(`}`)

	files = append(files, gen.TextFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.kt", serviceResponseInterfaceName(operation))),
		Content: w.String(),
	})
	return files
}

func generateResponsesImplementations(w *gen.Writer, response *spec.NamedResponse) {
	serviceResponseImplementationName := response.Name.PascalCase()
	w.Line(`class %s : %s {`, serviceResponseImplementationName, serviceResponseInterfaceName(response.Operation))
	if !response.Type.Definition.IsEmpty() {
		w.Line(`  private lateinit var body: %s`, KotlinType(&response.Type.Definition))
		w.EmptyLine()
		w.Line(`  constructor()`)
		w.EmptyLine()
		w.Line(`  constructor(body: %s) {`, KotlinType(&response.Type.Definition))
		w.Line(`    this.body = body`)
		w.Line(`  }`)
	}
	w.Line(`}`)
}
