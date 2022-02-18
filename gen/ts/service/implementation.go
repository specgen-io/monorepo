package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen/ts/modules"
	"github.com/specgen-io/specgen/v2/gen/ts/responses"
	"github.com/specgen-io/specgen/v2/gen/ts/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func generateServicesImplementations(specification *spec.Spec, generatedModule modules.Module, module modules.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, version := range specification.Versions {
		versionGeneratedModule := generatedModule.Submodule(version.Version.FlatCase())
		modelsModule := versionGeneratedModule.Submodule("models")
		for _, api := range version.Http.Apis {
			apiModule := versionGeneratedModule.Submodule(serviceName(&api)) //TODO: This logic is duplicated, other place is where API module is generated
			implModule := module.Submodule(version.Version.FlatCase()).Submodule(api.Name.SnakeCase() + "_service")
			files = append(files, *generateServiceImplementation(&api, apiModule, modelsModule, implModule))
		}
	}
	return files
}

func generateServiceImplementation(api *spec.Api, apiModule modules.Module, modelsModule modules.Module, module modules.Module) *sources.CodeFile {
	w := writer.NewTsWriter()

	w.Line("import * as service from '%s'", apiModule.GetImport(module))
	w.Line("import * as models from '%s'", modelsModule.GetImport(module))
	w.EmptyLine()
	w.Line("export let %sService = (): service.%s => {", api.Name.CamelCase(), serviceInterfaceName(api)) //TODO: remove services

	operations := []string{}
	for _, operation := range api.Operations {
		operations = append(operations, operation.Name.CamelCase())
		params := ""
		if operation.Body != nil || operation.HasParams() {
			params = fmt.Sprintf(`params: service.%s`, operationParamsTypeName(&operation))
		}
		w.Line("  let %s = async (%s): Promise<%s> => {", operation.Name.CamelCase(), params, responses.ResponseType(&operation, "service"))
		w.Line("    throw new Error('Not Implemented')")
		w.Line("  }")
		w.EmptyLine()
	}
	w.Line("  return {%s}", strings.Join(operations, ", "))
	w.Line("}")
	return &sources.CodeFile{module.GetPath(), w.String()}
}
