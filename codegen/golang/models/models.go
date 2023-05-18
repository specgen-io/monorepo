package models

import (
	"generator"
	"spec"
)

func GenerateModels(specification *spec.Spec, jsonmode string, moduleName string, generatePath string) *generator.Sources {
	sources := generator.NewSources()

	modules := NewModules(moduleName, generatePath, specification)
	generator := NewGenerator(jsonmode, modules)

	sources.AddGenerated(generator.EnumsHelperFunctions())

	for _, version := range specification.Versions {
		sources.AddGeneratedAll(generator.Models(&version))
	}
	return sources
}
