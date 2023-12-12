package client

import (
	"fmt"
	"java/types"
	"spec"
	"strings"
)

func operationSignature(types *types.Types, operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%s %s(%s)`,
		operationReturnType(types, operation),
		operation.Name.CamelCase(),
		strings.Join(operationParameters(operation, types), ", "),
	)
}

func operationReturnType(types *types.Types, operation *spec.NamedOperation) string {
	if len(operation.Responses.Success()) == 1 {
		return types.ResponseBodyJavaType(&operation.Responses.Success()[0].Body)
	}
	return responseInterfaceName(operation)
}

func operationParameters(operation *spec.NamedOperation, types *types.Types) []string {
	params := []string{}
	if operation.Body.IsText() || operation.Body.IsJson() {
		params = append(params, fmt.Sprintf("%s body", types.Java(&operation.Body.Type.Definition)))
	}
	if operation.Body.IsBodyFormData() {
		for _, param := range operation.Body.FormData {
			params = append(params, fmt.Sprintf("%s %s", types.Java(&param.Type.Definition), param.Name.CamelCase()))
		}
	}
	if operation.Body.IsBodyFormUrlEncoded() {
		for _, param := range operation.Body.FormUrlEncoded {
			params = append(params, fmt.Sprintf("%s %s", types.Java(&param.Type.Definition), param.Name.CamelCase()))
		}
	}
	for _, param := range operation.QueryParams {
		params = append(params, fmt.Sprintf("%s %s", types.Java(&param.Type.Definition), param.Name.CamelCase()))
	}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf("%s %s", types.Java(&param.Type.Definition), param.Name.CamelCase()))
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, fmt.Sprintf("%s %s", types.Java(&param.Type.Definition), param.Name.CamelCase()))
	}
	return params
}
