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

func generateParamsMembers(w *writer.Writer, params []spec.NamedParam) {
	for _, param := range params {
		paramName := param.Name.CamelCase()
		paramType := param.Type.Definition
		if paramType.IsNullable() {
			paramName = paramName + "?"
		}
		w.Line("  %s: %s,", paramName, types.TsType(&paramType))
	}
}

func generateOperationParams(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line("export interface %s {", operationParamsTypeName(operation))
	generateParamsMembers(w, operation.HeaderParams)
	generateParamsMembers(w, operation.Endpoint.UrlParams)
	generateParamsMembers(w, operation.QueryParams)
	if !operation.Body.IsEmpty() {
		w.Line("  body: %s,", types.RequestBodyTsType(&operation.Body))
	}
	w.Line("}")
}
