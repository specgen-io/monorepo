package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
	"strings"
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

	imports := []string{}
	imports = append(imports, `import "errors"`)
	if operationHasType(api, spec.TypeDate) {
		imports = append(imports, fmt.Sprintf(`import "%s"`, "cloud.google.com/go/civil"))
	}
	if operationHasType(api, spec.TypeJson) {
		imports = append(imports, fmt.Sprintf(`import "%s"`, "encoding/json"))
	}
	if operationHasType(api, spec.TypeUuid) {
		imports = append(imports, fmt.Sprintf(`import "%s"`, "github.com/google/uuid"))
	}
	if operationHasType(api, spec.TypeDecimal) {
		imports = append(imports, fmt.Sprintf(`import "%s"`, "github.com/shopspring/decimal"))
	}

	if len(imports) > 0 {
		w.EmptyLine()
		w.Line(`%s`, strings.Join(imports, "\n"))
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
	w.Line(`type %sService struct{}`, api.Name.PascalCase())
}

func generateFunction(w *gen.Writer, api spec.Api) {
	for _, operation := range api.Operations {
		w.Line(`func (service *%sService) %s(%s) (*%sResponse, error) {`, api.Name.PascalCase(), operation.Name.PascalCase(), strings.Join(addParams(operation), ", "), operation.Name.PascalCase())
		w.Line(`  return nil, errors.New("implementation has not added yet")`)
		w.Line(`}`)
	}
}
