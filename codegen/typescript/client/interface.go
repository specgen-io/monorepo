package client

import (
	"fmt"
	"spec"
	"typescript/types"
	"typescript/writer"
)

func operationSignature(operation *spec.NamedOperation) string {
	params := ""
	if operation.Body.IsText() || operation.Body.IsJson() || operation.HasParams() {
		params = fmt.Sprintf(`parameters: %s`, operationParamsTypeName(operation))
	}
	return fmt.Sprintf(`%s: async (%s): Promise<%s>`, operation.Name.CamelCase(), params, responseType(operation))
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
	if !operation.Body.IsEmpty() {
		w.Line("  body: %s,", types.RequestBodyTsType(&operation.Body))
	}
	w.Line("}")
}
