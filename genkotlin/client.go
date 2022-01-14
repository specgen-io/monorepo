package genkotlin

import (
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateClient(specification *spec.Spec, packageName string, generatePath string) *sources.Sources {
	sources := sources.NewSources()

	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}
	mainPackage := Package(generatePath, packageName)

	sources.AddGenerated(generateClientException(mainPackage))

	utilsPackage := mainPackage.Subpackage("utils")
	sources.AddGeneratedAll(generateUtils(utilsPackage))

	jsonPackage := mainPackage.Subpackage("json")
	sources.AddGenerated(generateJson(jsonPackage))

	for _, version := range specification.Versions {
		versionPackage := mainPackage.Subpackage(version.Version.FlatCase())

		modelsVersionPackage := versionPackage.Subpackage("models")
		sources.AddGenerated(generateVersionModels(&version, modelsVersionPackage))

		clientVersionPackage := versionPackage.Subpackage("clients")
		sources.AddGeneratedAll(generateClientsImplementations(&version, clientVersionPackage, modelsVersionPackage, jsonPackage, utilsPackage, mainPackage))
	}

	return sources
}
