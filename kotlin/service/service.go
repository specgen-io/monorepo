package service

import (
	"generator"
	"openapi"
	"spec"
)

func Generate(specification *spec.Spec, jsonlib, server, packageName, swaggerPath, generatePath, servicesPath string) *generator.Sources {
	sources := generator.NewSources()

	packages := NewPackages(packageName, generatePath, servicesPath, specification)
	generator := NewGenerator(jsonlib, server, packages)

	sources.AddGeneratedAll(generator.Server.ContentType())
	sources.AddGeneratedAll(generator.Server.Errors())
	sources.AddGeneratedAll(generator.Models.ErrorModels(specification.HttpErrors))

	for _, version := range specification.Versions {
		sources.AddGeneratedAll(generator.Models.Models(&version))
		sources.AddGeneratedAll(generator.ServicesInterfaces(&version))
		sources.AddGeneratedAll(generator.Server.ServicesControllers(&version))
	}
	sources.AddGenerated(generator.Server.ExceptionController(&specification.HttpErrors.Responses))

	sources.AddGeneratedAll(generator.Server.JsonHelpers())

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
