package service

import (
	"github.com/specgen-io/specgen/v2/gen/golang/models"
	"github.com/specgen-io/specgen/v2/gen/golang/module"
	"github.com/specgen-io/specgen/v2/gen/golang/types"
	"github.com/specgen-io/specgen/v2/gen/openapi"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateService(specification *spec.Spec, moduleName string, swaggerPath string, generatePath string, servicesPath string) *sources.Sources {
	sources := sources.NewSources()

	rootModule := module.New(moduleName, generatePath)

	emptyModule := rootModule.Submodule("empty")
	sources.AddGenerated(types.GenerateEmpty(emptyModule))
	sources.AddGenerated(generateSpecRouting(specification, rootModule))

	for _, version := range specification.Versions {
		versionModule := rootModule.Submodule(version.Version.FlatCase())
		modelsModule := versionModule.Submodule(types.ModelsPackage)

		sources.AddGenerated(generateParamsParser(versionModule, modelsModule))
		sources.AddGeneratedAll(generateRoutings(&version, versionModule, modelsModule))
		sources.AddGeneratedAll(generateServicesInterfaces(&version, versionModule, modelsModule, emptyModule))
		sources.AddGeneratedAll(models.GenerateVersionModels(&version, modelsModule))
	}

	if swaggerPath != "" {
		sources.AddGenerated(openapi.GenerateOpenapi(specification, swaggerPath))
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
