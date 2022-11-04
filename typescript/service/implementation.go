package service

import (
	"fmt"
	"strings"
	"typescript/types"

	"generator"
	"spec"
	"typescript/modules"
	"typescript/responses"
	"typescript/writer"
)

func generateServicesImplementations(specification *spec.Spec, generatedModule modules.Module, errorsModule modules.Module, module modules.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, version := range specification.Versions {
		versionGeneratedModule := generatedModule.Submodule(version.Name.FlatCase())
		modelsModule := versionGeneratedModule.Submodule("models")
		for _, api := range version.Http.Apis {
			apiModule := versionGeneratedModule.Submodule(serviceName(&api)) //TODO: This logic is duplicated, other place is where API module is generated
			implModule := module.Submodule(version.Name.FlatCase()).Submodule(api.Name.SnakeCase() + "_service")
			files = append(files, *generateServiceImplementation(&api, apiModule, modelsModule, errorsModule, implModule))
		}
	}
	return files
}

func generateServiceImplementation(api *spec.Api, apiModule modules.Module, modelsModule modules.Module, errorsModule modules.Module, implModule modules.Module) *generator.CodeFile {
	w := writer.New(implModule)
	w.Line("import * as service from '%s'", apiModule.GetImport(implModule))
	w.Line("import * as %s from '%s'", types.ModelsPackage, modelsModule.GetImport(implModule))
	w.Line("import * as %s from '%s'", types.ErrorsPackage, errorsModule.GetImport(implModule))
	w.EmptyLine()
	w.Line("export const %sService = (): service.%s => {", api.Name.CamelCase(), serviceInterfaceName(api)) //TODO: remove services

	operations := []string{}
	for _, operation := range api.Operations {
		operations = append(operations, operation.Name.CamelCase())
		params := ""
		if operation.Body != nil || operation.HasParams() {
			params = fmt.Sprintf(`params: service.%s`, operationParamsTypeName(&operation))
		}
		w.Line("  const %s = async (%s): Promise<%s> => {", operation.Name.CamelCase(), params, responses.ResponseType(&operation, "service"))
		w.Line("    throw new Error('Not Implemented')")
		w.Line("  }")
		w.EmptyLine()
	}
	w.Line("  return {%s}", strings.Join(operations, ", "))
	w.Line("}")
	return w.ToCodeFile()
}
