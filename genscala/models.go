package genscala

import (
	spec "github.com/specgen-io/spec.v2"
	"specgen/gen"
	"specgen/static"
)

func GenerateServiceModels(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}

	modelsPackage := modelsPackageName(specification.Name)

	scalaStaticCode := static.ScalaStaticCode{ PackageName: modelsPackage }

	scalaCirceFiles, err := static.RenderTemplate("scala-circe", generatePath, scalaStaticCode)
	if err != nil {
		return err
	}


	modelsFiles := GenerateCirceModels(specification, modelsPackage, generatePath)

	sourceManaged := append(scalaCirceFiles, modelsFiles...)

	err = gen.WriteFiles(sourceManaged, true)
	if err != nil {
		return err
	}

	return nil
}

func modelsPackageName(name spec.Name) string {
	return name.FlatCase() + ".models"
}