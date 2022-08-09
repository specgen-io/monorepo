package service

import (
	"generator"
	"golang/common"
	"golang/imports"
	"golang/module"
	"golang/responses"
	"golang/types"
	"golang/writer"
	"spec"
)

func generateServiceInterfaces(version *spec.Version, versionModule, modelsModule, emptyModule module.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		apiModule := versionModule.Submodule(api.Name.SnakeCase())
		files = append(files, *generateServiceInterface(&api, apiModule, modelsModule, emptyModule))
	}
	return files
}

func generateServiceInterface(api *spec.Api, apiModule, modelsModule, emptyModule module.Module) *generator.CodeFile {
	w := writer.NewGoWriter()
	w.Line("package %s", apiModule.Name)

	imports := imports.New()
	imports.AddApiTypes(api)
	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 && types.OperationHasType(&operation, spec.TypeEmpty) {
			imports.Add(emptyModule.Package)
		}
	}
	//TODO - potential bug, could be unused import
	imports.Add(modelsModule.Package)
	imports.Write(w)

	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			w.EmptyLine()
			responses.GenerateOperationResponseStruct(w, &operation)
		}
	}
	w.EmptyLine()
	w.Line(`type %s interface {`, serviceInterfaceName)
	for _, operation := range api.Operations {
		w.Line(`  %s`, common.OperationSignature(&operation, nil))
	}
	w.Line(`}`)
	return &generator.CodeFile{
		Path:    apiModule.GetPath("service.go"),
		Content: w.String(),
	}
}

const serviceInterfaceName = "Service"
