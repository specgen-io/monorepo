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
	if walkers.ApiHasType(api, spec.TypeDate) {
		w.Imports.Add("cloud.google.com/go/civil")
	}
	if walkers.ApiHasType(api, spec.TypeJson) {
		w.Imports.Add("encoding/json")
	}
	if walkers.ApiHasType(api, spec.TypeUuid) {
		w.Imports.Add("github.com/google/uuid")
	}
	if walkers.ApiHasType(api, spec.TypeDecimal) {
		w.Imports.Add("github.com/shopspring/decimal")
	}
	if walkers.ApiHasBodyOfKind(api, spec.BodyFile) {
		w.Imports.Module(g.Modules.HttpFile)
	}
	if walkers.ApiHasNonSingleResponse(api) {
		w.Imports.Module(g.Modules.ServicesApi(api))
	}
	if walkers.ApiIsUsingModels(api) {
		w.Imports.Module(g.Modules.Models(api.InHttp.InVersion))
	}

	w.Line(`type %s struct{}`, serviceTypeName(api))
	w.EmptyLine()
	apiPackage := api.Name.SnakeCase()
	for _, operation := range api.Operations {
		w.Line(`func (service %s) %s {`, serviceTypeName(api), g.operationSignature(&operation, &apiPackage))
		singleEmptyResponse := len(operation.Responses) == 1 && operation.Responses[0].Body.IsEmpty()
		if singleEmptyResponse {
			w.Line(`  return errors.New("implementation has not added yet")`)
		} else {
			w.Line(`  return nil, errors.New("implementation has not added yet")`)
		}
		w.Line(`}`)
		w.EmptyLine()
	}

	return w.ToCodeFile()
}

func serviceTypeName(api *spec.Api) string {
	return fmt.Sprintf(`%sService`, api.Name.PascalCase())
}
