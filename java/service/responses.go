package service

import (
	"fmt"
	"generator"
	"java/types"
	"java/writer"
	"spec"
	"strings"
)

func operationSignature(types *types.Types, operation *spec.NamedOperation) string {
	if len(operation.Responses) == 1 {
		for _, response := range operation.Responses {
			return fmt.Sprintf(`%s %s(%s)`, types.Java(&response.Type.Definition), operation.Name.CamelCase(), strings.Join(parameters(operation, types), ", "))
		}
	}
	if len(operation.Responses) > 1 {
		return fmt.Sprintf(`%s %s(%s)`, responseInterfaceName(operation), operation.Name.CamelCase(), strings.Join(parameters(operation, types), ", "))
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

func (g *Generator) responseInterface(operation *spec.NamedOperation) *generator.CodeFile {
	packages := g.Packages.Version(operation.InApi.InHttp.InVersion)
	apiPackage := packages.ServicesApi(operation.InApi)

	w := writer.NewJavaWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	w.Line(`import %s;`, packages.Models.PackageStar)
	w.Line(`import %s;`, g.Packages.ErrorsModels.PackageStar)
	w.EmptyLine()
	w.Line(`public interface %s {`, responseInterfaceName(operation))
	for index, response := range operation.Responses {
		if index > 0 {
			w.EmptyLine()
		}
		responseImpl(w.Indented(), g.Types, &response)
	}
	w.Line(`}`)

	return &generator.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", responseInterfaceName(operation))),
		Content: w.String(),
	}
}

func responseImpl(w *generator.Writer, types *types.Types, response *spec.OperationResponse) {
	serviceResponseImplementationName := response.Name.PascalCase()
	w.Line(`class %s implements %s {`, serviceResponseImplementationName, responseInterfaceName(response.Operation))
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

func responseInterfaceName(operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%sResponse`, operation.Name.PascalCase())
}

func getResponseBody(response *spec.OperationResponse, responseVarName string) string {
	return fmt.Sprintf(`((%s.%s) %s).body`, responseInterfaceName(response.Operation), response.Name.PascalCase(), responseVarName)
}
