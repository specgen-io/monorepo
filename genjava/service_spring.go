package genjava

import (
	"github.com/specgen-io/specgen/v2/genopenapi"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateService(specification *spec.Spec, jsonlib string, packageName string, swaggerPath string, generatePath string, servicesPath string) *sources.Sources {
	newSources := sources.NewSources()

	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}

	mainPackage := Package(generatePath, packageName)

	generator := NewGenerator(jsonlib)

	modelsPackage := mainPackage.Subpackage("models")
	newSources.AddGeneratedAll(generator.Models.SetupLibrary(modelsPackage))

	for _, version := range specification.Versions {
		versionPackage := mainPackage.Subpackage(version.Version.FlatCase())

		modelsVersionPackage := versionPackage.Subpackage("models")
		newSources.AddGeneratedAll(generator.generateVersionModels(&version, modelsVersionPackage))

		serviceVersionPackage := versionPackage.Subpackage("services")
		newSources.AddGeneratedAll(generator.generateServicesInterfaces(&version, serviceVersionPackage, modelsVersionPackage))

		controllerVersionPackage := versionPackage.Subpackage("controllers")
		newSources.AddGeneratedAll(generator.generateServicesControllers(&version, controllerVersionPackage, modelsPackage, modelsVersionPackage, serviceVersionPackage))
	}

	if swaggerPath != "" {
		newSources.AddGenerated(genopenapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		servicesImplPackage := Package(servicesPath, packageName)
		for _, version := range specification.Versions {
			servicesImplVersionPath := servicesImplPackage.Subpackage("services")
			serviceImplVersionPackage := servicesImplVersionPath.Subpackage(version.Version.FlatCase())

			versionPackage := mainPackage.Subpackage(version.Version.FlatCase())
			modelsVersionPackage := versionPackage.Subpackage("models")
			serviceVersionPackage := versionPackage.Subpackage("services")

			newSources.AddScaffoldedAll(generator.generateServicesImplementations(&version, serviceImplVersionPackage, modelsVersionPackage, serviceVersionPackage))
		}
	}

	return newSources
}
