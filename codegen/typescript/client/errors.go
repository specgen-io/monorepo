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
	w.Imports.Star(g.Modules.Validation, "t")
	w.Imports.LibNames("axios", "AxiosResponse")
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
		if response.BodyIs(spec.BodyEmpty) {
			w.Line(`  constructor() {`)
			w.Line(`    super('Error response with status code %s')`, spec.HttpStatusCode(response.Name))
			w.Line(`  }`)
		} else {
			w.Line(`  public body: %s`, types.TsType(&response.Type.Definition))
			w.EmptyLine()
			w.Line(`  constructor(body: %s) {`, types.TsType(&response.Type.Definition))
			w.Line(`    super('Error response with status code %s')`, spec.HttpStatusCode(response.Name))
			w.Line(`    this.body = body`)
			w.Line(`  }`)
		}
		w.Line(`}`)
	}
	w.EmptyLine()
	w.Line(`export const checkRequiredErrors = (response: AxiosResponse<any>): AxiosResponse<any> => {`)
	w.Line(`  switch (response.status) {`)
	for _, response := range httpErrors.Responses.Required() {
		w.Line(`    case %s:`, spec.HttpStatusCode(response.Name))
		w.Line(`      throw new %s(%s)`, errorExceptionName(&response.Response), g.responseBody(&response.Response))
	}
	w.Line(`  }`)
	w.Line(`  return response`)
	w.Line(`}`)
	return w.ToCodeFile()
}
