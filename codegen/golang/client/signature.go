package client

import (
	"fmt"
	"golang/module"
	"golang/types"
	"spec"
	"strings"
)

func operationSignature(types *types.Types, operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%s(%s) %s`,
		operation.Name.PascalCase(),
		strings.Join(operationParams(types, operation), ", "),
		operationReturn(types, operation),
	)
}

func operationReturn(types *types.Types, operation *spec.NamedOperation) string {
	successResponses := operation.Responses.Success()
	if len(successResponses) == 1 {
		successResponse := successResponses[0]
		if successResponse.Body.IsEmpty() {
			return `error`
		} else {
			return fmt.Sprintf(`(%s, error)`, types.ResponseBodyGoType(&successResponse.Body))
		}
	} else {
		return fmt.Sprintf(`(*%s, error)`, responseTypeName(operation))
	}
}

func operationError(operation *spec.NamedOperation, errorVar string) string {
	successResponses := operation.Responses.Success()
	if len(successResponses) == 1 && successResponses[0].Body.IsEmpty() {
		return errorVar
	} else {
		return fmt.Sprintf(`nil, %s`, errorVar)
	}
}

func resultSuccess(response *spec.OperationResponse, resultVar string) string {
	successResponses := response.Operation.Responses.Success()
	if len(successResponses) == 1 {
		if successResponses[0].Body.IsEmpty() {
			return `nil`
		} else {
			return fmt.Sprintf(`%s, nil`, resultVar)
		}
	} else {
		return fmt.Sprintf(`%s, nil`, newResponse(response, resultVar))
	}
}

func resultError(response *spec.OperationResponse, errorsModules module.Module, resultVar string) string {
	errorBody := ``
	if !response.Body.IsEmpty() {
		errorBody = resultVar
	}
	result := fmt.Sprintf(`&%s{%s}`, errorsModules.Get(response.Name.PascalCase()), errorBody)
	successResponses := response.Operation.Responses.Success()
	if len(successResponses) == 1 && successResponses[0].Body.IsEmpty() {
		return result
	} else {
		return fmt.Sprintf(`nil, %s`, result)
	}
}

func operationParams(types *types.Types, operation *spec.NamedOperation) []string {
	params := []string{}
	if operation.Body.IsText() || operation.Body.IsBinary() || operation.Body.IsJson() {
		params = append(params, fmt.Sprintf("body %s", types.RequestBodyGoType(&operation.Body)))
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
		params = append(params, fmt.Sprintf("%s %s", param.Name.CamelCase(), types.ParamGoType(&param)))
	}
	return params
}
