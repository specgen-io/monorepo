package service

import (
	"fmt"
	"generator"
	"golang/writer"
	"spec"
)

func ContentType(operation *spec.NamedOperation) string {
	if operation.Body.IsEmpty() {
		return ""
	} else if operation.Body.IsText() {
		return "text/plain"
	} else if operation.Body.IsJson() {
		return "application/json"
	} else if operation.Body.IsBinary() {
		return "application/octet-stream"
	} else if operation.Body.IsBodyFormData() {
		return "multipart/form-data"
	} else if operation.Body.IsBodyFormUrlEncoded() {
		return "application/x-www-form-urlencoded"
	} else {
		panic(fmt.Sprintf("Unknown Contet Type"))
	}
}

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
