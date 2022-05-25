package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen/golang/common"
	imports2 "github.com/specgen-io/specgen/v2/gen/golang/imports"
	"github.com/specgen-io/specgen/v2/gen/golang/module"
	"github.com/specgen-io/specgen/v2/gen/golang/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func generateServicesImplementations(version *spec.Version, versionModule, modelsModule, targetModule module.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		apiModule := versionModule.Submodule(api.Name.SnakeCase())
		files = append(files, *generateServiceImplementation(&api, apiModule, modelsModule, targetModule))
	}
	return files
}

func generateServiceImplementation(api *spec.Api, apiModule, modelsModule, targetModule module.Module) *sources.CodeFile {
	w := writer.NewGoWriter()
	w.Line("package %s", targetModule.Name)

	imports := imports2.Imports()
	imports.Add("errors")
	imports.AddApiTypes(api)
	imports.Add(apiModule.Package)
	if paramsContainsModel(api) {
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

	return &sources.CodeFile{
		Path:    targetModule.GetPath(fmt.Sprintf("%s.go", api.Name.SnakeCase())),
		Content: w.String(),
	}
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
