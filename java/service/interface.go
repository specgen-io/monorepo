package service

import (
	"fmt"

	"generator"
	"java/imports"
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
	version := api.InHttp.InVersion
	apiPackage := g.Packages.Version(version).ServicesApi(api)

	files := []generator.CodeFile{}

	w := writer.NewJavaWriter()
	w.Line(`package %s;`, apiPackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.Types.Imports()...)
	imports.Add(g.Packages.Models(version).PackageStar)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`public interface %s {`, serviceInterfaceName(api))
	for _, operation := range api.Operations {
		w.Line(`  %s;`, operationSignature(g.Types, &operation))
	}
	w.Line(`}`)

	for _, operation := range api.Operations {
		if len(operation.Responses) > 1 {
			files = append(files, *g.responseInterface(&operation))
		}
	}

	files = append(files, generator.CodeFile{
		Path:    apiPackage.GetPath(fmt.Sprintf("%s.java", serviceInterfaceName(api))),
		Content: w.String(),
	})

	return files
}
