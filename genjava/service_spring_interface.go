package genjava

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func (g *Generator) generateServicesInterfaces(version *spec.Version, thePackage Module, modelsVersionPackage Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		apiPackage := thePackage.Subpackage(api.Name.SnakeCase())
		files = append(files, g.generateInterface(&api, apiPackage, modelsVersionPackage)...)
	}
	return files
}

func (g *Generator) generateInterface(api *spec.Api, apiPackage Module, modelsVersionPackage Module) []sources.CodeFile {
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
		w.Line(`  %s;`, g.generateResponsesSignatures(&operation))
	}
	w.Line(`}`)

	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			files = append(files, g.generateResponseInterface(&operation, apiPackage, modelsVersionPackage)...)
		}
	}

	files = append(files, sources.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", serviceInterfaceName(api))),
		Content: w.String(),
	})

	return files
}

func (g *Generator) generateResponsesSignatures(operation *spec.NamedOperation) string {
	if len(operation.Responses) == 1 {
		for _, response := range operation.Responses {
			return fmt.Sprintf(`%s %s(%s)`, g.Types.JavaType(&response.Type.Definition), operation.Name.CamelCase(), JoinDelimParams(addOperationResponseParams(operation, g.Types)))
		}
	}
	if len(operation.Responses) > 1 {
		return fmt.Sprintf(`%s %s(%s)`, serviceResponseInterfaceName(operation), operation.Name.CamelCase(), JoinDelimParams(addOperationResponseParams(operation, g.Types)))
	}
	return ""
}

func addOperationResponseParams(operation *spec.NamedOperation, types *Types) []string {
	params := []string{}
	if operation.Body != nil {
		params = append(params, fmt.Sprintf("%s body", types.JavaType(&operation.Body.Type.Definition)))
	}
	for _, param := range operation.QueryParams {
		params = append(params, fmt.Sprintf("%s %s", types.JavaType(&param.Type.Definition), param.Name.CamelCase()))
	}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf("%s %s", types.JavaType(&param.Type.Definition), param.Name.CamelCase()))
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, fmt.Sprintf("%s %s", types.JavaType(&param.Type.Definition), param.Name.CamelCase()))
	}
	return params
}

func (g *Generator) generateResponseInterface(operation *spec.NamedOperation, apiPackage Module, modelsVersionPackage Module) []sources.CodeFile {
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
		g.generateResponsesImplementations(w.Indented(), &response)
	}
	w.Line(`}`)

	files = append(files, sources.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", serviceResponseInterfaceName(operation))),
		Content: w.String(),
	})
	return files
}

func (g *Generator) generateResponsesImplementations(w *sources.Writer, response *spec.NamedResponse) {
	serviceResponseImplementationName := response.Name.PascalCase()
	w.Line(`class %s implements %s {`, serviceResponseImplementationName, serviceResponseInterfaceName(response.Operation))
	if !response.Type.Definition.IsEmpty() {
		w.Line(`  public %s body;`, g.Types.JavaType(&response.Type.Definition))
		w.EmptyLine()
		w.Line(`  public %s() {`, serviceResponseImplementationName)
		w.Line(`  }`)
		w.EmptyLine()
		w.Line(`  public %s(%s body) {`, serviceResponseImplementationName, g.Types.JavaType(&response.Type.Definition))
		w.Line(`    this.body = body;`)
		w.Line(`  }`)
	}
	w.Line(`}`)
}
