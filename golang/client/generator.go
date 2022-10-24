package client

import (
	"generator"
	"golang/models"
	"golang/types"
	"spec"
)

type ClientGenerator interface {
	GenerateClientsImplementations(version *spec.Version) []generator.CodeFile
	GenerateResponseFunctions() *generator.CodeFile
}

type Generator struct {
	models.Generator
	ClientGenerator
	Types   *types.Types
	Modules *Modules
}

func NewGenerator(modules *Modules) *Generator {
	types := types.NewTypes()
	return &Generator{
		models.NewGenerator(&(modules.Modules)),
		NewNetHttpGenerator(modules, types),
		types,
		modules,
	}
}

func (g *Generator) GenerateEmptyType() *generator.CodeFile {
	return types.GenerateEmpty(g.Modules.Empty)
}

func (g *Generator) AllStaticFiles() []generator.CodeFile {
	return []generator.CodeFile{
		*g.GenerateEnumsHelperFunctions(),
		*g.GenerateEmptyType(),
		*g.GenerateConverter(),
		*g.GenerateResponseFunctions(),
	}
}
