package service

import (
	"fmt"
	"spec"
	"strings"
)

func getApiCallParamsObject(operation *spec.NamedOperation) string {
	if (operation.BodyIs(spec.RequestBodyEmpty) || operation.BodyIs(spec.RequestBodyFormData) || operation.BodyIs(spec.RequestBodyFormUrlEncoded)) && !operation.HasParams() {
		return ""
	}

	apiCallParams := []string{}
	if operation.BodyIs(spec.RequestBodyString) || operation.BodyIs(spec.RequestBodyJson) {
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
	if len(operation.Responses) == 1 && operation.Responses[0].Body.IsEmpty() {
		return fmt.Sprintf("await service.%s(%s)", operation.Name.CamelCase(), paramsObject)
	}
	return fmt.Sprintf("const result = await service.%s(%s)", operation.Name.CamelCase(), paramsObject)
}
