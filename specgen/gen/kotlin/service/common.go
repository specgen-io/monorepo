package service

import (
	"fmt"
	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/spec/v2"
	"github.com/specgen-io/specgen/v2/gen/kotlin/modules"
	"strings"
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
	serviceCall := fmt.Sprintf(`%s.%s(%s)`, serviceVarName(operation.Api), operation.Name.CamelCase(), joinParams(addServiceMethodParams(operation, bodyStringVar, bodyJsonVar)))
	if len(operation.Responses) == 1 && operation.Responses[0].BodyIs(spec.BodyEmpty) {
		w.Line(serviceCall)
	} else {
		w.Line(`val %s = %s`, resultVarName, serviceCall)
	}
}

func checkContentType(w *generator.Writer) {
	w.Lines(`
private fun checkContentType(request: HttpRequest<*>, expectedContentType: String) {
	val contentType = request.headers.contentType
	if (!(contentType.isPresent && contentType.get().contains(expectedContentType))) {
		throw ContentTypeMismatchException(expectedContentType, if (contentType.isPresent) contentType.get() else null )
	}
}
`)
}

func contentTypeMismatchException(thePackage modules.Module) *generator.CodeFile {
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
