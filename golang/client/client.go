package client

import (
	"generator"
	"golang/models"
	"golang/module"
	"golang/types"
	"spec"
)

func GenerateClient(specification *spec.Spec, moduleName string, generatePath string) *generator.Sources {
	sources := generator.NewSources()

	modules := models.NewModules(moduleName, generatePath, specification)
	clientGenerator := NewGenerator(modules)

	rootModule := module.New(moduleName, generatePath)

	sources.AddGenerated(clientGenerator.Models.GenerateEnumsHelperFunctions())

	emptyModule := rootModule.Submodule("empty")
	sources.AddGenerated(types.GenerateEmpty(emptyModule))

	convertModule := rootModule.Submodule("convert")
	sources.AddGenerated(generateConverter(convertModule))

	responseModule := rootModule.Submodule("response")
	sources.AddGenerated(generateResponseFunctions(responseModule))

	errorsModule := rootModule.Submodule("httperrors")
	errorsModelsModule := errorsModule.Submodule(types.ErrorsModelsPackage)
	sources.AddGenerated(clientGenerator.Models.GenerateErrorModels(specification.HttpErrors))
	sources.AddGenerated(httpErrors(errorsModule, errorsModelsModule, &specification.HttpErrors.Responses))

	for _, version := range specification.Versions {
		versionModule := rootModule.Submodule(version.Name.FlatCase())
		modelsModule := versionModule.Submodule(types.VersionModelsPackage)
		sources.AddGenerated(clientGenerator.Models.GenerateVersionModels(&version))
		sources.AddGeneratedAll(clientGenerator.Client.GenerateClientsImplementations(&version, versionModule, convertModule, emptyModule, errorsModule, errorsModelsModule, modelsModule, responseModule))
	}
	return sources
}
