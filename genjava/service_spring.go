package genjava

import (
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/genopenapi"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateService(specification *spec.Spec, packageName string, swaggerPath string, generatePath string, servicesPath string) *gen.Sources {
	sources := gen.NewSources()

	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}

	mainPackage := Package(generatePath, packageName)

	modelsPackage := mainPackage.Subpackage("models")
	sources.AddGenerated(generateJson(modelsPackage))

	for _, version := range specification.Versions {
		versionPackage := mainPackage.Subpackage(version.Version.FlatCase())

		modelsVersionPackage := versionPackage.Subpackage("models")
		sources.AddGeneratedAll(generateVersionModels(&version, modelsVersionPackage))

		serviceVersionPackage := versionPackage.Subpackage("services")
		sources.AddGeneratedAll(generateServicesInterfaces(&version, serviceVersionPackage, modelsVersionPackage))

		controllerVersionPackage := versionPackage.Subpackage("controllers")
		sources.AddGeneratedAll(generateServicesControllers(&version, controllerVersionPackage, modelsPackage, modelsVersionPackage, serviceVersionPackage))
	}

	if swaggerPath != "" {
		sources.AddGenerated(genopenapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		servicesImplPackage := Package(servicesPath, packageName)
		for _, version := range specification.Versions {
			servicesImplVersionPath := servicesImplPackage.Subpackage("services")
			serviceImplVersionPackage := servicesImplVersionPath.Subpackage(version.Version.FlatCase())

			versionPackage := mainPackage.Subpackage(version.Version.FlatCase())
			modelsVersionPackage := versionPackage.Subpackage("models")
			serviceVersionPackage := versionPackage.Subpackage("services")

			sources.AddScaffoldedAll(generateServicesImplementations(&version, serviceImplVersionPackage, modelsVersionPackage, serviceVersionPackage))
		}
	}

	return sources
}
