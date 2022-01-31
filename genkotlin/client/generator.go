package client

import (
	"github.com/specgen-io/specgen/v2/genkotlin/models"
	"github.com/specgen-io/specgen/v2/genkotlin/types"
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
