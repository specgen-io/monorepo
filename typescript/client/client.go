package client

import (
	"generator"
	"spec"
	"typescript/module"
	"typescript/validations"
	"typescript/validations/modules"
)

func GenerateClient(specification *spec.Spec, generatePath string, client string, validationName string) *generator.Sources {
	modules := modules.NewModules(validationName, generatePath, specification)
	validation := validations.New(validationName, modules)
	g := NewClientGenerator(client, validation)

	sources := generator.NewSources()
	rootModule := module.New(generatePath)

	validationModule := rootModule.Submodule(validationName)
	sources.AddGenerated(validation.SetupLibrary())
	paramsModule := rootModule.Submodule("params")
	sources.AddGenerated(generateParamsBuilder(paramsModule))
	for _, version := range specification.Versions {
		versionModule := rootModule.Submodule(version.Name.FlatCase())
		modelsModule := versionModule.Submodule("models")
		sources.AddGenerated(validation.Models(&version))
		for _, api := range version.Http.Apis {
			apiModule := versionModule.Submodule(api.Name.SnakeCase())
			sources.AddGenerated(g.ApiClient(api, validationModule, modelsModule, paramsModule, apiModule))
		}
	}

	return sources
}
