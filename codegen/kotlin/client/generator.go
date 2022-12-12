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
	Utils() []generator.CodeFile
	Exceptions(errors *spec.ErrorResponses) []generator.CodeFile
}

type Generator struct {
	ClientGenerator
	models.Generator
	Jsonlib  string
	Types    *types.Types
	Packages *Packages
}

func NewGenerator(jsonlib, client string, packages *Packages) *Generator {
	types := models.NewTypes(jsonlib)
	models := models.NewGenerator(jsonlib, &(packages.Packages))

	var clientGenerator ClientGenerator = nil
	switch client {
	case OkHttp:
		clientGenerator = NewOkHttpGenerator(types, models, packages)
		break
	case Micronaut:
		clientGenerator = NewMicronautGenerator(types, models, packages)
		break
	case MicronautDecl:
		clientGenerator = NewMicronautDeclGenerator(types, models, packages)
		break
	default:
		panic(fmt.Sprintf(`Unsupported client: %s`, client))
	}

	return &Generator{
		clientGenerator,
		models,
		jsonlib,
		types,
		packages,
	}
}
