package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
)

func generateServicesInterfaces(version *spec.Version, packageName string, generatePath string, fileName string) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", packageName)

	imports := []string{}
	imports = generateVersionImports(version, imports)
	if len(imports) > 0 {
		w.EmptyLine()
		for _, imp := range imports {
			w.Line(`import %s`, imp)
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
		generateInterface(w, api, fileName)
	}

	return &gen.TextFile{
		Path:    filepath.Join(generatePath, fileName),
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

func generateInterface(w *gen.Writer, api spec.Api, fileName string) {
	var interfaceName string
	if isServiceOrClient(fileName) {
		interfaceName = serviceInterfaceTypeName(&api)
	} else {
		interfaceName = clientInterfaceTypeName(&api)
	}
	w.Line(`type %s interface {`, interfaceName)
	for _, operation := range api.Operations {
		w.Line(`  %s(%s) (*%s, error)`, operation.Name.PascalCase(), JoinDelimParams(addMethodParams(operation)), responseTypeName(&operation))
	}
	w.Line(`}`)
}

func isServiceOrClient(fileName string) bool {
	if fileName == "services.go" {
		return true
	} else {
		return false
	}
}
