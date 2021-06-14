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
	for _, version := range specification.Versions {
		versionPath := filepath.Join(generatePath, versionedFolder(version.Version, "spec"))
		versionPackageName := versionedPackage(version.Version, "spec")

		generatedFiles = append(generatedFiles, *generateParamsParser(versionPackageName, filepath.Join(versionPath, "parsing.go")))
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

func versionedFolder(version spec.Name, folder string) string {
	if version.Source != "" {
		return fmt.Sprintf(`%s_%s`, folder, version.FlatCase())
	}
	return folder
}

func versionedPackage(version spec.Name, packageName string) string {
	if version.Source != "" {
		return fmt.Sprintf("%s_%s", packageName, version.FlatCase())
	}
	return packageName
}
