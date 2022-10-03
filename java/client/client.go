package client

import (
	"generator"
	"java/models"
	"java/packages"
	"spec"
)

func Generate(specification *spec.Spec, jsonlib string, packageName string, generatePath string) *generator.Sources {
	sources := generator.NewSources()

	thepackages := models.NewPackages(packageName, generatePath, specification)
	generator := NewGenerator(jsonlib, thepackages)

	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}

	mainPackage := packages.New(generatePath, packageName)

	sources.AddGenerated(clientException(mainPackage))

	utilsPackage := mainPackage.Subpackage("utils")
	sources.AddGeneratedAll(generateUtils(utilsPackage))

	jsonPackage := mainPackage.Subpackage("json")

	errorsPackage := mainPackage.Subpackage("errors")
	errorsModelsPackage := errorsPackage.Subpackage("models")

	sources.AddGeneratedAll(generator.Models.ErrorModels(specification.HttpErrors))

	for _, version := range specification.Versions {
		sources.AddGeneratedAll(generator.Models.Models(&version))

		versionPackage := mainPackage.Subpackage(version.Name.FlatCase())
		versionModelsPackage := versionPackage.Subpackage("models")
		versionClientsPackage := versionPackage.Subpackage("clients")
		sources.AddGeneratedAll(generator.Clients(&version, versionClientsPackage, versionModelsPackage, errorsModelsPackage, jsonPackage, utilsPackage, mainPackage))
	}

	sources.AddGeneratedAll(generator.Models.SetupLibrary())

	return sources
}
