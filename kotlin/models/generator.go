package models

import (
	"fmt"

	"generator"
	"kotlin/modules"
	"kotlin/types"
	"spec"
)

type Generator interface {
	Models(models []*spec.NamedModel, thePackage, jsonPackage modules.Module) []generator.CodeFile
	ModelsDefinitionsImports() []string
	ModelsUsageImports() []string
	SetupImport(jsonPackage modules.Module) string
	SetupLibrary(thePackage modules.Module) []generator.CodeFile
	JsonHelpersMethods() string
	ValidationErrorsHelpers(thePackage, errorsModelsPackage, jsonPackage modules.Module) *generator.CodeFile
	CreateJsonMapperField(annotation string) string
	InitJsonMapper(w *generator.Writer)

	JsonRead(varJson string, typ *spec.TypeDef) string
	JsonWrite(varData string, typ *spec.TypeDef) string

	ReadJson(jsonStr string, typ *spec.TypeDef) (string, string)
	WriteJson(varData string, typ *spec.TypeDef) (string, string)
}

func NewGenerator(jsonlib string) Generator {
	types := NewTypes(jsonlib)
	if jsonlib == Jackson {
		return NewJacksonGenerator(types)
	}
	if jsonlib == Moshi {
		return NewMoshiGenerator(types)
	}
	panic(fmt.Sprintf(`Unsupported jsonlib: %s`, jsonlib))
}

func NewTypes(jsonlib string) *types.Types {
	if jsonlib == Jackson {
		return &types.Types{"JsonNode"}
	}
	if jsonlib == Moshi {
		return &types.Types{"Map<String, Any>"}
	}
	panic(fmt.Sprintf(`Unsupported jsonlib: %s`, jsonlib))
}
