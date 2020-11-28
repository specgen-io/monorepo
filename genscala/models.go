package genscala

import (
	"github.com/ModaOperandi/spec"
	"specgen/gen"
	"specgen/static"
)

func GenerateServiceModels(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}

	scalaStaticCode := static.ScalaStaticCode{ PackageName: "models" }

	scalaCirceFiles, err := static.RenderTemplate("scala-circe", generatePath, scalaStaticCode)
	if err != nil {
		return err
	}


	modelsFile := GenerateCirceModels(specification, "models", generatePath)

	sourceManaged := append(scalaCirceFiles, *modelsFile)

	err = gen.WriteFiles(sourceManaged, true)
	if err != nil {
		return err
	}

	return nil
}
