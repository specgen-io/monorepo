package models

import (
	"generator"
	"spec"
)

func GenerateModels(specification *spec.Spec, moduleName string, generatePath string) *generator.Sources {
	sources := generator.NewSources()

	modules := NewModules(moduleName, generatePath, specification)
	generator := NewGenerator(modules)

	sources.AddGenerated(generator.GenerateEnumsHelperFunctions())

	for _, version := range specification.Versions {
		sources.AddGenerated(generator.GenerateVersionModels(&version))
	}
	return sources
}
