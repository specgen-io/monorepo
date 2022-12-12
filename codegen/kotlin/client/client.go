package client

import (
	"generator"
	"spec"
)

func Generate(specification *spec.Spec, jsonlib, client, packageName, generatePath string) *generator.Sources {
	sources := generator.NewSources()

	packages := NewPackages(packageName, generatePath, specification)
	generator := NewGenerator(jsonlib, client, packages)

	sources.AddGeneratedAll(generator.ErrorModels(specification.HttpErrors))
	sources.AddGeneratedAll(generator.Exceptions(&specification.HttpErrors.Responses))
	sources.AddGeneratedAll(generator.Utils())

	for _, version := range specification.Versions {
		sources.AddGeneratedAll(generator.Models(&version))
		sources.AddGeneratedAll(generator.Clients(&version))
	}
	sources.AddGeneratedAll(generator.JsonHelpers())

	return sources
}
