package service

import (
	"github.com/specgen-io/specgen/v2/gen/openapi"
	"github.com/specgen-io/specgen/v2/gen/ts/modules"
	"github.com/specgen-io/specgen/v2/gen/ts/validations"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateService(specification *spec.Spec, swaggerPath string, generatePath string, servicesPath string, server string, validationName string) *sources.Sources {
	validation := validations.New(validationName)
	generator := NewServiceGenerator(server, validation)

	sources := sources.NewSources()

	rootModule := modules.New(generatePath)

	validationModule := rootModule.Submodule(validationName)
	sources.AddGenerated(validation.SetupLibrary(validationModule))
	paramsModule := rootModule.Submodule("params")
	sources.AddGenerated(generateParamsStaticCode(paramsModule))

	for _, version := range specification.Versions {
		versionModule := rootModule.Submodule(version.Version.FlatCase())
		modelsModule := versionModule.Submodule("models")
		sources.AddGenerated(validation.VersionModels(&version, validationModule, modelsModule))
		sources.AddGeneratedAll(generateServiceApis(&version, modelsModule, versionModule))
		routingModule := versionModule.Submodule("routing")
		sources.AddGenerated(generator.VersionRouting(&version, validationModule, paramsModule, routingModule))
	}
	specRouterModule := rootModule.Submodule("spec_router")
	sources.AddGenerated(generator.SpecRouter(specification, rootModule, specRouterModule))

	if swaggerPath != "" {
		sources.AddGenerated(openapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		sources.AddScaffoldedAll(generateServicesImplementations(specification, rootModule, modules.New(servicesPath)))
	}

	return sources
}
