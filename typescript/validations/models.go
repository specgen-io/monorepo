package validations

import (
	"generator"
	"spec"
	"typescript/modules"
)

func GenerateModels(specification *spec.Spec, validation string, generatePath string) *generator.Sources {
	sources := generator.NewSources()

	generator := New(validation)

	module := modules.New(generatePath)
	validationModule := module.Submodule(validation)
	validationFile := generator.SetupLibrary(validationModule)
	sources.AddGenerated(validationFile)
	for _, version := range specification.Versions {
		versionModule := module.Submodule(version.Name.FlatCase())
		modelsModule := versionModule.Submodule("models")
		sources.AddGenerated(generator.Models(version.ResolvedModels, validationModule, modelsModule))
	}
	return sources
}
