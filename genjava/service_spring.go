package genjava

import (
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
)

func GenerateService(serviceFile string, packageName string, swaggerPath string, generatePath string, servicesPath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}

	sourcesOverride := []gen.TextFile{}
	sourcesScaffold := []gen.TextFile{}

	modelsPackage := Module(generatePath, packageName)
	sourcesOverride = append(sourcesOverride, *generateJsoner(modelsPackage))
	for _, version := range specification.Versions {
		versionPackage := modelsPackage.Submodule(version.Version.FlatCase())
		sourcesOverride = append(sourcesOverride, generateVersionModels(&version, versionPackage)...)
	}

	err = gen.WriteFiles(sourcesOverride, true)
	if err != nil {
		return err
	}
	err = gen.WriteFiles(sourcesScaffold, false)
	return err
}
