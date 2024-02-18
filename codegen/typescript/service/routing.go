package service

import (
	"fmt"
	"spec"
	"typescript/writer"
)

func parameterAssignment(source string, param *spec.NamedParam) string {
	return fmt.Sprintf("%s: %s['%s']", param.Name.CamelCase(), source, param.Name.Source)
}

func serviceCall(w *writer.Writer, operation *spec.NamedOperation) {
	parameters := ""
	if operation.Body.IsText() || operation.Body.IsJson() || operation.HasParams() {
		w.Line("const parameters = {")
		if operation.Body.IsText() || operation.Body.IsJson() {
			w.Line("  body,")
		}
		for _, param := range operation.Endpoint.UrlParams {
			w.Line("  %s,", parameterAssignment("urlParams.value", &param))
		}
		for _, param := range operation.HeaderParams {
			w.Line("  %s,", parameterAssignment("headerParams.value", &param))
		}
		for _, param := range operation.QueryParams {
			w.Line("  %s,", parameterAssignment("queryParams.value", &param))
		}
		w.Line("}")
		parameters = "parameters"
	}

	if len(operation.Responses) == 1 && operation.Responses[0].Body.IsEmpty() {
		w.Line(fmt.Sprintf("await service.%s(%s)", operation.Name.CamelCase(), parameters))
	} else {
		w.Line(fmt.Sprintf("const result = await service.%s(%s)", operation.Name.CamelCase(), parameters))
	}
}
