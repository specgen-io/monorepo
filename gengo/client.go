package gengo

import (
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
)

func GenerateGoClient(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}
	generatedFiles := []gen.TextFile{}

	packageName := "spec"
	for _, version := range specification.Versions {
		versionPackageName := versionedPackage(version.Version, packageName)
		generatedPath := filepath.Join(generatePath, packageName)
		versionPath := versionedFolder(version.Version, generatedPath)

		generatedFiles = append(generatedFiles, *generateConverter(versionPackageName, filepath.Join(versionPath, "converter.go")))
		generatedFiles = append(generatedFiles, generateVersionModels(&version, versionPackageName, versionPath)...)
		generatedFiles = append(generatedFiles, generateClientsImplementations(&version, versionPackageName, versionPath)...)
		generatedFiles = append(generatedFiles, *generateServicesInterfaces(&version, versionPackageName, versionPath, "responses.go"))
	}
	err = gen.WriteFiles(generatedFiles, true)
	return err
}
