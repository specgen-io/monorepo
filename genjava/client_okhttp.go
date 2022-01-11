package genjava

import (
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateClient(specification *spec.Spec, jsonlib string, packageName string, generatePath string) *sources.Sources {
	newSources := sources.NewSources()

	mainPackage := Package(generatePath, packageName)

	newSources.AddGenerated(generateClientException(mainPackage))

	utilsPackage := mainPackage.Subpackage("utils")
	newSources.AddGeneratedAll(generateUtils(utilsPackage))

	modelsPackage := mainPackage.Subpackage("models")
	newSources.AddGenerated(generateJson(modelsPackage))

	for _, version := range specification.Versions {
		versionPackage := mainPackage.Subpackage(version.Version.FlatCase())

		modelsVersionPackage := versionPackage.Subpackage("models")
		newSources.AddGeneratedAll(generateVersionModels(&version, modelsVersionPackage, jsonlib))

		clientVersionPackage := versionPackage.Subpackage("clients")
		newSources.AddGeneratedAll(generateClientsImplementations(&version, clientVersionPackage, modelsVersionPackage, modelsPackage, utilsPackage, mainPackage, jsonlib))
	}

	return newSources
}
