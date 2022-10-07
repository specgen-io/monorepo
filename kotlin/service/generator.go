package service

import (
	"fmt"

	"generator"
	"kotlin/models"
	"kotlin/types"
	"spec"
)

type ServerGenerator interface {
	ServiceImports() []string
	ServicesControllers(version *spec.Version) []generator.CodeFile
	ServiceImplAnnotation(api *spec.Api) (annotationImport, annotation string)
	ExceptionController(responses *spec.Responses) *generator.CodeFile
	Errors() []generator.CodeFile
	ContentType() []generator.CodeFile
	JsonHelpers() []generator.CodeFile
}

type Generator struct {
	ServerGenerator
	models.Generator
	Jsonlib  string
	Types    *types.Types
	Packages *Packages
}

func NewGenerator(jsonlib, server string, packages *Packages) *Generator {
	types := models.NewTypes(jsonlib)
	models := models.NewGenerator(jsonlib, &(packages.Packages))

	if server == Spring {
		return &Generator{
			NewSpringGenerator(types, models, packages),
			models,
			jsonlib,
			types,
			packages,
		}
	}
	if server == Micronaut {
		return &Generator{
			NewMicronautGenerator(types, models, packages),
			models,
			jsonlib,
			types,
			packages,
		}
	}

	panic(fmt.Sprintf(`Unsupported server: %s`, server))
}
