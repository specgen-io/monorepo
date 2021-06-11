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
		folder := versionedFolder("spec", &version)
		packageName := folder
		generatedFiles = append(generatedFiles, *generateParamsParser(folder, filepath.Join(generatePath, folder, "parsing.go")))
		generatedFiles = append(generatedFiles, *generateRouting(&version, packageName, generatePath, folder))
		generatedFiles = append(generatedFiles, *generateServicesInterfaces(&version, packageName, generatePath, folder))
		implFiles = append(implFiles, generateServicesImplementations(&version, packageName, generatePath, folder)...)
	}
	err = gen.WriteFiles(generatedFiles, true)
	if err != nil { return err }
	err = gen.WriteFiles(implFiles, false)
	return err
}

func versionedFolder(folder string, version *spec.Version) string {
	if version.Version.Source != "" {
		return fmt.Sprintf(`%s_%s`, folder, version.Version.FlatCase())
	}
	return folder
}