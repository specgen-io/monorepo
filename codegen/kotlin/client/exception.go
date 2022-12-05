package client

import (
	"fmt"
	"generator"
	"kotlin/packages"
	"kotlin/types"
	"kotlin/writer"
	"spec"
)

func clientException(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `ClientException`)
	w.Lines(`
import java.lang.RuntimeException

open class [[.ClassName]] : RuntimeException {
	constructor() : super()
	constructor(message: String) : super(message)
	constructor(cause: Throwable) : super(cause)
	constructor(message: String, cause: Throwable) : super(message, cause)
}
`)
	return w.ToCodeFile()
}

func responseException(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `ResponseException`)
	w.Lines(`
import java.lang.RuntimeException

open class [[.ClassName]](message: String) : RuntimeException(message)
`)
	return w.ToCodeFile()
}

func inheritedClientException(thePackage, errorsModelsPackage packages.Package, types *types.Types, error *spec.Response) *generator.CodeFile {
	className := fmt.Sprintf(`%sException`, error.Name.PascalCase())
	w := writer.New(thePackage, className)
	w.Imports.PackageStar(errorsModelsPackage)
	errorBody := ""
	if !error.BodyIs(spec.BodyEmpty) {
		errorBody = fmt.Sprintf(`(val body: %s)`, error.Type.Definition)
	}
	w.Line(`class [[.ClassName]]%s : ResponseException("Error response with status code %s")`, errorBody, spec.HttpStatusCode(error.Name))

	return w.ToCodeFile()
}
