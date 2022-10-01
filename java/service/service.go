package service

import (
	"generator"
	"openapi"
	"spec"
)

func Generate(specification *spec.Spec, jsonlib, server, packageName, swaggerPath, generatePath, servicesPath string) *generator.Sources {
	sources := generator.NewSources()

	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}

	servicePackages := NewServicePackages(packageName, generatePath, servicesPath)

	generator := NewGenerator(jsonlib, server, packageName, generatePath, servicesPath)

	sources.AddGeneratedAll(generator.Server.ContentType())
	sources.AddGeneratedAll(generator.Server.Errors())
	sources.AddGeneratedAll(generator.Server.JsonHelpers())
	sources.AddGeneratedAll(generator.Models.Models(specification.HttpErrors.ResolvedModels, servicePackages.ErrorsModels, servicePackages.Json))

	for _, version := range specification.Versions {
		versionPackages := servicePackages.Version(&version)
		sources.AddGeneratedAll(generator.Models.Models(version.ResolvedModels, versionPackages.Models, servicePackages.Json))
		sources.AddGeneratedAll(generator.ServicesInterfaces(&version))
		sources.AddGeneratedAll(generator.Server.ServicesControllers(&version))
	}
	sources.AddGenerated(generator.Server.ExceptionController(&specification.HttpErrors.Responses))

	if swaggerPath != "" {
		sources.AddGenerated(openapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		for _, version := range specification.Versions {
			sources.AddScaffoldedAll(generator.ServicesImplementations(&version))
		}
	}

	return sources
}
