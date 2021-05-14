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
	files := []gen.TextFile{}
	for _, version := range specification.Versions {
		w := NewGoWriter()
		generateModels(w, &version, "spec")
		folder := "spec"
		if version.Version.Source != "" {
			folder = folder + "_" + version.Version.FlatCase()
		}
		files = append(files, gen.TextFile{Path: filepath.Join(generatePath, folder, "models.go"), Content: w.String()})
		files = append(files, *generateHelperFunctions("spec", filepath.Join(generatePath, "spec", "helpers.go")))
	}
	return gen.WriteFiles(files, true)
}

func versionedPackage(version spec.Name, packageName string) string {
	if version.Source != "" {
		return fmt.Sprintf("%s_%s", packageName, version.FlatCase())
	}
	return packageName
}

func generateImport(version *spec.Version, typ string, importStr string) string {
	for _, model := range version.ResolvedModels {
		if model.IsObject() {
			for _, field := range model.Object.Fields {
				if field.Type.Definition.Plain == typ {
					return fmt.Sprintf(`import "%s"`, importStr)
				}
			}
		}
	}
	return ""
}

func generateModels(w *gen.Writer, version *spec.Version, packageName string) {
	w.Line("package %s", versionedPackage(version.Version, packageName))
	w.Line("")
	w.Line(generateImport(version, spec.TypeDecimal, "cloud.google.com/go/civil"))
	w.Line(generateImport(version, spec.TypeJson, "encoding/json"))
	w.Line(generateImport(version, spec.TypeUuid, "github.com/google/uuid"))
	w.Line(generateImport(version, spec.TypeDecimal, "github.com/shopspring/decimal"))
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
	modelName := model.Name.PascalCase()
	choiceValuesStringsParams := []string{}
	choiceValuesParams := []string{}
	for _, enumItem := range model.Enum.Items {
		enumItemName := enumItem.Name.PascalCase()
		w.Line("  %s %s = \"%s\"", enumItemName, modelName, enumItem.Value)
		choiceValuesStringsParams = append(choiceValuesStringsParams, fmt.Sprintf("string(%s)", enumItemName))
		choiceValuesParams = append(choiceValuesParams, fmt.Sprintf("%s", enumItemName))
	}
	w.Line(")")
	w.Line("")
	w.Line("var %sValuesStrings = []string{%s}", modelName, strings.Join(choiceValuesStringsParams, ", "))
	w.Line("var %sValues = []%s{%s}", modelName, modelName, strings.Join(choiceValuesParams, ", "))
	w.Line("")
	w.Line("func (self *%s) UnmarshalJSON(b []byte) error {", modelName)
	w.Line("  str, err := readEnumStringValue(b, %sValuesStrings)", modelName)
	w.Line("  if err != nil { return err }")
	w.Line("  *self = %s(str)", modelName)
	w.Line("  return nil")
	w.Line("}")
}

func generateOneOfModel(w *gen.Writer, model *spec.NamedModel) {
	w.Line("type %s struct {", model.Name.PascalCase())
	for _, item := range model.OneOf.Items {
		typ := GoType(&item.Type.Definition)
		w.Line("  %s *%s `json:\"%s,omitempty\"`", item.Name.PascalCase(), typ, item.Name.Source)
	}
	w.Line("}")
}
