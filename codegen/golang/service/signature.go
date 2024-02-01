package service

import (
	"fmt"
	"golang/types"
	"spec"
	"strings"
)

func (g *Generator) operationSignature(operation *spec.NamedOperation, apiPackage *string) string {
	return fmt.Sprintf(`%s(%s) %s`,
		operation.Name.PascalCase(),
		strings.Join(operationParams(g.Types, operation), ", "),
		g.operationReturn(operation, apiPackage),
	)
}

func (g *Generator) operationReturn(operation *spec.NamedOperation, responsePackageName *string) string {
	if len(operation.Responses) == 1 {
		response := operation.Responses[0]
		if response.Body.IsEmpty() {
			return `error`
		} else {
			return fmt.Sprintf(`(%s, error)`, g.Types.ResponseBodyGoType(&response.Body))
		}
	}
	responseType := responseTypeName(operation)
	if responsePackageName != nil {
		responseType = *responsePackageName + "." + responseType
	}
	return fmt.Sprintf(`(*%s, error)`, responseType)
}

func operationParams(types *types.Types, operation *spec.NamedOperation) []string {
	params := []string{}
	if operation.Body.IsText() || operation.Body.IsBinary() || operation.Body.IsJson() {
		params = append(params, fmt.Sprintf("body %s", types.RequestBodyGoType(&operation.Body)))
	}
	if operation.Body.IsBodyFormData() {
		for _, param := range operation.Body.FormData {
			params = appendParam(types, params, param)
		}
	}
	if operation.Body.IsBodyFormUrlEncoded() {
		for _, param := range operation.Body.FormUrlEncoded {
			params = appendParam(types, params, param)
		}
	}
	for _, param := range operation.QueryParams {
		params = appendParam(types, params, param)
	}
	for _, param := range operation.HeaderParams {
		params = appendParam(types, params, param)
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = appendParam(types, params, param)
	}
	return params
}

func appendParam(types *types.Types, params []string, param spec.NamedParam) []string {
	return append(params, fmt.Sprintf("%s %s", param.Name.CamelCase(), types.GoType(&param.Type.Definition)))
}
