package service

import (
	"fmt"
	"kotlin/writer"
	"spec"
	"strings"
)

func joinParams(params []string) string {
	return strings.Join(params, ", ")
}

func addServiceMethodParams(operation *spec.NamedOperation, bodyStringVar, bodyJsonVar string) []string {
	methodParams := []string{}
	if operation.BodyIs(spec.RequestBodyString) {
		methodParams = append(methodParams, bodyStringVar)
	}
	if operation.BodyIs(spec.RequestBodyJson) {
		methodParams = append(methodParams, bodyJsonVar)
	}
	if operation.BodyIs(spec.RequestBodyFormData) {
		for _, param := range operation.Body.FormData {
			methodParams = append(methodParams, param.Name.CamelCase())
		}
	}
	if operation.BodyIs(spec.RequestBodyFormUrlEncoded) {
		for _, param := range operation.Body.FormUrlEncoded {
			methodParams = append(methodParams, param.Name.CamelCase())
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

func serviceCall(w *writer.Writer, operation *spec.NamedOperation, bodyStringVar, bodyJsonVar, resultVarName string) {
	serviceCall := fmt.Sprintf(`%s.%s(%s)`, serviceVarName(operation.InApi), operation.Name.CamelCase(), joinParams(addServiceMethodParams(operation, bodyStringVar, bodyJsonVar)))
	if len(operation.Responses) == 1 && operation.Responses[0].Body.Is(spec.ResponseBodyEmpty) {
		w.Line(serviceCall)
	} else {
		w.Line(`val %s = %s`, resultVarName, serviceCall)
	}
}
