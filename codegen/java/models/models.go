package models

import (
	"generator"
	"spec"
)

func Generate(specification *spec.Spec, jsonlib string, packageName string, generatePath string) *generator.Sources {
	sources := generator.NewSources()

	packages := NewPackages(packageName, generatePath, specification)
	generator := NewGenerator(jsonlib, packages)

	for _, version := range specification.Versions {
		sources.AddGeneratedAll(generator.Models(&version))
	}
	sources.AddGeneratedAll(generator.JsonHelpers())

	return sources
}
