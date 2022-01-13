package genjava

import (
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateModels(specification *spec.Spec, jsonlib string, packageName string, generatePath string) *sources.Sources {
	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}
	mainPackage := Package(generatePath, packageName)

	generator := NewGenerator(jsonlib)

	modelsPackage := mainPackage.Subpackage("models")

	newSources := sources.NewSources()
	newSources.AddGeneratedAll(generator.generateModels(specification, modelsPackage))

	return newSources
}

type ModelsGenerator interface {
	SetupLibrary(thePackage Module) []sources.CodeFile
	ObjectModel(model *spec.NamedModel, thePackage Module) *sources.CodeFile
	OneOfModel(model *spec.NamedModel, thePackage Module) *sources.CodeFile
	EnumModel(model *spec.NamedModel, thePackage Module) *sources.CodeFile
}

func NewModelsGenerator(jsonlib string) ModelsGenerator {
	types := NewTypes(jsonlib)
	if jsonlib == Jackson {
		return NewJacksonGenerator(types)
	}
	if jsonlib == Moshi {
		return NewMoshiGenerator(types)
	}
	return nil
}

func (g *Generator) generateModels(specification *spec.Spec, thePackage Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, version := range specification.Versions {
		versionPackage := thePackage.Subpackage(version.Version.FlatCase())
		files = append(files, g.generateVersionModels(&version, versionPackage)...)
	}

	files = append(files, g.Models.SetupLibrary(thePackage)...)
	return files
}

func (g *Generator) generateVersionModels(version *spec.Version, thePackage Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, model := range version.ResolvedModels {
		if model.IsObject() {
			files = append(files, *g.Models.ObjectModel(model, thePackage))
		} else if model.IsOneOf() {
			files = append(files, *g.Models.OneOfModel(model, thePackage))
		} else if model.IsEnum() {
			files = append(files, *g.Models.EnumModel(model, thePackage))
		}
	}
	return files
}
