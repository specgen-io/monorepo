package genjava

import (
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateClient(specification *spec.Spec, packageName string, generatePath string) *gen.Sources {
	sources := gen.NewSources()

	mainPackage := Package(generatePath, packageName)

	sources.AddGenerated(generateClientException(mainPackage))

	utilsPackage := mainPackage.Subpackage("utils")
	sources.AddGeneratedAll(generateUtils(utilsPackage))

	modelsPackage := mainPackage.Subpackage("models")
	sources.AddGenerated(generateJson(modelsPackage))

	for _, version := range specification.Versions {
		versionPackage := mainPackage.Subpackage(version.Version.FlatCase())

		modelsVersionPackage := versionPackage.Subpackage("models")
		sources.AddGeneratedAll(generateVersionModels(&version, modelsVersionPackage))

		clientVersionPackage := versionPackage.Subpackage("clients")
		sources.AddGeneratedAll(generateClientsImplementations(&version, clientVersionPackage, modelsVersionPackage, modelsPackage, utilsPackage, mainPackage))
	}

	return sources
}
