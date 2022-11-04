package service

import (
	"generator"
	"openapi"
	"spec"
)

func GenerateService(specification *spec.Spec, swaggerPath string, generatePath string, servicesPath string, server string, validationName string) *generator.Sources {
	sources := generator.NewSources()
	modules := NewModules(validationName, generatePath, servicesPath, specification)
	generator := NewServiceGenerator(server, validationName, modules)

	sources.AddGenerated(generator.SetupLibrary())
	sources.AddGenerated(generator.ParamsStaticCode())
	sources.AddGenerated(generator.ErrorModels(specification.HttpErrors))
	sources.AddGenerated(generator.Responses())

	for _, version := range specification.Versions {
		sources.AddGenerated(generator.Models(&version))
		sources.AddGeneratedAll(generator.ServiceApis(&version))
		sources.AddGenerated(generator.VersionRouting(&version))
	}
	sources.AddGenerated(generator.SpecRouter(specification))

	if swaggerPath != "" {
		sources.AddGenerated(openapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		sources.AddScaffoldedAll(generator.ServicesImpls(specification))
	}

	return sources
}
