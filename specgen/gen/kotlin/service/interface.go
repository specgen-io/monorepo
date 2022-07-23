package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen/kotlin/imports"
	"github.com/specgen-io/specgen/v2/gen/kotlin/modules"
	"github.com/specgen-io/specgen/v2/gen/kotlin/responses"
	"github.com/specgen-io/specgen/v2/gen/kotlin/writer"
	"github.com/specgen-io/specgen/v2/generator"
	"github.com/specgen-io/specgen/v2/spec"
)

func (g *Generator) ServicesInterfaces(version *spec.Version, thePackage, modelsVersionPackage modules.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		apiPackage := thePackage.Subpackage(api.Name.SnakeCase())
		files = append(files, g.serviceInterface(&api, apiPackage, modelsVersionPackage)...)
	}
	return files
}

func (g *Generator) serviceInterface(api *spec.Api, apiPackage, modelsVersionPackage modules.Module) []generator.CodeFile {
	files := []generator.CodeFile{}

	w := writer.NewKotlinWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.Types.Imports()...)
	imports.Add(modelsVersionPackage.PackageStar)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`interface %s {`, serviceInterfaceName(api))
	for _, operation := range api.Operations {
		w.Line(`  fun %s`, responses.Signature(g.Types, &operation))
	}
	w.Line(`}`)

	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			files = append(files, responses.Interfaces(g.Types, &operation, apiPackage, modelsVersionPackage)...)
		}
	}

	files = append(files, generator.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.kt", serviceInterfaceName(api))),
		Content: w.String(),
	})

	return files
}
