package common

import (
	"spec"
	"strconv"
)

func IsSuccessfulStatusCode(statusCodeStr string) bool {
	statusCode, _ := strconv.Atoi(statusCodeStr)
	if statusCode >= 200 && statusCode <= 299 {
		return true
	}
	return false
}

func ResponsesNumber(operation *spec.NamedOperation) int {
	count := 0
	for _, response := range operation.Responses {
		if IsSuccessfulStatusCode(spec.HttpStatusCode(response.Name)) {
			count++
		}
	}
	return count
}
