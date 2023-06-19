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
		}
		return fmt.Sprintf(`(*%s, error)`, g.Types.GoType(&response.Body.Type.Definition))
	}
	responseType := responseTypeName(operation)
	if responsePackageName != nil {
		responseType = *responsePackageName + "." + responseType
	}
	return fmt.Sprintf(`(*%s, error)`, responseType)
}

func operationParams(types *types.Types, operation *spec.NamedOperation) []string {
	params := []string{}
	if operation.BodyIs(spec.RequestBodyString) {
		params = append(params, fmt.Sprintf("body %s", types.GoType(&operation.Body.Type.Definition)))
	}
	if operation.BodyIs(spec.RequestBodyJson) {
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
