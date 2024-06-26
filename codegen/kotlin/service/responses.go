package service

import (
	"fmt"
	"generator"
	"kotlin/types"
	"kotlin/writer"
	"spec"
)

func (g *Generator) responseInterface(operation *spec.NamedOperation) *generator.CodeFile {
	w := writer.New(g.Packages.ServicesApi(operation.InApi), responseInterfaceName(operation))
	w.Imports.PackageStar(g.Packages.Models(operation.InApi.InHttp.InVersion))
	w.Imports.PackageStar(g.Packages.ErrorsModels)
	w.Line(`interface [[.ClassName]] {`)
	for index, response := range operation.Responses {
		if index > 0 {
			w.EmptyLine()
		}
		responseImpl(w.Indented(), g.Types, &response)
	}
	w.Line(`}`)
	return w.ToCodeFile()
}

func responseImpl(w *writer.Writer, types *types.Types, response *spec.OperationResponse) {
	serviceResponseImplementationName := response.Name.PascalCase()
	if !response.Body.IsEmpty() {
		w.Line(`class %s(var body: %s) : %s`, serviceResponseImplementationName, types.Kotlin(&response.Body.Type.Definition), responseInterfaceName(response.Operation))
	} else {
		w.Line(`class %s : %s`, serviceResponseImplementationName, responseInterfaceName(response.Operation))
	}
}

func responseInterfaceName(operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%sResponse`, operation.Name.PascalCase())
}

func getResponseBody(varName string) string {
	return fmt.Sprintf(`%s.body`, varName)
}
