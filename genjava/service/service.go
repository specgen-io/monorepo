package service

import (
	"github.com/specgen-io/specgen/v2/genjava/packages"
	"github.com/specgen-io/specgen/v2/genopenapi"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func Generate(specification *spec.Spec, jsonlib string, packageName string, swaggerPath string, generatePath string, servicesPath string) *sources.Sources {
	newSources := sources.NewSources()

	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}

	mainPackage := packages.New(generatePath, packageName)

	generator := NewGenerator(jsonlib)

	jsonPackage := mainPackage.Subpackage("json")
	newSources.AddGeneratedAll(generator.Models.SetupLibrary(jsonPackage))

	for _, version := range specification.Versions {
		versionPackage := mainPackage.Subpackage(version.Version.FlatCase())

		modelsVersionPackage := versionPackage.Subpackage("models")
		newSources.AddGeneratedAll(generator.Models.VersionModels(&version, modelsVersionPackage))

		serviceVersionPackage := versionPackage.Subpackage("services")
		newSources.AddGeneratedAll(generator.ServicesInterfaces(&version, serviceVersionPackage, modelsVersionPackage))

		controllerVersionPackage := versionPackage.Subpackage("controllers")
		newSources.AddGeneratedAll(generator.ServicesControllers(&version, controllerVersionPackage, jsonPackage, modelsVersionPackage, serviceVersionPackage))
	}

	if swaggerPath != "" {
		newSources.AddGenerated(genopenapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		servicesImplPackage := packages.New(servicesPath, packageName)
		for _, version := range specification.Versions {
			servicesImplVersionPath := servicesImplPackage.Subpackage("services")
			serviceImplVersionPackage := servicesImplVersionPath.Subpackage(version.Version.FlatCase())

			versionPackage := mainPackage.Subpackage(version.Version.FlatCase())
			modelsVersionPackage := versionPackage.Subpackage("models")
			serviceVersionPackage := versionPackage.Subpackage("services")

			newSources.AddScaffoldedAll(generator.ServicesImplementations(&version, serviceImplVersionPackage, modelsVersionPackage, serviceVersionPackage))
		}
	}

	return newSources
}
