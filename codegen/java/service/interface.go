package service

import (
	"fmt"
	"strings"

	"generator"
	"java/types"
	"java/writer"
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
	w.Imports.Star(g.Packages.Models(api.InHttp.InVersion))
	w.EmptyLine()
	w.Line(`public interface [[.ClassName]] {`)
	for _, operation := range api.Operations {
		w.Line(`  %s;`, operationSignature(g.Types, &operation))
	}
	w.Line(`}`)
	return w.ToCodeFile()
}

func operationSignature(types *types.Types, operation *spec.NamedOperation) string {
	if len(operation.Responses) == 1 {
		for _, response := range operation.Responses {
			return fmt.Sprintf(`%s %s(%s)`, types.Java(&response.Type.Definition), operation.Name.CamelCase(), strings.Join(operationParameters(operation, types), ", "))
		}
	}
	if len(operation.Responses) > 1 {
		return fmt.Sprintf(`%s %s(%s)`, responseInterfaceName(operation), operation.Name.CamelCase(), strings.Join(operationParameters(operation, types), ", "))
	}
	return ""
}

func operationParameters(operation *spec.NamedOperation, types *types.Types) []string {
	params := []string{}
	if operation.Body != nil {
		params = append(params, fmt.Sprintf("%s body", types.Java(&operation.Body.Type.Definition)))
	}
	for _, param := range operation.QueryParams {
		params = append(params, fmt.Sprintf("%s %s", types.Java(&param.Type.Definition), param.Name.CamelCase()))
	}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf("%s %s", types.Java(&param.Type.Definition), param.Name.CamelCase()))
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, fmt.Sprintf("%s %s", types.Java(&param.Type.Definition), param.Name.CamelCase()))
	}
	return params
}
