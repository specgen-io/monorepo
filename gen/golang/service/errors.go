package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen/golang/types"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func generateErrors(w *sources.Writer, version *spec.Version) {
	badRequest := version.Http.Errors.Get(spec.HttpStatusBadRequest)
	w.Line(`respondBadRequest := func(logFields log.Fields, res http.ResponseWriter, error *%s) {`, types.GoType(&badRequest.Type.Definition))
	w.Line(`  log.WithFields(logFields).Warn(error.Message)`)
	generateResponseWriting(w.Indented(), `logFields`, badRequest, `error`)
	w.Line(`}`)
	w.Line(`_ = respondBadRequest`)

	notFound := version.Http.Errors.Get(spec.HttpStatusNotFound)
	w.Line(`respondNotFound := func(logFields log.Fields, res http.ResponseWriter, error *%s) {`, types.GoType(&notFound.Type.Definition))
	w.Line(`  log.WithFields(logFields).Warn(error.Message)`)
	generateResponseWriting(w.Indented(), `logFields`, notFound, `error`)
	w.Line(`}`)
	w.Line(`_ = respondNotFound`)

	internalServerError := version.Http.Errors.Get(spec.HttpStatusInternalServerError)
	w.EmptyLine()
	w.Line(`respondInternalServerError := func(logFields log.Fields, res http.ResponseWriter, error *%s) {`, types.GoType(&internalServerError.Type.Definition))
	w.Line(`  log.WithFields(logFields).Warn(error.Message)`)
	generateResponseWriting(w.Indented(), `logFields`, internalServerError, `error`)
	w.Line(`}`)
	w.Line(`_ = respondInternalServerError`)
}

func respondNotFound(w *sources.Writer, operation *spec.NamedOperation, message string) {
	badRequest := operation.Api.Apis.Errors.Get(spec.HttpStatusNotFound)
	error := fmt.Sprintf(`%s{Message: %s}`, types.GoType(&badRequest.Type.Definition), message)
	w.Line(`respondNotFound(%s, res, &%s)`, logFieldsName(operation), error)
	w.Line(`return`)
}

func respondBadRequest(w *sources.Writer, operation *spec.NamedOperation, location string, message string, params string) {
	badRequest := operation.Api.Apis.Errors.Get(spec.HttpStatusBadRequest)
	error := fmt.Sprintf(`%s{Location: "%s", Message: %s, Errors: %s}`, types.GoType(&badRequest.Type.Definition), location, message, params)
	w.Line(`respondBadRequest(%s, res, &%s)`, logFieldsName(operation), error)
	w.Line(`return`)
}

func respondInternalServerError(w *sources.Writer, operation *spec.NamedOperation, message string) {
	internalServerError := operation.Api.Apis.Errors.Get(spec.HttpStatusInternalServerError)
	error := fmt.Sprintf(`%s{Message: %s}`, types.GoType(&internalServerError.Type.Definition), message)
	w.Line(`respondInternalServerError(%s, res, &%s)`, logFieldsName(operation), error)
	w.Line(`return`)
}
