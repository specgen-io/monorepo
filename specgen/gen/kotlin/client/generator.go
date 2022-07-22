package client

import (
	"fmt"
	"github.com/specgen-io/specgen/spec/v2"
	"github.com/specgen-io/specgen/specgen/v2/gen/kotlin/models"
	"github.com/specgen-io/specgen/specgen/v2/gen/kotlin/modules"
	"github.com/specgen-io/specgen/specgen/v2/gen/kotlin/types"
	"github.com/specgen-io/specgen/specgen/v2/generator"
)

type ClientGenerator interface {
	ClientImplementation(version *spec.Version, thePackage modules.Module, modelsVersionPackage modules.Module, jsonPackage modules.Module, mainPackage modules.Module) []generator.CodeFile
}

type Generator struct {
	Types  *types.Types
	Models models.Generator
	Client ClientGenerator
}

func NewGenerator(jsonlib, client string) *Generator {
	types := models.NewTypes(jsonlib)
	models := models.NewGenerator(jsonlib)

	if client == OkHttp {
		return &Generator{
			types,
			models,
			NewOkHttpGenerator(types, models),
		}
	}
	if client == MicronautDecl {
		return &Generator{
			types,
			models,
			NewMicronautDeclGenerator(types, models),
		}
	}
	if client == MicronautLow {
		return &Generator{
			types,
			models,
			NewMicronautLowGenerator(types, models),
		}
	}

	panic(fmt.Sprintf(`Unsupported client: %s`, client))
}
