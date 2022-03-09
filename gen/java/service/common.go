package service

import (
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func addServiceMethodParams(operation *spec.NamedOperation, bodyStringVar, bodyJsonVar string) []string {
	methodParams := []string{}
	if operation.BodyIs(spec.BodyString) {
		methodParams = append(methodParams, bodyStringVar)
	}
	if operation.BodyIs(spec.BodyJson) {
		methodParams = append(methodParams, bodyJsonVar)
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

func joinParams(params []string) string {
	return strings.Join(params, ", ")
}
