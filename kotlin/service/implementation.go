package service

import (
	"fmt"

	"generator"
	"kotlin/imports"
	"kotlin/writer"
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
	w := writer.NewKotlinWriter()
	w.Line(`package %s`, g.Packages.Version(api.InHttp.InVersion).ServicesImpl.PackageName)
	w.EmptyLine()
	annotationImport, annotation := g.Server.ServiceImplAnnotation(api)
	imports := imports.New()
	imports.Add(annotationImport)
	imports.Add(g.Packages.Models(api.InHttp.InVersion).PackageStar)
	imports.Add(g.Packages.Version(api.InHttp.InVersion).ServicesApi(api).PackageStar)
	imports.Add(g.Types.Imports()...)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`@%s`, annotation)
	w.Line(`class %s : %s {`, serviceImplName(api), serviceInterfaceName(api))
	for _, operation := range api.Operations {
		w.Line(`  override fun %s {`, operationSignature(g.Types, &operation))
		w.Line(`    throw UnsupportedOperationException("Implementation has not added yet")`)
		w.Line(`  }`)
	}
	w.Line(`}`)

	return &generator.CodeFile{
		Path:    g.Packages.Version(api.InHttp.InVersion).ServicesImpl.GetPath(fmt.Sprintf("%s.kt", serviceImplName(api))),
		Content: w.String(),
	}
}
