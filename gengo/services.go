package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
	"strings"
)

func generateServices(version *spec.Version, packageName, generatePath, folder string) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", versionedPackage(version.Version, packageName))
	w.EmptyLine()
	addImport(w, version, spec.TypeDate, "cloud.google.com/go/civil")
	addImport(w, version, spec.TypeJson, "encoding/json")
	addImport(w, version, spec.TypeUuid, "github.com/google/uuid")
	addImport(w, version, spec.TypeDecimal, "github.com/shopspring/decimal")
	if strings.Contains(w.String(), "import") {
		w.EmptyLine()
	}
	w.Line(`type EmptyDef struct{}`)
	w.EmptyLine()
	w.Line(`var Empty = EmptyDef{}`)
	w.EmptyLine()

	for _, api := range version.Http.Apis {
		for _, operation := range api.Operations {
			generateOperationResponseStruct(w, operation)
		}
		generateInterface(w, api)
	}

	return &gen.TextFile{
		Path:    filepath.Join(generatePath, folder, "services.go"),
		Content: w.String(),
	}
}

func addImport(w *gen.Writer, version *spec.Version, typ string, importStr string) *gen.Writer {
	for _, api := range version.Http.Apis {
		for _, operation := range api.Operations {
			checkParam(w, operation.QueryParams, importStr, typ)
			checkParam(w, operation.HeaderParams, importStr, typ)
			checkParam(w, operation.Endpoint.UrlParams, importStr, typ)
		}
	}
	return nil
}

func checkParam(w *gen.Writer, namedParams []spec.NamedParam, importStr, typ string) *gen.Writer {
	if namedParams != nil && len(namedParams) > 0 {
		for _, param := range namedParams {
			if checkType(&param.Type.Definition, typ) {
				w.Line(`import "%s"`, importStr)
				return w
			}
		}
	}
	return nil
}

func addResponseParams(response spec.NamedResponse) []string {
	responseParams := []string{}
	responseType := GoType(&response.Type.Definition)
	if response.Type.Definition.IsEmpty() {
		responseType = "EmptyDef"
	}
	responseParams = append(responseParams, fmt.Sprintf(`%s *%s`, response.Name.PascalCase(), responseType))

	return responseParams
}

func generateOperationResponseStruct(w *gen.Writer, operation spec.NamedOperation) {
	w.Line(`type %sResponse struct {`, operation.Name.PascalCase())
	for _, response := range operation.Responses {
		w.Line(`  %s`, strings.Join(addResponseParams(response), "\n    "))
	}
	w.Line(`}`)
	w.EmptyLine()
}

func addParams(operation spec.NamedOperation) []string {
	urlParams := []string{}

	if operation.Body != nil {
		urlParams = append(urlParams, fmt.Sprintf("body *%s", GoType(&operation.Body.Type.Definition)))
	}
	for _, param := range operation.QueryParams {
		urlParams = append(urlParams, fmt.Sprintf("%s %s", param.Name.CamelCase(), GoType(&param.Type.Definition)))
	}
	for _, param := range operation.HeaderParams {
		urlParams = append(urlParams, fmt.Sprintf("%s %s", param.Name.CamelCase(), GoType(&param.Type.Definition)))
	}
	for _, param := range operation.Endpoint.UrlParams {
		urlParams = append(urlParams, fmt.Sprintf("%s %s", param.Name.CamelCase(), GoType(&param.Type.Definition)))
	}

	return urlParams
}

func generateInterface(w *gen.Writer, api spec.Api) {
	w.Line(`type I%sService interface {`, api.Name.PascalCase())
	for _, operation := range api.Operations {
		w.Line(`  %s(%s) (*%sResponse, error)`, operation.Name.PascalCase(), strings.Join(addParams(operation), ", "), operation.Name.PascalCase())
	}
	w.Line(`}`)
	w.EmptyLine()
}
