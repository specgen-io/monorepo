package models

import (
	"generator"
	"java/packages"
	"spec"
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
