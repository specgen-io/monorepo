package service

import (
	"fmt"
	"generator"
	"kotlin/models"
	"kotlin/types"
	"spec"
)

type ServerGenerator interface {
	FilesImports() []string
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
	var types *types.Types
	modelsGenerator := models.NewGenerator(jsonlib, &(packages.Packages))
	
	var serverGenerator ServerGenerator = nil
	switch server {
	case Spring:
		types = models.NewTypes(jsonlib, "Resource", "Resource", "MultipartFile", "Resource")
		serverGenerator = NewSpringGenerator(types, modelsGenerator, packages)
		break
	case Micronaut:
		types = models.NewTypes(jsonlib, "ByteArray", "ByteArray", "CompletedFileUpload", "StreamedFile")
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
