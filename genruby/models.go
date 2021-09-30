package genruby

import (
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
)

func GenerateModels(serviceFile string, generatePath string) error {
	result, err := spec.ReadSpecFile(serviceFile)
	if err != nil {
		return err
	}

	specification := result.Spec
	fileName := specification.Name.SnakeCase() + "_models.rb"
	moduleName := specification.Name.PascalCase()
	modelsPath := filepath.Join(generatePath, fileName)
	models := generateModels(specification, moduleName, modelsPath)

	err = gen.WriteFile(models, true)
	return err
}

func generateModels(specification *spec.Spec, moduleName string, generatePath string) *gen.TextFile {
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

	return &gen.TextFile{Path: generatePath, Content: w.String()}
}

func generateVersionModelsModule(w *gen.Writer, version *spec.Version, moduleName string) {
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

func generateObjectModel(w *gen.Writer, model *spec.NamedModel) {
	w.Line(`class %s`, model.Name.PascalCase())
	w.Line(`  include DataClass`)
	for _, field := range model.Object.Fields {
		typ := RubyType(&field.Type.Definition)
		w.Line(`  val :%s, %s`, field.Name.Source, typ)
	}
	w.Line(`end`)
}

func generateEnumModel(w *gen.Writer, model *spec.NamedModel) {
	w.Line(`class %s`, model.Name.PascalCase())
	w.Line(`  include Enum`)
	for _, enumItem := range model.Enum.Items {
		w.Line(`  define :%s, '%s'`, enumItem.Name.Source, enumItem.Value)
	}
	w.Line(`end`)
}

func generateOneOfModel(w *gen.Writer, model *spec.NamedModel) {
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
