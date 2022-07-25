package service

import (
	"fmt"

	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/spec/v2"
	"github.com/specgen-io/specgen/v2/gen/golang/common"
	"github.com/specgen-io/specgen/v2/gen/golang/imports"
	"github.com/specgen-io/specgen/v2/gen/golang/module"
	"github.com/specgen-io/specgen/v2/gen/golang/types"
	"github.com/specgen-io/specgen/v2/gen/golang/writer"
)

func generateServiceImplementations(version *spec.Version, versionModule, modelsModule, targetModule module.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		apiModule := versionModule.Submodule(api.Name.SnakeCase())
		files = append(files, *generateServiceImplementation(&api, apiModule, modelsModule, targetModule))
	}
	return files
}

func generateServiceImplementation(api *spec.Api, apiModule, modelsModule, targetModule module.Module) *generator.CodeFile {
	w := writer.NewGoWriter()
	w.Line("package %s", targetModule.Name)

	imports := imports.New()
	imports.Add("errors")
	imports.AddApiTypes(api)
	if types.ApiHasBody(api) {
		imports.Add(apiModule.Package)
	}
	if isContainsModel(api) {
		imports.Add(modelsModule.Package)
	}
	imports.Write(w)

	w.EmptyLine()
	w.Line(`type %s struct{}`, serviceTypeName(api))
	w.EmptyLine()
	apiPackage := api.Name.SnakeCase()
	for _, operation := range api.Operations {
		w.Line(`func (service *%s) %s {`, serviceTypeName(api), common.OperationSignature(&operation, &apiPackage))
		singleEmptyResponse := len(operation.Responses) == 1 && operation.Responses[0].Type.Definition.IsEmpty()
		if singleEmptyResponse {
			w.Line(`  return errors.New("implementation has not added yet")`)
		} else {
			w.Line(`  return nil, errors.New("implementation has not added yet")`)
		}
		w.Line(`}`)
	}

	return &generator.CodeFile{
		Path:    targetModule.GetPath(fmt.Sprintf("%s.go", api.Name.SnakeCase())),
		Content: w.String(),
	}
}

func isContainsModel(api *spec.Api) bool {
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
		for _, response := range operation.Responses {
			if isModel(&response.Type.Definition) {
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
