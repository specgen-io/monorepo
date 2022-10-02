package service

import (
	"java/models"
	"java/packages"
	"spec"
)

type Packages struct {
	models.Packages
	versions    map[string]*VersionPackages
	ContentType packages.Package
	Converters  packages.Package
	Controllers packages.Package
}

func NewPackages(packageName, generatePath, servicesPath string, specification *spec.Spec) *Packages {
	generated := packages.New(generatePath, packageName)
	contenttype := generated.Subpackage("contenttype")
	converters := generated.Subpackage("converters")
	controllers := generated.Subpackage("controllers")
	implementations := packages.New(servicesPath, packageName)

	versions := map[string]*VersionPackages{}
	for _, version := range specification.Versions {
		versions[version.Name.Source] = newVersionPackages(generated, implementations, &version)
	}

	return &Packages{
		*models.NewPackages(packageName, generatePath, specification),
		map[string]*VersionPackages{},
		contenttype,
		converters,
		controllers,
	}
}

func (p *Packages) Version(version *spec.Version) *VersionPackages {
	return p.versions[version.Name.Source]
}

type VersionPackages struct {
	Services     packages.Package
	Controllers  packages.Package
	ServicesImpl packages.Package
}

func (p *VersionPackages) ServicesApi(api *spec.Api) packages.Package {
	return p.Services.Subpackage(api.Name.SnakeCase())
}

func newVersionPackages(generated, implementations packages.Package, version *spec.Version) *VersionPackages {
	main := generated.Subpackage(version.Name.FlatCase())
	services := main.Subpackage("services")
	controllers := main.Subpackage("controllers")
	servicesImpl := implementations.Subpackage("services").Subpackage(version.Name.FlatCase())

	return &VersionPackages{
		services,
		controllers,
		servicesImpl,
	}
}
