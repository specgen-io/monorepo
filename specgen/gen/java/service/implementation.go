package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen/java/imports"
	"github.com/specgen-io/specgen/v2/gen/java/packages"
	"github.com/specgen-io/specgen/v2/gen/java/responses"
	"github.com/specgen-io/specgen/v2/gen/java/writer"
	"github.com/specgen-io/specgen/v2/generator"
	"github.com/specgen-io/specgen/v2/spec"
)

func (g *Generator) ServicesImplementations(version *spec.Version, thePackage, modelsVersionPackage, servicesVersionPackage packages.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, api := range version.Http.Apis {
		serviceVersionSubpackage := servicesVersionPackage.Subpackage(api.Name.SnakeCase())
		files = append(files, *g.serviceImplementation(&api, thePackage, modelsVersionPackage, serviceVersionSubpackage))
	}
	return files
}

func (g *Generator) serviceImplementation(api *spec.Api, thePackage, modelsVersionPackage, serviceVersionSubpackage packages.Module) *generator.CodeFile {
	w := writer.NewJavaWriter()
	w.Line(`package %s;`, thePackage.PackageName)
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
	w.Line(`public class %s implements %s {`, serviceImplName(api), serviceInterfaceName(api))
	for _, operation := range api.Operations {
		w.Line(`  @Override`)
		w.Line(`  public %s {`, responses.Signature(g.Types, &operation))
		w.Line(`    throw new UnsupportedOperationException("Implementation has not added yet");`)
		w.Line(`  }`)
	}
	w.Line(`}`)

	return &generator.CodeFile{
		Path:    thePackage.GetPath(fmt.Sprintf("%s.java", serviceImplName(api))),
		Content: w.String(),
	}
}
