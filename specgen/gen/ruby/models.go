package ruby

import (
	"github.com/specgen-io/specgen/spec/v2"
	"github.com/specgen-io/specgen/v2/generator"
	"path/filepath"
)

func GenerateModels(specification *spec.Spec, generatePath string) *generator.Sources {
	sources := generator.NewSources()
	fileName := specification.Name.SnakeCase() + "_models.rb"
	moduleName := specification.Name.PascalCase()
	modelsPath := filepath.Join(generatePath, fileName)
	models := generateModels(specification, moduleName, modelsPath)

	sources.AddGenerated(models)
	return sources
}

func generateModels(specification *spec.Spec, moduleName string, generatePath string) *generator.CodeFile {
	w := NewRubyWriter()
	w.Line(`require "date"`)
	w.Line(`require "emery"`)

	for _, version := range specification.Versions {
		if version.Version.Source == "" {
			generateVersionModelsModule(w, &version, moduleName)
		}
	}

	for _, version := range specification.Versions {
		if version.Version.Source != "" {
			generateVersionModelsModule(w, &version, moduleName)
		}
	}

	return &generator.CodeFile{Path: generatePath, Content: w.String()}
}

func generateVersionModelsModule(w *generator.Writer, version *spec.Version, moduleName string) {
	w.EmptyLine()
	w.Line(`module %s`, versionedModule(moduleName, version.Version))
	for index, model := range version.ResolvedModels {
		if index != 0 {
			w.EmptyLine()
		}
		if model.IsObject() {
			generateObjectModel(w.Indented(), model)
		} else if model.IsOneOf() {
			generateOneOfModel(w.Indented(), model)
		} else if model.IsEnum() {
			generateEnumModel(w.Indented(), model)
		}
	}
	w.Line(`end`)
}

func generateObjectModel(w *generator.Writer, model *spec.NamedModel) {
	w.Line(`class %s`, model.Name.PascalCase())
	w.Line(`  include DataClass`)
	for _, field := range model.Object.Fields {
		typ := RubyType(&field.Type.Definition)
		w.Line(`  val :%s, %s`, field.Name.Source, typ)
	}
	w.Line(`end`)
}

func generateEnumModel(w *generator.Writer, model *spec.NamedModel) {
	w.Line(`class %s`, model.Name.PascalCase())
	w.Line(`  include Enum`)
	for _, enumItem := range model.Enum.Items {
		w.Line(`  define :%s, '%s'`, enumItem.Name.Source, enumItem.Value)
	}
	w.Line(`end`)
}

func generateOneOfModel(w *generator.Writer, model *spec.NamedModel) {
	w.Line(`class %s`, model.Name.PascalCase())
	w.Line(`  include TaggedUnion`)
	if model.OneOf.Discriminator != nil {
		w.Line(`  with_discriminator "%s"`, *model.OneOf.Discriminator)
	}
	for _, field := range model.OneOf.Items {
		typ := RubyType(&field.Type.Definition)
		w.Line(`  tag :%s, %s`, field.Name.Source, typ)
	}
	w.Line(`end`)

}
