package service

import (
	"fmt"
	"strings"
	"typescript/types"

	"generator"
	"spec"
	"typescript/responses"
	"typescript/writer"
)

func (g *Generator) ServicesImpls(specification *spec.Spec) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			files = append(files, *g.serviceImpl(&api))
		}
	}
	return files
}

func (g *Generator) serviceImpl(api *spec.Api) *generator.CodeFile {
	implModule := g.Modules.ServiceImpl(api)
	w := writer.New(implModule)
	w.Line("import * as service from '%s'", g.Modules.ServiceApi(api).GetImport(implModule))
	w.Line("import * as %s from '%s'", types.ModelsPackage, g.Modules.Models(api.InHttp.InVersion).GetImport(implModule))
	w.Line("import * as %s from '%s'", types.ErrorsPackage, g.Modules.Errors.GetImport(implModule))
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
