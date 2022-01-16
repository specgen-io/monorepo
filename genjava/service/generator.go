package service

import (
	"github.com/specgen-io/specgen/v2/genjava/models"
	"github.com/specgen-io/specgen/v2/genjava/types"
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
