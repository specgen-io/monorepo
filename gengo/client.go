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
	for _, version := range specification.Versions {
		versionPath := filepath.Join(generatePath, versionedFolder(version.Version, "spec"))
		versionPackageName := versionedPackage(version.Version, "spec")

		generatedFiles = append(generatedFiles, *generateConverter(versionPackageName, filepath.Join(versionPath, "converter.go")))
		generatedFiles = append(generatedFiles, generateVersionModels(&version, versionPackageName, versionPath)...)
		generatedFiles = append(generatedFiles, generateClientsImplementations(&version, versionPackageName, versionPath)...)
	}
	err = gen.WriteFiles(generatedFiles, true)
	return err
}
