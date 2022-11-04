package service

import (
	"fmt"

	"generator"
	"spec"
	"typescript/common"
	"typescript/responses"
	"typescript/types"
	"typescript/writer"
)

func (g *Generator) ServiceApis(version *spec.Version) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *g.serviceApi(&api))
	}
	return files
}

func (g *Generator) serviceApi(api *spec.Api) *generator.CodeFile {
	apiModule := g.Modules.ServiceApi(api)
	w := writer.New(apiModule)
	w.Line("import * as %s from '%s'", types.ModelsPackage, g.Modules.Models(api.InHttp.InVersion).GetImport(apiModule))
	w.Line("import * as %s from '%s'", types.ErrorsPackage, g.Modules.Errors.GetImport(apiModule))
	for _, operation := range api.Operations {
		if operation.Body != nil || operation.HasParams() {
			w.EmptyLine()
			generateOperationParams(w, &operation)
		}
		if len(operation.Responses) > 1 {
			w.EmptyLine()
			responses.GenerateOperationResponse(w, &operation)
		}
	}
	w.EmptyLine()
	w.Line("export interface %s {", serviceInterfaceName(api))
	for _, operation := range api.Operations {
		params := ""
		if operation.Body != nil || operation.HasParams() {
			params = fmt.Sprintf(`params: %s`, operationParamsTypeName(&operation))
		}
		w.Line("  %s(%s): Promise<%s>", operation.Name.CamelCase(), params, responses.ResponseType(&operation, ""))
	}
	w.Line("}")
	return w.ToCodeFile()
}

func serviceInterfaceName(api *spec.Api) string {
	return api.Name.PascalCase() + "Service"
}

func serviceInterfaceNameVersioned(api *spec.Api) string {
	result := serviceInterfaceName(api)
	version := api.InHttp.InVersion.Name
	if version.Source != "" {
		result = result + version.PascalCase()
	}
	return result
}

func serviceName(api *spec.Api) string {
	return api.Name.SnakeCase()
}

func operationParamsTypeName(operation *spec.NamedOperation) string {
	return operation.Name.PascalCase() + "Params"
}

func addServiceParam(w generator.Writer, paramName string, typ *spec.TypeDef) {
	if typ.IsNullable() {
		paramName = paramName + "?"
	}
	w.Line("  %s: %s,", paramName, types.TsType(typ))
}

func generateServiceParams(w generator.Writer, params []spec.NamedParam) {
	for _, param := range params {
		addServiceParam(w, common.TSIdentifier(param.Name.Source), &param.Type.Definition)
	}
}

func generateOperationParams(w generator.Writer, operation *spec.NamedOperation) {
	w.Line("export interface %s {", operationParamsTypeName(operation))
	generateServiceParams(w, operation.HeaderParams)
	generateServiceParams(w, operation.Endpoint.UrlParams)
	generateServiceParams(w, operation.QueryParams)
	if operation.Body != nil {
		addServiceParam(w, "body", &operation.Body.Type.Definition)
	}
	w.Line("}")
}
