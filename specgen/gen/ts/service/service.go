package service

import (
	"github.com/specgen-io/specgen/spec/v2"
	"github.com/specgen-io/specgen/specgen/v2/gen/openapi"
	"github.com/specgen-io/specgen/specgen/v2/gen/ts/modules"
	"github.com/specgen-io/specgen/specgen/v2/gen/ts/validations"
	"github.com/specgen-io/specgen/specgen/v2/generator"
)

func GenerateService(specification *spec.Spec, swaggerPath string, generatePath string, servicesPath string, server string, validationName string) *generator.Sources {
	validation := validations.New(validationName)
	g := NewServiceGenerator(server, validation)

	sources := generator.NewSources()

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
		sources.AddGenerated(g.VersionRouting(&version, validationModule, paramsModule, routingModule))
	}
	specRouterModule := rootModule.Submodule("spec_router")
	sources.AddGenerated(g.SpecRouter(specification, rootModule, specRouterModule))

	if swaggerPath != "" {
		sources.AddGenerated(openapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		sources.AddScaffoldedAll(generateServicesImplementations(specification, rootModule, modules.New(servicesPath)))
	}

	return sources
}
