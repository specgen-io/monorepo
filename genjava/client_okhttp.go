package genjava

import (
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateClient(specification *spec.Spec, packageName string, generatePath string) error {
	generatedFiles := []gen.TextFile{}

	mainPackage := Package(generatePath, packageName)

	generatedFiles = append(generatedFiles, *generateClientException(mainPackage))

	utilsPackage := mainPackage.Subpackage("utils")
	generatedFiles = append(generatedFiles, generateUtils(utilsPackage)...)

	modelsPackage := mainPackage.Subpackage("models")
	generatedFiles = append(generatedFiles, *generateJson(modelsPackage))

	for _, version := range specification.Versions {
		versionPackage := mainPackage.Subpackage(version.Version.FlatCase())

		modelsVersionPackage := versionPackage.Subpackage("models")
		generatedFiles = append(generatedFiles, generateVersionModels(&version, modelsVersionPackage)...)

		clientVersionPackage := versionPackage.Subpackage("clients")
		generatedFiles = append(generatedFiles, generateClientsImplementations(&version, clientVersionPackage, modelsVersionPackage, utilsPackage, mainPackage)...)
	}
	err := gen.WriteFiles(generatedFiles, true)

	return err
}
