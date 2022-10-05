package service

import (
	"fmt"
	"generator"
	"java/types"
	"java/writer"
	"spec"
)

func (g *Generator) responseInterface(operation *spec.NamedOperation) *generator.CodeFile {
	version := operation.InApi.InHttp.InVersion
	apiPackage := g.Packages.Version(version).ServicesApi(operation.InApi)

	w := writer.NewJavaWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	w.Line(`import %s;`, g.Packages.Models(version).PackageStar)
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
