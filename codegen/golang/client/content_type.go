package client

import (
	"fmt"
	"spec"
)

func ContentType(operation *spec.NamedOperation) string {
	if operation.BodyIs(spec.RequestBodyEmpty) {
		return ""
	} else if operation.BodyIs(spec.RequestBodyString) {
		return `"text/plain"`
	} else if operation.BodyIs(spec.RequestBodyJson) {
		return `"application/json"`
	} else if operation.BodyIs(spec.RequestBodyFormData) {
		return `writer.FormDataContentType()`
	} else if operation.BodyIs(spec.RequestBodyFormUrlEncoded) {
		return `"application/x-www-form-urlencoded"`
	} else {
		panic(fmt.Sprintf("Unknown Contet Type"))
	}
}
