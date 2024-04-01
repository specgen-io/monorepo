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
	w.Imports.Add(g.FilesImports()...)
	w.Imports.PackageStar(g.Packages.Models(api.InHttp.InVersion))
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
				return fmt.Sprintf(`%s(%s): %s`, operation.Name.CamelCase(), strings.Join(operationParameters(operation, types), ", "), types.ResponseBodyKotlinType(&response.Body))
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
	if operation.Body.IsText() || operation.Body.IsBinary() || operation.Body.IsJson() {
		params = append(params, fmt.Sprintf("body: %s", types.RequestBodyKotlinType(&operation.Body)))
	}
	if operation.Body.IsBodyFormData() {
		params = appendParams(types, params, operation.Body.FormData)
	}
	if operation.Body.IsBodyFormUrlEncoded() {
		params = appendParams(types, params, operation.Body.FormUrlEncoded)
	}
	params = appendParams(types, params, operation.QueryParams)
	params = appendParams(types, params, operation.HeaderParams)
	params = appendParams(types, params, operation.Endpoint.UrlParams)
	return params
}

func appendParams(types *types.Types, params []string, namedParams []spec.NamedParam) []string {
	for _, param := range namedParams {
		params = append(params, fmt.Sprintf("%s: %s", param.Name.CamelCase(), types.ParamKotlinType(&param)))
	}
	return params
}
