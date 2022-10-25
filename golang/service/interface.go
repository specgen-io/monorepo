package service

import (
	"generator"
	"golang/imports"
	"golang/types"
	"golang/writer"
	"spec"
)

func (g *Generator) generateServiceInterfaces(version *spec.Version) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *g.generateServiceInterface(&api))
	}
	return files
}

func (g *Generator) generateServiceInterface(api *spec.Api) *generator.CodeFile {
	w := writer.New(g.Modules.ServicesApi(api), "service.go")

	imports := imports.New()
	imports.AddApiTypes(api)
	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 && types.OperationHasType(&operation, spec.TypeEmpty) {
			imports.Module(g.Modules.Empty)
		}
	}
	//TODO - potential bug, could be unused import
	imports.Module(g.Modules.Models(api.InHttp.InVersion))
	if usingErrorModels(api) {
		imports.Module(g.Modules.HttpErrorsModels)
	}
	imports.Write(w)

	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			w.EmptyLine()
			generateResponseStruct(w, g.Types, &operation)
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
