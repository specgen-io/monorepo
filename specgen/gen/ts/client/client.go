package client

import (
	"github.com/specgen-io/specgen/v2/gen/ts/modules"
	"github.com/specgen-io/specgen/v2/gen/ts/validations"
	"github.com/specgen-io/specgen/v2/generator"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateClient(specification *spec.Spec, generatePath string, client string, validationName string) *generator.Sources {
	validation := validations.New(validationName)
	g := NewClientGenerator(client, validation)

	sources := generator.NewSources()
	rootModule := modules.New(generatePath)

	validationModule := rootModule.Submodule(validationName)
	sources.AddGenerated(validation.SetupLibrary(validationModule))
	paramsModule := rootModule.Submodule("params")
	sources.AddGenerated(generateParamsBuilder(paramsModule))
	for _, version := range specification.Versions {
		versionModule := rootModule.Submodule(version.Version.FlatCase())
		modelsModule := versionModule.Submodule("models")
		sources.AddGenerated(validation.VersionModels(&version, validationModule, modelsModule))
		for _, api := range version.Http.Apis {
			apiModule := versionModule.Submodule(api.Name.SnakeCase())
			sources.AddGenerated(g.ApiClient(api, validationModule, modelsModule, paramsModule, apiModule))
		}
	}

	return sources
}
