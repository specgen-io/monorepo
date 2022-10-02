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

	generator := NewGenerator(jsonlib, server, packageName, generatePath, servicesPath, specification)

	sources.AddGeneratedAll(generator.ContentType())
	sources.AddGeneratedAll(generator.JsonHelpers())
	sources.AddGenerated(generator.ModelsValidation())
	sources.AddGeneratedAll(generator.ErrorModels(specification.HttpErrors))

	sources.AddGenerated(generator.ErrorsHelpers())
	sources.AddGenerated(generator.ExceptionController(&specification.HttpErrors.Responses))
	for _, version := range specification.Versions {
		sources.AddGeneratedAll(generator.Models(&version))
		sources.AddGeneratedAll(generator.ServicesInterfaces(&version))
		sources.AddGeneratedAll(generator.ServicesControllers(&version))
	}

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
