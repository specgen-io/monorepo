package golang

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"path/filepath"
)

func generateServicesImplementations(moduleName string, versionModulePath string, version *spec.Version, generatePath string) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *generateServiceImplementation(moduleName, versionModulePath, &api, generatePath))
	}
	return files
}

func generateServiceImplementation(moduleName string, versionModulePath string, api *spec.Api, generatePath string) *sources.CodeFile {
	w := NewGoWriter()
	w.Line("package %s", getShortPackageName(generatePath))

	imports := Imports()
	imports.Add("errors")
	imports.AddApiTypes(api)
	imports.Add(createPackageName(moduleName, versionModulePath, api.Name.SnakeCase()))
	if paramsContainsModel(api) {
		imports.Add(createPackageName(moduleName, versionModulePath, modelsPackage))
	}
	imports.Write(w)

	w.EmptyLine()
	w.Line(`type %s struct{}`, serviceTypeName(api))
	w.EmptyLine()
	apiPackage := api.Name.SnakeCase()
	for _, operation := range api.Operations {
		w.Line(`func (service *%s) %s(%s) %s {`,
			serviceTypeName(api),
			operation.Name.PascalCase(),
			JoinDelimParams(addVersionedMethodParams(&operation)),
			operationReturn(&operation, &apiPackage),
		)
		singleEmptyResponse := len(operation.Responses) == 1 && operation.Responses[0].Type.Definition.IsEmpty()
		if singleEmptyResponse {
			w.Line(`  return errors.New("implementation has not added yet")`)
		} else {
			w.Line(`  return nil, errors.New("implementation has not added yet")`)
		}
		w.Line(`}`)
	}

	return &sources.CodeFile{
		Path:    filepath.Join(generatePath, fmt.Sprintf("%s.go", api.Name.SnakeCase())),
		Content: w.String(),
	}
}

func addVersionedMethodParams(operation *spec.NamedOperation) []string {
	params := []string{}
	if operation.BodyIs(spec.BodyString) {
		params = append(params, fmt.Sprintf("body %s", GoType(&operation.Body.Type.Definition)))
	}
	if operation.BodyIs(spec.BodyJson) {
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

func paramsContainsModel(api *spec.Api) bool {
	for _, operation := range api.Operations {
		if operation.Body != nil {
			if isModel(&operation.Body.Type.Definition) {
				return true
			}
		}
		for _, param := range operation.QueryParams {
			if isModel(&param.Type.Definition) {
				return true
			}
		}
		for _, param := range operation.HeaderParams {
			if isModel(&param.Type.Definition) {
				return true
			}
		}
		for _, param := range operation.Endpoint.UrlParams {
			if isModel(&param.Type.Definition) {
				return true
			}
		}
	}
	return false
}

func isModel(def *spec.TypeDef) bool {
	return def.Info.Model != nil
}

func serviceTypeName(api *spec.Api) string {
	return fmt.Sprintf(`%sService`, api.Name.PascalCase())
}
