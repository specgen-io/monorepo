package client

import (
	"generator"
	"spec"
)

func Generate(specification *spec.Spec, jsonlib string, packageName string, generatePath string) *generator.Sources {
	sources := generator.NewSources()

	packages := NewPackages(packageName, generatePath, specification)
	generator := NewGenerator(jsonlib, packages)

	sources.AddGenerated(generator.Exceptions())
	sources.AddGeneratedAll(generator.Utils())
	sources.AddGeneratedAll(generator.ErrorModels(specification.HttpErrors))

	for _, version := range specification.Versions {
		sources.AddGeneratedAll(generator.Models(&version))
		sources.AddGeneratedAll(generator.Clients(&version))
	}

	sources.AddGeneratedAll(generator.SetupLibrary())

	return sources
}
