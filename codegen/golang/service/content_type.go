package service

import (
	"generator"
	"golang/writer"
)

func (g *Generator) CheckContentType() *generator.CodeFile {
	w := writer.New(g.Modules.ContentType, `check.go`)
	w.Template(
		map[string]string{
			`ErrorsPackage`:       g.Modules.HttpErrors.Package,
			`ErrorsModelsPackage`: g.Modules.HttpErrorsModels.Package,
		}, `
import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"[[.ErrorsPackage]]"
	"[[.ErrorsModelsPackage]]"
	"net/http"
	"strings"
)

func Check(logFields log.Fields, expectedContentType string, req *http.Request, res http.ResponseWriter) bool {
	contentType := req.Header.Get("Content-Type")
	if !strings.Contains(contentType, expectedContentType) {
		message := fmt.Sprintf("Expected Content-Type header: '%s' was not provided, found: '%s'", expectedContentType, contentType)
		httperrors.RespondBadRequest(logFields, res, &errmodels.BadRequestError{Location: "header", Message: "Failed to parse header", Errors: []errmodels.ValidationError{{Path: "Content-Type", Code: "missing", Message: &message}}})
		return false
	}
	return true
}
`)
	return w.ToCodeFile()
}
