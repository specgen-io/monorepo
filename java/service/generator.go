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
	Errors(models []*spec.NamedModel) []generator.CodeFile
	ContentType() []generator.CodeFile
	JsonHelpers() []generator.CodeFile
}

type Generator struct {
	ServerGenerator
	jsonlib         string
	Types           *types.Types
	ModelsGenerator models.Generator
	Packages        *ServicePackages
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
	return g.ModelsGenerator.Models(version.ResolvedModels, g.Packages.Version(version).Models, g.Packages.Json)
}
