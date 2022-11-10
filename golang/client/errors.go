package client

import (
	"generator"
	"golang/writer"
	"spec"
)

func (g *Generator) Errors(errors *spec.Responses) []generator.CodeFile {
	files := []generator.CodeFile{}

	files = append(files, *g.httpErrors(errors))
	files = append(files, *g.httpErrorsHandler(errors))

	return files
}

func (g *Generator) httpErrors(errors *spec.Responses) *generator.CodeFile {
	w := writer.New(g.Modules.HttpErrors, "errors.go")

	w.Imports.Add("fmt")
	w.Imports.Module(g.Modules.HttpErrorsModels)

	for _, errorResponse := range *errors {
		w.EmptyLine()
		w.Line(`type %s struct {`, errorResponse.Name.PascalCase())
		w.Line(`	Body %s`, g.Types.GoType(&errorResponse.Type.Definition))
		w.Line(`}`)
		w.EmptyLine()
		w.Line(`func (obj *%s) Error() string {`, errorResponse.Name.PascalCase())
		w.Line(`	return fmt.Sprintf("Body:  PERCENT_v", obj.Body)`)
		w.Line(`}`)
	}

	return w.ToCodeFile()
}

func (g *Generator) httpErrorsHandler(errors *spec.Responses) *generator.CodeFile {
	w := writer.New(g.Modules.HttpErrors, `errors_handler.go`)

	w.Imports.Module(g.Modules.HttpErrorsModels)
	w.Imports.Module(g.Modules.Response)
	w.Imports.Add("net/http")
	w.Imports.AddAliased("github.com/sirupsen/logrus", "log")

	w.EmptyLine()
	w.Line(`func HandleErrors(resp *http.Response, log log.Fields) error {`)
	for _, errorResponse := range *errors {
		w.Line(`  if resp.StatusCode == %s {`, spec.HttpStatusCode(errorResponse.Name))
		w.Line(`    var result %s`, g.Types.GoType(&errorResponse.Type.Definition))
		w.Line(`    err := response.Json(log, resp, &result)`)
		w.Line(`    if err != nil {`)
		w.Line(`      return err`)
		w.Line(`    }`)
		w.Line(`    return &%s{Body: result}`, errorResponse.Name.PascalCase())
		w.Line(`  }`)
	}
	w.Line(`  return nil`)
	w.Line(`}`)

	return w.ToCodeFile()
}
