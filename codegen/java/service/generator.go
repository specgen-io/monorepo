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
	ExceptionController(responses *spec.ErrorResponses) *generator.CodeFile
	ErrorsHelpers() *generator.CodeFile
	ContentType() []generator.CodeFile
}

type Generator struct {
	ServerGenerator
	models.Generator
	jsonlib  string
	Types    *types.Types
	Packages *Packages
}

func NewGenerator(jsonlib, server string, packages *Packages) *Generator {
	types := models.NewTypes(jsonlib)
	models := models.NewGenerator(jsonlib, &(packages.Packages))

	var serverGenerator ServerGenerator = nil
	switch server {
	case Spring:
		serverGenerator = NewSpringGenerator(types, models, packages)
		break
	case Micronaut:
		serverGenerator = NewMicronautGenerator(types, models, packages)
		break
	default:
		panic(fmt.Sprintf(`Unsupported server: %s`, server))
	}

	return &Generator{
		serverGenerator,
		models,
		jsonlib,
		types,
		packages,
	}
}
