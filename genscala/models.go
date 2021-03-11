package genscala

import (
	"github.com/specgen-io/spec"
	"specgen/gen"
	"specgen/static"
)

func GenerateServiceModels(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}

	modelsPackage := modelsPackageName(specification.ServiceName)

	scalaStaticCode := static.ScalaStaticCode{ PackageName: modelsPackage }

	scalaCirceFiles, err := static.RenderTemplate("scala-circe", generatePath, scalaStaticCode)
	if err != nil {
		return err
	}


	modelsFile := GenerateCirceModels(specification, modelsPackage, generatePath)

	sourceManaged := append(scalaCirceFiles, *modelsFile)

	err = gen.WriteFiles(sourceManaged, true)
	if err != nil {
		return err
	}

	return nil
}

func modelsPackageName(name spec.Name) string {
	return name.FlatCase() + ".models"
}