package client

import (
	"generator"
	"golang/writer"
	"spec"
)

func (g *Generator) Errors(errors *spec.ErrorResponses) *generator.CodeFile {
	w := writer.New(g.Modules.HttpErrors, "errors.go")

	w.Imports.Add("fmt")
	w.Imports.Module(g.Modules.HttpErrorsModels)

	for _, response := range *errors {
		w.EmptyLine()
		w.Line(`type %s struct {`, response.Name.PascalCase())
		if !response.BodyIs(spec.ResponseBodyEmpty) {
			w.Line(`	Body %s`, g.Types.GoType(&response.ResponseBody.Type.Definition))
		}
		w.Line(`}`)
		w.EmptyLine()
		w.Line(`func (obj *%s) Error() string {`, response.Name.PascalCase())
		if response.BodyIs(spec.ResponseBodyEmpty) {
			w.Line(`	return "%s"`, response.Name.PascalCase())
		} else {
			w.Line(`	return fmt.Sprintf("%s - Body:  PERCENT_v", obj.Body)`, response.Name.PascalCase())
		}
		w.Line(`}`)
	}

	return w.ToCodeFile()
}
