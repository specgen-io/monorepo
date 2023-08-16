package client

import (
	"fmt"
	"kotlin/types"
	"spec"
	"strings"
)

func operationSignature(types *types.Types, operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%s(%s): %s`,
		operation.Name.CamelCase(),
		strings.Join(operationParameters(operation, types), ", "),
		operationReturnType(types, operation),
	)
}

func operationReturnType(types *types.Types, operation *spec.NamedOperation) string {
	if len(operation.Responses.Success()) == 1 {
		return types.ResponseBodyType(&operation.Responses.Success()[0].Body)
	}
	return responseInterfaceName(operation)
}

func operationParameters(operation *spec.NamedOperation, types *types.Types) []string {
	params := []string{}
	if operation.BodyIs(spec.RequestBodyString) || operation.BodyIs(spec.RequestBodyJson) {
		params = append(params, fmt.Sprintf("body: %s", types.Kotlin(&operation.Body.Type.Definition)))
	}
	for _, param := range operation.QueryParams {
		params = append(params, fmt.Sprintf("%s: %s", param.Name.CamelCase(), types.Kotlin(&param.Type.Definition)))
	}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf("%s: %s", param.Name.CamelCase(), types.Kotlin(&param.Type.Definition)))
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, fmt.Sprintf("%s: %s", param.Name.CamelCase(), types.Kotlin(&param.Type.Definition)))
	}
	return params
}
