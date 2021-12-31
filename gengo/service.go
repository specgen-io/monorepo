package gengo

import (
	"github.com/specgen-io/specgen/v2/genopenapi"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateService(specification *spec.Spec, moduleName string, swaggerPath string, generatePath string, servicesPath string) *sources.Sources {
	sources := sources.NewSources()

	rootModule := Module(moduleName, generatePath)

	emptyModule := rootModule.Submodule("empty")
	sources.AddGenerated(generateEmpty(emptyModule))
	sources.AddGenerated(generateSpecRouting(specification, rootModule))

	for _, version := range specification.Versions {
		versionModule := rootModule.Submodule(version.Version.FlatCase())
		modelsModule := versionModule.Submodule(modelsPackage)

		sources.AddGenerated(generateParamsParser(versionModule))
		sources.AddGeneratedAll(generateRoutings(&version, versionModule, modelsModule))
		sources.AddGeneratedAll(generateServicesInterfaces(&version, versionModule, modelsModule, emptyModule))
		sources.AddGeneratedAll(generateVersionModels(&version, modelsModule))
	}

	if swaggerPath != "" {
		sources.AddGenerated(genopenapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		for _, version := range specification.Versions {
			versionPath := createPath(generatePath, version.Version.FlatCase())
			servicesVersionPath := createPath(servicesPath, version.Version.FlatCase())
			sources.AddScaffoldedAll(generateServicesImplementations(moduleName, versionPath, &version, servicesVersionPath))
		}
	}

	return sources
}
