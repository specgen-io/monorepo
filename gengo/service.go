package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
)

func GenerateService(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}
	generatedFiles := []gen.TextFile{}
	implFiles := []gen.TextFile{}

	packageName := fmt.Sprintf("%s_service", specification.Name.SnakeCase())
	generatePath = filepath.Join(generatePath, packageName)
	for _, version := range specification.Versions {
		versionPackageName := versionedPackage(version.Version, packageName)
		versionPath := versionedFolder(version.Version, generatePath)

		generatedFiles = append(generatedFiles, *generateParamsParser(versionPackageName, filepath.Join(versionPath, "params_parsing.go")))
		generatedFiles = append(generatedFiles, *generateRouting(&version, versionPackageName, versionPath))
		generatedFiles = append(generatedFiles, *generateServicesInterfaces(&version, versionPackageName, versionPath))
		generatedFiles = append(generatedFiles, generateVersionModels(&version, versionPackageName, versionPath)...)

		implFiles = append(implFiles, generateServicesImplementations(&version, versionPackageName, versionPath)...)
	}
	err = gen.WriteFiles(generatedFiles, true)
	if err != nil {
		return err
	}
	err = gen.WriteFiles(implFiles, false)
	return err
}
