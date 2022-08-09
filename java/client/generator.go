package client

import (
	"java/models"
	"java/types"
)

type Generator struct {
	Jsonlib string
	Types   *types.Types
	Models  models.Generator
}

func NewGenerator(jsonlib string) *Generator {
	return &Generator{
		jsonlib,
		models.NewTypes(jsonlib),
		models.NewGenerator(jsonlib),
	}
}
