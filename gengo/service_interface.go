package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
)

func generateServicesInterfaces(version *spec.Version, packageName string, generatePath string) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", packageName)

	imports := []string{}
	if versionHasType(version, spec.TypeDate) {
		imports = append(imports, `import "cloud.google.com/go/civil"`)
	}
	if versionHasType(version, spec.TypeJson) {
		imports = append(imports, `import "encoding/json"`)
	}
	if versionHasType(version, spec.TypeUuid) {
		imports = append(imports, `import "github.com/google/uuid"`)
	}
	if versionHasType(version, spec.TypeDecimal) {
		imports = append(imports, `import "github.com/shopspring/decimal"`)
	}

	if len(imports) > 0 {
		w.EmptyLine()
		for _, imp := range imports {
			w.Line(imp)
		}
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

func generateOperationResponseStruct(w *gen.Writer, operation spec.NamedOperation) {
	w.Line(`type %s struct {`, responseTypeName(&operation))
	for _, response := range operation.Responses {
		w.Line(`  %s *%s`, response.Name.PascalCase(), GoType(&response.Type.Definition))
	}
	w.Line(`}`)
}

func addMethodParams(operation spec.NamedOperation) []string {
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
	w.Line(`type %s interface {`, serviceInterfaceTypeName(&api))
	for _, operation := range api.Operations {
		w.Line(`  %s(%s) (*%s, error)`, operation.Name.PascalCase(), JoinDelimParams(addMethodParams(operation)), responseTypeName(&operation))
	}
	w.Line(`}`)
}
