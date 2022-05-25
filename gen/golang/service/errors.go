package service

import (
	"github.com/specgen-io/specgen/v2/gen/golang/imports"
	"github.com/specgen-io/specgen/v2/gen/golang/module"
	"github.com/specgen-io/specgen/v2/gen/golang/types"
	"github.com/specgen-io/specgen/v2/gen/golang/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func generateErrors(version *spec.Version, versionModule module.Module, modelsModule module.Module) *sources.CodeFile {
	w := writer.NewGoWriter()
	w.Line("package %s", versionModule.Name)

	imports := imports.Imports()
	imports.Add("encoding/json")
	imports.AddAlias("github.com/sirupsen/logrus", "log")
	imports.Add("net/http")
	imports.Add("fmt")
	imports.Add("strings")
	imports.Add(modelsModule.Package)
	imports.Write(w)

	badRequest := version.Http.Errors.Get(spec.HttpStatusBadRequest)
	w.EmptyLine()
	w.Line(`func BadRequest(logFields log.Fields, res http.ResponseWriter, error *%s) {`, types.GoType(&badRequest.Type.Definition))
	w.Line(`  log.WithFields(logFields).Warn(error.Message)`)
	generateResponseWriting(w.Indented(), `logFields`, badRequest, `error`)
	w.Line(`}`)

	internalServerError := version.Http.Errors.Get(spec.HttpStatusInternalServerError)
	w.EmptyLine()
	w.Line(`func InternalServerError(logFields log.Fields, res http.ResponseWriter, error *%s) {`, types.GoType(&internalServerError.Type.Definition))
	w.Line(`  log.WithFields(logFields).Warn(error.Message)`)
	generateResponseWriting(w.Indented(), `logFields`, internalServerError, `error`)
	w.Line(`}`)

	w.EmptyLine()
	w.Line(`func CheckContentType(logFields log.Fields, res http.ResponseWriter, req *http.Request, expectedContentType string) bool {`)
	w.Line(`  contentType := req.Header.Get("Content-Type")`)
	w.Line(`  if !strings.Contains(contentType, expectedContentType) {`)
	w.Line(`    error := models.BadRequestError{%s, nil}`, genFmtSprintf("Wrong Content-type: %s", "contentType"))
	w.Line(`    BadRequest(logFields, res, &error)`)
	w.Line(`    return false`)
	w.Line(`  }`)
	w.Line(`  return true`)
	w.Line(`}`)

	return &sources.CodeFile{
		Path:    versionModule.GetPath("errors.go"),
		Content: w.String(),
	}
}

func generateBadRequestResponse(w *sources.Writer, operation *spec.NamedOperation, message string, params string) {
	badRequest := operation.Api.Apis.Errors.Get(spec.HttpStatusBadRequest)
	w.Line(`error := %s{Message: %s, Params: %s}`, types.GoType(&badRequest.Type.Definition), message, params)
	w.Line(`BadRequest(%s, res, &error)`, logFieldsName(operation))
	w.Line(`return`)
}

func generateInternalServerErrorResponse(w *sources.Writer, operation *spec.NamedOperation, message string) {
	internalServerError := operation.Api.Apis.Errors.Get(spec.HttpStatusInternalServerError)
	w.Line(`error := %s{Message: %s}`, types.GoType(&internalServerError.Type.Definition), message)
	w.Line(`InternalServerError(%s, res, &error)`, logFieldsName(operation))
	w.Line(`return`)
}

func checkRequestContentType(w *sources.Writer, operation *spec.NamedOperation, contentType string) {
	w.Line(`if !CheckContentType(%s, res, req, "%s") {`, logFieldsName(operation), contentType)
	w.Line(`  return`)
	w.Line(`}`)
}
