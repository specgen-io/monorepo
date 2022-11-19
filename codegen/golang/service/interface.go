package service

import (
	"generator"
	"golang/walkers"
	"golang/writer"
	"spec"
)

func (g *Generator) ServicesInterfaces(version *spec.Version) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *g.serviceInterface(&api))
	}
	return files
}

func (g *Generator) serviceInterface(api *spec.Api) *generator.CodeFile {
	w := writer.New(g.Modules.ServicesApi(api), "service.go")

	w.Imports.AddApiTypes(api)
	if walkers.ApiHasMultiResponsesWithEmptyBody(api) {
		w.Imports.Module(g.Modules.Empty)
	}
	if walkers.ApiIsUsingModels(api) {
		w.Imports.Module(g.Modules.Models(api.InHttp.InVersion))
	}
	if walkers.ApiIsUsingErrorModels(api) {
		w.Imports.Module(g.Modules.HttpErrorsModels)
	}

	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			w.EmptyLine()
			Response(w, g.Types, &operation)
		}
	}
	w.EmptyLine()
	w.Line(`type %s interface {`, serviceInterfaceName)
	for _, operation := range api.Operations {
		w.Line(`  %s`, g.operationSignature(&operation, nil))
	}
	w.Line(`}`)

	return w.ToCodeFile()
}

const serviceInterfaceName = "Service"
