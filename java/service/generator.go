package service

import (
	"fmt"

	"generator"
	"java/models"
	"java/packages"
	"java/types"
	"spec"
)

type ServerGenerator interface {
	ServiceImports() []string
	ServiceImplAnnotation(api *spec.Api) (annotationImport, annotation string)
	ServicesControllers(version *spec.Version, mainPackage, thePackage, contentTypePackage, jsonPackage, modelsVersionPackage, serviceVersionPackage packages.Module) []generator.CodeFile
	ExceptionController(responses *spec.Responses, thePackage, errorsPackage, errorsModelsPackage, jsonPackage packages.Module) *generator.CodeFile
	Errors(thePackage, errorsModelsPackage, contentTypePackage, jsonPackage packages.Module) []generator.CodeFile
	ContentType(thePackage packages.Module) []generator.CodeFile
	JsonHelpers(thePackage packages.Module) []generator.CodeFile
}

type Generator struct {
	Jsonlib string
	Types   *types.Types
	Models  models.Generator
	Server  ServerGenerator
}

func NewGenerator(jsonlib, server string) *Generator {
	types := models.NewTypes(jsonlib)
	models := models.NewGenerator(jsonlib)

	if server == Spring {
		return &Generator{
			jsonlib,
			types,
			models,
			NewSpringGenerator(types, models),
		}
	}
	if server == Micronaut {
		return &Generator{
			jsonlib,
			types,
			models,
			NewMicronautGenerator(types, models),
		}
	}

	panic(fmt.Sprintf(`Unsupported server: %s`, server))
}
