package client

import (
	"generator"
	"golang/models"
	"golang/module"
	"golang/types"
	"spec"
)

type ClientGenerator interface {
	GenerateClientsImplementations(version *spec.Version, versionModule, convertModule, emptyModule, errorsModule, errorsModelsModule, modelsModule, respondModule module.Module) []generator.CodeFile
}

type Generator struct {
	Types  *types.Types
	Models models.Generator
	Client ClientGenerator
}

func NewGenerator(modules *models.Modules) *Generator {
	types := types.NewTypes()
	return &Generator{
		types,
		models.NewGenerator(modules),
		NewNetHttpGenerator(types),
	}
}
