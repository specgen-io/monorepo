package genjava

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
)

func generateServicesImplementations(version *spec.Version, thePackage Module, modelsVersionPackage Module, servicesVersionPackage Module) []gen.TextFile {
	files := []gen.TextFile{}
	for _, api := range version.Http.Apis {
		serviceVersionSubpackage := servicesVersionPackage.Subpackage(api.Name.SnakeCase())
		files = append(files, *generateServiceImplementation(version, &api, thePackage, modelsVersionPackage, serviceVersionSubpackage))
	}
	return files
}

func generateServiceImplementation(version *spec.Version, api *spec.Api, thePackage Module, modelsVersionPackage Module, serviceVersionSubpackage Module) *gen.TextFile {
	w := NewJavaWriter()
	w.Line(`package %s;`, thePackage.PackageName)
	w.EmptyLine()
	w.Line(`import java.math.BigDecimal;`)
	w.Line(`import java.time.*;`)
	w.Line(`import java.util.*;`)
	w.EmptyLine()
	w.Line(`import org.springframework.stereotype.Service;`)
	w.EmptyLine()
	w.Line(`import %s;`, modelsVersionPackage.PackageStar)
	w.Line(`import %s;`, serviceVersionSubpackage.PackageStar)
	w.EmptyLine()
	w.Line(`@Service("%s")`, versionServiceName(serviceName(api), version))
	w.Line(`public class %s implements %s {`, serviceName(api), serviceInterfaceName(api))
	for _, operation := range api.Operations {
		w.Line(`  @Override`)
		w.Line(`  public %s {`, generateResponsesSignatures(operation))
		w.Line(`    throw new UnsupportedOperationException("Implementation has not added yet");`)
		w.Line(`  }`)
	}
	w.Line(`}`)

	return &gen.TextFile{
		Path:    thePackage.GetPath(fmt.Sprintf("%s.java", serviceName(api))),
		Content: w.String(),
	}
}