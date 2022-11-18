package generators

import (
	"path/filepath"
	"ruby/writer"

	"generator"
	"spec"
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
	w := writer.New(generatePath)
	w.Line(`require "date"`)
	w.Line(`require "emery"`)

	for _, version := range specification.Versions {
		if version.Name.Source == "" {
			generateVersionModelsModule(w, &version, moduleName)
		}
	}

	for _, version := range specification.Versions {
		if version.Name.Source != "" {
			generateVersionModelsModule(w, &version, moduleName)
		}
	}

	return w.ToCodeFile()
}

func generateVersionModelsModule(w *writer.Writer, version *spec.Version, moduleName string) {
	w.EmptyLine()
	w.Line(`module %s`, versionedModule(moduleName, version.Name))
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

func generateObjectModel(w *writer.Writer, model *spec.NamedModel) {
	w.Line(`class %s`, model.Name.PascalCase())
	w.Line(`  include DataClass`)
	for _, field := range model.Object.Fields {
		typ := RubyType(&field.Type.Definition)
		w.Line(`  val :%s, %s`, field.Name.Source, typ)
	}
	w.Line(`end`)
}

func generateEnumModel(w *writer.Writer, model *spec.NamedModel) {
	w.Line(`class %s`, model.Name.PascalCase())
	w.Line(`  include Enum`)
	for _, enumItem := range model.Enum.Items {
		w.Line(`  define :%s, '%s'`, enumItem.Name.Source, enumItem.Value)
	}
	w.Line(`end`)
}

func generateOneOfModel(w *writer.Writer, model *spec.NamedModel) {
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
