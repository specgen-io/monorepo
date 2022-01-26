package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/genjava/packages"
	"github.com/specgen-io/specgen/v2/genjava/responses"
	"github.com/specgen-io/specgen/v2/genjava/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func (g *Generator) ServicesInterfaces(version *spec.Version, thePackage packages.Module, modelsVersionPackage packages.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		apiPackage := thePackage.Subpackage(api.Name.SnakeCase())
		files = append(files, g.serviceInterface(&api, apiPackage, modelsVersionPackage)...)
	}
	return files
}

func (g *Generator) serviceInterface(api *spec.Api, apiPackage packages.Module, modelsVersionPackage packages.Module) []sources.CodeFile {
	files := []sources.CodeFile{}

	w := writer.NewJavaWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	GenerateImports(w, GeneralImports())
	w.EmptyLine()
	w.Line(`import %s;`, modelsVersionPackage.PackageStar)
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

	files = append(files, sources.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", serviceInterfaceName(api))),
		Content: w.String(),
	})

	return files
}
