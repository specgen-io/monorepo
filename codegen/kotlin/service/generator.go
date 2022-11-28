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
	ExceptionController(responses *spec.ErrorResponses) *generator.CodeFile
	ErrorsHelpers() *generator.CodeFile
	ContentType() []generator.CodeFile
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
