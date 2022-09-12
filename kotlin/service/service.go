package service

import (
	"generator"
	"kotlin/modules"
	"openapi"
	"spec"
)

func Generate(specification *spec.Spec, jsonlib, server, packageName, swaggerPath, generatePath, servicesPath string) *generator.Sources {
	sources := generator.NewSources()

	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}

	mainPackage := modules.Package(generatePath, packageName)

	generator := NewGenerator(jsonlib, server)

	contentTypePackage := mainPackage.Subpackage("contenttype")
	sources.AddGeneratedAll(generator.Server.ContentType(contentTypePackage))

	jsonPackage := mainPackage.Subpackage("json")

	errorsPackage := mainPackage.Subpackage("errors")
	errorsModelsPackage := errorsPackage.Subpackage("models")

	sources.AddGeneratedAll(generator.Server.Errors(errorsPackage, errorsModelsPackage, contentTypePackage, jsonPackage))
	sources.AddGeneratedAll(generator.Models.Models(specification.HttpErrors.ResolvedModels, errorsModelsPackage, jsonPackage))

	for _, version := range specification.Versions {
		versionPackage := mainPackage.Subpackage(version.Name.FlatCase())

		modelsVersionPackage := versionPackage.Subpackage("models")
		sources.AddGeneratedAll(generator.Models.Models(version.ResolvedModels, modelsVersionPackage, jsonPackage))

		serviceVersionPackage := versionPackage.Subpackage("services")
		sources.AddGeneratedAll(generator.ServicesInterfaces(&version, serviceVersionPackage, modelsVersionPackage, errorsModelsPackage))

		controllerVersionPackage := versionPackage.Subpackage("controllers")
		sources.AddGeneratedAll(generator.Server.ServicesControllers(&version, mainPackage, controllerVersionPackage, contentTypePackage, jsonPackage, modelsVersionPackage, errorsModelsPackage, serviceVersionPackage))
	}
	controllerPackage := mainPackage.Subpackage("controllers")
	sources.AddGenerated(generator.Server.ExceptionController(&specification.HttpErrors.Responses, controllerPackage, errorsPackage, errorsModelsPackage, jsonPackage))

	sources.AddGeneratedAll(generator.Server.JsonHelpers(jsonPackage))

	if swaggerPath != "" {
		sources.AddGenerated(openapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		servicesImplPackage := modules.Package(servicesPath, packageName)
		for _, version := range specification.Versions {
			servicesImplVersionPath := servicesImplPackage.Subpackage("services")
			serviceImplVersionPackage := servicesImplVersionPath.Subpackage(version.Name.FlatCase())

			versionPackage := mainPackage.Subpackage(version.Name.FlatCase())
			modelsVersionPackage := versionPackage.Subpackage("models")
			serviceVersionPackage := versionPackage.Subpackage("services")

			sources.AddScaffoldedAll(generator.ServicesImplementations(&version, serviceImplVersionPackage, modelsVersionPackage, serviceVersionPackage))
		}
	}

	return sources
}
