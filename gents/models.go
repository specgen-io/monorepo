package gents

import (
	"fmt"
	"github.com/specgen-io/spec"
	"specgen/gen"
)

func GenerateModels(serviceFile string, generatePath string, validation string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}
	files := generateModels(specification, validation, generatePath)
	return gen.WriteFiles(files, true)
}

func generateModels(specification *spec.Spec, validation string, generatePath string) []gen.TextFile {
	if validation == Superstruct {
		return generateSuperstructModels(specification, generatePath)
	}
	if validation == IoTs {
		return generateIoTsModels(specification, generatePath)
	}
	panic(fmt.Sprintf("Unknown validation: %s", validation))
}
