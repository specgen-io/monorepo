package client

import (
	"fmt"
	"golang/common"
	"golang/types"
	"spec"
	"strings"
)

func operationSignature(operation *spec.NamedOperation, types *types.Types, apiPackage *string) string {
	return fmt.Sprintf(`%s(%s) %s`,
		operation.Name.PascalCase(),
		strings.Join(operationParams(types, operation), ", "),
		operationReturn(operation, types, apiPackage),
	)
}

func operationReturn(operation *spec.NamedOperation, types *types.Types, responsePackageName *string) string {
	if common.ResponsesNumber(operation) == 1 {
		response := operation.Responses[0]
		return fmt.Sprintf(`(*%s, error)`, types.GoType(&response.Type.Definition))
	}
	responseType := responseTypeName(operation)
	if responsePackageName != nil {
		responseType = *responsePackageName + "." + responseType
	}
	return fmt.Sprintf(`(*%s, error)`, responseType)
}

func operationParams(types *types.Types, operation *spec.NamedOperation) []string {
	params := []string{}
	if operation.BodyIs(spec.BodyString) {
		params = append(params, fmt.Sprintf("body %s", types.GoType(&operation.Body.Type.Definition)))
	}
	if operation.BodyIs(spec.BodyJson) {
		params = append(params, fmt.Sprintf("body *%s", types.GoType(&operation.Body.Type.Definition)))
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
