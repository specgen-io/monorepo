package golang

import (
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateClient(specification *spec.Spec, moduleName string, generatePath string) *sources.Sources {
	sources := sources.NewSources()

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
