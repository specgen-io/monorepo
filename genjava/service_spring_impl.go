package genjava

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func generateServicesImplementations(version *spec.Version, thePackage Module, modelsVersionPackage Module, servicesVersionPackage Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, api := range version.Http.Apis {
		serviceVersionSubpackage := servicesVersionPackage.Subpackage(api.Name.SnakeCase())
		files = append(files, *generateServiceImplementation(&api, thePackage, modelsVersionPackage, serviceVersionSubpackage))
	}
	return files
}

func generateServiceImplementation(api *spec.Api, thePackage Module, modelsVersionPackage Module, serviceVersionSubpackage Module) *sources.CodeFile {
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
	w.Line(`@Service("%s")`, versionServiceName(serviceName(api), api.Apis.Version))
	w.Line(`public class %s implements %s {`, serviceImplName(api), serviceInterfaceName(api))
	for _, operation := range api.Operations {
		w.Line(`  @Override`)
		w.Line(`  public %s {`, generateResponsesSignatures(&operation))
		w.Line(`    throw new UnsupportedOperationException("Implementation has not added yet");`)
		w.Line(`  }`)
	}
	w.Line(`}`)

	return &sources.CodeFile{
		Path:    thePackage.GetPath(fmt.Sprintf("%s.java", serviceImplName(api))),
		Content: w.String(),
	}
}
