package service

import (
	"fmt"
	"generator"
	"golang/writer"
	"spec"
)

func responseTypeName(operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%sResponse`, operation.Name.PascalCase())
}

func respondJson(logFields, resVar, statusCode, dataVar string) string {
	return fmt.Sprintf(`respond.Json(%s, %s, %s, %s)`, logFields, resVar, statusCode, dataVar)
}

func respondText(logFields, resVar, statusCode, dataVar string) string {
	return fmt.Sprintf(`respond.Text(%s, %s, %s, %s)`, logFields, resVar, statusCode, dataVar)
}

func respondBinary(logFields, resVar, statusCode, dataVar string) string {
	return fmt.Sprintf(`respond.Binary(%s, %s, %s, %s)`, logFields, resVar, statusCode, dataVar)
}

func respondEmpty(logFields, resVar, statusCode string) string {
	return fmt.Sprintf(`respond.Empty(%s, %s, %s)`, logFields, resVar, statusCode)
}

func writeResponse(w *writer.Writer, logFieldsName string, response *spec.Response, responseVar string) {
	if response.Body.IsEmpty() {
		w.Line(respondEmpty(logFieldsName, `res`, spec.HttpStatusCode(response.Name)))
	}
	if response.Body.IsText() {
		w.Line(respondText(logFieldsName, `res`, spec.HttpStatusCode(response.Name), `*`+responseVar))
	}
	if response.Body.IsJson() {
		w.Line(respondJson(logFieldsName, `res`, spec.HttpStatusCode(response.Name), responseVar))
	}
	if response.Body.IsBinary() {
		w.Line(fmt.Sprintf(`err = %s`, respondBinary(logFieldsName, `res`, spec.HttpStatusCode(response.Name), responseVar)))
		w.Line(`if err != nil {`)
		w.Line(`  httperrors.RespondInternalServerError(%s, res, &errmodels.InternalServerError{Message: "Error sending data"})`, logFieldsName)
		w.Line(`  return`)
		w.Line(`}`)
	}
}

func (g *Generator) Response(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line(`type %s struct {`, responseTypeName(operation))
	w.Indent()
	for _, response := range operation.Responses {
		w.LineAligned(`%s *%s`, response.Name.PascalCase(), g.Types.ResponseBodyGoType(&response.Body))
	}
	w.Unindent()
	w.Line(`}`)
}

func (g *Generator) ResponseHelperFunctions() *generator.CodeFile {
	w := writer.New(g.Modules.Respond, `respond.go`)
	w.Lines(`
import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func Json(logFields log.Fields, res http.ResponseWriter, statusCode int, data interface{}) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	json.NewEncoder(res).Encode(data)
	log.WithFields(logFields).WithField("status", statusCode).Info("Completed request")
}

func Text(logFields log.Fields, res http.ResponseWriter, statusCode int, data string) {
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(statusCode)
	res.Write([]byte(data))
	log.WithFields(logFields).WithField("status", statusCode).Info("Completed request")
}

func Binary(logFields log.Fields, res http.ResponseWriter, statusCode int, data io.ReadCloser) error {
	res.Header().Set("Content-Type", "application/octet-stream")
	_, err := io.Copy(res, data)
	res.WriteHeader(statusCode)
	log.WithFields(logFields).WithField("status", statusCode).Info("Completed request")
	return err
}

func Empty(logFields log.Fields, res http.ResponseWriter, statusCode int) {
	res.WriteHeader(statusCode)
	log.WithFields(logFields).WithField("status", statusCode).Info("Completed request")
}
`)
	return w.ToCodeFile()
}
