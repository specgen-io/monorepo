package gents

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
)

func GenerateModels(serviceFile string, generatePath string, validation string) error {
	result, err := spec.ReadSpecFile(serviceFile)
	if err != nil {
		return err
	}

	specification := result.Spec
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
