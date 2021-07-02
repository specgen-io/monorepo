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

	imports := []string{}
	if apiHasType(version, spec.TypeDate) {
		imports = append(imports, fmt.Sprintf(`import "cloud.google.com/go/civil"`))
	}
	if apiHasType(version, spec.TypeJson) {
		imports = append(imports, fmt.Sprintf(`import "encoding/json"`))
	}
	if apiHasType(version, spec.TypeUuid) {
		imports = append(imports, fmt.Sprintf(`import "github.com/google/uuid"`))
	}
	if apiHasType(version, spec.TypeDecimal) {
		imports = append(imports, fmt.Sprintf(`import "github.com/shopspring/decimal"`))
	}

	if len(imports) > 0 {
		w.EmptyLine()
		w.Line(`%s`, strings.Join(imports, "\n"))
	}

	w.EmptyLine()
	w.Line(`type EmptyDef struct{}`)
	w.EmptyLine()
	w.Line(`var Empty = EmptyDef{}`)

	for _, api := range version.Http.Apis {
		for _, operation := range api.Operations {
			w.EmptyLine()
			generateOperationResponseStruct(w, operation)
		}
		w.EmptyLine()
		generateInterface(w, api)
	}

	return &gen.TextFile{
		Path:    filepath.Join(generatePath, "services.go"),
		Content: w.String(),
	}
}

func addResponseParams(response spec.NamedResponse) []string {
	responseParams := []string{}
	//TODO: Check if it's possible to return EmptyDef from GoType - what other parts of codegen would be affected?
	responseType := GoType(&response.Type.Definition)
	if response.Type.Definition.IsEmpty() {
		responseType = "EmptyDef"
	}
	responseParams = append(responseParams, fmt.Sprintf(`%s *%s`, response.Name.PascalCase(), responseType))

	return responseParams
}

func generateOperationResponseStruct(w *gen.Writer, operation spec.NamedOperation) {
	//TODO: Make helper function responseTypeName(operation *spec.NamedOperation) string - use it everywhere instead of %sResponse
	w.Line(`type %sResponse struct {`, operation.Name.PascalCase())
	for _, response := range operation.Responses {
		w.Line(`  %s`, strings.Join(addResponseParams(response), "\n    "))
	}
	w.Line(`}`)
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

//TODO: Make helper function serviceInterfaceTypeName(api *spec.Api) string - use it everywhere instead of I%sService
func generateInterface(w *gen.Writer, api spec.Api) {
	w.Line(`type I%sService interface {`, api.Name.PascalCase())
	for _, operation := range api.Operations {
		w.Line(`  %s(%s) (*%sResponse, error)`, operation.Name.PascalCase(), strings.Join(addParams(operation), ", "), operation.Name.PascalCase())
	}
	w.Line(`}`)
}
