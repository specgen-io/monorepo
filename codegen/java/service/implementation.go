package service

import (
	"generator"
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
	w := writer.New(g.Packages.ServicesImpl(api.InHttp.InVersion), serviceImplName(api))
	annotationImport, annotation := g.ServiceImplAnnotation(api)
	w.Imports.Add(annotationImport)
	w.Imports.Add(g.FilesImports()...)
	w.Imports.Star(g.Packages.Models(api.InHttp.InVersion))
	w.Imports.Star(g.Packages.ServicesApi(api))
	w.Imports.Add(g.Types.Imports()...)
	w.Line(`@%s`, annotation)
	w.Line(`public class [[.ClassName]] implements %s {`, serviceInterfaceName(api))
	for _, operation := range api.Operations {
		w.Line(`  @Override`)
		w.Line(`  public %s {`, operationSignature(g.Types, &operation))
		w.Line(`    throw new UnsupportedOperationException("Implementation has not added yet");`)
		w.Line(`  }`)
	}
	w.Line(`}`)
	return w.ToCodeFile()
}
