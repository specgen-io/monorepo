package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/genopenapi"
	"path/filepath"
)

func GenerateService(serviceFile string, swaggerPath string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}
	generatedFiles := []gen.TextFile{}
	implFiles := []gen.TextFile{}

	packageName := "spec"
	generatePath = filepath.Join(generatePath, packageName)

	generatedFiles = append(generatedFiles, *generateRoutes(specification, packageName, generatePath))

	for _, version := range specification.Versions {
		versionPackageName := versionedPackage(version.Version, packageName)
		versionPath := versionedFolder(version.Version, generatePath)

		generatedFiles = append(generatedFiles, *generateParamsParser(versionPackageName, filepath.Join(versionPath, "params_parsing.go")))
		generatedFiles = append(generatedFiles, *generateRouting(&version, versionPackageName, versionPath))
		generatedFiles = append(generatedFiles, *generateServicesInterfaces(&version, versionPackageName, versionPath))
		generatedFiles = append(generatedFiles, generateVersionModels(&version, versionPackageName, versionPath)...)

		implFiles = append(implFiles, generateServicesImplementations(&version, versionPackageName, versionPath)...)
	}

	if swaggerPath != "" {
		generatedFiles = append(generatedFiles, *genopenapi.GenerateOpenapi(specification, swaggerPath))
	}

	err = gen.WriteFiles(generatedFiles, true)
	if err != nil {
		return err
	}
	err = gen.WriteFiles(implFiles, false)
	return err
}

func generateRoutes(specification *spec.Spec, packageName string, generatePath string) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", packageName)
	w.EmptyLine()
	w.Line("import (")
	w.Line(`  "github.com/husobee/vestigo"`)
	w.Line(`  "github.com/specgen-io/test-service/go/spec/v2"`)
	w.Line(`)`)
	w.EmptyLine()
	w.Line(`func AddRoutes(router *vestigo.Router) {`)
	for _, version := range specification.Versions {
		for _, api := range version.Http.Apis {
			versionedPackageName := addVersionedPackage(version.Version, packageName)
			w.Line(`  %sAdd%sRoutes(router, &%s%s{})`, versionedPackageName, api.Name.PascalCase(), versionedPackageName, serviceTypeName(&api))
		}
	}
	w.Line(`}`)

	return &gen.TextFile{
		Path:    filepath.Join(generatePath, "routes.go"),
		Content: w.String(),
	}
}

func addVersionedPackage(version spec.Name, packageName string) string {
	if version.Source != "" {
		return fmt.Sprintf(`%s.`, versionedPackage(version, packageName))
	}
	return ""
}
