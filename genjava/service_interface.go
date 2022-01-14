package genjava

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func (g *Generator) ServicesInterfaces(version *spec.Version, thePackage Module, modelsVersionPackage Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		apiPackage := thePackage.Subpackage(api.Name.SnakeCase())
		files = append(files, g.serviceInterface(&api, apiPackage, modelsVersionPackage)...)
	}
	return files
}

func (g *Generator) serviceInterface(api *spec.Api, apiPackage Module, modelsVersionPackage Module) []sources.CodeFile {
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
		w.Line(`  %s;`, generateResponsesSignatures(g.Types, &operation))
	}
	w.Line(`}`)

	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			files = append(files, g.ResponsesInterfaces(&operation, apiPackage, modelsVersionPackage)...)
		}
	}

	files = append(files, sources.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", serviceInterfaceName(api))),
		Content: w.String(),
	})

	return files
}
