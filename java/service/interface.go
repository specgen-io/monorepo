package service

import (
	"fmt"

	"generator"
	"java/imports"
	"java/responses"
	"java/writer"
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
	packages := g.Packages.Version(api.InHttp.InVersion)
	apiPackage := packages.ServicesApi(api)

	files := []generator.CodeFile{}

	w := writer.NewJavaWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.Types.Imports()...)
	imports.Add(packages.Models.PackageStar)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`public interface %s {`, serviceInterfaceName(api))
	for _, operation := range api.Operations {
		w.Line(`  %s;`, responses.Signature(g.Types, &operation))
	}
	w.Line(`}`)

	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			files = append(files, responses.Interfaces(g.Types, &operation, apiPackage, packages.Models, g.Packages.ErrorsModels)...)
		}
	}

	files = append(files, generator.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", serviceInterfaceName(api))),
		Content: w.String(),
	})

	return files
}
