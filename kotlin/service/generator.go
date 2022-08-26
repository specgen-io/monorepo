package service

import (
	"fmt"

	"generator"
	"kotlin/models"
	"kotlin/modules"
	"kotlin/types"
	"spec"
)

type ServerGenerator interface {
	ServiceImports() []string
	ServicesControllers(version *spec.Version, mainPackage, thePackage, contentTypePackage, jsonPackage, modelsVersionPackage, serviceVersionPackage modules.Module) []generator.CodeFile
	ServiceImplAnnotation(api *spec.Api) (annotationImport, annotation string)
	ExceptionController(responses *spec.Responses, thePackage, errorsPackage, errorsModelsPackage, jsonPackage modules.Module) *generator.CodeFile
	Errors(thePackage, errorsModelsPackage, contentTypePackage, jsonPackage modules.Module) []generator.CodeFile
	ContentType(thePackage modules.Module) []generator.CodeFile
	JsonHelpers(thePackage modules.Module) []generator.CodeFile
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
