package service

import (
	"fmt"
	"strings"

	"generator"
	"kotlin/packages"
	"spec"
)

func joinParams(params []string) string {
	return strings.Join(params, ", ")
}

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

func serviceCall(w *generator.Writer, operation *spec.NamedOperation, bodyStringVar, bodyJsonVar, resultVarName string) {
	serviceCall := fmt.Sprintf(`%s.%s(%s)`, serviceVarName(operation.InApi), operation.Name.CamelCase(), joinParams(addServiceMethodParams(operation, bodyStringVar, bodyJsonVar)))
	if len(operation.Responses) == 1 && operation.Responses[0].BodyIs(spec.BodyEmpty) {
		w.Line(serviceCall)
	} else {
		w.Line(`val %s = %s`, resultVarName, serviceCall)
	}
}

func contentTypeMismatchException(thePackage packages.Package) *generator.CodeFile {
	code := `
package [[.PackageName]]

class ContentTypeMismatchException(expected: String, actual: String?) :
    RuntimeException(
        String.format(
            "Expected Content-Type header: '%s' was not provided, found: '%s'",
            expected,
            actual
        )
    )
`
	code, _ = generator.ExecuteTemplate(code, struct {
		PackageName string
	}{
		thePackage.PackageName,
	})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("ContentTypeMismatchException.kt"),
		Content: strings.TrimSpace(code),
	}
}

func jsonParseException(thePackage packages.Package) *generator.CodeFile {
	code := `
package [[.PackageName]]

class JsonParseException(exception: Throwable) :
    RuntimeException("Failed to parse body: " + exception.message, exception)
`
	code, _ = generator.ExecuteTemplate(code, struct {
		PackageName string
	}{
		thePackage.PackageName,
	})
	return &generator.CodeFile{
		Path:    thePackage.GetPath("JsonParseException.kt"),
		Content: strings.TrimSpace(code),
	}
}
