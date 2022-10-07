package service

import (
	"fmt"
	"generator"
	"kotlin/types"
	"kotlin/writer"
	"spec"
)

func (g *Generator) responseInterface(operation *spec.NamedOperation) *generator.CodeFile {
	apiPackage := g.Packages.Version(operation.InApi.InHttp.InVersion).ServicesApi(operation.InApi)
	w := writer.NewKotlinWriter()
	w.Line(`package %s`, apiPackage.PackageName)
	w.EmptyLine()
	w.Line(`import %s`, g.Packages.Models(operation.InApi.InHttp.InVersion).PackageStar)
	w.Line(`import %s`, g.Packages.ErrorsModels.PackageStar)
	w.EmptyLine()
	w.Line(`interface %s {`, responseInterfaceName(operation))
	for index, response := range operation.Responses {
		if index > 0 {
			w.EmptyLine()
		}
		responseImpl(w.Indented(), g.Types, &response)
	}
	w.Line(`}`)

	return &generator.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.kt", responseInterfaceName(operation))),
		Content: w.String(),
	}
}

func responseImpl(w *generator.Writer, types *types.Types, response *spec.OperationResponse) {
	serviceResponseImplementationName := response.Name.PascalCase()
	if !response.Type.Definition.IsEmpty() {
		w.Line(`class %s(var body: %s) : %s`, serviceResponseImplementationName, types.Kotlin(&response.Type.Definition), responseInterfaceName(response.Operation))
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
