package models

import (
	"generator"
	"spec"
)

type Generator interface {
	GenerateVersionModels(version *spec.Version) *generator.CodeFile
	GenerateErrorModels(httperrors *spec.HttpErrors) *generator.CodeFile
	EnumValuesStrings(model *spec.NamedModel) string
	GenerateEnumsHelperFunctions() *generator.CodeFile
}

func NewGenerator(modules *Modules) Generator {
	return NewEncodingJsonGenerator(modules)
}
