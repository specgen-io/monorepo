package service

import (
	"fmt"
	"generator"
	"java/models"
	"java/types"
	"spec"
)

type ServerGenerator interface {
	FilesImports() []string
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
	var types *types.Types
	modelsGenerator := models.NewGenerator(jsonlib, &(packages.Packages))

	var serverGenerator ServerGenerator = nil
	switch server {
	case Spring:
		types = models.NewTypes(jsonlib, "Resource", "Resource")
		serverGenerator = NewSpringGenerator(types, modelsGenerator, packages)
		break
	case Micronaut:
		types = models.NewTypes(jsonlib, "byte[]", "byte[]")
		serverGenerator = NewMicronautGenerator(types, modelsGenerator, packages)
		break
	default:
		panic(fmt.Sprintf(`Unsupported server: %s`, server))
	}

	return &Generator{
		serverGenerator,
		modelsGenerator,
		jsonlib,
		types,
		packages,
	}
}
