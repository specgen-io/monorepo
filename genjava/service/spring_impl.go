package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/genjava/imports"
	"github.com/specgen-io/specgen/v2/genjava/packages"
	"github.com/specgen-io/specgen/v2/genjava/responses"
	"github.com/specgen-io/specgen/v2/genjava/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func (g *Generator) ServicesImplementations(version *spec.Version, thePackage packages.Module, modelsVersionPackage packages.Module, servicesVersionPackage packages.Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		serviceVersionSubpackage := servicesVersionPackage.Subpackage(api.Name.SnakeCase())
		files = append(files, *g.serviceImplementation(&api, thePackage, modelsVersionPackage, serviceVersionSubpackage))
	}
	return files
}

func (g *Generator) serviceImplementation(api *spec.Api, thePackage packages.Module, modelsVersionPackage packages.Module, serviceVersionSubpackage packages.Module) *sources.CodeFile {
	w := writer.NewJavaWriter()
	w.Line(`package %s;`, thePackage.PackageName)
	w.EmptyLine()
	imports := imports.New()
	imports.Add(g.Types.Imports()...)
	imports.Add(`org.springframework.stereotype.Service`)
	imports.Add(modelsVersionPackage.PackageStar)
	imports.Add(serviceVersionSubpackage.PackageStar)
	imports.Write(w)
	w.EmptyLine()
	w.Line(`@Service("%s")`, versionServiceName(serviceName(api), api.Apis.Version))
	w.Line(`public class %s implements %s {`, serviceImplName(api), serviceInterfaceName(api))
	for _, operation := range api.Operations {
		w.Line(`  @Override`)
		w.Line(`  public %s {`, responses.Signature(g.Types, &operation))
		w.Line(`    throw new UnsupportedOperationException("Implementation has not added yet");`)
		w.Line(`  }`)
	}
	w.Line(`}`)

	return &sources.CodeFile{
		Path:    thePackage.GetPath(fmt.Sprintf("%s.java", serviceImplName(api))),
		Content: w.String(),
	}
}
