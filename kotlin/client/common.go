package client

import (
	"fmt"
	"strconv"

	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/spec/v2"
)

func addBuilderParam(param *spec.NamedParam) string {
	if param.Type.Definition.IsNullable() {
		return fmt.Sprintf(`%s!!`, param.Name.CamelCase())
	}
	return param.Name.CamelCase()
}

func generateTryCatch(w *generator.Writer, valName string, exceptionObject string, codeBlock func(w *generator.Writer), exceptionHandler func(w *generator.Writer)) {
	w.Line(`val %s = try {`, valName)
	codeBlock(w.Indented())
	w.Line(`} catch (%s) {`, exceptionObject)
	exceptionHandler(w.Indented())
	w.Line(`}`)
}

func generateClientTryCatch(w *generator.Writer, valName string, statement string, exceptionType, exceptionVar, errorMessage string) {
	generateTryCatch(w, valName, exceptionVar+`: `+exceptionType,
		func(w *generator.Writer) {
			w.Line(statement)
		},
		func(w *generator.Writer) {
			generateThrowClientException(w, errorMessage, exceptionVar)
		})
}

func generateThrowClientException(w *generator.Writer, errorMessage string, wrapException string) {
	w.Line(`val errorMessage = %s`, errorMessage)
	w.Line(`logger.error(errorMessage)`)
	params := "errorMessage"
	if wrapException != "" {
		params += ", " + wrapException
	}
	w.Line(`throw ClientException(%s)`, params)
}

func isSuccessfulStatusCode(statusCodeStr string) bool {
	statusCode, _ := strconv.Atoi(statusCodeStr)
	if statusCode >= 200 && statusCode <= 299 {
		return true
	}
	return false
}

func responsesNumber(operation *spec.NamedOperation) int {
	count := 0
	for _, response := range operation.Responses {
		if isSuccessfulStatusCode(spec.HttpStatusCode(response.Name)) {
			count++
		}
	}
	return count
}
