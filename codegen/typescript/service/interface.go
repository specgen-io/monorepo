package service

import (
	"fmt"
	"generator"
	"spec"
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
	w := writer.New(g.Modules.ServiceApi(api))
	w.Imports.Star(g.Modules.Models(api.InHttp.InVersion), types.ModelsPackage)
	w.Imports.Star(g.Modules.ErrorsModels, types.ErrorsPackage)
	for _, operation := range api.Operations {
		if operation.Body.IsText() || operation.Body.IsJson() || operation.Body.IsBodyFormData() || operation.Body.IsBodyFormUrlEncoded() || operation.HasParams() {
			w.EmptyLine()
			generateOperationParams(w, &operation)
		}
		if len(operation.Responses) > 1 {
			w.EmptyLine()
			GenerateOperationResponse(w, &operation)
		}
	}
	w.EmptyLine()
	w.Line("export interface %s {", serviceInterfaceName(api))
	for _, operation := range api.Operations {
		params := ""
		if operation.Body.IsText() || operation.Body.IsJson() || operation.Body.IsBodyFormData() || operation.Body.IsBodyFormUrlEncoded() || operation.HasParams() {
			params = fmt.Sprintf(`params: %s`, operationParamsTypeName(&operation))
		}
		w.Line("  %s(%s): Promise<%s>", operation.Name.CamelCase(), params, ResponseType(&operation, ""))
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

func operationParamsTypeName(operation *spec.NamedOperation) string {
	return operation.Name.PascalCase() + "Params"
}

func generateOperationParams(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line("export interface %s {", operationParamsTypeName(operation))
	for _, param := range operation.HeaderParams {
		w.Line("  %s,", types.ParamTsDeclaration(&param))
	}
	for _, param := range operation.Endpoint.UrlParams {
		w.Line("  %s,", types.ParamTsDeclaration(&param))
	}
	for _, param := range operation.QueryParams {
		w.Line("  %s,", types.ParamTsDeclaration(&param))
	}
	if operation.Body.IsBodyFormData() {
		for _, param := range operation.Body.FormData {
			w.Line("  %s,", types.ParamTsDeclaration(&param))
		}
	}
	if operation.Body.IsBodyFormUrlEncoded() {
		for _, param := range operation.Body.FormUrlEncoded {
			w.Line("  %s,", types.ParamTsDeclaration(&param))
		}
	}
	if operation.Body.IsText() || operation.Body.IsJson() {
		w.Line("  body: %s,", types.RequestBodyTsType(&operation.Body))
	}
	w.Line("}")
}
