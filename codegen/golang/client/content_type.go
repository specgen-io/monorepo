package client

import (
	"fmt"
	"spec"
)

func ContentType(operation *spec.NamedOperation) string {
	if operation.Body.IsEmpty() {
		return ""
	} else if operation.Body.IsText() {
		return `"text/plain"`
	} else if operation.Body.IsJson() {
		return `"application/json"`
	} else if operation.Body.IsBodyFormData() {
		return `writer.FormDataContentType()`
	} else if operation.Body.IsBodyFormUrlEncoded() {
		return `"application/x-www-form-urlencoded"`
	} else {
		panic(fmt.Sprintf("Unknown Contet Type"))
	}
}
