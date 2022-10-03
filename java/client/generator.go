package client

import (
	"java/models"
	"java/types"
)

type Generator struct {
	models.Generator
	Jsonlib  string
	Types    *types.Types
	Packages *Packages
}

func NewGenerator(jsonlib string, packages *Packages) *Generator {
	return &Generator{
		models.NewGenerator(jsonlib, &(packages.Packages)),
		jsonlib,
		models.NewTypes(jsonlib),
		packages,
	}
}
