package client

import (
	"fmt"

	"generator"
	"kotlin/models"
	"kotlin/packages"
	"kotlin/types"
	"spec"
)

type ClientGenerator interface {
	Clients(version *spec.Version, thePackage packages.Package, modelsVersionPackage packages.Package, errorModelsPackage packages.Package, jsonPackage packages.Package, mainPackage packages.Package) []generator.CodeFile
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
