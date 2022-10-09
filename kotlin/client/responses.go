package client

import (
	"fmt"
	"generator"
	"kotlin/packages"
	"kotlin/types"
	"kotlin/writer"
	"spec"
)

func responseCreate(response *spec.OperationResponse, resultVar string) string {
	if len(response.Operation.Responses) > 1 {
		return fmt.Sprintf(`%s.%s(%s)`, reponseInterfaceName(response.Operation), response.Name.PascalCase(), resultVar)
	}
	return resultVar
}

func responses(api *spec.Api, types *types.Types, apiPackage packages.Package, modelsVersionPackage packages.Package, errorModelsPackage packages.Package) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			files = append(files, *responseInterface(types, &operation, apiPackage, modelsVersionPackage, errorModelsPackage))
		}
	}
	return files
}

func responseInterface(types *types.Types, operation *spec.NamedOperation, apiPackage packages.Package, modelsVersionPackage packages.Package, errorModelsPackage packages.Package) *generator.CodeFile {
	w := writer.New(apiPackage, reponseInterfaceName(operation))
	w.Line(`import %s`, modelsVersionPackage.PackageStar)
	w.Line(`import %s`, errorModelsPackage.PackageStar)
	w.EmptyLine()
	w.Line(`interface [[.ClassName]] {`)
	for index, response := range operation.Responses {
		if index > 0 {
			w.EmptyLine()
		}
		responseImpl(w.Indented(), types, &response)
	}
	w.Line(`}`)
	return w.ToCodeFile()
}

func responseImpl(w generator.Writer, types *types.Types, response *spec.OperationResponse) {
	serviceResponseImplementationName := response.Name.PascalCase()
	if !response.Type.Definition.IsEmpty() {
		w.Line(`class %s(var body: %s) : %s`, serviceResponseImplementationName, types.Kotlin(&response.Type.Definition), reponseInterfaceName(response.Operation))
	} else {
		w.Line(`class %s : %s`, serviceResponseImplementationName, reponseInterfaceName(response.Operation))
	}
}

func reponseInterfaceName(operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%sResponse`, operation.Name.PascalCase())
}
