package service

import (
	"fmt"
	"generator"
	"golang/empty"
	"golang/models"
	"golang/types"
	"spec"
)

type ServerGenerator interface {
	RootRouting(specification *spec.Spec) *generator.CodeFile
	HttpErrors(responses *spec.ErrorResponses) []generator.CodeFile
	CheckContentType() *generator.CodeFile
	Routings(version *spec.Version) []generator.CodeFile
	ResponseHelperFunctions() *generator.CodeFile
}

type Generator struct {
	ServerGenerator
	models.Generator
	Types   *types.Types
	Modules *Modules
}

func NewGenerator(server string, modules *Modules) *Generator {
	types := types.NewTypes()
	models := models.NewGenerator(&(modules.Modules))

	var serverGenerator ServerGenerator = nil
	switch server {
	case Vestigo:
		serverGenerator = NewVestigoGenerator(types, models, modules)
		break
	default:
		panic(fmt.Sprintf(`Unsupported server: %s`, server))
	}

	return &Generator{
		serverGenerator,
		models,
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
