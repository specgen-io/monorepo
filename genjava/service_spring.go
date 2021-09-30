package genjava

import (
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/genopenapi"
)

func GenerateService(serviceFile string, packageName string, swaggerPath string, generatePath string, servicesPath string) error {
	result, err := spec.ReadSpecFile(serviceFile)
	if err != nil {
		return err
	}

	specification := result.Spec

	sourcesOverride := []gen.TextFile{}
	sourcesScaffold := []gen.TextFile{}

	mainPackage := Package(generatePath, packageName)

	modelsPackage := mainPackage.Subpackage("models")
	sourcesOverride = append(sourcesOverride, *generateJsoner(modelsPackage))

	for _, version := range specification.Versions {
		versionPackage := mainPackage.Subpackage(version.Version.FlatCase())

		modelsVersionPackage := versionPackage.Subpackage("models")
		sourcesOverride = append(sourcesOverride, generateVersionModels(&version, modelsVersionPackage)...)

		serviceVersionPackage := versionPackage.Subpackage("services")
		sourcesOverride = append(sourcesOverride, generateServicesInterfaces(&version, serviceVersionPackage, modelsVersionPackage)...)

		controllerVersionPackage := versionPackage.Subpackage("controllers")
		sourcesOverride = append(sourcesOverride, generateServicesControllers(&version, controllerVersionPackage, modelsPackage, modelsVersionPackage, serviceVersionPackage)...)
	}

	if swaggerPath != "" {
		sourcesOverride = append(sourcesOverride, *genopenapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		servicesImplPackage := Package(servicesPath, packageName)
		for _, version := range specification.Versions {
			servicesImplVersionPath := servicesImplPackage.Subpackage("services_implementations")
			serviceImplVersionPackage := servicesImplVersionPath.Subpackage(version.Version.FlatCase())

			versionPackage := mainPackage.Subpackage(version.Version.FlatCase())
			modelsVersionPackage := versionPackage.Subpackage("models")
			serviceVersionPackage := versionPackage.Subpackage("services")

			sourcesScaffold = append(sourcesScaffold, generateServicesImplementations(&version, serviceImplVersionPackage, modelsVersionPackage, serviceVersionPackage)...)
		}
	}

	err = gen.WriteFiles(sourcesOverride, true)
	if err != nil {
		return err
	}
	err = gen.WriteFiles(sourcesScaffold, false)
	return err
}
