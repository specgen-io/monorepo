package service

import (
	"fmt"

	"generator"
	"golang/imports"
	"golang/module"
	"golang/types"
	"golang/writer"
	"spec"
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
	w := writer.New(targetModule, fmt.Sprintf("%s.go", api.Name.SnakeCase()))

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
		w.Line(`func (service *%s) %s {`, serviceTypeName(api), OperationSignature(&operation, &apiPackage))
		singleEmptyResponse := len(operation.Responses) == 1 && operation.Responses[0].Type.Definition.IsEmpty()
		if singleEmptyResponse {
			w.Line(`  return errors.New("implementation has not added yet")`)
		} else {
			w.Line(`  return nil, errors.New("implementation has not added yet")`)
		}
		w.Line(`}`)
	}

	return w.ToCodeFile()
}

func isContainsModel(api *spec.Api) bool {
	for _, operation := range api.Operations {
		if operation.Body != nil {
			if types.IsModel(&operation.Body.Type.Definition) {
				return true
			}
		}
		for _, param := range operation.QueryParams {
			if types.IsModel(&param.Type.Definition) {
				return true
			}
		}
		for _, param := range operation.HeaderParams {
			if types.IsModel(&param.Type.Definition) {
				return true
			}
		}
		for _, param := range operation.Endpoint.UrlParams {
			if types.IsModel(&param.Type.Definition) {
				return true
			}
		}
		for _, response := range operation.Responses {
			if types.IsModel(&response.Type.Definition) {
				return true
			}
		}
	}
	return false
}

func serviceTypeName(api *spec.Api) string {
	return fmt.Sprintf(`%sService`, api.Name.PascalCase())
}
