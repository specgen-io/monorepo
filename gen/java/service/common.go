package service

import (
	"github.com/specgen-io/specgen/v2/sources"
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

func generateServiceTryCatch(w *sources.Writer, statement, exceptionType, exceptionVar, errorMessage, responseInstance string) {
	generateTryCatch(w, exceptionType+` `+exceptionVar,
		func(w *sources.Writer) {
			w.Line(statement)
		},
		func(w *sources.Writer) {
			generateThrowClientException(w, errorMessage, responseInstance)
		})
}

func generateTryCatch(w *sources.Writer, exceptionObject string, codeBlock func(w *sources.Writer), exceptionHandler func(w *sources.Writer)) {
	w.Line(`try {`)
	codeBlock(w.Indented())
	w.Line(`} catch (%s) {`, exceptionObject)
	exceptionHandler(w.Indented())
	w.Line(`}`)
}

func generateThrowClientException(w *sources.Writer, errorMessage, responseInstance string) {
	w.Line(`logger.error(%s);`, errorMessage)
	w.Line(`return %s;`, responseInstance)
}
