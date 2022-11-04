package service

import (
	"generator"
	"openapi"
	"spec"
	"typescript/module"
	"typescript/validations"
	"typescript/validations/modules"
)

func GenerateService(specification *spec.Spec, swaggerPath string, generatePath string, servicesPath string, server string, validationName string) *generator.Sources {
	modules := modules.NewModules(validationName, generatePath, specification)
	validation := validations.New(validationName, modules)
	g := NewServiceGenerator(server, validation)

	sources := generator.NewSources()

	rootModule := module.New(generatePath)

	validationModule := rootModule.Submodule(validationName)
	sources.AddGenerated(validation.SetupLibrary())
	paramsModule := rootModule.Submodule("params")
	sources.AddGenerated(generateParamsStaticCode(paramsModule))
	errorsModule := rootModule.Submodule("errors")
	sources.AddGenerated(validation.ErrorModels(specification.HttpErrors))
	responsesModule := rootModule.Submodule("responses")
	sources.AddGenerated(g.Responses(responsesModule, validationModule, errorsModule))

	for _, version := range specification.Versions {
		versionModule := rootModule.Submodule(version.Name.FlatCase())
		modelsModule := versionModule.Submodule("models")
		sources.AddGenerated(validation.Models(&version))
		sources.AddGeneratedAll(generateServiceApis(&version, modelsModule, errorsModule, versionModule))
		routingModule := versionModule.Submodule("routing")
		sources.AddGenerated(g.VersionRouting(&version, routingModule, modelsModule, validationModule, paramsModule, errorsModule, responsesModule))
	}
	specRouterModule := rootModule.Submodule("spec_router")
	sources.AddGenerated(g.SpecRouter(specification, rootModule, specRouterModule))

	if swaggerPath != "" {
		sources.AddGenerated(openapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		sources.AddScaffoldedAll(generateServicesImplementations(specification, rootModule, errorsModule, module.New(servicesPath)))
	}

	return sources
}
