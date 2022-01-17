package service

import (
	"github.com/specgen-io/specgen/v2/genopenapi"
	"github.com/specgen-io/specgen/v2/gents/modules"
	"github.com/specgen-io/specgen/v2/gents/validations"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateService(specification *spec.Spec, swaggerPath string, generatePath string, servicesPath string, server string, validationName string) *sources.Sources {
	validation := validations.New(validationName)
	generator := NewServiceGenerator(server, validation)

	sources := sources.NewSources()

	generateModule := modules.NewModule(generatePath)

	validationModule := generateModule.Submodule(validationName)
	sources.AddGenerated(validation.SetupLibrary(validationModule))
	paramsModule := generateModule.Submodule("params")
	sources.AddGenerated(generateParamsStaticCode(paramsModule))

	for _, version := range specification.Versions {
		versionModule := generateModule.Submodule(version.Version.FlatCase())
		modelsModule := versionModule.Submodule("models")
		sources.AddGenerated(validation.VersionModels(&version, validationModule, modelsModule))
		sources.AddGeneratedAll(generateServiceApis(&version, modelsModule, versionModule))
		routingModule := versionModule.Submodule("routing")
		sources.AddGenerated(generator.VersionRouting(&version, validationModule, paramsModule, routingModule))
	}
	specRouterModule := generateModule.Submodule("spec_router")
	sources.AddGenerated(generator.SpecRouter(specification, generateModule, specRouterModule))

	if swaggerPath != "" {
		sources.AddGenerated(genopenapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		sources.AddScaffoldedAll(generateServicesImplementations(specification, generateModule, modules.NewModule(servicesPath)))
	}

	return sources
}
