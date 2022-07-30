package client

import (
	"github.com/specgen-io/specgen/java/v2/models"
	"github.com/specgen-io/specgen/java/v2/types"
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
