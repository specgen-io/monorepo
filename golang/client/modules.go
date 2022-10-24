package client

import (
	"golang/models"
	"golang/module"
	"spec"
)

type Modules struct {
	models.Modules
	clients  map[string]map[string]module.Module
	Root     module.Module
	Empty    module.Module
	Convert  module.Module
	Response module.Module
}

func NewModules(moduleName string, generatePath string, specification *spec.Spec) *Modules {
	root := module.New(moduleName, generatePath)
	empty := root.Submodule("empty")
	convert := root.Submodule("convert")
	response := root.Submodule("response")

	clients := map[string]map[string]module.Module{}
	for _, version := range specification.Versions {
		versionModule := root.Submodule(version.Name.FlatCase())
		for _, api := range version.Http.Apis {
			clients[version.Name.Source][api.Name.Source] = versionModule.Submodule(api.Name.SnakeCase())
		}
	}

	return &Modules{
		*models.NewModules(moduleName, generatePath, specification),
		clients,
		root,
		empty,
		convert,
		response,
	}
}

func (p *Modules) Client(api *spec.Api) module.Module {
	return p.clients[api.InHttp.InVersion.Name.Source][api.Name.Source]
}
