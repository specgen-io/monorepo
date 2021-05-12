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

	modelsPath := filepath.Join(generatePath, "models.go")
	models := generateModels(specification, "spec", modelsPath)

	err = gen.WriteFiles(models, true)
	return err
}

func versionedPackage(version spec.Name, packageName string) string {
	if version.Source != "" {
		return fmt.Sprintf("%s_%s", packageName, version.FlatCase())
	}
	return packageName
}

func generateModels(specification *spec.Spec, packageName string, generatePath string) []gen.TextFile {
	files := []gen.TextFile{}

	for _, version := range specification.Versions {
		w := NewGoWriter()
		w.Line("package %s", versionedPackage(version.Version, packageName))
		w.Line("")
		w.Line("import (")
		w.Line(`  "encoding/json"`)
		w.Line(`  "github.com/google/uuid"`)
		w.Line(`  "github.com/shopspring/decimal"`)
		w.Line(`  "time"`)
		w.Line(")")
		for _, model := range version.ResolvedModels {
			w.Line("")
			if model.IsObject() {
				generateObjectModel(w, model)
			} else if model.IsOneOf() {
				generateOneOfModel(w, model)
			} else if model.IsEnum() {
				generateEnumModel(w, model)
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

func generateObjectModel(w *gen.Writer, model *spec.NamedModel) {
	w.Line("type %s struct {", model.Name.PascalCase())
	for _, field := range model.Object.Fields {
		typ := GoType(&field.Type.Definition)
		w.Line("  %s %s `json:\"%s\"`", field.Name.PascalCase(), typ, field.Name.Source)
	}
	w.Line("}")
}

func generateEnumModel(w *gen.Writer, model *spec.NamedModel) {
	w.Line("type %s %s", model.Name.PascalCase(), "string")
	w.Line("")
	w.Line("const (")
	var modelName string
	choiceValuesStringsParams := []string{}
	choiceValuesParams := []string{}
	for _, enumItem := range model.Enum.Items {
		modelName = model.Name.PascalCase()
		w.Line("  %s %s = \"%s\"", enumItem.Name.PascalCase(), modelName, enumItem.Value)
		choiceValuesStringsParams = append(choiceValuesStringsParams, fmt.Sprintf("string(%s)", enumItem.Name.PascalCase()))
		choiceValuesParams = append(choiceValuesParams, fmt.Sprintf("%s", enumItem.Name.PascalCase()))
	}
	w.Line(")")
	w.Line("")
	w.Line("var %sValuesStrings = []string{%s}", modelName, strings.Join(choiceValuesStringsParams, ", "))
	w.Line("var %sValues = []%s{%s}", modelName, modelName, strings.Join(choiceValuesParams, ", "))
}

func generateOneOfModel(w *gen.Writer, model *spec.NamedModel) {
	w.Line("type %s struct {", model.Name.PascalCase())
	for _, item := range model.OneOf.Items {
		typ := GoType(&item.Type.Definition)
		w.Line("  %s %s `json:\"%s\"`", item.Name.PascalCase(), typ, item.Name.Source)
	}
	w.Line("}")
}
