package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
)

func generateServicesInterfaces(version *spec.Version, rootPackage, generatePath string) []gen.TextFile {
	files := []gen.TextFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *generateInterface(&api, rootPackage, generatePath))
	}
	return files
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

func generateInterface(api *spec.Api, rootPackage string, generatePath string) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", api.Name.SnakeCase())
	imports := []string{}
	imports = generateApiImports(api, imports)
	imports = append(imports, createPackageName(rootPackage, generatePath, modelsPackage))

	w.Line(`import (`)
	for _, imp := range imports {
		w.Line(`  %s`, imp)
	}
	w.Line(`)`)

	w.EmptyLine()
	w.Line(`type EmptyDef struct{}`)
	w.EmptyLine()
	w.Line(`var Empty = EmptyDef{}`)
	for _, operation := range api.Operations {
		w.EmptyLine()
		generateOperationResponseStruct(w, operation)
	}
	w.EmptyLine()
	w.Line(`type Service interface {`)
	for _, operation := range api.Operations {
		w.Line(`  %s(%s) (*%s, error)`, operation.Name.PascalCase(), JoinDelimParams(addMethodParams(operation)), responseTypeName(&operation))
	}
	w.Line(`}`)
	return &gen.TextFile{
		Path:    filepath.Join(generatePath, api.Name.SnakeCase(), "service.go"),
		Content: w.String(),
	}
}
