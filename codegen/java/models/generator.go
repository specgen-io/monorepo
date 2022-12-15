package models

import (
	"fmt"
	"generator"
	"java/types"
	"spec"
)

type Generator interface {
	Models(version *spec.Version) []generator.CodeFile
	ErrorModels(httperrors *spec.HttpErrors) []generator.CodeFile
	ModelsUsageImports() []string
	ValidationErrorsHelpers() *generator.CodeFile
	JsonRead(varJson string, typ *spec.TypeDef) string
	JsonWrite(varData string, typ *spec.TypeDef) string
	JsonHelpers() []generator.CodeFile

	JsonMapperInit() string
	JsonMapperType() string
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
		return &types.Types{RawJsonType: "JsonNode"}
	}
	if jsonlib == Moshi {
		return &types.Types{RawJsonType: "Map<String, Object>"}
	}
	panic(fmt.Sprintf(`Unsupported jsonlib: %s`, jsonlib))
}
