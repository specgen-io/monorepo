package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen/golang/module"
	"github.com/specgen-io/specgen/v2/generator"
	"strings"
)

func respondJson(logFields, resVar, statusCode, dataVar string) string {
	return fmt.Sprintf(`respondJson(%s, %s, %s, %s)`, logFields, resVar, statusCode, dataVar)
}

func respondText(logFields, resVar, statusCode, dataVar string) string {
	return fmt.Sprintf(`respondText(%s, %s, %s, %s)`, logFields, resVar, statusCode, dataVar)
}

func respondEmpty(logFields, resVar, statusCode string) string {
	return fmt.Sprintf(`respondEmpty(%s, %s, %s)`, logFields, resVar, statusCode)
}

func generateRespondFunctions(module module.Module) *generator.CodeFile {
	data := struct {
		PackageName string
	}{
		module.Name,
	}
	code := `
package [[.PackageName]]

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func respondJson(logFields log.Fields, res http.ResponseWriter, statusCode int, data interface{}) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	json.NewEncoder(res).Encode(data)
	log.WithFields(logFields).WithField("status", statusCode).Info("Completed request")
}

func respondText(logFields log.Fields, res http.ResponseWriter, statusCode int, data string) {
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(statusCode)
	res.Write([]byte(data))
	log.WithFields(logFields).WithField("status", statusCode).Info("Completed request")
}

func respondEmpty(logFields log.Fields, res http.ResponseWriter, statusCode int) {
	res.WriteHeader(statusCode)
	log.WithFields(logFields).WithField("status", statusCode).Info("Completed request")
}
`

	code, _ = generator.ExecuteTemplate(code, data)
	return &generator.CodeFile{module.GetPath("responses.go"), strings.TrimSpace(code)}
}
