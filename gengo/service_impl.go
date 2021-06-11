package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
	"strings"
)

func generateServicesImplementations(version *spec.Version, packageName string, generatePath string) []gen.TextFile {
	versionPackageName := versionedPackage(version.Version, packageName)
	files := []gen.TextFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *generateServiceImplementation(&api, versionPackageName, generatePath))
	}
	return files
}

func generateServiceImplementation(api *spec.Api, packageName string, generatePath string) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", packageName)
	w.EmptyLine()
	addImportApi(w, api, spec.TypeDate, "cloud.google.com/go/civil")
	addImportApi(w, api, spec.TypeJson, "encoding/json")
	addImportApi(w, api, spec.TypeUuid, "github.com/google/uuid")
	addImportApi(w, api, spec.TypeDecimal, "github.com/shopspring/decimal")
	if strings.Contains(w.String(), "import") {
		w.EmptyLine()
	}
	generateTypeStruct(w, *api)
	generateFunction(w, *api)

	return &gen.TextFile{
		Path:    filepath.Join(generatePath, fmt.Sprintf("%s_service.go", api.Name.SnakeCase())),
		Content: w.String(),
	}
}

func generateTypeStruct(w *gen.Writer, api spec.Api) {
	w.Line(`type %sService struct{}`, api.Name.PascalCase())
	w.EmptyLine()
}

func generateOperationResponse(response spec.NamedResponse, operation spec.NamedOperation) string {
	if operation.Body != nil {
		return fmt.Sprintf(`%s: %s`, response.Name.PascalCase(), "body")
	}
	if response.Type.Definition.IsEmpty() {
		return fmt.Sprintf(`%s: &%s`, response.Name.PascalCase(), ToPascalCase(GoType(&response.Type.Definition)))
	}
	return fmt.Sprintf(`%s: &%s{%s}`, response.Name.PascalCase(), GoType(&response.Type.Definition), strings.Join(addOperationTypeParams(operation), ", "))
}

func addOperationTypeParams(operation spec.NamedOperation) []string {
	params := []string{}
	for _, param := range operation.QueryParams {
		params = append(params, fmt.Sprintf(`%s`, param.Name.CamelCase()))
	}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf(`%s`, param.Name.CamelCase()))
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, fmt.Sprintf(`%s`, param.Name.CamelCase()))

	}
	return params
}

func generateFunction(w *gen.Writer, api spec.Api) {
	for _, operation := range api.Operations {
		w.Line(`func (service *%sService) %s(%s) (*%sResponse, error) {`, api.Name.PascalCase(), operation.Name.PascalCase(), strings.Join(addParams(operation), ", "), operation.Name.PascalCase())
		for _, response := range operation.Responses {
			w.Line(`  return &%sResponse{%s}, nil`, operation.Name.PascalCase(), generateOperationResponse(response, operation))
		}
		w.Line(`}`)
		w.EmptyLine()
	}

}
