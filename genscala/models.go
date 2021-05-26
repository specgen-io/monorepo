package genscala

import (
	spec "github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/gen"
	"path/filepath"
)

func GenerateServiceModels(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}

	modelsPackage := modelsPackageName(specification.Name)

	scalaCirceFile := generateJson(modelsPackage, filepath.Join(generatePath, "Json.scala"))

	modelsFiles := GenerateCirceModels(specification, modelsPackage, generatePath)

	sourceManaged := append(modelsFiles, *scalaCirceFile)

	err = gen.WriteFiles(sourceManaged, true)
	if err != nil {
		return err
	}

	return nil
}

func modelsPackageName(name spec.Name) string {
	return name.FlatCase() + ".models"
}