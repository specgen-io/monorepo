package client

import (
	"generator"
	"golang/imports"
	"golang/module"
	"golang/types"
	"golang/writer"
	"spec"
)

func httpErrors(module, errorsModelsModule module.Module, errors *spec.Responses) *generator.CodeFile {
	w := writer.New(module, "errors.go")

	imports := imports.New()
	imports.Add("fmt")
	imports.ModuleAliased(errorsModelsModule)
	imports.Write(w)

	badRequestError := errors.GetByStatusName(spec.HttpStatusBadRequest)
	getError(w, badRequestError)
	notFoundError := errors.GetByStatusName(spec.HttpStatusNotFound)
	getError(w, notFoundError)
	internalServerError := errors.GetByStatusName(spec.HttpStatusInternalServerError)
	getError(w, internalServerError)

	return w.ToCodeFile()
}

func getError(w generator.Writer, response *spec.Response) {
	w.EmptyLine()
	w.Line(`type %s struct {`, response.Name.PascalCase())
	w.Line(`	Body %s`, types.GoType(&response.Type.Definition))
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`func (obj *%s) Error() string {`, response.Name.PascalCase())
	w.Line(`	return fmt.Sprintf("Body:  PERCENT_v", obj.Body)`)
	w.Line(`}`)
}
