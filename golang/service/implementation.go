package service

import (
	"fmt"
	"generator"
	"golang/imports"
	"golang/types"
	"golang/writer"
	"spec"
)

func (g *Generator) generateServiceImplementations(version *spec.Version) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *g.generateServiceImplementation(&api))
	}
	return files
}

func (g *Generator) generateServiceImplementation(api *spec.Api) *generator.CodeFile {
	w := writer.New(g.Modules.ServicesImpl(api), fmt.Sprintf("%s.go", api.Name.SnakeCase()))

	imports := imports.New()
	imports.Add("errors")
	imports.AddApiTypes(api)
	if types.ApiHasBody(api) {
		imports.Module(g.Modules.ServicesApi(api))
	}
	if isContainsModel(api) {
		imports.Module(g.Modules.Models(api.InHttp.InVersion))
	}
	imports.Write(w)

	w.EmptyLine()
	w.Line(`type %s struct{}`, serviceTypeName(api))
	w.EmptyLine()
	apiPackage := api.Name.SnakeCase()
	for _, operation := range api.Operations {
		w.Line(`func (service *%s) %s {`, serviceTypeName(api), g.operationSignature(&operation, &apiPackage))
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
