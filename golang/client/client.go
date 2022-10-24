package client

import (
	"generator"
	"spec"
)

func GenerateClient(specification *spec.Spec, moduleName string, generatePath string) *generator.Sources {
	sources := generator.NewSources()

	modules := NewModules(moduleName, generatePath, specification)
	generator := NewGenerator(modules)

	sources.AddGeneratedAll(generator.AllStaticFiles())

	sources.AddGenerated(generator.ErrorModels(specification.HttpErrors))
	sources.AddGeneratedAll(generator.Errors(&specification.HttpErrors.Responses))

	for _, version := range specification.Versions {
		sources.AddGenerated(generator.Models(&version))
		sources.AddGeneratedAll(generator.Clients(&version))
	}
	return sources
}
