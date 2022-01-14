package genjava

import (
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateClient(specification *spec.Spec, jsonlib string, packageName string, generatePath string) *sources.Sources {
	newSources := sources.NewSources()

	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}

	mainPackage := Package(generatePath, packageName)

	generator := NewGenerator(jsonlib)

	newSources.AddGenerated(ClientException(mainPackage))

	utilsPackage := mainPackage.Subpackage("utils")
	newSources.AddGeneratedAll(generateUtils(utilsPackage))

	jsonPackage := mainPackage.Subpackage("json")
	newSources.AddGeneratedAll(generator.Models.SetupLibrary(jsonPackage))

	for _, version := range specification.Versions {
		versionPackage := mainPackage.Subpackage(version.Version.FlatCase())

		modelsVersionPackage := versionPackage.Subpackage("models")
		newSources.AddGeneratedAll(generator.Models.VersionModels(&version, modelsVersionPackage))

		clientVersionPackage := versionPackage.Subpackage("clients")
		newSources.AddGeneratedAll(generator.Clients(&version, clientVersionPackage, modelsVersionPackage, jsonPackage, utilsPackage, mainPackage))
	}

	return newSources
}
