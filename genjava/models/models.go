package models

import (
	"github.com/specgen-io/specgen/v2/genjava/packages"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func Generate(specification *spec.Spec, jsonlib string, packageName string, generatePath string) *sources.Sources {
	sources := sources.NewSources()

	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}

	mainPackage := packages.New(generatePath, packageName)

	generator := NewGenerator(jsonlib)

	jsonPackage := mainPackage.Subpackage("json")
	sources.AddGeneratedAll(generator.SetupLibrary(jsonPackage))

	for _, version := range specification.Versions {
		versionPackage := mainPackage.Subpackage(version.Version.FlatCase())

		versionModelsPackage := versionPackage.Subpackage("models")
		sources.AddGeneratedAll(generator.VersionModels(&version, versionModelsPackage))
	}

	return sources
}
