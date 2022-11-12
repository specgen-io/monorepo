package service

import (
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
	w := writer.New(g.Packages.ServicesImpl(api.InHttp.InVersion), serviceImplName(api))
	annotationImport, annotation := g.ServiceImplAnnotation(api)
	imports := imports.New()
	imports.Add(annotationImport)
	imports.Add(g.Packages.Models(api.InHttp.InVersion).PackageStar)
	imports.Add(g.Packages.ServicesApi(api).PackageStar)
	imports.Add(g.Types.Imports()...)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`@%s`, annotation)
	w.Line(`class [[.ClassName]] : %s {`, serviceInterfaceName(api))
	for _, operation := range api.Operations {
		w.Line(`  override fun %s {`, operationSignature(g.Types, &operation))
		w.Line(`    throw UnsupportedOperationException("Implementation has not added yet")`)
		w.Line(`  }`)
	}
	w.Line(`}`)
	return w.ToCodeFile()
}
