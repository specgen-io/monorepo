package client

import (
	"generator"
	"kotlin/modules"
	"spec"
)

func Generate(specification *spec.Spec, jsonlib, client, packageName, generatePath string) *generator.Sources {
	sources := generator.NewSources()

	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}

	mainPackage := modules.Package(generatePath, packageName)

	generator := NewGenerator(jsonlib, client)

	jsonPackage := mainPackage.Subpackage("json")

	errorsPackage := mainPackage.Subpackage("errors")
	errorsModelsPackage := errorsPackage.Subpackage("models")

	sources.AddGeneratedAll(generator.Models.Models(specification.HttpErrors.ResolvedModels, errorsModelsPackage, jsonPackage))

	for _, version := range specification.Versions {
		versionPackage := mainPackage.Subpackage(version.Name.FlatCase())

		modelsVersionPackage := versionPackage.Subpackage("models")
		sources.AddGeneratedAll(generator.Models.Models(version.ResolvedModels, modelsVersionPackage, jsonPackage))

		clientVersionPackage := versionPackage.Subpackage("clients")

		sources.AddGeneratedAll(generator.Client.ClientImplementation(&version, clientVersionPackage, modelsVersionPackage, errorsModelsPackage, jsonPackage, mainPackage))
	}

	sources.AddGeneratedAll(generator.Models.SetupLibrary(jsonPackage))

	return sources
}
