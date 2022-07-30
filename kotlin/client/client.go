package client

import (
	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/kotlin/v2/modules"
	"github.com/specgen-io/specgen/spec/v2"
)

func Generate(specification *spec.Spec, jsonlib, client, packageName, generatePath string) *generator.Sources {
	sources := generator.NewSources()

	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}

	mainPackage := modules.Package(generatePath, packageName)

	generator := NewGenerator(jsonlib, client)

	jsonPackage := mainPackage.Subpackage("json")

	for _, version := range specification.Versions {
		versionPackage := mainPackage.Subpackage(version.Version.FlatCase())

		modelsVersionPackage := versionPackage.Subpackage("models")
		sources.AddGeneratedAll(generator.Models.VersionModels(&version, modelsVersionPackage, jsonPackage))

		clientVersionPackage := versionPackage.Subpackage("clients")

		sources.AddGeneratedAll(generator.Client.ClientImplementation(&version, clientVersionPackage, modelsVersionPackage, jsonPackage, mainPackage))
	}

	sources.AddGeneratedAll(generator.Models.SetupLibrary(jsonPackage))

	return sources
}
