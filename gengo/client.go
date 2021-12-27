package gengo

import (
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateGoClient(specification *spec.Spec, moduleName string, generatePath string) *gen.Sources {
	sources := gen.NewSources()

	rootModule := Module(moduleName, generatePath)

	emptyModule := rootModule.Submodule("empty")
	sources.AddGenerated(generateEmpty(emptyModule))

	for _, version := range specification.Versions {
		versionModule := rootModule.Submodule(version.Version.FlatCase())
		modelsModule := versionModule.Submodule(modelsPackage)

		sources.AddGeneratedAll(generateVersionModels(&version, modelsModule))
		sources.AddGeneratedAll(generateClientsImplementations(&version, versionModule, modelsModule, emptyModule))
	}
	return sources
}
