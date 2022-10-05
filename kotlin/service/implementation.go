package service

import (
	"fmt"

	"generator"
	"kotlin/imports"
	"kotlin/packages"
	"kotlin/writer"
	"spec"
)

func (g *Generator) ServicesImplementations(version *spec.Version, thePackage, modelsVersionPackage, servicesVersionPackage packages.Package) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		serviceVersionSubpackage := servicesVersionPackage.Subpackage(api.Name.SnakeCase())
		files = append(files, *g.serviceImplementation(&api, thePackage, modelsVersionPackage, serviceVersionSubpackage))
	}
	return files
}

func (g *Generator) serviceImplementation(api *spec.Api, thePackage, modelsVersionPackage, serviceVersionSubpackage packages.Package) *generator.CodeFile {
	w := writer.NewKotlinWriter()
	w.Line(`package %s`, thePackage.PackageName)
	w.EmptyLine()
	annotationImport, annotation := g.Server.ServiceImplAnnotation(api)
	imports := imports.New()
	imports.Add(annotationImport)
	imports.Add(modelsVersionPackage.PackageStar)
	imports.Add(serviceVersionSubpackage.PackageStar)
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
		Path:    thePackage.GetPath(fmt.Sprintf("%s.kt", serviceImplName(api))),
		Content: w.String(),
	}
}
