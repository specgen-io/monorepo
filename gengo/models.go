package gengo

import (
	"fmt"
	spec "github.com/specgen-io/spec.v2"
	"path/filepath"
	"specgen/gen"
	"strings"
)

func GenerateModels(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}

	fileName := specification.Name.SnakeCase() + "_models.go"
	modelsPath := filepath.Join(generatePath, fileName)
	models := generateModels(specification.ResolvedModels, "spec", modelsPath)

	err = gen.WriteFiles(models, true)
	return err
}

func versionedPackage(version spec.Name, packageName string) string {
	if version.Source != "" {
		return fmt.Sprintf("%s_%s", packageName, version.FlatCase())
	}
	return packageName
}

func generateModels(versionedModels spec.VersionedModels, packageName string, generatePath string) []gen.TextFile {
	files := []gen.TextFile{}

	for _, models := range versionedModels {
		w := NewGoWriter()
		w.Line("package %s", versionedPackage(models.Version, packageName))
		for _, model := range models.Models {
			w.Line("")
			if model.IsObject() {
				generateObjectModel(w, model)
			} else if model.IsOneOf() {
				generateEnumModel(w, model)
			} else if model.IsEnum() {
				generateOneOfModel(w, model)
			}
		}
		file := gen.TextFile{
			Path:    generatePath,
			Content: w.String(),
		}
		files = append(files, file)
	}

	return files
}

func generateObjectModel(w *gen.Writer, model spec.NamedModel) {
	w.Line("type %s struct {", model.Name.PascalCase())
	for _, field := range model.Object.Fields {
		typ := GoType(&field.Type.Definition)
		w.Line("  %s %s `json:\"%s\"`", field.Name.PascalCase(), typ, field.Name.Source)
	}
	w.Line("}")
}

func generateEnumModel(w *gen.Writer, model spec.NamedModel) {
	w.Line("type %s %s", model.Name.PascalCase(), "string")
	w.Line("")
	w.Line("const (")
	for _, enumItem := range model.Enum.Items {
		w.Line("  %s %s = %s", enumItem.Name.PascalCase(), model.Name.PascalCase(), enumItem.Value)
	}
	w.Line(")")
}

func generateOneOfModel(w *gen.Writer, model spec.NamedModel) {
	w.Line("")
	params := []string{}
	for _, item := range model.OneOf.Items {
		params = append(params, fmt.Sprintf("%s(%s)", GoType(&item.Type.Definition), item.Name.PascalCase()))
	}
	w.Line("var %s = []string{%s}", model.Name.PascalCase(), strings.Join(params, ", "))
}
