package responses

import (
	"fmt"
	"strings"

	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/spec/v2"
	"github.com/specgen-io/specgen/v2/gen/java/packages"
	"github.com/specgen-io/specgen/v2/gen/java/types"
	"github.com/specgen-io/specgen/v2/gen/java/writer"
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

func CreateResponse(response *spec.OperationResponse, resultVar string) string {
	if len(response.Operation.Responses) > 1 {
		return fmt.Sprintf(`return new %s.%s(%s);`, InterfaceName(response.Operation), response.Name.PascalCase(), resultVar)
	} else {
		if resultVar != "" {
			return fmt.Sprintf(`return %s;`, resultVar)
		} else {
			return `return;`
		}
	}
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

func Interfaces(types *types.Types, operation *spec.NamedOperation, apiPackage packages.Module, modelsVersionPackage packages.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
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

	files = append(files, generator.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", InterfaceName(operation))),
		Content: w.String(),
	})
	return files
}

func implementation(w *generator.Writer, types *types.Types, response *spec.OperationResponse) {
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

func GetBody(response *spec.OperationResponse, responseVarName string) string {
	return fmt.Sprintf(`((%s.%s) %s).body`, InterfaceName(response.Operation), response.Name.PascalCase(), responseVarName)
}

func joinParams(params []string) string {
	return strings.Join(params, ", ")
}
