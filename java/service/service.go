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
	sources.AddGenerated(generator.ValidationErrorsHelpers())
	sources.AddGeneratedAll(generator.ErrorModels(specification.HttpErrors))

	sources.AddGenerated(generator.ErrorsHelpers())
	sources.AddGenerated(generator.ExceptionController(&specification.HttpErrors.Responses))
	for _, version := range specification.Versions {
		sources.AddGeneratedAll(generator.Models(&version))
		sources.AddGeneratedAll(generator.ServicesInterfaces(&version))
		sources.AddGeneratedAll(generator.ServicesControllers(&version))
	}
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
