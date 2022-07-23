package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen/ts/common"
	"github.com/specgen-io/specgen/v2/gen/ts/modules"
	"github.com/specgen-io/specgen/v2/gen/ts/responses"
	"github.com/specgen-io/specgen/v2/gen/ts/types"
	"github.com/specgen-io/specgen/v2/gen/ts/writer"
	"github.com/specgen-io/specgen/v2/generator"
	"github.com/specgen-io/specgen/v2/spec"
)

func generateServiceApis(version *spec.Version, modelsModule modules.Module, module modules.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		apiModule := module.Submodule(serviceName(&api))
		serviceFile := generateApiService(&api, modelsModule, apiModule)
		files = append(files, *serviceFile)
	}
	return files
}

func generateApiService(api *spec.Api, modelsModule modules.Module, module modules.Module) *generator.CodeFile {
	w := writer.NewTsWriter()
	w.Line("import * as %s from '%s'", types.ModelsPackage, modelsModule.GetImport(module))
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
	return &generator.CodeFile{module.GetPath(), w.String()}
}

func serviceInterfaceName(api *spec.Api) string {
	return api.Name.PascalCase() + "Service"
}

func serviceInterfaceNameVersioned(api *spec.Api) string {
	result := serviceInterfaceName(api)
	version := api.Http.Version.Version
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

func addServiceParam(w *generator.Writer, paramName string, typ *spec.TypeDef) {
	if typ.IsNullable() {
		paramName = paramName + "?"
	}
	w.Line("  %s: %s,", paramName, types.TsType(typ))
}

func generateServiceParams(w *generator.Writer, params []spec.NamedParam) {
	for _, param := range params {
		addServiceParam(w, common.TSIdentifier(param.Name.Source), &param.Type.Definition)
	}
}

func generateOperationParams(w *generator.Writer, operation *spec.NamedOperation) {
	w.Line("export interface %s {", operationParamsTypeName(operation))
	generateServiceParams(w, operation.HeaderParams)
	generateServiceParams(w, operation.Endpoint.UrlParams)
	generateServiceParams(w, operation.QueryParams)
	if operation.Body != nil {
		addServiceParam(w, "body", &operation.Body.Type.Definition)
	}
	w.Line("}")
}
