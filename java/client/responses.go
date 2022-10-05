package client

import (
	"fmt"
	"generator"
	"java/types"
	"java/writer"
	"spec"
)

func responseCreate(response *spec.OperationResponse, resultVar string) string {
	if len(response.Operation.Responses) > 1 {
		return fmt.Sprintf(`return new %s.%s(%s);`, responseInterfaceName(response.Operation), response.Name.PascalCase(), resultVar)
	} else {
		if resultVar != "" {
			return fmt.Sprintf(`return %s;`, resultVar)
		} else {
			return `return;`
		}
	}
}

func (g *Generator) responseInterface(types *types.Types, operation *spec.NamedOperation) []generator.CodeFile {
	clientPackage := g.Packages.Client(operation.InApi)
	files := []generator.CodeFile{}
	w := writer.NewJavaWriter()
	w.Line(`package %s;`, clientPackage.PackageName)
	w.EmptyLine()
	w.Line(`import %s;`, g.Packages.Models(operation.InApi.InHttp.InVersion).PackageStar)
	w.Line(`import %s;`, g.Packages.ErrorsModels.PackageStar)
	w.EmptyLine()
	w.Line(`public interface %s {`, responseInterfaceName(operation))
	for index, response := range operation.Responses {
		if index > 0 {
			w.EmptyLine()
		}
		responseImpl(w.Indented(), types, &response)
	}
	w.Line(`}`)

	files = append(files, generator.CodeFile{
		Path:    clientPackage.GetPath(fmt.Sprintf("%s.java", responseInterfaceName(operation))),
		Content: w.String(),
	})
	return files
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
