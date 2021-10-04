package genjava

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
)

func generateServicesInterfaces(version *spec.Version, thePackage Module, modelsVersionPackage Module) []gen.TextFile {
	files := []gen.TextFile{}
	for _, api := range version.Http.Apis {
		apiPackage := thePackage.Subpackage(api.Name.SnakeCase())
		files = append(files, generateInterface(&api, apiPackage, modelsVersionPackage)...)
	}
	return files
}

func generateInterface(api *spec.Api, apiPackage Module, modelsVersionPackage Module) []gen.TextFile {

	files := []gen.TextFile{}
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
		w.Line(`  %s;`, generateResponsesSignatures(operation))
	}
	w.Line(`}`)

	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			files = append(files, generateResponseInterface(operation, apiPackage, modelsVersionPackage)...)
		}
	}

	files = append(files, gen.TextFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", serviceInterfaceName(api))),
		Content: w.String(),
	})

	return files
}

func generateResponsesSignatures(operation spec.NamedOperation) string {
	if len(operation.Responses) == 1 {
		for _, response := range operation.Responses {
			return fmt.Sprintf(`%s %s(%s)`, JavaType(&response.Type.Definition), operation.Name.CamelCase(), JoinParams(addOperationResponseParams(operation)))
		}
	}
	if len(operation.Responses) > 1 {
		return fmt.Sprintf(`%s %s(%s)`, serviceResponseInterfaceName(operation), operation.Name.CamelCase(), JoinParams(addOperationResponseParams(operation)))
	}
	return ""
}

func addOperationResponseParams(operation spec.NamedOperation) []string {
	params := []string{}
	if operation.Body != nil {
		params = append(params, fmt.Sprintf("%s body", JavaType(&operation.Body.Type.Definition)))
	}
	for _, param := range operation.QueryParams {
		params = append(params, fmt.Sprintf("%s %s", JavaType(&param.Type.Definition), param.Name.CamelCase()))
	}
	for _, param := range operation.HeaderParams {
		params = append(params, fmt.Sprintf("%s %s", JavaType(&param.Type.Definition), param.Name.CamelCase()))
	}
	for _, param := range operation.Endpoint.UrlParams {
		params = append(params, fmt.Sprintf("%s %s", JavaType(&param.Type.Definition), param.Name.CamelCase()))
	}
	return params
}

func generateResponseInterface(operation spec.NamedOperation, apiPackage Module, modelsVersionPackage Module) []gen.TextFile {
	files := []gen.TextFile{}
	w := NewJavaWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	w.Line(`public interface %s {`, serviceResponseInterfaceName(operation))
	w.Line(`}`)

	for _, response := range operation.Responses {
		files = append(files, *generateResponsesImplementations(operation, response, apiPackage, modelsVersionPackage))
	}

	files = append(files, gen.TextFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", serviceResponseInterfaceName(operation))),
		Content: w.String(),
	})
	return files
}

func generateResponsesImplementations(operation spec.NamedOperation, response spec.NamedResponse, apiPackage Module, modelsVersionPackage Module) *gen.TextFile {
	w := NewJavaWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	w.Line(`import %s;`, modelsVersionPackage.PackageStar)
	w.EmptyLine()
	w.Line(`public class %s implements %s {`, serviceResponseImplName(operation, response), serviceResponseInterfaceName(operation))
	if !response.Type.Definition.IsEmpty() {
		w.Line(`  public %s %s;`, JavaType(&response.Type.Definition), response.Name.Source)
		w.EmptyLine()
		w.Line(`  public %s() {`, serviceResponseImplName(operation, response))
		w.Line(`  }`)
		w.EmptyLine()
		w.Line(`  public %s(%s %s) {`, serviceResponseImplName(operation, response), JavaType(&response.Type.Definition), response.Name.Source)
		w.Line(`    this.%s = %s;`, response.Name.Source, response.Name.Source)
		w.Line(`  }`)
	}
	w.Line(`}`)

	return &gen.TextFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", serviceResponseImplName(operation, response))),
		Content: w.String(),
	}
}
