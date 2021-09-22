package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
)

func generateServicesImplementations(moduleName string, importedPackage string, version *spec.Version, packageName string, generatePath string) []gen.TextFile {
	files := []gen.TextFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *generateServiceImplementation(moduleName, importedPackage, version, &api, packageName, generatePath))
	}
	return files
}

func generateServiceImplementation(moduleName string, importedPackage string, version *spec.Version, api *spec.Api, packageName string, generatePath string) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", packageName)

	imports := Imports()
	imports.Add("errors")
	imports.AddApiTypes(api)
	imports.Add(fmt.Sprintf(`%s/%s`, moduleName, versionedFolder(version.Version, importedPackage)))
	imports.Write(w)

	w.EmptyLine()
	generateTypeStruct(w, *api)
	w.EmptyLine()
	generateFunction(w, *api, version, importedPackage)

	return &gen.TextFile{
		Path:    filepath.Join(generatePath, fmt.Sprintf("%s_service.go", api.Name.SnakeCase())),
		Content: w.String(),
	}
}

func generateTypeStruct(w *gen.Writer, api spec.Api) {
	w.Line(`type %s struct{}`, serviceTypeName(&api))
}

func generateFunction(w *gen.Writer, api spec.Api, version *spec.Version, packageName string) {
	for _, operation := range api.Operations {
		w.Line(`func (service *%s) %s(%s) (*%s.%s, error) {`, serviceTypeName(&api), operation.Name.PascalCase(), JoinDelimParams(addVersionedMethodParams(operation, version, packageName)), versionedPackage(version.Version, packageName), responseTypeName(&operation))
		w.Line(`  return nil, errors.New("implementation has not added yet")`)
		w.Line(`}`)
	}
}

func addVersionedMethodParams(operation spec.NamedOperation, version *spec.Version, packageName string) []string {
	params := []string{}

	if operation.Body != nil {
		params = append(params, fmt.Sprintf("body *%s", versionedType(&operation.Body.Type.Definition, version, packageName)))
	}
	for _, param := range operation.QueryParams {
		params = append(params, fmt.Sprintf("%s %s", param.Name.CamelCase(), versionedType(&param.Type.Definition, version, packageName)))
	}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf("%s %s", param.Name.CamelCase(), versionedType(&param.Type.Definition, version, packageName)))
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, fmt.Sprintf("%s %s", param.Name.CamelCase(), versionedType(&param.Type.Definition, version, packageName)))
	}
	return params
}

func versionedType(def *spec.TypeDef, version *spec.Version, packageName string) string {
	typ := GoType(def)
	if def.Info.Model != nil {
		typ = fmt.Sprintf(`%s.%s`, versionedPackage(version.Version, packageName), typ)
	}
	return typ
}

func serviceTypeName(api *spec.Api) string {
	return fmt.Sprintf(`%sService`, api.Name.PascalCase())
}