package validations

import (
	"github.com/specgen-io/specgen/v2/gents/modules"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateModels(specification *spec.Spec, validation string, generatePath string) *sources.Sources {
	sources := sources.NewSources()

	generator := New(validation)

	module := modules.New(generatePath)
	validationModule := module.Submodule(validation)
	validationFile := generator.SetupLibrary(validationModule)
	sources.AddGenerated(validationFile)
	for _, version := range specification.Versions {
		versionModule := module.Submodule(version.Version.FlatCase())
		modelsModule := versionModule.Submodule("models")
		sources.AddGenerated(generator.VersionModels(&version, validationModule, modelsModule))
	}
	return sources
}