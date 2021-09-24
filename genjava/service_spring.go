package genjava

import (
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
)

func GenerateService(serviceFile string, packageName string, swaggerPath string, generatePath string, servicesPath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}

	generatePath = filepath.Join(generatePath, packageName)

	sourcesOverride := []gen.TextFile{}
	sourcesScaffold := []gen.TextFile{}

	sourcesOverride = append(sourcesOverride, *generateJsoner(packageName, generatePath))
	for _, version := range specification.Versions {
		versionedPackageName := versionedPackage(version.Version, packageName)
		versionPath := versionedPath(version.Version, generatePath)
		sourcesOverride = append(sourcesOverride, generateVersionModels(&version, versionedPackageName, versionPath)...)
	}

	err = gen.WriteFiles(sourcesOverride, true)
	if err != nil {
		return err
	}
	err = gen.WriteFiles(sourcesScaffold, false)
	return err
}
