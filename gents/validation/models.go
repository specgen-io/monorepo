package validation

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gents/modules"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateModels(specification *spec.Spec, validation string, generatePath string) *sources.Sources {
	sources := sources.NewSources()
	module := modules.NewModule(generatePath)
	validationModule := module.Submodule(validation)
	validationFile := GenerateValidation(validation, validationModule)
	sources.AddGenerated(validationFile)
	for _, version := range specification.Versions {
		versionModule := module.Submodule(version.Version.FlatCase())
		modelsModule := versionModule.Submodule("models")
		sources.AddGenerated(GenerateVersionModels(&version, validation, validationModule, modelsModule))
	}
	return sources
}

func GenerateVersionModels(version *spec.Version, validation string, validationModule modules.Module, module modules.Module) *sources.CodeFile {
	if validation == Superstruct {
		return generateSuperstructVersionModels(version, validationModule, module)
	}
	if validation == IoTs {
		return generateIoTsVersionModels(version, validationModule, module)
	}
	panic(fmt.Sprintf("Unknown validation: %s", validation))
}
