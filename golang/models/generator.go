package models

import (
	"generator"
	"golang/module"
	"spec"
)

type Generator interface {
	GenerateVersionModels(models []*spec.NamedModel, module, enumsModule module.Module) *generator.CodeFile
	EnumValuesStrings(model *spec.NamedModel) string
	GenerateEnumsHelperFunctions(module module.Module) *generator.CodeFile
}

func NewGenerator() Generator {
	return NewEncodingJsonGenerator()
}
