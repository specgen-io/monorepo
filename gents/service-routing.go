package gents

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func generateParametersParsing(w *sources.Writer, validation string, operation *spec.NamedOperation, body, rawBody string, headers string, urlParams string, query string, badRequestStatement string) string {
	if operation.Body == nil && !operation.HasParams() {
		return ""
	}

	apiCallParams := []string{}
	if operation.Body != nil {
		if operation.Body.Type.Definition.Plain == spec.TypeString {
			w.Line(`const body: string = %s`, rawBody)
		} else {
			w.Line("var body: %s", TsType(&operation.Body.Type.Definition))
		}
	}
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
	if operation.Body != nil {
		if operation.Body.Type.Definition.Plain != spec.TypeString {
			w.Line("  body = t.decode(%s.%s, %s)", modelsPackage, runtimeType(validation, &operation.Body.Type.Definition), body)
		}
		apiCallParams = append(apiCallParams, "body")
	}
	if len(operation.Endpoint.UrlParams) > 0 {
		w.Line("  urlParams = t.decode(%s, %s)", paramsRuntimeTypeName(paramsTypeName(operation, "UrlParams")), urlParams)
		apiCallParams = append(apiCallParams, "...urlParams")
	}
	if len(operation.HeaderParams) > 0 {
		w.Line("  headerParams = t.decode(%s, %s)", paramsRuntimeTypeName(paramsTypeName(operation, "HeaderParams")), headers)
		apiCallParams = append(apiCallParams, "...headerParams")
	}
	if len(operation.QueryParams) > 0 {
		w.Line("  queryParams = t.decode(%s, %s)", paramsRuntimeTypeName(paramsTypeName(operation, "QueryParams")), query)
		apiCallParams = append(apiCallParams, "...queryParams")
	}
	w.Line("} catch (error) {")
	w.Line("  %s", badRequestStatement)
	w.Line("  return")
	w.Line("}")

	return fmt.Sprintf(`{%s}`, strings.Join(apiCallParams, ", "))
}

func serviceCall(operation *spec.NamedOperation, paramsObject string) string {
	if len(operation.Responses) == 1 && operation.Responses[0].Definition.Type.Definition.IsEmpty() {
		return fmt.Sprintf("await service.%s(%s)", operation.Name.CamelCase(), paramsObject)
	}
	return fmt.Sprintf("let result = await service.%s(%s)", operation.Name.CamelCase(), paramsObject)
}
