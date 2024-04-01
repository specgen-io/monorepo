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
		return types.ResponseBodyKotlinType(&operation.Responses.Success()[0].Body)
	}
	return responseInterfaceName(operation)
}

func operationParameters(operation *spec.NamedOperation, types *types.Types) []string {
	params := []string{}
	if operation.Body.IsText() || operation.Body.IsBinary() || operation.Body.IsJson() {
		params = append(params, fmt.Sprintf("body: %s", types.RequestBodyKotlinType(&operation.Body)))
	}
	if operation.Body.IsBodyFormData() {
		params = appendParams(types, params, operation.Body.FormData)
	}
	if operation.Body.IsBodyFormUrlEncoded() {
		params = appendParams(types, params, operation.Body.FormUrlEncoded)
	}
	params = appendParams(types, params, operation.QueryParams)
	params = appendParams(types, params, operation.HeaderParams)
	params = appendParams(types, params, operation.Endpoint.UrlParams)
	return params
}

func appendParams(types *types.Types, params []string, namedParams []spec.NamedParam) []string {
	for _, param := range namedParams {
		if param.Type.Definition.String() == spec.TypeFile {
			params = append(params, "fileName: String")
		}
		params = append(params, fmt.Sprintf("%s: %s", param.Name.CamelCase(), types.ParamKotlinType(&param)))
	}
	return params
}
