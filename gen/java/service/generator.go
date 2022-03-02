package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen/java/models"
	"github.com/specgen-io/specgen/v2/gen/java/packages"
	"github.com/specgen-io/specgen/v2/gen/java/types"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

type ServerGenerator interface {
	ServicesControllers(version *spec.Version, thePackage packages.Module, jsonPackage packages.Module, modelsVersionPackage packages.Module, serviceVersionPackage packages.Module) []sources.CodeFile
	ServicesImplementations(version *spec.Version, thePackage packages.Module, modelsVersionPackage packages.Module, servicesVersionPackage packages.Module) []sources.CodeFile
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
	panic(fmt.Sprintf(`Unsupported server: %s`, server))
}
