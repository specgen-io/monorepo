package service

import (
	"fmt"
	"golang/imports"
	"golang/module"
	"golang/writer"

	"generator"
	"golang/types"
	"spec"
)

func generateErrors(module, errorsModelsModule, respondModule module.Module, responses *spec.Responses) *generator.CodeFile {
	w := writer.NewGoWriter()
	w.Line("package %s", module.Name)

	imports := imports.New()
	imports.AddAlias("github.com/sirupsen/logrus", "log")
	imports.Add("net/http")
	imports.AddAlias(errorsModelsModule.Package, "errmodels")
	imports.Add(respondModule.Package)
	imports.Write(w)

	w.EmptyLine()
	badRequest := responses.GetByStatusName(spec.HttpStatusBadRequest)
	w.Line(`func RespondBadRequest(logFields log.Fields, res http.ResponseWriter, error *%s) {`, types.GoErrType(&badRequest.Type.Definition))
	w.Line(`  log.WithFields(logFields).Warn(error.Message)`)
	generateResponseWriting(w.Indented(), `logFields`, badRequest, `error`)
	w.Line(`}`)

	w.EmptyLine()
	notFound := responses.GetByStatusName(spec.HttpStatusNotFound)
	w.Line(`func RespondNotFound(logFields log.Fields, res http.ResponseWriter, error *%s) {`, types.GoErrType(&notFound.Type.Definition))
	w.Line(`  log.WithFields(logFields).Warn(error.Message)`)
	generateResponseWriting(w.Indented(), `logFields`, notFound, `error`)
	w.Line(`}`)

	w.EmptyLine()
	internalServerError := responses.GetByStatusName(spec.HttpStatusInternalServerError)
	w.Line(`func RespondInternalServerError(logFields log.Fields, res http.ResponseWriter, error *%s) {`, types.GoErrType(&internalServerError.Type.Definition))
	w.Line(`  log.WithFields(logFields).Warn(error.Message)`)
	generateResponseWriting(w.Indented(), `logFields`, internalServerError, `error`)
	w.Line(`}`)

	return &generator.CodeFile{
		Path:    module.GetPath(fmt.Sprintf("responses.go")),
		Content: w.String(),
	}
}

func callCheckContentType(logFieldsVar, expectedContentType, requestVar, responseVar string) string {
	return fmt.Sprintf(`contenttype.Check(%s, %s, %s, %s)`, logFieldsVar, expectedContentType, requestVar, responseVar)
}

func respondNotFound(w *generator.Writer, operation *spec.NamedOperation, message string) {
	badRequest := operation.Api.Http.Errors.GetByStatusName(spec.HttpStatusNotFound)
	errorMessage := fmt.Sprintf(`%s{Message: %s}`, types.GoErrType(&badRequest.Type.Definition), message)
	w.Line(`httperrors.RespondNotFound(%s, res, &%s)`, logFieldsName(operation), errorMessage)
	w.Line(`return`)
}

func respondBadRequest(w *generator.Writer, operation *spec.NamedOperation, location string, message string, params string) {
	badRequest := operation.Api.Http.Errors.GetByStatusName(spec.HttpStatusBadRequest)
	errorMessage := fmt.Sprintf(`%s{Location: "%s", Message: %s, Errors: %s}`, types.GoErrType(&badRequest.Type.Definition), location, message, params)
	w.Line(`httperrors.RespondBadRequest(%s, res, &%s)`, logFieldsName(operation), errorMessage)
	w.Line(`return`)
}

func respondInternalServerError(w *generator.Writer, operation *spec.NamedOperation, message string) {
	internalServerError := operation.Api.Http.Errors.GetByStatusName(spec.HttpStatusInternalServerError)
	errorMessage := fmt.Sprintf(`%s{Message: %s}`, types.GoErrType(&internalServerError.Type.Definition), message)
	w.Line(`httperrors.RespondInternalServerError(%s, res, &%s)`, logFieldsName(operation), errorMessage)
	w.Line(`return`)
}
