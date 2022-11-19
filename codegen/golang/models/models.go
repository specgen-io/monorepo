package models

import (
	"generator"
	"spec"
)

func GenerateModels(specification *spec.Spec, moduleName string, generatePath string) *generator.Sources {
	sources := generator.NewSources()

	modules := NewModules(moduleName, generatePath, specification)
	generator := NewGenerator(modules)

	sources.AddGenerated(generator.EnumsHelperFunctions())

	for _, version := range specification.Versions {
		sources.AddGeneratedAll(generator.Models(&version))
	}
	return sources
}
