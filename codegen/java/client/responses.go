package client

import (
	"fmt"
	"generator"
	"java/types"
	"java/writer"
	"spec"
	"strconv"
)

func isSuccessfulStatusCode(statusCodeStr string) bool {
	statusCode, _ := strconv.Atoi(statusCodeStr)
	if statusCode >= 200 && statusCode <= 299 {
		return true
	}
	return false
}

func successfulResponsesNumber(operation *spec.NamedOperation) int {
	count := 0
	for _, response := range operation.Responses {
		if isSuccessfulStatusCode(spec.HttpStatusCode(response.Name)) {
			count++
		}
	}
	return count
}

func responseCreate(response *spec.OperationResponse, resultVar string) string {
	if successfulResponsesNumber(response.Operation) > 1 {
		return fmt.Sprintf(`return new %s.%s(%s);`, responseInterfaceName(response.Operation), response.Name.PascalCase(), resultVar)
	} else {
		if resultVar != "" {
			return fmt.Sprintf(`return %s;`, resultVar)
		} else {
			return `return;`
		}
	}
}

func (g *Generator) responseInterface(types *types.Types, operation *spec.NamedOperation) *generator.CodeFile {
	w := writer.New(g.Packages.Client(operation.InApi), responseInterfaceName(operation))
	w.Line(`import %s;`, g.Packages.Models(operation.InApi.InHttp.InVersion).PackageStar)
	w.Line(`import %s;`, g.Packages.ErrorsModels.PackageStar)
	w.EmptyLine()
	w.Line(`public interface [[.ClassName]] {`)
	for index, response := range operation.Responses {
		if index > 0 {
			w.EmptyLine()
		}
		responseImpl(w.Indented(), types, &response)
	}
	w.Line(`}`)
	return w.ToCodeFile()
}

func responseImpl(w *writer.Writer, types *types.Types, response *spec.OperationResponse) {
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
