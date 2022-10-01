package models

import (
	"fmt"

	"generator"
	"java/packages"
	"java/types"
	"spec"
)

type Generator interface {
	Models(models []*spec.NamedModel, thePackage packages.Package, jsonPackage packages.Package) []generator.CodeFile
	ModelsDefinitionsImports() []string
	ModelsUsageImports() []string
	SetupImport(jsonPackage packages.Package) string
	SetupLibrary(thePackage packages.Package) []generator.CodeFile
	JsonHelpersMethods() string
	JsonParseException(thePackage packages.Package) *generator.CodeFile
	ValidationErrorsHelpers(thePackage, errorsModelsPackage, jsonPackage packages.Package) *generator.CodeFile
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
