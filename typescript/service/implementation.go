package service

import (
	"fmt"
	"strings"
	"typescript/imports"
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
	w := writer.New(g.Modules.ServiceImpl(api))
	imports := imports.New(g.Modules.ServiceImpl(api))
	imports.Star(g.Modules.ServiceApi(api), `service`)
	imports.Star(g.Modules.Models(api.InHttp.InVersion), types.ModelsPackage)
	imports.Star(g.Modules.Errors, types.ErrorsPackage)
	imports.Write(w)
	w.EmptyLine()
	w.Line("export const %sService = (): service.%s => {", api.Name.CamelCase(), serviceInterfaceName(api))

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
