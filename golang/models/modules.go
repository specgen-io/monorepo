package models

import (
	"golang/module"
	"golang/types"
	"spec"
)

type Modules struct {
	models           map[string]module.Module
	Root             module.Module
	Enums            module.Module
	HttpErrors       module.Module
	HttpErrorsModels module.Module
}

func NewModules(moduleName string, generatePath string, specification *spec.Spec) *Modules {
	generated := module.New(moduleName, generatePath)
	enums := generated.Submodule("enums")
	httperrors := generated.Submodule("httperrors")
	httperrorsModels := httperrors.Submodule(types.ErrorsModelsPackage)

	models := map[string]module.Module{}
	for _, version := range specification.Versions {
		models[version.Name.Source] = generated.Submodule(version.Name.FlatCase()).Submodule(types.VersionModelsPackage)
	}

	return &Modules{
		models,
		generated,
		enums,
		httperrors,
		httperrorsModels,
	}
}

func (p *Modules) Models(version *spec.Version) module.Module {
	return p.models[version.Name.Source]
}
