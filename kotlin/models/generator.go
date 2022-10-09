package models

import (
	"fmt"

	"generator"
	"kotlin/types"
	"spec"
)

type Generator interface {
	Models(version *spec.Version) []generator.CodeFile
	ErrorModels(httperrors *spec.HttpErrors) []generator.CodeFile
	ModelsDefinitionsImports() []string
	ModelsUsageImports() []string
	SetupLibrary() []generator.CodeFile
	JsonHelpersMethods() string
	ValidationErrorsHelpers() *generator.CodeFile
	CreateJsonMapperField(annotation string) string
	InitJsonMapper(w generator.Writer)

	JsonRead(varJson string, typ *spec.TypeDef) string
	JsonWrite(varData string, typ *spec.TypeDef) string

	ReadJson(jsonStr string, typ *spec.TypeDef) (string, string)
	WriteJson(varData string, typ *spec.TypeDef) (string, string)
}

func NewGenerator(jsonlib string, packages *Packages) Generator {
	types := NewTypes(jsonlib)
	if jsonlib == Jackson {
		return NewJacksonGenerator(types, packages)
	}
	if jsonlib == Moshi {
		return NewMoshiGenerator(types, packages)
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
