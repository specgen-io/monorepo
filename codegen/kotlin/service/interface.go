package service

import (
	"fmt"
	"strings"

	"generator"
	"kotlin/types"
	"kotlin/writer"
	"spec"
)

func (g *Generator) ServicesInterfaces(version *spec.Version) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *g.serviceInterface(&api))
		for _, operation := range api.Operations {
			if len(operation.Responses) > 1 {
				files = append(files, *g.responseInterface(&operation))
			}
		}
	}
	return files
}

func (g *Generator) serviceInterface(api *spec.Api) *generator.CodeFile {
	w := writer.New(g.Packages.ServicesApi(api), serviceInterfaceName(api))
	w.Imports.Add(g.Types.Imports()...)
	w.Imports.PackageStar(g.Packages.Models(api.InHttp.InVersion))
	w.EmptyLine()
	w.Line(`interface [[.ClassName]] {`)
	for _, operation := range api.Operations {
		w.Line(`  fun %s`, operationSignature(g.Types, &operation))
	}
	w.Line(`}`)
	return w.ToCodeFile()
}

func operationSignature(types *types.Types, operation *spec.NamedOperation) string {
	if len(operation.Responses) == 1 {
		for _, response := range operation.Responses {
			if !response.Body.IsEmpty() {
				return fmt.Sprintf(`%s(%s): %s`, operation.Name.CamelCase(), strings.Join(operationParameters(operation, types), ", "), types.Kotlin(&response.Body.Type.Definition))
			} else {
				return fmt.Sprintf(`%s(%s)`, operation.Name.CamelCase(), strings.Join(operationParameters(operation, types), ", "))
			}
		}
	}
	if len(operation.Responses) > 1 {
		return fmt.Sprintf(`%s(%s): %s`, operation.Name.CamelCase(), strings.Join(operationParameters(operation, types), ", "), responseInterfaceName(operation))
	}
	return ""
}

func operationParameters(operation *spec.NamedOperation, types *types.Types) []string {
	params := []string{}
	if operation.Body != nil {
		params = append(params, fmt.Sprintf("body: %s", types.Kotlin(&operation.Body.Type.Definition)))
	}
	for _, param := range operation.QueryParams {
		params = append(params, fmt.Sprintf("%s: %s", param.Name.CamelCase(), types.Kotlin(&param.Type.Definition)))
	}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf("%s: %s", param.Name.CamelCase(), types.Kotlin(&param.Type.Definition)))
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, fmt.Sprintf("%s: %s", param.Name.CamelCase(), types.Kotlin(&param.Type.Definition)))
	}
	return params
}
