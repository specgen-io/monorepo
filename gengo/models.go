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
		w.Line(`  "cloud.google.com/go/civil"`)
		w.Line(`  "encoding/json"`)
		w.Line(`  "errors"`)
		w.Line(`  "fmt"`)
		w.Line(`  "github.com/google/uuid"`)
		w.Line(`  "github.com/shopspring/decimal"`)
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
		generateHelperFunctions(w)
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

func generateHelperFunctions(w *gen.Writer) {
	w.Line("")
	w.Line("func contains(lookFor string, arr []string) bool {")
	w.Line("  for _, value := range arr{")
	w.Line("    if lookFor == value {")
	w.Line("      return true")
	w.Line("    }")
	w.Line("  }")
	w.Line("  return false")
	w.Line("}")
	w.Line("")
	w.Line("func readEnumStringValue(b []byte, values []string) (string, error) {")
	w.Line("  var str string")
	w.Line("  if err := json.Unmarshal(b, &str); err != nil {")
	w.Line("    return \"\", err")
	w.Line("  }")
	w.Line("  if !contains(str, values) {")
	w.Line("    return \"\", errors.New(fmt.Sprintf(\"Unknown enum value: %ss\", str))", "%")
	w.Line("  }")
	w.Line("  return str, nil")
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
