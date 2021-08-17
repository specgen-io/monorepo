package gengo

import (
	"fmt"
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

	packageName := fmt.Sprintf("%s_client", specification.Name.SnakeCase())
	generatePath = filepath.Join(generatePath, packageName)
	for _, version := range specification.Versions {
		versionPackageName := versionedPackage(version.Version, packageName)
		versionPath := versionedFolder(version.Version, generatePath)

		generatedFiles = append(generatedFiles, *generateConverter(versionPackageName, filepath.Join(versionPath, "converter.go")))
		generatedFiles = append(generatedFiles, generateVersionModels(&version, versionPackageName, versionPath)...)
		generatedFiles = append(generatedFiles, generateClientsImplementations(&version, versionPackageName, versionPath)...)
	}
	err = gen.WriteFiles(generatedFiles, true)
	return err
}
