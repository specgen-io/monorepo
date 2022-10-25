package service

import (
	"golang/models"
	"golang/module"
	"spec"
)

type Modules struct {
	models.Modules
	servicesApis  map[string]map[string]module.Module
	servicesImpls map[string]map[string]module.Module
	routing       map[string]module.Module
	Root          module.Module
	Empty         module.Module
	ParamsParser  module.Module
	Respond       module.Module
	ContentType   module.Module
}

func NewModules(moduleName, generatePath, servicesPath string, specification *spec.Spec) *Modules {
	root := module.New(moduleName, generatePath)
	empty := root.Submodule("empty")
	paramsParser := root.Submodule("paramsparser")
	respond := root.Submodule("respond")
	contentType := root.Submodule("contenttype")

	servicesApis := map[string]map[string]module.Module{}
	servicesImpls := map[string]map[string]module.Module{}
	routing := map[string]module.Module{}

	rootImpls := module.New(moduleName, servicesPath)
	for _, version := range specification.Versions {
		routing[version.Name.Source] = root.Submodule(version.Name.FlatCase()).Submodule("routing")
		servicesApis[version.Name.Source] = map[string]module.Module{}
		servicesImpls[version.Name.Source] = map[string]module.Module{}
		for _, api := range version.Http.Apis {
			servicesApis[version.Name.Source][api.Name.Source] = root.Submodule(version.Name.FlatCase()).Submodule(api.Name.SnakeCase())
			servicesApis[version.Name.Source][api.Name.Source] = rootImpls.Submodule(version.Name.FlatCase()).Submodule(api.Name.SnakeCase())
		}
	}

	return &Modules{
		*models.NewModules(moduleName, generatePath, specification),
		servicesApis,
		servicesImpls,
		routing,
		root,
		empty,
		paramsParser,
		respond,
		contentType,
	}
}

func (m *Modules) ServicesApi(api *spec.Api) module.Module {
	return m.servicesApis[api.InHttp.InVersion.Name.Source][api.Name.Source]
}

func (m *Modules) Routing(version *spec.Version) module.Module {
	return m.routing[version.Name.Source]
}

func (m *Modules) ServicesImpl(api *spec.Api) module.Module {
	return m.servicesImpls[api.InHttp.InVersion.Name.Source][api.Name.Source]
}

//func (p *Modules) ServicesImpl(version *spec.Version) module.Module {
//	return p.versions[version.Name.Source].ServicesImpl
//}
