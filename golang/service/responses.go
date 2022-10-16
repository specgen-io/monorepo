package service

import (
	"fmt"
	"generator"
	"golang/module"
	"golang/writer"
)

func respondJson(logFields, resVar, statusCode, dataVar string) string {
	return fmt.Sprintf(`respond.Json(%s, %s, %s, %s)`, logFields, resVar, statusCode, dataVar)
}

func respondText(logFields, resVar, statusCode, dataVar string) string {
	return fmt.Sprintf(`respond.Text(%s, %s, %s, %s)`, logFields, resVar, statusCode, dataVar)
}

func respondEmpty(logFields, resVar, statusCode string) string {
	return fmt.Sprintf(`respond.Empty(%s, %s, %s)`, logFields, resVar, statusCode)
}

func generateRespondFunctions(respondModule module.Module) *generator.CodeFile {
	w := writer.New(respondModule, `respond.go`)
	w.Lines(`
import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
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

func Empty(logFields log.Fields, res http.ResponseWriter, statusCode int) {
	res.WriteHeader(statusCode)
	log.WithFields(logFields).WithField("status", statusCode).Info("Completed request")
}
`)
	return w.ToCodeFile()
}
