package service

import (
	"generator"
	"golang/empty"
	"golang/models"
	"golang/types"
	"spec"
)

type ServiceGenerator interface {
	RootRouting(specification *spec.Spec) *generator.CodeFile
	HttpErrors(responses *spec.Responses) []generator.CodeFile
	CheckContentType() *generator.CodeFile
	Routings(version *spec.Version) []generator.CodeFile
	ResponseHelperFunctions() *generator.CodeFile
}

type Generator struct {
	models.Generator
	ServiceGenerator
	Types   *types.Types
	Modules *Modules
}

func NewGenerator(modules *Modules) *Generator {
	modelsGenerator := models.NewGenerator(&(modules.Modules))
	types := types.NewTypes()
	return &Generator{
		modelsGenerator,
		NewVestigoGenerator(types, modelsGenerator, modules),
		types,
		modules,
	}
}

func (g *Generator) AllStaticFiles() []generator.CodeFile {
	return []generator.CodeFile{
		*g.EnumsHelperFunctions(),
		*empty.GenerateEmpty(g.Modules.Empty),
		*generateParamsParser(g.Modules.ParamsParser),
		*g.ResponseHelperFunctions(),
		*g.CheckContentType(),
	}
}
