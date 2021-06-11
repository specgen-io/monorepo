package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
	"strings"
)

func generateServicesInterfaces(version *spec.Version, packageName string, generatePath string) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", packageName)
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
		Path:    filepath.Join(generatePath, "services.go"),
		Content: w.String(),
	}
}

func addImport(w *gen.Writer, version *spec.Version, typ string, importStr string) *gen.Writer {
	for _, api := range version.Http.Apis {
		addImportApi(w, &api, typ, importStr)
	}
	return nil
}

func addImportApi(w *gen.Writer, api *spec.Api, typ string, importStr string) *gen.Writer {
	for _, operation := range api.Operations {
		checkParam(w, operation.QueryParams, importStr, typ)
		checkParam(w, operation.HeaderParams, importStr, typ)
		checkParam(w, operation.Endpoint.UrlParams, importStr, typ)
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
	params := []string{}

	if operation.Body != nil {
		params = append(params, fmt.Sprintf("body *%s", GoType(&operation.Body.Type.Definition)))
	}
	for _, param := range operation.QueryParams {
		params = append(params, fmt.Sprintf("%s %s", param.Name.CamelCase(), GoType(&param.Type.Definition)))
	}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf("%s %s", param.Name.CamelCase(), GoType(&param.Type.Definition)))
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, fmt.Sprintf("%s %s", param.Name.CamelCase(), GoType(&param.Type.Definition)))
	}

	return params
}

func generateInterface(w *gen.Writer, api spec.Api) {
	w.Line(`type I%sService interface {`, api.Name.PascalCase())
	for _, operation := range api.Operations {
		w.Line(`  %s(%s) (*%sResponse, error)`, operation.Name.PascalCase(), strings.Join(addParams(operation), ", "), operation.Name.PascalCase())
	}
	w.Line(`}`)
	w.EmptyLine()
}
