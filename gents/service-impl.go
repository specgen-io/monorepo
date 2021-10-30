package gents

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"strings"
)

func generateServicesImplementations(specification *spec.Spec, generatedModule module, module module) []gen.TextFile {
	files := []gen.TextFile{}
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			apiModule := generatedModule.Submodule(version.Version.FlatCase()).Submodule(serviceName(&api))  //TODO: This logic is duplicated, other place is where API module is generated
			implModule := module.Submodule(version.Version.FlatCase()).Submodule(api.Name.SnakeCase()+"_service")
			files = append(files, *generateServiceImplementation(&api, apiModule, implModule))
		}
	}
	return files
}

func generateServiceImplementation(api *spec.Api, apiModule module, module module) *gen.TextFile {
	w := NewTsWriter()

	w.Line("import * as service from '%s'", apiModule.GetImport(module))
	w.EmptyLine()
	w.Line("export let %sService = (): service.%s => {", api.Name.CamelCase(), serviceInterfaceName(api)) //TODO: remove services

	operations := []string{}
	for _, operation := range api.Operations {
		operations = append(operations, operation.Name.CamelCase())
		params := ""
		if operation.Body != nil || operation.HasParams() {
			params = fmt.Sprintf(`params: service.%s`, operationParamsTypeName(&operation))
		}
		w.Line("  let %s = async (%s): Promise<service.%s> => {", operation.Name.CamelCase(), params, responseTypeName(&operation))
		w.Line("    throw new Error('Not Implemented')")
		w.Line("  }")
		w.EmptyLine()
	}
	w.Line("  return {%s}", strings.Join(operations, ", "))
	w.Line("}")
	return &gen.TextFile{module.GetPath(), w.String()}
}
