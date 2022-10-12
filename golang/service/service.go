package service

import (
	"generator"
	"golang/models"
	"golang/module"
	"golang/types"
	"openapi"
	"spec"
)

func GenerateService(specification *spec.Spec, moduleName string, swaggerPath string, generatePath string, servicesPath string) *generator.Sources {
	sources := generator.NewSources()

	modules := models.NewModules(moduleName, generatePath, specification)
	modelsGenerator := models.NewGenerator(modules)

	rootModule := module.New(moduleName, generatePath)
	sources.AddGenerated(generateSpecRouting(specification, rootModule))

	sources.AddGenerated(modelsGenerator.GenerateEnumsHelperFunctions())

	emptyModule := rootModule.Submodule("empty")
	sources.AddGenerated(types.GenerateEmpty(emptyModule))

	paramsParserModule := rootModule.Submodule("paramsparser")
	sources.AddGenerated(generateParamsParser(paramsParserModule))

	respondModule := rootModule.Submodule("respond")
	sources.AddGenerated(generateRespondFunctions(respondModule))

	errorsModule := rootModule.Submodule("httperrors")
	errorsModelsModule := errorsModule.Submodule("models")
	sources.AddGenerated(modelsGenerator.GenerateErrorModels(specification.HttpErrors))
	sources.AddGeneratedAll(httpErrors(errorsModule, errorsModelsModule, paramsParserModule, respondModule, &specification.HttpErrors.Responses))

	contentTypeModule := rootModule.Submodule("contenttype")
	sources.AddGenerated(checkContentType(contentTypeModule, errorsModule, errorsModelsModule))

	for _, version := range specification.Versions {
		versionModule := rootModule.Submodule(version.Name.FlatCase())
		modelsModule := versionModule.Submodule(types.VersionModelsPackage)
		routingModule := versionModule.Submodule("routing")

		sources.AddGeneratedAll(generateRoutings(&version, versionModule, routingModule, contentTypeModule, errorsModule, errorsModelsModule, modelsModule, paramsParserModule, respondModule, modelsGenerator))
		sources.AddGeneratedAll(generateServiceInterfaces(&version, versionModule, modelsModule, errorsModelsModule, emptyModule))
		sources.AddGenerated(modelsGenerator.GenerateVersionModels(&version))
	}

	if swaggerPath != "" {
		sources.AddGenerated(openapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		rootServicesModule := module.New(moduleName, servicesPath)
		for _, version := range specification.Versions {
			versionServicesModule := rootServicesModule.Submodule(version.Name.FlatCase())
			versionModule := rootModule.Submodule(version.Name.FlatCase())
			modelsModule := versionModule.Submodule(types.VersionModelsPackage)
			sources.AddScaffoldedAll(generateServiceImplementations(&version, versionModule, modelsModule, versionServicesModule))
		}
	}

	return sources
}
