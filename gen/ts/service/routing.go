package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func getApiCallParamsObject(operation *spec.NamedOperation) string {
	if operation.Body == nil && !operation.HasParams() {
		return ""
	}

	apiCallParams := []string{}
	if operation.Body != nil {
		apiCallParams = append(apiCallParams, "body")
	}
	if len(operation.Endpoint.UrlParams) > 0 {
		apiCallParams = append(apiCallParams, "...urlParams")
	}
	if len(operation.HeaderParams) > 0 {
		apiCallParams = append(apiCallParams, "...headerParams")
	}
	if len(operation.QueryParams) > 0 {
		apiCallParams = append(apiCallParams, "...queryParams")
	}

	return fmt.Sprintf(`{%s}`, strings.Join(apiCallParams, ", "))

}

func serviceCall(operation *spec.NamedOperation, paramsObject string) string {
	if len(operation.Responses) == 1 && operation.Responses[0].Definition.Type.Definition.IsEmpty() {
		return fmt.Sprintf("await service.%s(%s)", operation.Name.CamelCase(), paramsObject)
	}
	return fmt.Sprintf("let result = await service.%s(%s)", operation.Name.CamelCase(), paramsObject)
}
