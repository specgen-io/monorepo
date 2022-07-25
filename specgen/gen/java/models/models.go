package models

import (
	"github.com/specgen-io/specgen/spec/v2"
	"github.com/specgen-io/specgen/v2/gen/java/packages"
	"github.com/specgen-io/specgen/v2/generator"
)

func Generate(specification *spec.Spec, jsonlib string, packageName string, generatePath string) *generator.Sources {
	sources := generator.NewSources()

	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}

	mainPackage := packages.New(generatePath, packageName)

	generator := NewGenerator(jsonlib)

	jsonPackage := mainPackage.Subpackage("json")

	for _, version := range specification.Versions {
		versionPackage := mainPackage.Subpackage(version.Version.FlatCase())

		versionModelsPackage := versionPackage.Subpackage("models")
		sources.AddGeneratedAll(generator.VersionModels(&version, versionModelsPackage, jsonPackage))
	}

	sources.AddGeneratedAll(generator.SetupLibrary(jsonPackage))

	return sources
}
