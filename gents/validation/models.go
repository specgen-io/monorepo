package validation

import (
	"github.com/specgen-io/specgen/v2/gents/modules"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateModels(specification *spec.Spec, validationName string, generatePath string) *sources.Sources {
	sources := sources.NewSources()

	validation := New(validationName)

	module := modules.NewModule(generatePath)
	validationModule := module.Submodule(validationName)
	validationFile := validation.SetupLibrary(validationModule)
	sources.AddGenerated(validationFile)
	for _, version := range specification.Versions {
		versionModule := module.Submodule(version.Version.FlatCase())
		modelsModule := versionModule.Submodule("models")
		sources.AddGenerated(validation.GenerateVersionModels(&version, validationModule, modelsModule))
	}
	return sources
}
