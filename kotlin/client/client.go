package client

import (
	"generator"
	"kotlin/packages"
	"spec"
)

func Generate(specification *spec.Spec, jsonlib, client, packageName, generatePath string) *generator.Sources {
	sources := generator.NewSources()

	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}

	mainPackage := packages.New(generatePath, packageName)

	thepackages := NewPackages(packageName, generatePath, specification)
	generator := NewGenerator(jsonlib, client, thepackages)

	jsonPackage := mainPackage.Subpackage("json")

	errorsPackage := mainPackage.Subpackage("errors")
	errorsModelsPackage := errorsPackage.Subpackage("models")

	sources.AddGeneratedAll(generator.Models.ErrorModels(specification.HttpErrors))

	for _, version := range specification.Versions {
		versionPackage := mainPackage.Subpackage(version.Name.FlatCase())

		modelsVersionPackage := versionPackage.Subpackage("models")
		sources.AddGeneratedAll(generator.Models.Models(&version))

		clientVersionPackage := versionPackage.Subpackage("clients")

		sources.AddGeneratedAll(generator.Client.Clients(&version, clientVersionPackage, modelsVersionPackage, errorsModelsPackage, jsonPackage, mainPackage))
	}

	sources.AddGeneratedAll(generator.Models.SetupLibrary())

	return sources
}
