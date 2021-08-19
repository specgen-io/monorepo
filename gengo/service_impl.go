package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
)

func generateServicesImplementations(version *spec.Version, packageName string, generatePath string) []gen.TextFile {
	files := []gen.TextFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *generateServiceImplementation(&api, packageName, generatePath))
	}
	return files
}

func generateServiceImplementation(api *spec.Api, packageName string, generatePath string) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", packageName)

	imports := []string{`import "errors"`}
	if apiHasType(api, spec.TypeDate) {
		imports = append(imports, `import "cloud.google.com/go/civil"`)
	}
	if apiHasType(api, spec.TypeJson) {
		imports = append(imports, `import "encoding/json"`)
	}
	if apiHasType(api, spec.TypeUuid) {
		imports = append(imports, `import "github.com/google/uuid"`)
	}
	if apiHasType(api, spec.TypeDecimal) {
		imports = append(imports, `import "github.com/shopspring/decimal"`)
	}

	w.EmptyLine()
	for _, imp := range imports {
		w.Line(imp)
	}

	w.EmptyLine()
	generateTypeStruct(w, *api)
	w.EmptyLine()
	generateFunction(w, *api)

	return &gen.TextFile{
		Path:    filepath.Join(generatePath, fmt.Sprintf("%s_service.go", api.Name.SnakeCase())),
		Content: w.String(),
	}
}

func generateTypeStruct(w *gen.Writer, api spec.Api) {
	w.Line(`type %s struct{}`, serviceTypeName(&api))
}

func generateFunction(w *gen.Writer, api spec.Api) {
	for _, operation := range api.Operations {
		w.Line(`func (service *%s) %s(%s) (*%s, error) {`, serviceTypeName(&api), operation.Name.PascalCase(), JoinDelimParams(addMethodParams(operation)), responseTypeName(&operation))
		w.Line(`  return nil, errors.New("implementation has not added yet")`)
		w.Line(`}`)
	}
}
