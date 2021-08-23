package gengo

import (
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/genopenapi"
	"path/filepath"
)

func GenerateService(moduleName string, serviceFile string, swaggerPath string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}
	generatedFiles := []gen.TextFile{}
	implFiles := []gen.TextFile{}

	packageName := "spec"
	generatedFiles = append(generatedFiles, *generateRoutes(moduleName, specification, packageName, filepath.Join(generatePath, packageName)))

	for _, version := range specification.Versions {
		versionPackageName := versionedPackage(version.Version, packageName)
		versionPath := versionedFolder(version.Version, filepath.Join(generatePath, packageName))

		generatedFiles = append(generatedFiles, *generateParamsParser(versionPackageName, filepath.Join(versionPath, "params_parsing.go")))
		generatedFiles = append(generatedFiles, *generateRouting(&version, versionPackageName, versionPath))
		generatedFiles = append(generatedFiles, *generateServicesInterfaces(&version, versionPackageName, versionPath))
		generatedFiles = append(generatedFiles, generateVersionModels(&version, versionPackageName, versionPath)...)

		implFiles = append(implFiles, generateServicesImplementations(moduleName, packageName, &version, versionedPackage(version.Version, "main"), versionedFolder(version.Version, generatePath))...)
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
