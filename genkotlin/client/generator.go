package client

import (
	"github.com/specgen-io/specgen/v2/genkotlin/models"
	"github.com/specgen-io/specgen/v2/genkotlin/types"
)

type Generator struct {
	Types  *types.Types
	Models models.Generator
}

func NewGenerator(jsonlib string) *Generator {
	return &Generator{
		models.NewTypes(jsonlib),
		models.NewGenerator(jsonlib),
	}
}
