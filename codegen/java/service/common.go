package service

import (
	"spec"
)

func addServiceMethodParams(operation *spec.NamedOperation, bodyStringVar, bodyJsonVar, bodyBinaryVar string) []string {
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
			methodParams = append(methodParams, param.Name.CamelCase())
		}
	}
	if operation.Body.IsBodyFormUrlEncoded() {
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
