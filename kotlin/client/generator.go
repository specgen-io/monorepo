package client

import (
	"fmt"

	"generator"
	"kotlin/models"
	"kotlin/types"
	"spec"
)

type ClientGenerator interface {
	Clients(version *spec.Version) []generator.CodeFile
}

type Generator struct {
	Types  *types.Types
	Models models.Generator
	Client ClientGenerator
}

func NewGenerator(jsonlib, client string, packages *Packages) *Generator {
	types := models.NewTypes(jsonlib)
	models := models.NewGenerator(jsonlib, &(packages.Packages))

	if client == OkHttp {
		return &Generator{
			types,
			models,
			NewOkHttpGenerator(types, models, packages),
		}
	}
	if client == MicronautDecl {
		return &Generator{
			types,
			models,
			NewMicronautDeclGenerator(types, models, packages),
		}
	}
	if client == MicronautLow {
		return &Generator{
			types,
			models,
			NewMicronautLowGenerator(types, models, packages),
		}
	}

	panic(fmt.Sprintf(`Unsupported client: %s`, client))
}
