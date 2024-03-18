package service

import (
	"fmt"
	"kotlin/writer"
	"spec"
	"strings"
)

func addServiceMethodParams(operation *spec.NamedOperation, bodyStringVar, bodyJsonVar, bodyBinaryVar string, isSupportDefaulted bool) []string {
	methodParams := []string{}
	if operation.Body.IsText() {
		methodParams = append(methodParams, bodyStringVar)
	}
	if operation.Body.IsJson() {
		methodParams = append(methodParams, bodyJsonVar)
	}
	if operation.Body.IsBinary() {
		methodParams = append(methodParams, bodyBinaryVar)
	}
	if operation.Body.IsBodyFormData() {
		for _, param := range operation.Body.FormData {
			if !isSupportDefaulted && param.DefinitionDefault.Default != nil {
				methodParams = append(methodParams, fmt.Sprintf(`%s ?: "%s"`, param.Name.CamelCase(), *param.DefinitionDefault.Default))
			} else {
				methodParams = append(methodParams, param.Name.CamelCase())
			}
		}
	}
	if operation.Body.IsBodyFormUrlEncoded() {
		for _, param := range operation.Body.FormUrlEncoded {
			if !isSupportDefaulted && param.DefinitionDefault.Default != nil {
				methodParams = append(methodParams, fmt.Sprintf(`%s ?: "%s"`, param.Name.CamelCase(), *param.DefinitionDefault.Default))
			} else {
				methodParams = append(methodParams, param.Name.CamelCase())
			}
		}
	}
	for _, param := range operation.QueryParams {
		methodParams = append(methodParams, param.Name.CamelCase())
	}
	for _, param := range operation.HeaderParams {
		methodParams = append(methodParams, param.Name.CamelCase())
	}
	for _, param := range operation.Endpoint.UrlParams {
		methodParams = append(methodParams, param.Name.CamelCase())
	}
	return methodParams
}

func serviceCall(w *writer.Writer, operation *spec.NamedOperation, bodyStringVar, bodyJsonVar, bodyBinaryVar, resultVarName string, isSupportDefaulted bool) {
	serviceCall := fmt.Sprintf(`%s.%s(%s)`, serviceVarName(operation.InApi), operation.Name.CamelCase(), strings.Join(addServiceMethodParams(operation, bodyStringVar, bodyJsonVar, bodyBinaryVar, isSupportDefaulted), ", "))
	if len(operation.Responses) == 1 && operation.Responses[0].Body.IsEmpty() {
		w.Line(serviceCall)
	} else {
		w.Line(`val %s = %s`, resultVarName, serviceCall)
	}
}
