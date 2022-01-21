package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gents/types"
	"github.com/specgen-io/specgen/v2/gents/validations"
	"github.com/specgen-io/specgen/v2/gents/validations/common"
	"github.com/specgen-io/specgen/v2/sources"
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

func generateParametersParsing(w *sources.Writer, operation *spec.NamedOperation, headers, urlParams, query string, badRequestStatement string) {
	if operation.HasParams() {
		if len(operation.Endpoint.UrlParams) > 0 {
			w.Line("var urlParams: %s", paramsTypeName(operation, "UrlParams"))
		}
		if len(operation.HeaderParams) > 0 {
			w.Line("var headerParams: %s", paramsTypeName(operation, "HeaderParams"))
		}
		if len(operation.QueryParams) > 0 {
			w.Line("var queryParams: %s", paramsTypeName(operation, "QueryParams"))
		}
		w.Line("try {")
		if len(operation.Endpoint.UrlParams) > 0 {
			w.Line("  urlParams = t.decode(%s, %s)", common.ParamsRuntimeTypeName(paramsTypeName(operation, "UrlParams")), urlParams)
		}
		if len(operation.HeaderParams) > 0 {
			w.Line("  headerParams = t.decode(%s, %s)", common.ParamsRuntimeTypeName(paramsTypeName(operation, "HeaderParams")), headers)
		}
		if len(operation.QueryParams) > 0 {
			w.Line("  queryParams = t.decode(%s, %s)", common.ParamsRuntimeTypeName(paramsTypeName(operation, "QueryParams")), query)
		}
		w.Line("} catch (error) {")
		w.Line("  %s", badRequestStatement)
		w.Line("  return")
		w.Line("}")
	}
}

func generateBodyParsing(w *sources.Writer, validation validations.Validation, operation *spec.NamedOperation, body, rawBody string, badRequestStatement string) {
	if operation.Body != nil {
		if operation.Body.Type.Definition.Plain == spec.TypeString {
			w.Line(`const body: string = %s`, rawBody)
		} else {
			w.Line("var body: %s", types.TsType(&operation.Body.Type.Definition))
			w.Line("try {")
			w.Line("  body = t.decode(%s, %s)", validation.RuntimeTypeFromPackage(types.ModelsPackage, &operation.Body.Type.Definition), body)
			w.Line("} catch (error) {")
			w.Line("  %s", badRequestStatement)
			w.Line("  return")
			w.Line("}")
		}
	}
}

func serviceCall(operation *spec.NamedOperation, paramsObject string) string {
	if len(operation.Responses) == 1 && operation.Responses[0].Definition.Type.Definition.IsEmpty() {
		return fmt.Sprintf("await service.%s(%s)", operation.Name.CamelCase(), paramsObject)
	}
	return fmt.Sprintf("let result = await service.%s(%s)", operation.Name.CamelCase(), paramsObject)
}
