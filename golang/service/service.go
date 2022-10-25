package service

import (
	"generator"
	"openapi"
	"spec"
)

func GenerateService(specification *spec.Spec, moduleName string, swaggerPath string, generatePath string, servicesPath string) *generator.Sources {
	sources := generator.NewSources()

	modules := NewModules(moduleName, generatePath, servicesPath, specification)
	serviceGenerator := NewGenerator(modules)

	sources.AddGenerated(serviceGenerator.GenerateSpecRouting(specification))
	sources.AddGeneratedAll(serviceGenerator.AllStaticFiles())
	sources.AddGenerated(serviceGenerator.ErrorModels(specification.HttpErrors))
	sources.AddGeneratedAll(serviceGenerator.HttpErrors(&specification.HttpErrors.Responses))
	for _, version := range specification.Versions {
		sources.AddGeneratedAll(serviceGenerator.GenerateRoutings(&version))
		sources.AddGeneratedAll(serviceGenerator.generateServiceInterfaces(&version))
		sources.AddGenerated(serviceGenerator.Models(&version))
	}

	if swaggerPath != "" {
		sources.AddGenerated(openapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		for _, version := range specification.Versions {
			sources.AddScaffoldedAll(serviceGenerator.generateServiceImplementations(&version))
		}
	}

	return sources
}
