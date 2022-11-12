package client

import (
	"spec"
	"typescript/module"
	"typescript/validations/modules"
)

type Modules struct {
	modules.Modules
	clients map[string]map[string]module.Module
	Params  module.Module
}

func NewModules(validationName, generatePath string, specification *spec.Spec) *Modules {
	root := module.New(generatePath)
	params := root.Submodule("params")
	clients := map[string]map[string]module.Module{}
	for _, version := range specification.Versions {
		clients[version.Name.Source] = map[string]module.Module{}
		versionModule := root.Submodule(version.Name.FlatCase())
		for _, api := range version.Http.Apis {
			clients[version.Name.Source][api.Name.Source] = versionModule.Submodule(api.Name.SnakeCase())

		}
	}

	return &Modules{
		*modules.NewModules(validationName, generatePath, specification),
		clients,
		params,
	}
}

func (m *Modules) Client(api *spec.Api) module.Module {
	return m.clients[api.InHttp.InVersion.Name.Source][api.Name.Source]
}
