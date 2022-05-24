package service

import (
	"github.com/specgen-io/specgen/v2/gen/golang/common"
	imports2 "github.com/specgen-io/specgen/v2/gen/golang/imports"
	"github.com/specgen-io/specgen/v2/gen/golang/module"
	"github.com/specgen-io/specgen/v2/gen/golang/responses"
	"github.com/specgen-io/specgen/v2/gen/golang/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func generateServicesInterfaces(version *spec.Version, versionModule, modelsModule, emptyModule module.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		apiModule := versionModule.Submodule(api.Name.SnakeCase())
		files = append(files, *generateInterface(&api, apiModule, modelsModule, emptyModule))
	}
	return files
}

func generateInterface(api *spec.Api, apiModule, modelsModule, emptyModule module.Module) *sources.CodeFile {
	w := writer.NewGoWriter()
	w.Line("package %s", apiModule.Name)

	imports := imports2.Imports()
	imports.AddApiTypes(api)
	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 && common.OperationHasType(&operation, spec.TypeEmpty) {
			imports.Add(emptyModule.Package)
		}
	}
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
		w.Line(`  %s(%s) %s`,
			operation.Name.PascalCase(),
			strings.Join(common.AddMethodParams(&operation), ", "),
			common.OperationReturn(&operation, nil),
		)
	}
	w.Line(`}`)
	return &sources.CodeFile{
		Path:    apiModule.GetPath("service.go"),
		Content: w.String(),
	}
}

const serviceInterfaceName = "Service"
