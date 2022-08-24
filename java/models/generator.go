package models

import (
	"fmt"

	"generator"
	"java/packages"
	"java/types"
	"spec"
)

type Generator interface {
	Models(models []*spec.NamedModel, thePackage packages.Module, jsonPackage packages.Module) []generator.CodeFile
	ModelsDefinitionsImports() []string
	ModelsUsageImports() []string
	SetupImport(jsonPackage packages.Module) string
	SetupLibrary(thePackage packages.Module) []generator.CodeFile
	JsonHelpersMethods() string
	JsonParseException(thePackage packages.Module) *generator.CodeFile
	ValidationErrorsHelpers(thePackage, errorsModelsPackage, jsonPackage packages.Module) *generator.CodeFile
	CreateJsonMapperField(w *generator.Writer, annotation string)
	InitJsonMapper(w *generator.Writer)

	JsonRead(varJson string, typ *spec.TypeDef) string
	JsonWrite(varData string, typ *spec.TypeDef) string

	ReadJson(jsonStr string, typ *spec.TypeDef) (string, string)
	WriteJson(varData string, typ *spec.TypeDef) (string, string)
	WriteJsonNoCheckedException(varData string, typ *spec.TypeDef) string
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
		return &types.Types{RawJsonType: "JsonNode"}
	}
	if jsonlib == Moshi {
		return &types.Types{RawJsonType: "Map<String, Object>"}
	}
	panic(fmt.Sprintf(`Unsupported jsonlib: %s`, jsonlib))
}
