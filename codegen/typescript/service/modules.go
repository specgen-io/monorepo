package service

import (
	"spec"
	"typescript/module"
	"typescript/validations/modules"
)

type Modules struct {
	modules.Modules
	routings     map[string]module.Module
	serviceApis  map[string]map[string]module.Module
	serviceImpls map[string]map[string]module.Module
	Root         module.Module
	Params       module.Module
	Errors       module.Module
	Responses    module.Module
	SpecRouter   module.Module
}

func NewModules(validationName, generatePath, servicesPath string, specification *spec.Spec) *Modules {
	root := module.New(generatePath)

	services := module.New(servicesPath)

	params := root.Submodule("params")
	errors := root.Submodule("errors")
	responses := root.Submodule("responses")
	specRouter := root.Submodule("spec_router")

	routings := map[string]module.Module{}
	serviceApis := map[string]map[string]module.Module{}
	serviceImpls := map[string]map[string]module.Module{}
	for _, version := range specification.Versions {
		serviceApis[version.Name.Source] = map[string]module.Module{}
		serviceImpls[version.Name.Source] = map[string]module.Module{}
		versionModule := root.Submodule(version.Name.FlatCase())
		routings[version.Name.Source] = versionModule.Submodule("routing")
		for _, api := range version.Http.Apis {
			serviceApis[version.Name.Source][api.Name.Source] = versionModule.Submodule(api.Name.SnakeCase())
			serviceImpls[version.Name.Source][api.Name.Source] = services.Submodule(version.Name.FlatCase()).Submodule(api.Name.SnakeCase())
		}
	}

	return &Modules{
		*modules.NewModules(validationName, generatePath, specification),
		routings,
		serviceApis,
		serviceImpls,
		root,
		params,
		errors,
		responses,
		specRouter,
	}
}

func (m *Modules) Routing(version *spec.Version) module.Module {
	return m.routings[version.Name.Source]
}

func (m *Modules) ServiceApi(api *spec.Api) module.Module {
	return m.serviceApis[api.InHttp.InVersion.Name.Source][api.Name.Source]
}

func (m *Modules) ServiceImpl(api *spec.Api) module.Module {
	return m.serviceImpls[api.InHttp.InVersion.Name.Source][api.Name.Source]
}
