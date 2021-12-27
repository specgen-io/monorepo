package gents

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateModels(specification *spec.Spec, validation string, generatePath string) *gen.Sources {
	sources := gen.NewSources()
	module := Module(generatePath)
	validationModule := module.Submodule(validation)
	validationFile := generateValidation(validation, validationModule)
	sources.AddGenerated(validationFile)
	for _, version := range specification.Versions {
		versionModule := module.Submodule(version.Version.FlatCase())
		modelsModule := versionModule.Submodule("models")
		sources.AddGenerated(generateVersionModels(&version, validation, validationModule, modelsModule))
	}
	return sources
}

func generateVersionModels(version *spec.Version, validation string, validationModule module, module module) *gen.TextFile {
	if validation == Superstruct {
		return generateSuperstructVersionModels(version, validationModule, module)
	}
	if validation == IoTs {
		return generateIoTsVersionModels(version, validationModule, module)
	}
	panic(fmt.Sprintf("Unknown validation: %s", validation))
}
