package genscala

import (
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateServiceModels(specification *spec.Spec, packageName string, generatePath string) error {
	if packageName == "" {
		packageName = specification.Name.FlatCase() + ".models"
	}
	modelsPackage := NewPackage(generatePath, packageName, "")

	modelsFiles := GenerateCirceModels(specification, modelsPackage)

	err := gen.WriteFiles(modelsFiles, true)
	if err != nil {
		return err
	}

	return nil
}
