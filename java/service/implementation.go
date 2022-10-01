package service

import (
	"fmt"

	"generator"
	"java/imports"
	"java/writer"
	"spec"
)

func (g *Generator) ServicesImplementations(version *spec.Version) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		files = append(files, *g.serviceImplementation(&api))
	}
	return files
}

func (g *Generator) serviceImplementation(api *spec.Api) *generator.CodeFile {
	packages := g.Packages.Version(api.InHttp.InVersion)
	w := writer.NewJavaWriter()
	w.Line(`package %s;`, packages.ServicesImpl.PackageName)
	w.EmptyLine()
	annotationImport, annotation := g.ServiceImplAnnotation(api)
	imports := imports.New()
	imports.Add(annotationImport)
	imports.Add(packages.Models.PackageStar)
	imports.Add(packages.ServicesApi(api).PackageStar)
	imports.Add(g.Types.Imports()...)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`@%s`, annotation)
	w.Line(`public class %s implements %s {`, serviceImplName(api), serviceInterfaceName(api))
	for _, operation := range api.Operations {
		w.Line(`  @Override`)
		w.Line(`  public %s {`, operationSignature(g.Types, &operation))
		w.Line(`    throw new UnsupportedOperationException("Implementation has not added yet");`)
		w.Line(`  }`)
	}
	w.Line(`}`)

	return &generator.CodeFile{
		Path:    packages.ServicesImpl.GetPath(fmt.Sprintf("%s.java", serviceImplName(api))),
		Content: w.String(),
	}
}
