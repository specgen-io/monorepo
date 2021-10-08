package gents

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
)

func GenerateModels(specification *spec.Spec, validation string, generatePath string) error {
	module := Module(generatePath)
	validationModule := module.Submodule(validation)
	files := generateModels(specification, validation, validationModule, module)
	return gen.WriteFiles(files, true)
}

func generateModels(specification *spec.Spec, validation string, validationModule module, module module) []gen.TextFile {
	files := []gen.TextFile{}
	for _, version := range specification.Versions {
		versionModule := module.Submodule(version.Version.FlatCase())
		modelsModule := versionModule.Submodule("models")
		files = append(files, *generateVersionModels(&version, validation, validationModule, modelsModule))
	}
	return files
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