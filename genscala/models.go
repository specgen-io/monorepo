package genscala

import (
	"github.com/ModaOperandi/spec"
	"specgen/gen"
)

func GenerateServiceModels(serviceFile string, generatePath string) (err error) {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return
	}

	jsonFile := GenerateJsonObject("models", generatePath)

	modelsFile := GenerateCirceModels(specification, "models", generatePath)

	sourceManaged := []gen.TextFile{*modelsFile, *jsonFile}

	err = gen.WriteFiles(sourceManaged, true)
	if err != nil {
		return
	}

	return
}
