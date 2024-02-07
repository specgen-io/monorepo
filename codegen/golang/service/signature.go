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
