package responses

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen/kotlin/modules"
	"github.com/specgen-io/specgen/v2/gen/kotlin/types"
	"github.com/specgen-io/specgen/v2/gen/kotlin/writer"
	"github.com/specgen-io/specgen/v2/generator"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func Signature(types *types.Types, operation *spec.NamedOperation) string {
	if len(operation.Responses) == 1 {
		for _, response := range operation.Responses {
			if !response.Type.Definition.IsEmpty() {
				return fmt.Sprintf(`%s(%s): %s`, operation.Name.CamelCase(), joinDelimParams(Parameters(operation, types)), types.Kotlin(&response.Type.Definition))
			} else {
				return fmt.Sprintf(`%s(%s)`, operation.Name.CamelCase(), joinDelimParams(Parameters(operation, types)))
			}
		}
	}
	if len(operation.Responses) > 1 {
		return fmt.Sprintf(`%s(%s): %s`, operation.Name.CamelCase(), joinDelimParams(Parameters(operation, types)), InterfaceName(operation))
	}
	return ""
}

func CreateResponse(response *spec.OperationResponse, resultVar string) string {
	if len(response.Operation.Responses) > 1 {
		return fmt.Sprintf(`%s.%s(%s)`, InterfaceName(response.Operation), response.Name.PascalCase(), resultVar)
	}
	return resultVar
}

func Parameters(operation *spec.NamedOperation, types *types.Types) []string {
	params := []string{}
	if operation.Body != nil {
		params = append(params, fmt.Sprintf("body: %s", types.Kotlin(&operation.Body.Type.Definition)))
	}
	for _, param := range operation.QueryParams {
		params = append(params, fmt.Sprintf("%s: %s", param.Name.CamelCase(), types.Kotlin(&param.Type.Definition)))
	}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf("%s: %s", param.Name.CamelCase(), types.Kotlin(&param.Type.Definition)))
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, fmt.Sprintf("%s: %s", param.Name.CamelCase(), types.Kotlin(&param.Type.Definition)))
	}
	return params
}

func Interfaces(types *types.Types, operation *spec.NamedOperation, apiPackage modules.Module, modelsVersionPackage modules.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
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
		implementations(w.Indented(), types, &response)
	}
	w.Line(`}`)

	files = append(files, generator.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.kt", InterfaceName(operation))),
		Content: w.String(),
	})
	return files
}

func implementations(w *generator.Writer, types *types.Types, response *spec.OperationResponse) {
	serviceResponseImplementationName := response.Name.PascalCase()
	if !response.Type.Definition.IsEmpty() {
		w.Line(`class %s(var body: %s) : %s`, serviceResponseImplementationName, types.Kotlin(&response.Type.Definition), InterfaceName(response.Operation))
	} else {
		w.Line(`class %s : %s`, serviceResponseImplementationName, InterfaceName(response.Operation))
	}
}

func GetBody(varName string) string {
	return fmt.Sprintf(`%s.body`, varName)
}

func InterfaceName(operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%sResponse`, operation.Name.PascalCase())
}

func joinDelimParams(params []string) string {
	return strings.Join(params, ", ")
}
