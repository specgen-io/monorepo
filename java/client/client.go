package client

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

	sources.AddGenerated(clientException(mainPackage))

	utilsPackage := mainPackage.Subpackage("utils")
	sources.AddGeneratedAll(generateUtils(utilsPackage))

	jsonPackage := mainPackage.Subpackage("json")

	for _, version := range specification.Versions {
		versionPackage := mainPackage.Subpackage(version.Name.FlatCase())

		modelsVersionPackage := versionPackage.Subpackage("models")
		sources.AddGeneratedAll(generator.Models.VersionModels(&version, modelsVersionPackage, jsonPackage))

		clientVersionPackage := versionPackage.Subpackage("clients")
		sources.AddGeneratedAll(generator.Clients(&version, clientVersionPackage, modelsVersionPackage, jsonPackage, utilsPackage, mainPackage))
	}

	sources.AddGeneratedAll(generator.Models.SetupLibrary(jsonPackage))

	return sources
}
