package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
	"strings"
)

func GenerateModels(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}
	files := generateModels(specification, generatePath)
	return gen.WriteFiles(files, true)
}

func generateModels(specification *spec.Spec, generatePath string) []gen.TextFile {
	files := []gen.TextFile{}
	for _, version := range specification.Versions {
		versionPath := filepath.Join(generatePath, versionedFolder(version.Version, "spec"))
		versionPackageName := versionedPackage(version.Version, "spec")
		files = append(files, generateVersionModels(&version, versionPackageName, versionPath)...)
	}
	return files
}

func generateVersionModels(version *spec.Version, packageName string, generatePath string) []gen.TextFile {
	return []gen.TextFile{
		*generateVersionModelsCode(version, packageName, generatePath),
		*generateHelperFunctions(packageName, filepath.Join(generatePath, "helpers.go")),
	}
}

func generateVersionModelsCode(version *spec.Version, packageName string, generatePath string) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", packageName)

	imports := []string{}
	if versionHasType(version, spec.TypeDate) {
		imports = append(imports, fmt.Sprintf(`import "cloud.google.com/go/civil"`))
	}
	if versionHasType(version, spec.TypeJson) {
		imports = append(imports, fmt.Sprintf(`import "encoding/json"`))
	}
	if versionHasType(version, spec.TypeUuid) {
		imports = append(imports, fmt.Sprintf(`import "github.com/google/uuid"`))
	}
	if versionHasType(version, spec.TypeDecimal) {
		imports = append(imports, fmt.Sprintf(`import "github.com/shopspring/decimal"`))
	}

	if len(imports) > 0 {
		w.EmptyLine()
		w.Line(`%s`, strings.Join(imports, "\n"))
	}

	for _, model := range version.ResolvedModels {
		w.EmptyLine()
		if model.IsObject() {
			generateObjectModel(w, model)
		} else if model.IsOneOf() {
			generateOneOfModel(w, model)
		} else if model.IsEnum() {
			generateEnumModel(w, model)
		}
	}
	return &gen.TextFile{Path: filepath.Join(generatePath, "models.go"), Content: w.String()}
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
	w.EmptyLine()
	w.Line("const (")
	modelName := model.Name.PascalCase()
	choiceValuesStringsParams := []string{}
	choiceValuesParams := []string{}
	for _, enumItem := range model.Enum.Items {
		enumConstName := modelName + enumItem.Name.PascalCase()
		w.Line("  %s %s = \"%s\"", enumConstName, modelName, enumItem.Value)
		choiceValuesStringsParams = append(choiceValuesStringsParams, fmt.Sprintf("string(%s)", enumConstName))
		choiceValuesParams = append(choiceValuesParams, fmt.Sprintf("%s", enumConstName))
	}
	w.Line(")")
	w.EmptyLine()
	w.Line("var %sValuesStrings = []string{%s}", modelName, strings.Join(choiceValuesStringsParams, ", "))
	w.Line("var %sValues = []%s{%s}", modelName, modelName, strings.Join(choiceValuesParams, ", "))
	w.EmptyLine()
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
