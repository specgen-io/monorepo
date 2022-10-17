package client

import (
	"generator"
	"golang/imports"
	"golang/module"
	"golang/writer"
	"spec"
)

func (g *Generator) HttpErrors(errorsModule, errorsModelsModule module.Module, errors *spec.Responses) *generator.CodeFile {
	w := writer.New(errorsModule, "errors.go")

	imports := imports.New()
	imports.Add("fmt")
	imports.Module(errorsModelsModule)
	imports.Write(w)

	badRequestError := errors.GetByStatusName(spec.HttpStatusBadRequest)
	g.getError(w, badRequestError)
	notFoundError := errors.GetByStatusName(spec.HttpStatusNotFound)
	g.getError(w, notFoundError)
	internalServerError := errors.GetByStatusName(spec.HttpStatusInternalServerError)
	g.getError(w, internalServerError)

	return w.ToCodeFile()
}

func (g *Generator) getError(w generator.Writer, response *spec.Response) {
	w.EmptyLine()
	w.Line(`type %s struct {`, response.Name.PascalCase())
	w.Line(`	Body %s`, g.Types.GoType(&response.Type.Definition))
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`func (obj *%s) Error() string {`, response.Name.PascalCase())
	w.Line(`	return fmt.Sprintf("Body:  PERCENT_v", obj.Body)`)
	w.Line(`}`)
}
