package responses

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/genjava/packages"
	"github.com/specgen-io/specgen/v2/genjava/types"
	"github.com/specgen-io/specgen/v2/genjava/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func Signature(types *types.Types, operation *spec.NamedOperation) string {
	if len(operation.Responses) == 1 {
		for _, response := range operation.Responses {
			return fmt.Sprintf(`%s %s(%s)`, types.Java(&response.Type.Definition), operation.Name.CamelCase(), joinParams(parameters(operation, types)))
		}
	}
	if len(operation.Responses) > 1 {
		return fmt.Sprintf(`%s %s(%s)`, InterfaceName(operation), operation.Name.CamelCase(), joinParams(parameters(operation, types)))
	}
	return ""
}

func parameters(operation *spec.NamedOperation, types *types.Types) []string {
	params := []string{}
	if operation.Body != nil {
		params = append(params, fmt.Sprintf("%s body", types.Java(&operation.Body.Type.Definition)))
	}
	for _, param := range operation.QueryParams {
		params = append(params, fmt.Sprintf("%s %s", types.Java(&param.Type.Definition), param.Name.CamelCase()))
	}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf("%s %s", types.Java(&param.Type.Definition), param.Name.CamelCase()))
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, fmt.Sprintf("%s %s", types.Java(&param.Type.Definition), param.Name.CamelCase()))
	}
	return params
}

func Interfaces(types *types.Types, operation *spec.NamedOperation, apiPackage packages.Module, modelsVersionPackage packages.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	w := writer.NewJavaWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	w.Line(`import %s;`, modelsVersionPackage.PackageStar)
	w.EmptyLine()
	w.Line(`public interface %s {`, InterfaceName(operation))
	for index, response := range operation.Responses {
		if index > 0 {
			w.EmptyLine()
		}
		implementation(w.Indented(), types, &response)
	}
	w.Line(`}`)

	files = append(files, sources.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", InterfaceName(operation))),
		Content: w.String(),
	})
	return files
}

func implementation(w *sources.Writer, types *types.Types, response *spec.NamedResponse) {
	serviceResponseImplementationName := response.Name.PascalCase()
	w.Line(`class %s implements %s {`, serviceResponseImplementationName, InterfaceName(response.Operation))
	if !response.Type.Definition.IsEmpty() {
		w.Line(`  public %s body;`, types.Java(&response.Type.Definition))
		w.EmptyLine()
		w.Line(`  public %s() {`, serviceResponseImplementationName)
		w.Line(`  }`)
		w.EmptyLine()
		w.Line(`  public %s(%s body) {`, serviceResponseImplementationName, types.Java(&response.Type.Definition))
		w.Line(`    this.body = body;`)
		w.Line(`  }`)
	}
	w.Line(`}`)
}

func InterfaceName(operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%sResponse`, operation.Name.PascalCase())
}

func joinParams(params []string) string {
	return strings.Join(params, ", ")
}
