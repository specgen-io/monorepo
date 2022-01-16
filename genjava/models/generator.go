package models

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/genjava/packages"
	"github.com/specgen-io/specgen/v2/genjava/types"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

type Generator interface {
	SetupLibrary(thePackage packages.Module) []sources.CodeFile
	VersionModels(version *spec.Version, thePackage packages.Module) []sources.CodeFile
	ReadJson(jsonStr string, javaType string) string
	WriteJson(varData string) string
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
