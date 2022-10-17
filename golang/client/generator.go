package client

import (
	"generator"
	"golang/models"
	"golang/module"
	"spec"
)

type ClientGenerator interface {
	GenerateClientsImplementations(version *spec.Version, versionModule, convertModule, emptyModule, errorsModule, errorsModelsModule, modelsModule, respondModule module.Module) []generator.CodeFile
}

type Generator struct {
	Models models.Generator
	Client ClientGenerator
}

func NewGenerator(modules *models.Modules) *Generator {
	return &Generator{
		models.NewGenerator(modules),
		NewNetHttpGenerator(),
	}
}
