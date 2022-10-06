package client

import (
	"kotlin/models"
	"kotlin/packages"
	"spec"
)

type Packages struct {
	models.Packages
	clients map[string]map[string]packages.Package
	Root    packages.Package
	Utils   packages.Package
}

func NewPackages(packageName, generatePath string, specification *spec.Spec) *Packages {
	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}

	generated := packages.New(generatePath, packageName)
	utils := generated.Subpackage("utils")

	clients := map[string]map[string]packages.Package{}
	for _, version := range specification.Versions {
		clients[version.Name.Source] = map[string]packages.Package{}
		versionPackage := generated.Subpackage(version.Name.FlatCase())
		versionClientsPackage := versionPackage.Subpackage("clients")
		for _, api := range version.Http.Apis {
			clients[version.Name.Source][api.Name.Source] = versionClientsPackage.Subpackage(api.Name.SnakeCase())
		}
	}

	return &Packages{
		*models.NewPackages(packageName, generatePath, specification),
		clients,
		generated,
		utils,
	}
}

func (p *Packages) Client(api *spec.Api) packages.Package {
	return p.clients[api.InHttp.InVersion.Name.Source][api.Name.Source]
}
