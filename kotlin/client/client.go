package client

import (
	"generator"
	"spec"
)

func Generate(specification *spec.Spec, jsonlib, client, packageName, generatePath string) *generator.Sources {
	sources := generator.NewSources()

	packages := NewPackages(packageName, generatePath, specification)
	generator := NewGenerator(jsonlib, client, packages)

	sources.AddGeneratedAll(generator.Models.ErrorModels(specification.HttpErrors))

	for _, version := range specification.Versions {
		sources.AddGeneratedAll(generator.Models.Models(&version))
		sources.AddGeneratedAll(generator.Client.Clients(&version))
	}

	sources.AddGeneratedAll(generator.Models.SetupLibrary())

	return sources
}
