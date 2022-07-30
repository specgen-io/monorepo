package service

import (
	"fmt"

	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/java/v2/models"
	"github.com/specgen-io/specgen/java/v2/packages"
	"github.com/specgen-io/specgen/java/v2/types"
	"github.com/specgen-io/specgen/spec/v2"
)

type ServerGenerator interface {
	ServicesControllers(version *spec.Version, mainPackage, thePackage, jsonPackage, modelsVersionPackage, serviceVersionPackage packages.Module) []generator.CodeFile
	ServiceImplAnnotation(api *spec.Api) (annotationImport, annotation string)
}

type Generator struct {
	Jsonlib string
	Types   *types.Types
	Models  models.Generator
	Server  ServerGenerator
}

func NewGenerator(jsonlib, server string) *Generator {
	types := models.NewTypes(jsonlib)
	models := models.NewGenerator(jsonlib)

	if server == Spring {
		return &Generator{
			jsonlib,
			types,
			models,
			NewSpringGenerator(types, models),
		}
	}
	if server == Micronaut {
		return &Generator{
			jsonlib,
			types,
			models,
			NewMicronautGenerator(types, models),
		}
	}

	panic(fmt.Sprintf(`Unsupported server: %s`, server))
}
