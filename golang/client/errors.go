package client

import (
	"generator"
	"golang/imports"
	"golang/module"
	"golang/writer"
	"spec"
)

func (g *Generator) Errors(errorsModule, errorsModelsModule, responseModule module.Module, errors *spec.Responses) []generator.CodeFile {
	files := []generator.CodeFile{}

	files = append(files, *g.httpErrors(errorsModule, errorsModelsModule, errors))
	files = append(files, *g.httpErrorsHandler(errorsModule, errorsModelsModule, responseModule, errors))

	return files
}

func (g *Generator) httpErrors(errorsModule, errorsModelsModule module.Module, errors *spec.Responses) *generator.CodeFile {
	w := writer.New(errorsModule, "errors.go")

	imports := imports.New()
	imports.Add("fmt")
	imports.Module(errorsModelsModule)
	imports.Write(w)

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

func (g *Generator) httpErrorsHandler(module, errorsModelsModule, responseModule module.Module, errors *spec.Responses) *generator.CodeFile {
	w := writer.New(module, `errors_handler.go`)

	imports := imports.New().
		Module(errorsModelsModule).
		Module(responseModule).
		Add("net/http").
		AddAliased("github.com/sirupsen/logrus", "log")
	imports.Write(w)

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
