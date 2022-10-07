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

	sources.AddGeneratedAll(generator.ContentType())
	sources.AddGeneratedAll(generator.Errors())
	sources.AddGeneratedAll(generator.ErrorModels(specification.HttpErrors))

	for _, version := range specification.Versions {
		sources.AddGeneratedAll(generator.Models(&version))
		sources.AddGeneratedAll(generator.ServicesInterfaces(&version))
		sources.AddGeneratedAll(generator.ServicesControllers(&version))
	}
	sources.AddGenerated(generator.ExceptionController(&specification.HttpErrors.Responses))

	sources.AddGeneratedAll(generator.JsonHelpers())

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
