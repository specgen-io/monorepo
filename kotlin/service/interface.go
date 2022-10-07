package service

import (
	"fmt"
	"strings"

	"generator"
	"kotlin/imports"
	"kotlin/types"
	"kotlin/writer"
	"spec"
)

func (g *Generator) ServicesInterfaces(version *spec.Version) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, g.serviceInterface(&api)...)
	}
	return files
}

func (g *Generator) serviceInterface(api *spec.Api) []generator.CodeFile {
	files := []generator.CodeFile{}

	w := writer.NewKotlinWriter()
	w.Line(`package %s`, g.Packages.Version(api.InHttp.InVersion).ServicesApi(api).PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.Types.Imports()...)
	imports.Add(g.Packages.Models(api.InHttp.InVersion).PackageStar)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`interface %s {`, serviceInterfaceName(api))
	for _, operation := range api.Operations {
		w.Line(`  fun %s`, operationSignature(g.Types, &operation))
	}
	w.Line(`}`)

	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			files = append(files, *g.responseInterface(&operation))
		}
	}

	files = append(files, generator.CodeFile{
		Path:    g.Packages.Version(api.InHttp.InVersion).ServicesApi(api).GetPath(fmt.Sprintf("%s.kt", serviceInterfaceName(api))),
		Content: w.String(),
	})

	return files
}

func operationSignature(types *types.Types, operation *spec.NamedOperation) string {
	if len(operation.Responses) == 1 {
		for _, response := range operation.Responses {
			if !response.Type.Definition.IsEmpty() {
				return fmt.Sprintf(`%s(%s): %s`, operation.Name.CamelCase(), strings.Join(operationParameters(operation, types), ", "), types.Kotlin(&response.Type.Definition))
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
