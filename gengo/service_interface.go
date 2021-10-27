package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
)

func generateServicesInterfaces(version *spec.Version, versionModule module, modelsModule module) []gen.TextFile {
	files := []gen.TextFile{}
	for _, api := range version.Http.Apis {
		apiModule := versionModule.Submodule(api.Name.SnakeCase())
		files = append(files, *generateInterface(&api, apiModule, modelsModule))
	}
	return files
}

func generateInterface(api *spec.Api, apiModule module, modelsModule module) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", apiModule.Name)

	imports := Imports()
	imports.AddApiTypes(api)
	imports.Add(modelsModule.Package)
	imports.Write(w)

	w.EmptyLine()
	w.Line(`type EmptyDef struct{}`)
	w.EmptyLine()
	w.Line(`var Empty = EmptyDef{}`)
	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			w.EmptyLine()
			generateOperationResponseStruct(w, operation)
		}
	}
	w.EmptyLine()
	w.Line(`type Service interface {`)
	for _, operation := range api.Operations {
		w.Line(`  %s(%s) (*%s, error)`, operation.Name.PascalCase(), JoinDelimParams(addMethodParams(operation)), generateResponsesSignatures(&operation))
	}
	w.Line(`}`)
	return &gen.TextFile{
		Path:    apiModule.GetPath("service.go"),
		Content: w.String(),
	}
}

func generateResponsesSignatures(operation *spec.NamedOperation) string {
	if len(operation.Responses) == 1 {
		for _, response := range operation.Responses {
			return GoType(&response.Type.Definition)
		}
	}
	if len(operation.Responses) > 1 {
		return responseTypeName(operation)
	}
	return ""
}

func generateOperationResponseStruct(w *gen.Writer, operation spec.NamedOperation) {
	w.Line(`type %s struct {`, responseTypeName(&operation))
	for _, response := range operation.Responses {
		w.Line(`  %s *%s`, response.Name.PascalCase(), GoType(&response.Type.Definition))
	}
	w.Line(`}`)
}

func addMethodParams(operation spec.NamedOperation) []string {
	params := []string{}
	if operation.Body != nil {
		params = append(params, fmt.Sprintf("body *%s", GoType(&operation.Body.Type.Definition)))
	}
	for _, param := range operation.QueryParams {
		params = append(params, fmt.Sprintf("%s %s", param.Name.CamelCase(), GoType(&param.Type.Definition)))
	}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf("%s %s", param.Name.CamelCase(), GoType(&param.Type.Definition)))
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, fmt.Sprintf("%s %s", param.Name.CamelCase(), GoType(&param.Type.Definition)))
	}
	return params
}
