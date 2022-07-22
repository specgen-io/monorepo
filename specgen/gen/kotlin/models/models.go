package models

import (
	"github.com/specgen-io/specgen/spec/v2"
	"github.com/specgen-io/specgen/specgen/v2/gen/kotlin/modules"
	"github.com/specgen-io/specgen/specgen/v2/generator"
)

func Generate(specification *spec.Spec, jsonlib string, packageName string, generatePath string) *generator.Sources {
	sources := generator.NewSources()

	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}

	mainPackage := modules.Package(generatePath, packageName)

	generator := NewGenerator(jsonlib)

	jsonPackage := mainPackage.Subpackage("json")

	for _, version := range specification.Versions {
		versionPackage := mainPackage.Subpackage(version.Version.FlatCase())

		modelsVersionPackage := versionPackage.Subpackage("models")
		sources.AddGeneratedAll(generator.VersionModels(&version, modelsVersionPackage, jsonPackage))
	}

	sources.AddGeneratedAll(generator.SetupLibrary(jsonPackage))

	return sources
}
