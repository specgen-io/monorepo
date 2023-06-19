package service

import (
	"fmt"
	"generator"
	"golang/types"
	"golang/writer"
	"spec"
)

func respondNotFound(w *writer.Writer, operation *spec.NamedOperation, types *types.Types, message string) {
	specification := operation.InApi.InHttp.InVersion.InSpec
	badRequest := specification.HttpErrors.Responses.GetByStatusName(spec.HttpStatusNotFound)
	errorMessage := fmt.Sprintf(`%s{Message: %s}`, types.GoType(&badRequest.ResponseBody.Type.Definition), message)
	w.Line(`httperrors.RespondNotFound(%s, res, &%s)`, logFieldsName(operation), errorMessage)
	w.Line(`return`)
}

func respondBadRequest(w *writer.Writer, operation *spec.NamedOperation, types *types.Types, location string, message string, params string) {
	specification := operation.InApi.InHttp.InVersion.InSpec
	badRequest := specification.HttpErrors.Responses.GetByStatusName(spec.HttpStatusBadRequest)
	errorMessage := fmt.Sprintf(`%s{Location: "%s", Message: %s, Errors: %s}`, types.GoType(&badRequest.ResponseBody.Type.Definition), location, message, params)
	w.Line(`httperrors.RespondBadRequest(%s, res, &%s)`, logFieldsName(operation), errorMessage)
	w.Line(`return`)
}

func respondInternalServerError(w *writer.Writer, operation *spec.NamedOperation, types *types.Types, message string) {
	specification := operation.InApi.InHttp.InVersion.InSpec
	internalServerError := specification.HttpErrors.Responses.GetByStatusName(spec.HttpStatusInternalServerError)
	errorMessage := fmt.Sprintf(`%s{Message: %s}`, types.GoType(&internalServerError.ResponseBody.Type.Definition), message)
	w.Line(`httperrors.RespondInternalServerError(%s, res, &%s)`, logFieldsName(operation), errorMessage)
	w.Line(`return`)
}

func (g *Generator) HttpErrors(responses *spec.ErrorResponses) []generator.CodeFile {
	files := []generator.CodeFile{}

	files = append(files, *g.errorsModelsConverter())
	files = append(files, *g.ErrorResponses(responses))

	return files
}

func (g *Generator) errorsModelsConverter() *generator.CodeFile {
	w := writer.New(g.Modules.HttpErrors, `converter.go`)
	w.Template(
		map[string]string{
			`ErrorsModelsPackage`: g.Modules.HttpErrorsModels.Package,
			`ParamsParserModule`:  g.Modules.ParamsParser.Package,
		}, `
import (
	"[[.ErrorsModelsPackage]]"
	"[[.ParamsParserModule]]"
)

func Convert(parsingErrors []paramsparser.ParsingError) []errmodels.ValidationError {
	var validationErrors []errmodels.ValidationError

	for _, parsingError := range parsingErrors {
		validationError := errmodels.ValidationError(parsingError)
		validationErrors = append(validationErrors, validationError)
	}

	return validationErrors
}
`)
	return w.ToCodeFile()
}

func (g *Generator) ErrorResponses(errors *spec.ErrorResponses) *generator.CodeFile {
	w := writer.New(g.Modules.HttpErrors, "responses.go")

	w.Imports.AddAliased("github.com/sirupsen/logrus", "log")
	w.Imports.Add("net/http")
	w.Imports.Module(g.Modules.HttpErrorsModels)
	w.Imports.Module(g.Modules.Respond)

	for _, response := range *errors {
		if response.BodyIs(spec.ResponseBodyEmpty) {
			w.Line(`func Respond%s(logFields log.Fields, res http.ResponseWriter) {`, response.Name.PascalCase())
			w.Line(`  log.WithFields(logFields).Warn("")`)
		} else {
			w.Line(`func Respond%s(logFields log.Fields, res http.ResponseWriter, error *%s) {`, response.Name.PascalCase(), g.Types.GoType(&response.ResponseBody.Type.Definition))
			w.Line(`  log.WithFields(logFields).Warn(error.Message)`)
		}
		writeResponse(w.Indented(), `logFields`, &response.Response, `error`)
		w.Line(`}`)
		w.EmptyLine()
	}

	return w.ToCodeFile()
}
