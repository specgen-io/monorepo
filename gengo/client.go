package gengo

import (
	"github.com/specgen-io/specgen/v2/spec"
	"github.com/specgen-io/specgen/v2/gen"
)

func GenerateGoClient(specification *spec.Spec, moduleName string, generatePath string) error {
	generatedFiles := []gen.TextFile{}

	for _, version := range specification.Versions {
		versionModule := Module(moduleName, createPath(generatePath, version.Version.FlatCase()))
		modelsModule := Module(moduleName, createPath(generatePath, version.Version.FlatCase(), modelsPackage))

		generatedFiles = append(generatedFiles, generateVersionModels(&version, modelsModule)...)
		generatedFiles = append(generatedFiles, generateClientsImplementations(&version, versionModule, modelsModule)...)
	}
	err := gen.WriteFiles(generatedFiles, true)
	return err
}