package gengo

import (
	spec "github.com/specgen-io/spec.v2"
	"path/filepath"
	"specgen/gen"
)

func GenerateModels(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}

	fileName := specification.Name.SnakeCase() + "_models.go"
	moduleName := specification.Name.PascalCase()
	modelsPath := filepath.Join(generatePath, fileName)
	models := generateModels(specification.ResolvedModels, moduleName, modelsPath)

	err = gen.WriteFile(models, true)
	return err
}

func generateModels(versionedModels spec.VersionedModels, moduleName string, generatePath string) *gen.TextFile {
	return nil
}
