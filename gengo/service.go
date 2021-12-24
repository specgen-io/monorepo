package gengo

import (
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/genopenapi"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateService(specification *spec.Spec, moduleName string, swaggerPath string, generatePath string, servicesPath string) error {
	sourcesOverride := []gen.TextFile{}
	sourcesScaffold := []gen.TextFile{}

	emptyModule := Module(moduleName, "empty")
	sourcesOverride = append(sourcesOverride, *generateEmpty(emptyModule))

	sourcesOverride = append(sourcesOverride, *generateSpecRouting(specification, moduleName, generatePath))

	for _, version := range specification.Versions {
		versionModule := Module(moduleName, createPath(generatePath, version.Version.FlatCase()))
		modelsModule := Module(moduleName, createPath(generatePath, version.Version.FlatCase(), modelsPackage))

		sourcesOverride = append(sourcesOverride, *generateParamsParser(versionModule))
		sourcesOverride = append(sourcesOverride, generateRoutings(&version, versionModule, modelsModule)...)
		sourcesOverride = append(sourcesOverride, generateServicesInterfaces(&version, versionModule, modelsModule, emptyModule)...)
		sourcesOverride = append(sourcesOverride, generateVersionModels(&version, modelsModule, emptyModule)...)
	}

	if swaggerPath != "" {
		sourcesOverride = append(sourcesOverride, *genopenapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		for _, version := range specification.Versions {
			versionPath := createPath(generatePath, version.Version.FlatCase())
			servicesVersionPath := createPath(servicesPath, version.Version.FlatCase())
			sourcesScaffold = append(sourcesScaffold, generateServicesImplementations(moduleName, versionPath, &version, servicesVersionPath)...)
		}
	}

	err := gen.WriteFiles(sourcesOverride, true)
	if err != nil {
		return err
	}
	err = gen.WriteFiles(sourcesScaffold, false)
	return err
}
