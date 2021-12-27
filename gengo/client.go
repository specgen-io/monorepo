package gengo

import (
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateGoClient(specification *spec.Spec, moduleName string, generatePath string) error {
	generatedFiles := []gen.TextFile{}

	rootModule := Module(moduleName, generatePath)

	emptyModule := rootModule.Submodule("empty")
	generatedFiles = append(generatedFiles, *generateEmpty(emptyModule))

	for _, version := range specification.Versions {
		versionModule := rootModule.Submodule(version.Version.FlatCase())
		modelsModule := versionModule.Submodule(modelsPackage)

		generatedFiles = append(generatedFiles, generateVersionModels(&version, modelsModule)...)
		generatedFiles = append(generatedFiles, generateClientsImplementations(&version, versionModule, modelsModule, emptyModule)...)
	}
	err := gen.WriteFiles(generatedFiles, true)
	return err
}
