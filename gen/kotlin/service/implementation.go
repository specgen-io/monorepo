package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen/kotlin/imports"
	"github.com/specgen-io/specgen/v2/gen/kotlin/modules"
	"github.com/specgen-io/specgen/v2/gen/kotlin/responses"
	"github.com/specgen-io/specgen/v2/gen/kotlin/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func (g *Generator) ServicesImplementations(version *spec.Version, thePackage, modelsVersionPackage, servicesVersionPackage modules.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		serviceVersionSubpackage := servicesVersionPackage.Subpackage(api.Name.SnakeCase())
		files = append(files, *g.serviceImplementation(&api, thePackage, modelsVersionPackage, serviceVersionSubpackage))
	}
	return files
}

func (g *Generator) serviceImplementation(api *spec.Api, thePackage, modelsVersionPackage, serviceVersionSubpackage modules.Module) *sources.CodeFile {
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
		w.Line(`  override fun %s {`, responses.Signature(g.Types, &operation))
		w.Line(`    throw UnsupportedOperationException("Implementation has not added yet")`)
		w.Line(`  }`)
	}
	w.Line(`}`)

	return &sources.CodeFile{
		Path:    thePackage.GetPath(fmt.Sprintf("%s.kt", serviceImplName(api))),
		Content: w.String(),
	}
}
