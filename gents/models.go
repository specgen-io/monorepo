package gents

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/spec"
	"github.com/specgen-io/specgen/v2/gen"
)

func GenerateModels(specification *spec.Spec, validation string, generatePath string) error {
	module := Module(generatePath)
	sources := []gen.TextFile{}
	validationModule := module.Submodule(validation)
	validationFile := generateValidation(validation, validationModule)
	sources = append(sources, *validationFile)
	for _, version := range specification.Versions {
		versionModule := module.Submodule(version.Version.FlatCase())
		modelsModule := versionModule.Submodule("models")
		sources = append(sources, *generateVersionModels(&version, validation, validationModule, modelsModule))
	}
	return gen.WriteFiles(sources, true)
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