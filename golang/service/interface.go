package service

import (
	"generator"
	"golang/imports"
	"golang/module"
	"golang/responses"
	"golang/types"
	"golang/writer"
	"spec"
)

func generateServiceInterfaces(version *spec.Version, versionModule, modelsModule, errorsModelsModule, emptyModule module.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		apiModule := versionModule.Submodule(api.Name.SnakeCase())
		files = append(files, *generateServiceInterface(&api, apiModule, modelsModule, errorsModelsModule, emptyModule))
	}
	return files
}

func generateServiceInterface(api *spec.Api, apiModule, modelsModule, errorsModelsModule, emptyModule module.Module) *generator.CodeFile {
	w := writer.New(apiModule, "service.go")

	imports := imports.New()
	imports.AddApiTypes(api)
	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 && types.OperationHasType(&operation, spec.TypeEmpty) {
			imports.Add(emptyModule.Package)
		}
	}
	//TODO - potential bug, could be unused import
	imports.Add(modelsModule.Package)
	if usingErrorModels(api) {
		imports.AddAlias(errorsModelsModule.Package, types.ErrorsModelsPackage)
	}
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
		w.Line(`  %s`, OperationSignature(&operation, nil))
	}
	w.Line(`}`)

	return w.ToCodeFile()
}

const serviceInterfaceName = "Service"

func usingErrorModels(api *spec.Api) bool {
	foundErrorModels := false
	walk := spec.NewWalker().
		OnTypeDef(func(typ *spec.TypeDef) {
			if typ.Info.Model != nil && typ.Info.Model.InHttpErrors != nil {
				foundErrorModels = true
			}
		})
	walk.Api(api)
	return foundErrorModels
}
