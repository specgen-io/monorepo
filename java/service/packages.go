package service

import (
	"java/packages"
	"spec"
)

type Packages struct {
	generated       packages.Package
	implementations packages.Package
	versions        map[string]*VersionPackages
	ContentType     packages.Package
	Json            packages.Package
	Errors          packages.Package
	ErrorsModels    packages.Package
	Converters      packages.Package
	Controllers     packages.Package
}

type VersionPackages struct {
	Models       packages.Package
	Services     packages.Package
	Controllers  packages.Package
	ServicesImpl packages.Package
}

func (p *VersionPackages) ServicesApi(api *spec.Api) packages.Package {
	return p.Services.Subpackage(api.Name.SnakeCase())
}

func newVersionPackages(generated, implementations packages.Package, version *spec.Version) *VersionPackages {
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

func NewServicePackages(packageName, generatePath, servicesPath string) *Packages {
	generated := packages.New(generatePath, packageName)
	contenttype := generated.Subpackage("contenttype")
	json := generated.Subpackage("json")
	errors := generated.Subpackage("errors")
	errorsModels := errors.Subpackage("models")
	converters := generated.Subpackage("converters")
	controllers := generated.Subpackage("controllers")
	implementations := packages.New(servicesPath, packageName)

	return &Packages{
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

func (p *Packages) Version(version *spec.Version) *VersionPackages {
	versionPackages, found := p.versions[version.Name.Source]
	if !found {
		versionPackages = newVersionPackages(p.generated, p.implementations, version)
	}
	return versionPackages
}
