package client

import (
	"github.com/specgen-io/specgen/v2/gents/modules"
	"github.com/specgen-io/specgen/v2/gents/validation"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateClient(specification *spec.Spec, generatePath string, client string, validationName string) *sources.Sources {
	validation := validation.New(validationName)
	generator := newClientGenerator(client, validation)

	sources := sources.NewSources()
	module := modules.NewModule(generatePath)

	validationModule := module.Submodule(validationName)
	sources.AddGenerated(validation.SetupLibrary(validationModule))
	paramsModule := module.Submodule("params")
	sources.AddGenerated(generateParamsBuilder(paramsModule))
	for _, version := range specification.Versions {
		versionModule := module.Submodule(version.Version.FlatCase())
		modelsModule := versionModule.Submodule("models")
		sources.AddGenerated(validation.GenerateVersionModels(&version, validationModule, modelsModule))
		for _, api := range version.Http.Apis {
			apiModule := versionModule.Submodule(api.Name.SnakeCase())
			sources.AddGenerated(generator.generateApiClient(api, validationModule, modelsModule, paramsModule, apiModule))
		}
	}

	return sources
}
