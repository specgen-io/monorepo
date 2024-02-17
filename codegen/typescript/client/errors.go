package client

import (
	"fmt"
	"generator"
	"spec"
	"typescript/types"
	"typescript/writer"
)

func errorExceptionName(response *spec.Response) string {
	return fmt.Sprintf(`%sException`, response.Name.PascalCase())
}

func (g *CommonGenerator) Errors(httpErrors *spec.HttpErrors) *generator.CodeFile {
	w := writer.New(g.Modules.Errors)
	w.Imports.Star(g.Modules.ErrorsModels, "errors")
	w.Line(`export * from '%s'`, g.Modules.ErrorsModels.GetImport(g.Modules.Errors))
	w.EmptyLine()
	w.Lines(`
export class ResponseException extends Error {
  constructor(message: string) {
    super(message)
  }
}
`)
	for _, response := range httpErrors.Responses {
		w.EmptyLine()
		w.Line(`export class %s extends ResponseException {`, errorExceptionName(&response.Response))
		if response.Body.IsEmpty() {
			w.Line(`  constructor() {`)
			w.Line(`    super('Error response with status code %s')`, spec.HttpStatusCode(response.Name))
			w.Line(`  }`)
		} else {
			w.Line(`  public body: %s`, types.ResponseBodyTsType(&response.Body))
			w.EmptyLine()
			w.Line(`  constructor(body: %s) {`, types.ResponseBodyTsType(&response.Body))
			w.Line(`    super('Error response with status code %s')`, spec.HttpStatusCode(response.Name))
			w.Line(`    this.body = body`)
			w.Line(`  }`)
		}
		w.Line(`}`)
	}
	return w.ToCodeFile()
}
