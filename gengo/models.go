package gengo

import (
	"fmt"
	"github.com/specgen-io/spec.v2"
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
			folder += "_" + version.Version.FlatCase()
		}
		files = append(files, gen.TextFile{Path: filepath.Join(generatePath, folder, "models.go"), Content: w.String()})
		files = append(files, *generateHelperFunctions(folder, filepath.Join(generatePath, folder, "helpers.go")))
	}
	return gen.WriteFiles(files, true)
}

func versionedPackage(version spec.Name, packageName string) string {
	if version.Source != "" {
		return fmt.Sprintf("%s_%s", packageName, version.FlatCase())
	}
	return packageName
}

func generateImport(w *gen.Writer, version *spec.Version, typ string, importStr string) *gen.Writer {
	for _, model := range version.ResolvedModels {
		if model.IsObject() {
			for _, field := range model.Object.Fields {
 				if checkType(&field.Type.Definition, typ) {
					w.Line(`import "%s"`, importStr)
					return w
				}
			}
		}
	}
	return nil
}

func checkType(fieldType *spec.TypeDef, typ string) bool {
	if fieldType.Plain != typ {
		switch fieldType.Node {
		case spec.PlainType:
			return false
		case spec.NullableType:
			return checkType(fieldType.Child, typ)
		case spec.ArrayType:
			return checkType(fieldType.Child, typ)
		case spec.MapType:
			return checkType(fieldType.Child, typ)
		default:
			panic(fmt.Sprintf("Unknown type: %v", typ))
		}
	}
	return true
}

func generateModels(w *gen.Writer, version *spec.Version, packageName string) {
	w.Line("package %s", versionedPackage(version.Version, packageName))
	w.EmptyLine()
	generateImport(w, version, spec.TypeDate, "cloud.google.com/go/civil")
	generateImport(w, version, spec.TypeJson, "encoding/json")
	generateImport(w, version, spec.TypeUuid, "github.com/google/uuid")
	generateImport(w, version, spec.TypeDecimal, "github.com/shopspring/decimal")
	for _, model := range version.ResolvedModels {
		if strings.Contains(w.String(), "import") {
			w.EmptyLine()
		}
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
	w.EmptyLine()
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
