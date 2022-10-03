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
	ErrorsHelpers() *generator.CodeFile
	ContentType() []generator.CodeFile
	JsonHelpers() []generator.CodeFile
}

type Generator struct {
	ServerGenerator
	jsonlib         string
	Types           *types.Types
	ModelsGenerator models.Generator
	Packages        *Packages
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
		jsonlib,
		types,
		models,
		packages,
	}
}

func (g *Generator) Models(version *spec.Version) []generator.CodeFile {
	return g.ModelsGenerator.Models(version, g.Packages.Models(version), g.Packages.Json)
}

func (g *Generator) ErrorModels(httperrors *spec.HttpErrors) []generator.CodeFile {
	return g.ModelsGenerator.ErrorModels(httperrors, g.Packages.ErrorsModels, g.Packages.Json)
}

func (g *Generator) ModelsValidation() *generator.CodeFile {
	return g.ModelsGenerator.ModelsValidation(g.Packages.Errors, g.Packages.ErrorsModels, g.Packages.Json)
}
