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

	generator := NewGenerator(jsonlib, server, packageName, generatePath, servicesPath)

	sources.AddGeneratedAll(generator.ContentType())
	sources.AddGeneratedAll(generator.JsonHelpers())
	sources.AddGeneratedAll(generator.Errors(specification.HttpErrors.ResolvedModels))

	for _, version := range specification.Versions {
		sources.AddGeneratedAll(generator.Models(&version))
		sources.AddGeneratedAll(generator.ServicesInterfaces(&version))
		sources.AddGeneratedAll(generator.ServicesControllers(&version))
	}
	sources.AddGenerated(generator.ExceptionController(&specification.HttpErrors.Responses))

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
