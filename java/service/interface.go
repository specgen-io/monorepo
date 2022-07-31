package service

import (
	"fmt"

	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/java/v2/imports"
	"github.com/specgen-io/specgen/java/v2/packages"
	"github.com/specgen-io/specgen/java/v2/responses"
	"github.com/specgen-io/specgen/java/v2/writer"
	"github.com/specgen-io/specgen/spec/v2"
)

func (g *Generator) ServicesInterfaces(version *spec.Version, thePackage, modelsVersionPackage packages.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		apiPackage := thePackage.Subpackage(api.Name.SnakeCase())
		files = append(files, g.serviceInterface(&api, apiPackage, modelsVersionPackage)...)
	}
	return files
}

func (g *Generator) serviceInterface(api *spec.Api, apiPackage, modelsVersionPackage packages.Module) []generator.CodeFile {
	files := []generator.CodeFile{}

	w := writer.NewJavaWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.Types.Imports()...)
	imports.Add(modelsVersionPackage.PackageStar)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`public interface %s {`, serviceInterfaceName(api))
	for _, operation := range api.Operations {
		w.Line(`  %s;`, responses.Signature(g.Types, &operation))
	}
	w.Line(`}`)

	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			files = append(files, responses.Interfaces(g.Types, &operation, apiPackage, modelsVersionPackage)...)
		}
	}

	files = append(files, generator.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", serviceInterfaceName(api))),
		Content: w.String(),
	})

	return files
}