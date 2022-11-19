package service

import (
	"fmt"
	"generator"
	"golang/walkers"
	"golang/writer"
	"spec"
)

func (g *Generator) ServicesImpls(version *spec.Version) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *g.serviceImpl(&api))
	}
	return files
}

func (g *Generator) serviceImpl(api *spec.Api) *generator.CodeFile {
	w := writer.New(g.Modules.ServicesImpl(api.InHttp.InVersion), fmt.Sprintf("%s.go", api.Name.SnakeCase()))

	w.Imports.Add("errors")
	w.Imports.AddApiTypes(api)
	if walkers.ApiHasNonSingleResponse(api) {
		w.Imports.Module(g.Modules.ServicesApi(api))
	}
	if walkers.ApiIsUsingModels(api) {
		w.Imports.Module(g.Modules.Models(api.InHttp.InVersion))
	}

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

func serviceTypeName(api *spec.Api) string {
	return fmt.Sprintf(`%sService`, api.Name.PascalCase())
}
