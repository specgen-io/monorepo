package models

import (
	"fmt"
	"github.com/specgen-io/specgen/spec/v2"
	"github.com/specgen-io/specgen/specgen/v2/gen/kotlin/modules"
	"github.com/specgen-io/specgen/specgen/v2/gen/kotlin/types"
	"github.com/specgen-io/specgen/specgen/v2/generator"
)

type Generator interface {
	SetupLibrary(thePackage modules.Module) []generator.CodeFile
	SetupImport(jsonPackage modules.Module) string
	VersionModels(version *spec.Version, thePackage modules.Module, jsonPackage modules.Module) []generator.CodeFile
	ReadJson(jsonStr string, typ *spec.TypeDef) (string, string)
	WriteJson(varData string, typ *spec.TypeDef) (string, string)
	CreateJsonMapperField() string
	InitJsonMapper(w *generator.Writer)
	ModelsDefinitionsImports() []string
	ModelsUsageImports() []string

	GenerateJsonParseException(thePackage, modelsPackage modules.Module) *generator.CodeFile
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
