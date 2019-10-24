package genscala

import (
	"github.com/ModaOperandi/spec"
	"specgen/gen"
)

func GenerateServiceModels(serviceFile string, sourceManagedPath string) (err error) {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return
	}

	modelsFile := GenerateCirceModels(specification, "models", sourceManagedPath)

	sourceManaged := []gen.TextFile{*modelsFile}

	err = gen.WriteFiles(sourceManaged, true)
	if err != nil {
		return
	}

	return
}
