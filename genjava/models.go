package genjava

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateModels(specification *spec.Spec, jsonlib string, packageName string, generatePath string) *sources.Sources {
	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}
	mainPackage := Package(generatePath, packageName)

	generator := NewModelsGenerator(jsonlib)

	modelsPackage := mainPackage.Subpackage("models")

	newSources := sources.NewSources()

	newSources.AddGeneratedAll(generator.SetupLibrary(modelsPackage))
	for _, version := range specification.Versions {
		versionPackage := modelsPackage.Subpackage(version.Version.FlatCase())
		newSources.AddGeneratedAll(generator.VersionModels(&version, versionPackage))
	}

	return newSources
}

type ModelsGenerator interface {
	SetupLibrary(thePackage Module) []sources.CodeFile
	VersionModels(version *spec.Version, thePackage Module) []sources.CodeFile
	ReadJson(jsonStr string, javaType string) string
	WriteJson(varData string) string
}

func NewModelsGenerator(jsonlib string) ModelsGenerator {
	types := NewTypes(jsonlib)
	if jsonlib == Jackson {
		return NewJacksonGenerator(types)
	}
	if jsonlib == Moshi {
		return NewMoshiGenerator(types)
	}
	panic(fmt.Sprintf(`Unsupported jsonlib: %s`, jsonlib))
}

func NewTypes(jsonlib string) *Types {
	if jsonlib == Jackson {
		return &Types{"JsonNode"}
	}
	if jsonlib == Moshi {
		return &Types{"Map<String, Object>"}
	}
	panic(fmt.Sprintf(`Unsupported jsonlib: %s`, jsonlib))
}
