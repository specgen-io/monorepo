package service

import (
	"fmt"
	"generator"
	"golang/imports"
	"golang/writer"
	"spec"
)

func (g *VestigoGenerator) generateErrors(errors *spec.Responses) *generator.CodeFile {
	w := writer.New(g.Modules.HttpErrors, "responses.go")

	imports := imports.New()
	imports.AddAliased("github.com/sirupsen/logrus", "log")
	imports.Add("net/http")
	imports.Module(g.Modules.HttpErrorsModels)
	imports.Module(g.Modules.Respond)
	imports.Write(w)

	for _, errorResponse := range *errors {
		w.EmptyLine()
		w.Line(`func Respond%s(logFields log.Fields, res http.ResponseWriter, error *%s) {`, errorResponse.Name.PascalCase(), g.Types.GoType(&errorResponse.Type.Definition))
		w.Line(`  log.WithFields(logFields).Warn(error.Message)`)
		g.generateResponseWriting(w.Indented(), `logFields`, &errorResponse, `error`)
		w.Line(`}`)
	}

	return w.ToCodeFile()
}

func callCheckContentType(logFieldsVar, expectedContentType, requestVar, responseVar string) string {
	return fmt.Sprintf(`contenttype.Check(%s, %s, %s, %s)`, logFieldsVar, expectedContentType, requestVar, responseVar)
}

func (g *VestigoGenerator) respondNotFound(w generator.Writer, operation *spec.NamedOperation, message string) {
	specification := operation.InApi.InHttp.InVersion.InSpec
	badRequest := specification.HttpErrors.Responses.GetByStatusName(spec.HttpStatusNotFound)
	errorMessage := fmt.Sprintf(`%s{Message: %s}`, g.Types.GoType(&badRequest.Type.Definition), message)
	w.Line(`httperrors.RespondNotFound(%s, res, &%s)`, logFieldsName(operation), errorMessage)
	w.Line(`return`)
}

func (g *VestigoGenerator) respondBadRequest(w generator.Writer, operation *spec.NamedOperation, location string, message string, params string) {
	specification := operation.InApi.InHttp.InVersion.InSpec
	badRequest := specification.HttpErrors.Responses.GetByStatusName(spec.HttpStatusBadRequest)
	errorMessage := fmt.Sprintf(`%s{Location: "%s", Message: %s, Errors: %s}`, g.Types.GoType(&badRequest.Type.Definition), location, message, params)
	w.Line(`httperrors.RespondBadRequest(%s, res, &%s)`, logFieldsName(operation), errorMessage)
	w.Line(`return`)
}

func (g *VestigoGenerator) respondInternalServerError(w generator.Writer, operation *spec.NamedOperation, message string) {
	specification := operation.InApi.InHttp.InVersion.InSpec
	internalServerError := specification.HttpErrors.Responses.GetByStatusName(spec.HttpStatusInternalServerError)
	errorMessage := fmt.Sprintf(`%s{Message: %s}`, g.Types.GoType(&internalServerError.Type.Definition), message)
	w.Line(`httperrors.RespondInternalServerError(%s, res, &%s)`, logFieldsName(operation), errorMessage)
	w.Line(`return`)
}
