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
		} else if successResponse.Body.IsBinary() {
			return fmt.Sprintf(`(%s, error)`, types.ResponseBodyGoType(&successResponse.Body))
		} else {
			return fmt.Sprintf(`(*%s, error)`, types.GoType(&successResponse.Body.Type.Definition))
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
		} else if successResponses[0].Body.IsBinary() {
			return fmt.Sprintf(`%s, nil`, resultVar)
		} else {
			return fmt.Sprintf(`&%s, nil`, resultVar)
		}
	} else {
		return fmt.Sprintf(`&%s, nil`, newResponse(response, resultVar))
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
	if operation.Body.IsText() {
		params = append(params, fmt.Sprintf("body %s", types.GoType(&operation.Body.Type.Definition)))
	}
	if operation.Body.IsJson() {
		params = append(params, fmt.Sprintf("body *%s", types.GoType(&operation.Body.Type.Definition)))
	}
	if operation.Body.IsBinary() {
		params = append(params, fmt.Sprintf("body %s", types.RequestBodyGoType(&operation.Body)))
	}
	if operation.Body.IsBodyFormData() {
		for _, param := range operation.Body.FormData {
			params = append(params, fmt.Sprintf("%s %s", param.Name.CamelCase(), types.GoType(&param.Type.Definition)))
		}
	}
	if operation.Body.IsBodyFormUrlEncoded() {
		for _, param := range operation.Body.FormUrlEncoded {
			params = append(params, fmt.Sprintf("%s %s", param.Name.CamelCase(), types.GoType(&param.Type.Definition)))
		}
	}
	for _, param := range operation.QueryParams {
		params = append(params, fmt.Sprintf("%s %s", param.Name.CamelCase(), types.GoType(&param.Type.Definition)))
	}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf("%s %s", param.Name.CamelCase(), types.GoType(&param.Type.Definition)))
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, fmt.Sprintf("%s %s", param.Name.CamelCase(), types.GoType(&param.Type.Definition)))
	}
	return params
}
