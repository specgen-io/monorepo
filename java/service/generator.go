package service

import (
	"fmt"

	"generator"
	"java/models"
	"java/types"
	"spec"
)

type ServerGenerator interface {
	ServiceImports() []string
	ServiceImplAnnotation(api *spec.Api) (annotationImport, annotation string)
	ServicesControllers(version *spec.Version) []generator.CodeFile
	ExceptionController(responses *spec.Responses) *generator.CodeFile
	Errors() []generator.CodeFile
	ContentType() []generator.CodeFile
	JsonHelpers() []generator.CodeFile
}

type Generator struct {
	Jsonlib  string
	Types    *types.Types
	Models   models.Generator
	Packages *ServicePackages
	Server   ServerGenerator
}

func NewGenerator(jsonlib, server, packageName, generatePath, servicesPath string) *Generator {
	types := models.NewTypes(jsonlib)
	models := models.NewGenerator(jsonlib)

	servicePackages := NewServicePackages(packageName, generatePath, servicesPath)

	if server == Spring {
		return &Generator{
			jsonlib,
			types,
			models,
			servicePackages,
			NewSpringGenerator(types, models, servicePackages),
		}
	}
	if server == Micronaut {
		return &Generator{
			jsonlib,
			types,
			models,
			servicePackages,
			NewMicronautGenerator(types, models, servicePackages),
		}
	}

	panic(fmt.Sprintf(`Unsupported server: %s`, server))
}
