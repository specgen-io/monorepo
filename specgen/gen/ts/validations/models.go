package validations

import (
	"github.com/specgen-io/specgen/spec/v2"
	"github.com/specgen-io/specgen/specgen/v2/gen/ts/modules"
	"github.com/specgen-io/specgen/specgen/v2/generator"
)

func GenerateModels(specification *spec.Spec, validation string, generatePath string) *generator.Sources {
	sources := generator.NewSources()

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
