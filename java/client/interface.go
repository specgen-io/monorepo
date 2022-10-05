package client

import (
	"fmt"
	"java/types"
	"spec"
	"strings"
)

func operationSignature(types *types.Types, operation *spec.NamedOperation) string {
	if len(operation.Responses) == 1 {
		for _, response := range operation.Responses {
			return fmt.Sprintf(`%s %s(%s)`, types.Java(&response.Type.Definition), operation.Name.CamelCase(), strings.Join(operationParameters(operation, types), ", "))
		}
	}
	if len(operation.Responses) > 1 {
		return fmt.Sprintf(`%s %s(%s)`, responseInterfaceName(operation), operation.Name.CamelCase(), strings.Join(operationParameters(operation, types), ", "))
	}
	return ""
}

func operationParameters(operation *spec.NamedOperation, types *types.Types) []string {
	params := []string{}
	if operation.Body != nil {
		params = append(params, fmt.Sprintf("%s body", types.Java(&operation.Body.Type.Definition)))
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
