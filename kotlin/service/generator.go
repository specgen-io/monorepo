package service

import (
	"fmt"

	"generator"
	"kotlin/models"
	"kotlin/packages"
	"kotlin/types"
	"spec"
)

type ServerGenerator interface {
	ServiceImports() []string
	ServicesControllers(version *spec.Version, mainPackage, thePackage, contentTypePackage, jsonPackage, modelsVersionPackage, errorModelsPackage, serviceVersionPackage packages.Package) []generator.CodeFile
	ServiceImplAnnotation(api *spec.Api) (annotationImport, annotation string)
	ExceptionController(responses *spec.Responses, thePackage, errorsPackage, errorsModelsPackage, jsonPackage packages.Package) *generator.CodeFile
	Errors(thePackage, errorsModelsPackage, contentTypePackage, jsonPackage packages.Package) []generator.CodeFile
	ContentType(thePackage packages.Package) []generator.CodeFile
	JsonHelpers(thePackage packages.Package) []generator.CodeFile
}

type Generator struct {
	Jsonlib string
	Types   *types.Types
	Models  models.Generator
	Server  ServerGenerator
}

func NewGenerator(jsonlib, server string, packages *models.Packages) *Generator {
	types := models.NewTypes(jsonlib)
	models := models.NewGenerator(jsonlib, packages)

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
