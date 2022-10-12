package models

import (
	"generator"
	"golang/module"
	"golang/types"
	"spec"
)

func GenerateModels(specification *spec.Spec, moduleName string, generatePath string) *generator.Sources {
	sources := generator.NewSources()

	modules := NewModules(moduleName, generatePath, specification)
	generator := NewGenerator(modules)

	rootModule := module.New(moduleName, generatePath)

	enumsModule := rootModule.Submodule("enums")
	sources.AddGenerated(generator.GenerateEnumsHelperFunctions(enumsModule))

	for _, version := range specification.Versions {
		versionModule := rootModule.Submodule(version.Name.FlatCase())
		modelsModule := versionModule.Submodule(types.VersionModelsPackage)
		sources.AddGenerated(generator.GenerateVersionModels(version.ResolvedModels, modelsModule, enumsModule))
	}
	return sources
}
