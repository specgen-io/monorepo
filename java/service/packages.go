package service

import (
	"java/packages"
	"spec"
)

type ServicePackages struct {
	generated       packages.Module
	implementations packages.Module
	versions        map[string]*VersionPackages
	ContentType     packages.Module
	Json            packages.Module
	Errors          packages.Module
	ErrorsModels    packages.Module
	Converters      packages.Module
	Controllers     packages.Module
}

type VersionPackages struct {
	Models       packages.Module
	Services     packages.Module
	Controllers  packages.Module
	ServicesImpl packages.Module
}

func (p *VersionPackages) ServicesApi(api *spec.Api) packages.Module {
	return p.Services.Subpackage(api.Name.SnakeCase())
}

func newVersionPackages(generated, implementations packages.Module, version *spec.Version) *VersionPackages {
	main := generated.Subpackage(version.Name.FlatCase())
	models := main.Subpackage("models")
	services := main.Subpackage("services")
	controllers := main.Subpackage("controllers")
	servicesImpl := implementations.Subpackage("services").Subpackage(version.Name.FlatCase())

	return &VersionPackages{
		models,
		services,
		controllers,
		servicesImpl,
	}
}

func NewServicePackages(packageName, generatePath, servicesPath string) *ServicePackages {
	generated := packages.New(generatePath, packageName)
	contenttype := generated.Subpackage("contenttype")
	json := generated.Subpackage("json")
	errors := generated.Subpackage("errors")
	errorsModels := errors.Subpackage("models")
	converters := generated.Subpackage("converters")
	controllers := generated.Subpackage("controllers")
	implementations := packages.New(servicesPath, packageName)

	return &ServicePackages{
		generated,
		implementations,
		map[string]*VersionPackages{},
		contenttype,
		json,
		errors,
		errorsModels,
		converters,
		controllers,
	}
}

func (p *ServicePackages) Version(version *spec.Version) *VersionPackages {
	versionPackages, found := p.versions[version.Name.Source]
	if !found {
		versionPackages = newVersionPackages(p.generated, p.implementations, version)
	}
	return versionPackages
}
