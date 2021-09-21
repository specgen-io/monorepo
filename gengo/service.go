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
	generatedFiles = append(generatedFiles, *generateRoutes(moduleName, specification, moduleName, filepath.Join(generatePath, packageName)))

	for _, version := range specification.Versions {
		versionPath := versionedFolder(version.Version, filepath.Join(generatePath, packageName))

		generatedFiles = append(generatedFiles, *generateParamsParser(versionPath))
		generatedFiles = append(generatedFiles, *generateRouting(&version, moduleName, versionPath))
		generatedFiles = append(generatedFiles, generateServicesInterfaces(&version, moduleName, versionPath)...)
		generatedFiles = append(generatedFiles, generateVersionModels(&version, versionPath)...)

		servicesPackageName := "services"
		servicesVersionPackageName := versionedPackage(version.Version, servicesPackageName)
		servicesVersionPath := versionedFolder(version.Version, filepath.Join(generatePath, servicesPackageName))
		implFiles = append(implFiles, generateServicesImplementations(moduleName, packageName, &version, servicesVersionPackageName, servicesVersionPath)...)
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
