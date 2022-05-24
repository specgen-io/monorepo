package common

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen/golang/responses"
	"github.com/specgen-io/specgen/v2/gen/golang/types"
	"github.com/specgen-io/specgen/v2/spec"
)

func OperationReturn(operation *spec.NamedOperation, responsePackageName *string) string {
	if len(operation.Responses) == 1 {
		response := operation.Responses[0]
		if response.Type.Definition.IsEmpty() {
			return `error`
		}
		return fmt.Sprintf(`(*%s, error)`, types.GoType(&response.Type.Definition))
	}
	responseType := responses.ResponseTypeName(operation)
	if responsePackageName != nil {
		responseType = *responsePackageName + "." + responseType
	}
	return fmt.Sprintf(`(*%s, error)`, responseType)
}

func AddMethodParams(operation *spec.NamedOperation) []string {
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
