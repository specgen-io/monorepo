package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen/golang/types"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func generateErrors(w *sources.Writer, version *spec.Version) {
	badRequest := version.Http.Errors.GetByStatusName(spec.HttpStatusBadRequest)
	w.Line(`respondBadRequest := func(logFields log.Fields, res http.ResponseWriter, error *%s) {`, types.GoType(&badRequest.Type.Definition))
	w.Line(`  log.WithFields(logFields).Warn(error.Message)`)
	generateResponseWriting(w.Indented(), `logFields`, badRequest, `error`)
	w.Line(`}`)
	w.Line(`_ = respondBadRequest`)

	notFound := version.Http.Errors.GetByStatusName(spec.HttpStatusNotFound)
	w.Line(`respondNotFound := func(logFields log.Fields, res http.ResponseWriter, error *%s) {`, types.GoType(&notFound.Type.Definition))
	w.Line(`  log.WithFields(logFields).Warn(error.Message)`)
	generateResponseWriting(w.Indented(), `logFields`, notFound, `error`)
	w.Line(`}`)
	w.Line(`_ = respondNotFound`)

	internalServerError := version.Http.Errors.GetByStatusName(spec.HttpStatusInternalServerError)
	w.EmptyLine()
	w.Line(`respondInternalServerError := func(logFields log.Fields, res http.ResponseWriter, error *%s) {`, types.GoType(&internalServerError.Type.Definition))
	w.Line(`  log.WithFields(logFields).Warn(error.Message)`)
	generateResponseWriting(w.Indented(), `logFields`, internalServerError, `error`)
	w.Line(`}`)
	w.Line(`_ = respondInternalServerError`)
}

func generateCheckContentType(w *sources.Writer) {
	w.Lines(`
checkContentType := func(logFields log.Fields, expectedContentType string, req *http.Request, res http.ResponseWriter) bool {
	contentType := req.Header.Get("Content-Type")
	if !strings.Contains(contentType, expectedContentType) {
		message := fmt.Sprintf("Expected Content-Type header: '%s' was not provided, found: '%s'", expectedContentType, contentType)
		respondBadRequest(logFields, res, &models.BadRequestError{Location: "header", Message: "Failed to parse header", Errors: []models.ValidationError{{Path: "Content-Type", Code: "missing", Message: &message}}})
		return false
	}
	return true
}
_ = checkContentType
`)
}

func callCheckContentType(logFieldsVar, expectedContentType, requestVar, responseVar string) string {
	return fmt.Sprintf(`checkContentType(%s, %s, %s, %s)`, logFieldsVar, expectedContentType, requestVar, responseVar)
}

func respondNotFound(w *sources.Writer, operation *spec.NamedOperation, message string) {
	badRequest := operation.Api.Http.Errors.GetByStatusName(spec.HttpStatusNotFound)
	error := fmt.Sprintf(`%s{Message: %s}`, types.GoType(&badRequest.Type.Definition), message)
	w.Line(`respondNotFound(%s, res, &%s)`, logFieldsName(operation), error)
	w.Line(`return`)
}

func respondBadRequest(w *sources.Writer, operation *spec.NamedOperation, location string, message string, params string) {
	badRequest := operation.Api.Http.Errors.GetByStatusName(spec.HttpStatusBadRequest)
	error := fmt.Sprintf(`%s{Location: "%s", Message: %s, Errors: %s}`, types.GoType(&badRequest.Type.Definition), location, message, params)
	w.Line(`respondBadRequest(%s, res, &%s)`, logFieldsName(operation), error)
	w.Line(`return`)
}

func respondInternalServerError(w *sources.Writer, operation *spec.NamedOperation, message string) {
	internalServerError := operation.Api.Http.Errors.GetByStatusName(spec.HttpStatusInternalServerError)
	error := fmt.Sprintf(`%s{Message: %s}`, types.GoType(&internalServerError.Type.Definition), message)
	w.Line(`respondInternalServerError(%s, res, &%s)`, logFieldsName(operation), error)
	w.Line(`return`)
}
