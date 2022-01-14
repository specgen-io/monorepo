package genjava

import (
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateClient(specification *spec.Spec, jsonlib string, packageName string, generatePath string) *sources.Sources {
	newSources := sources.NewSources()

	mainPackage := Package(generatePath, packageName)

	newSources.AddGenerated(ClientException(mainPackage))

	utilsPackage := mainPackage.Subpackage("utils")
	newSources.AddGeneratedAll(generateUtils(utilsPackage))

	modelsPackage := mainPackage.Subpackage("models")

	generator := NewGenerator(jsonlib)

	newSources.AddGeneratedAll(generator.Models.SetupLibrary(modelsPackage))

	for _, version := range specification.Versions {
		versionPackage := mainPackage.Subpackage(version.Version.FlatCase())

		modelsVersionPackage := versionPackage.Subpackage("models")
		newSources.AddGeneratedAll(generator.Models.VersionModels(&version, modelsVersionPackage))

		clientVersionPackage := versionPackage.Subpackage("clients")
		newSources.AddGeneratedAll(generator.Clients(&version, clientVersionPackage, modelsVersionPackage, modelsPackage, utilsPackage, mainPackage))
	}

	return newSources
}
