package genjava

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func generateServicesInterfaces(version *spec.Version, thePackage Module, modelsVersionPackage Module, jsonlib string) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		apiPackage := thePackage.Subpackage(api.Name.SnakeCase())
		files = append(files, generateInterface(&api, apiPackage, modelsVersionPackage, jsonlib)...)
	}
	return files
}

func generateInterface(api *spec.Api, apiPackage Module, modelsVersionPackage Module, jsonlib string) []sources.CodeFile {
	files := []sources.CodeFile{}

	w := NewJavaWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	w.Line(`import java.math.BigDecimal;`)
	w.Line(`import java.time.*;`)
	w.Line(`import java.util.*;`)
	w.EmptyLine()
	w.Line(`import %s;`, modelsVersionPackage.PackageStar)
	w.EmptyLine()
	w.Line(`public interface %s {`, serviceInterfaceName(api))
	for _, operation := range api.Operations {
		w.Line(`  %s;`, generateResponsesSignatures(&operation, jsonlib))
	}
	w.Line(`}`)

	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			files = append(files, generateResponseInterface(&operation, apiPackage, modelsVersionPackage, jsonlib)...)
		}
	}

	files = append(files, sources.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", serviceInterfaceName(api))),
		Content: w.String(),
	})

	return files
}

func generateResponsesSignatures(operation *spec.NamedOperation, jsonlib string) string {
	if len(operation.Responses) == 1 {
		for _, response := range operation.Responses {
			return fmt.Sprintf(`%s %s(%s)`, JavaType(&response.Type.Definition, jsonlib), operation.Name.CamelCase(), JoinDelimParams(addOperationResponseParams(operation, jsonlib)))
		}
	}
	if len(operation.Responses) > 1 {
		return fmt.Sprintf(`%s %s(%s)`, serviceResponseInterfaceName(operation), operation.Name.CamelCase(), JoinDelimParams(addOperationResponseParams(operation, jsonlib)))
	}
	return ""
}

func addOperationResponseParams(operation *spec.NamedOperation, jsonlib string) []string {
	params := []string{}
	if operation.Body != nil {
		params = append(params, fmt.Sprintf("%s body", JavaType(&operation.Body.Type.Definition, jsonlib)))
	}
	for _, param := range operation.QueryParams {
		params = append(params, fmt.Sprintf("%s %s", JavaType(&param.Type.Definition, jsonlib), param.Name.CamelCase()))
	}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf("%s %s", JavaType(&param.Type.Definition, jsonlib), param.Name.CamelCase()))
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, fmt.Sprintf("%s %s", JavaType(&param.Type.Definition, jsonlib), param.Name.CamelCase()))
	}
	return params
}

func generateResponseInterface(operation *spec.NamedOperation, apiPackage Module, modelsVersionPackage Module, jsonlib string) []sources.CodeFile {
	files := []sources.CodeFile{}
	w := NewJavaWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	w.Line(`import %s;`, modelsVersionPackage.PackageStar)
	w.EmptyLine()
	w.Line(`public interface %s {`, serviceResponseInterfaceName(operation))
	for index, response := range operation.Responses {
		if index > 0 {
			w.EmptyLine()
		}
		generateResponsesImplementations(w.Indented(), jsonlib, &response)
	}
	w.Line(`}`)

	files = append(files, sources.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", serviceResponseInterfaceName(operation))),
		Content: w.String(),
	})
	return files
}

func generateResponsesImplementations(w *sources.Writer, jsonlib string, response *spec.NamedResponse) {
	serviceResponseImplementationName := response.Name.PascalCase()
	w.Line(`class %s implements %s {`, serviceResponseImplementationName, serviceResponseInterfaceName(response.Operation))
	if !response.Type.Definition.IsEmpty() {
		w.Line(`  public %s body;`, JavaType(&response.Type.Definition, jsonlib))
		w.EmptyLine()
		w.Line(`  public %s() {`, serviceResponseImplementationName)
		w.Line(`  }`)
		w.EmptyLine()
		w.Line(`  public %s(%s body) {`, serviceResponseImplementationName, JavaType(&response.Type.Definition, jsonlib))
		w.Line(`    this.body = body;`)
		w.Line(`  }`)
	}
	w.Line(`}`)
}
