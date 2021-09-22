package gengo

import (
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/genopenapi"
)

func GenerateService(moduleName string, serviceFile string, swaggerPath string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}
	generatedFiles := []gen.TextFile{}
	implFiles := []gen.TextFile{}

	generatedFiles = append(generatedFiles, *generateSpecRouting(specification, moduleName, generatePath))

	for _, version := range specification.Versions {
		versionPath := createPath(generatePath, version.Version.FlatCase())

		generatedFiles = append(generatedFiles, *generateParamsParser(versionPath))
		generatedFiles = append(generatedFiles, *generateRouting(&version, moduleName, versionPath))
		generatedFiles = append(generatedFiles, generateServicesInterfaces(&version, moduleName, versionPath)...)
		generatedFiles = append(generatedFiles, generateVersionModels(&version, createPath(versionPath, modelsPackage))...)

		servicesPackageName := "services"
		servicesVersionPath := createPath(generatePath, servicesPackageName, version.Version.FlatCase())
		implFiles = append(implFiles, generateServicesImplementations(moduleName, versionPath, &version, servicesVersionPath)...)
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
