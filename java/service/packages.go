package service

import (
	"java/packages"
	"spec"
)

type ServicePackages struct {
	Root             packages.Module
	ContentType      packages.Module
	Json             packages.Module
	Errors           packages.Module
	ErrorsModels     packages.Module
	ErrorsController packages.Module
	ServicesImpl     packages.Module
	Converters       packages.Module
	Controllers      packages.Module
	versions         map[string]*VersionPackages
}

type VersionPackages struct {
	Main        packages.Module
	Models      packages.Module
	Services    packages.Module
	Controllers packages.Module
}

func (p *VersionPackages) ServicesApi(api *spec.Api) packages.Module {
	return p.Services.Subpackage(api.Name.SnakeCase())
}

func newVersionPackages(rootPackage packages.Module, version *spec.Version) *VersionPackages {
	mainPackage := rootPackage.Subpackage(version.Name.FlatCase())
	modelsPackage := mainPackage.Subpackage("models")
	servicePackage := mainPackage.Subpackage("services")
	controllersPackage := mainPackage.Subpackage("controllers")

	return &VersionPackages{
		mainPackage,
		modelsPackage,
		servicePackage,
		controllersPackage,
	}
}

func NewServicePackages(packageName, generatePath, servicesPath string) *ServicePackages {
	rootPackage := packages.New(generatePath, packageName)
	contentTypePackage := rootPackage.Subpackage("contenttype")
	jsonPackage := rootPackage.Subpackage("json")
	errorsPackage := rootPackage.Subpackage("errors")
	errorsModelsPackage := errorsPackage.Subpackage("models")
	errorsControllerPackage := rootPackage.Subpackage("controllers")
	servicesImplPackage := packages.New(servicesPath, packageName)
	convertersPackage := rootPackage.Subpackage("converters")
	controllerPackage := rootPackage.Subpackage("controllers")

	return &ServicePackages{
		rootPackage,
		contentTypePackage,
		jsonPackage,
		errorsPackage,
		errorsModelsPackage,
		errorsControllerPackage,
		servicesImplPackage,
		convertersPackage,
		controllerPackage,
		map[string]*VersionPackages{},
	}
}

func (p *ServicePackages) Version(version *spec.Version) *VersionPackages {
	versionPackages, found := p.versions[version.Name.Source]
	if !found {
		versionPackages = newVersionPackages(p.Root, version)
	}
	return versionPackages
}
