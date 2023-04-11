package service

import (
	"fmt"
	"generator"
	"golang/models"
	"golang/types"
	"spec"
)

type ServerGenerator interface {
	RootRouting(specification *spec.Spec) *generator.CodeFile
	Routings(version *spec.Version) []generator.CodeFile
	GenerateUrlParamsCtor() *generator.CodeFile
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
	case HttpRouter:
		serverGenerator = NewHttpRouterGenerator(types, models, modules)
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

//func (g *Generator) AllStaticFiles() []generator.CodeFile {
//	return []generator.CodeFile{
//		*empty.GenerateEmpty(g.Modules.Empty),
//		*g.EnumsHelperFunctions(),
//		*g.ResponseHelperFunctions(),
//		*g.CheckContentType(),
//		*g.GenerateParamsParser(),
//	}
//}
