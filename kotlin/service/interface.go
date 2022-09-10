package service

import (
	"fmt"

	"generator"
	"kotlin/imports"
	"kotlin/modules"
	"kotlin/responses"
	"kotlin/writer"
	"spec"
)

func (g *Generator) ServicesInterfaces(version *spec.Version, thePackage, modelsVersionPackage modules.Module, errorModelsPackage modules.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		apiPackage := thePackage.Subpackage(api.Name.SnakeCase())
		files = append(files, g.serviceInterface(&api, apiPackage, modelsVersionPackage, errorModelsPackage)...)
	}
	return files
}

func (g *Generator) serviceInterface(api *spec.Api, apiPackage, modelsVersionPackage modules.Module, errorModelsPackage modules.Module) []generator.CodeFile {
	files := []generator.CodeFile{}

	w := writer.NewKotlinWriter()
	w.Line(`package %s`, apiPackage.PackageName)
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
			files = append(files, responses.Interfaces(g.Types, &operation, apiPackage, modelsVersionPackage, errorModelsPackage)...)
		}
	}

	files = append(files, generator.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.kt", serviceInterfaceName(api))),
		Content: w.String(),
	})

	return files
}
