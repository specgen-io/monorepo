package modules

import (
	"spec"
	"typescript/module"
)

type Modules struct {
	models       map[string]module.Module
	ErrorModules module.Module
	Validation   module.Module
}

func NewModules(validationName, generatePath string, specification *spec.Spec) *Modules {
	root := module.New(generatePath)
	errorModels := root.Submodule("errors")
	validation := root.Submodule(validationName)

	models := map[string]module.Module{}
	for _, version := range specification.Versions {
		versionModule := root.Submodule(version.Name.FlatCase())
		models[version.Name.Source] = versionModule.Submodule("models")
	}

	return &Modules{
		models,
		errorModels,
		validation,
	}
}

func (m *Modules) Models(version *spec.Version) module.Module {
	return m.models[version.Name.Source]
}
