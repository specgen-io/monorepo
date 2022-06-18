package models

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen/java/packages"
	"github.com/specgen-io/specgen/v2/gen/java/types"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

type Generator interface {
	SetupLibrary(thePackage packages.Module) []sources.CodeFile
	SetupImport(jsonPackage packages.Module) string
	VersionModels(version *spec.Version, thePackage packages.Module, jsonPackage packages.Module) []sources.CodeFile
	ReadJson(jsonStr string, typ *spec.TypeDef) (string, string)
	WriteJson(varData string, typ *spec.TypeDef) (string, string)
	WriteJsonNoCheckedException(varData string, typ *spec.TypeDef) string
	CreateJsonMapperField(w *sources.Writer, annotation string)
	InitJsonMapper(w *sources.Writer)
	JsonImports() []string

	GenerateBodyBadRequestErrorCreator(thePackage, modelsPackage packages.Module) *sources.CodeFile
	CreateBodyBadRequestError(exceptionVar string) string
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
		return &types.Types{"Map<String, Object>"}
	}
	panic(fmt.Sprintf(`Unsupported jsonlib: %s`, jsonlib))
}
