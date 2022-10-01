package service

import (
	"generator"
	"java/packages"
	"openapi"
	"spec"
)

func Generate(specification *spec.Spec, jsonlib, server, packageName, swaggerPath, generatePath, servicesPath string) *generator.Sources {
	sources := generator.NewSources()

	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}

	servicePackages := NewServicePackages(packageName, generatePath, servicesPath)

	mainPackage := packages.New(generatePath, packageName)

	generator := NewGenerator(jsonlib, server, servicePackages)

	sources.AddGeneratedAll(generator.Server.ContentType())

	sources.AddGeneratedAll(generator.Server.Errors())
	sources.AddGeneratedAll(generator.Models.Models(specification.HttpErrors.ResolvedModels, servicePackages.ErrorsModels, servicePackages.Json))

	for _, version := range specification.Versions {
		versionPackages := servicePackages.Version(&version)
		sources.AddGeneratedAll(generator.Models.Models(version.ResolvedModels, versionPackages.Models, servicePackages.Json))
		sources.AddGeneratedAll(generator.ServicesInterfaces(&version))
		sources.AddGeneratedAll(generator.Server.ServicesControllers(&version))
	}
	sources.AddGenerated(generator.Server.ExceptionController(&specification.HttpErrors.Responses))

	sources.AddGeneratedAll(generator.Server.JsonHelpers())

	if swaggerPath != "" {
		sources.AddGenerated(openapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		servicesImplPackage := packages.New(servicesPath, packageName)
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
