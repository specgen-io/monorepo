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

func NewGenerator(jsonlib, server, packageName, generatePath, servicesPath string) *Generator {
	types := models.NewTypes(jsonlib)
	models := models.NewGenerator(jsonlib)

	servicePackages := NewServicePackages(packageName, generatePath, servicesPath)

	var serverGenerator ServerGenerator = nil
	switch server {
	case Spring:
		serverGenerator = NewSpringGenerator(types, models, servicePackages)
		break
	case Micronaut:
		serverGenerator = NewMicronautGenerator(types, models, servicePackages)
		break
	default:
		panic(fmt.Sprintf(`Unsupported server: %s`, server))
	}

	return &Generator{
		serverGenerator,
		jsonlib,
		types,
		models,
		servicePackages,
	}
}

func (g *Generator) Models(version *spec.Version) []generator.CodeFile {
	return g.ModelsGenerator.Models(version, g.Packages.Version(version).Models, g.Packages.Json)
}

func (g *Generator) ErrorModels(httperrors *spec.HttpErrors) []generator.CodeFile {
	return g.ModelsGenerator.ErrorModels(httperrors, g.Packages.ErrorsModels, g.Packages.Json)
}

func (g *Generator) ModelsValidation() *generator.CodeFile {
	return g.ModelsGenerator.ModelsValidation(g.Packages.Errors, g.Packages.ErrorsModels, g.Packages.Json)
}
