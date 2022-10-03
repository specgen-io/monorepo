package models

import (
	"generator"
	"spec"
)

func Generate(specification *spec.Spec, jsonlib string, packageName string, generatePath string) *generator.Sources {
	sources := generator.NewSources()

	thepackages := NewPackages(packageName, generatePath, specification)
	generator := NewGenerator(jsonlib, thepackages)

	sources.AddGeneratedAll(generator.SetupLibrary())
	for _, version := range specification.Versions {
		sources.AddGeneratedAll(generator.Models(&version))
	}

	return sources
}
