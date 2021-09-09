package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"specgen/gen"
	"path/filepath"
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

	packageName := fmt.Sprintf("%s_models", specification.Name.SnakeCase())
	generatePath = filepath.Join(generatePath, packageName)
	for _, version := range specification.Versions {
		versionPackageName := versionedPackage(version.Version, packageName)
		versionPath := versionedFolder(version.Version, generatePath)
		files = append(files, generateVersionModels(&version, versionPackageName, versionPath)...)
	}
	return files
}

func generateVersionModels(version *spec.Version, packageName string, generatePath string) []gen.TextFile {
	return []gen.TextFile{
		*generateVersionModelsCode(version, packageName, generatePath),
		*generateEnumsHelperFunctions(packageName, filepath.Join(generatePath, "enums_helpers.go")),
	}
}

func generateVersionModelsCode(version *spec.Version, packageName string, generatePath string) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", packageName)

	imports := []string{}
	if versionModelsHasType(version, spec.TypeDate) {
		imports = append(imports, `"cloud.google.com/go/civil"`)
	}
	if versionModelsHasType(version, spec.TypeJson) {
		imports = append(imports, `"encoding/json"`)
	}
	if versionModelsHasType(version, spec.TypeUuid) {
		imports = append(imports, `"github.com/google/uuid"`)
	}
	if versionModelsHasType(version, spec.TypeDecimal) {
		imports = append(imports, `"github.com/shopspring/decimal"`)
	}
	if isOneOfModel(version) {
		imports = append(imports, `"errors"`, `"fmt"`)
	}

	if len(imports) > 0 {
		w.EmptyLine()
		for _, imp := range imports {
			w.Line(`import %s`, imp)
		}
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
		w.Line("  %s %s `json:\"%s\"`", field.Name.PascalCase(), GoType(&field.Type.Definition), field.Name.Source)
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
	w.Line("var %sValuesStrings = []string{%s}", modelName, JoinDelimParams(choiceValuesStringsParams))
	w.Line("var %sValues = []%s{%s}", modelName, modelName, JoinDelimParams(choiceValuesParams))
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
		w.Line("  %s *%s `json:\"%s,omitempty\"`", item.Name.PascalCase(), GoType(&item.Type.Definition), item.Name.Source)
	}
	w.Line("}")
	if model.OneOf.Discriminator != nil {
		w.EmptyLine()
		w.Line(`func (u %s) MarshalJSON() ([]byte, error) {`, model.Name.PascalCase())
		for _, item := range model.OneOf.Items {
			w.Line(`  if u.%s != nil {`, item.Name.PascalCase())
			w.Line(`    return json.Marshal(&struct {`)
			w.Line("      Discriminator string `json:\"%s\"`", *model.OneOf.Discriminator)
			w.Line(`      *%s`, GoType(&item.Type.Definition))
			w.Line(`    }{`)
			w.Line(`      Discriminator: "%s",`, item.Name.Source)
			w.Line(`      %s: u.%s,`, GoType(&item.Type.Definition), item.Name.PascalCase())
			w.Line(`    })`)
			w.Line(`  }`)
		}
		w.Line(`  return nil, errors.New("union case not set")`)
		w.Line(`}`)
		w.EmptyLine()
		w.Line(`func (u *%s) UnmarshalJSON(data []byte) error {`, model.Name.PascalCase())
		w.Line(`  var discriminator struct {`)
		w.Line("    Value string `json:\"%s\"`", *model.OneOf.Discriminator)
		w.Line(`  }`)
		w.Line(`  err := json.Unmarshal(data, &discriminator)`)
		w.Line(`  if err != nil { return err }`)
		w.EmptyLine()
		w.Line(`  switch discriminator.Value {`)
		for _, item := range model.OneOf.Items {
			w.Line(`    case "%s":`, item.Name.Source)
			w.Line(`      unionCase := %s{}`, GoType(&item.Type.Definition))
			w.Line(`      err := json.Unmarshal(data, &unionCase)`)
			w.Line(`      if err != nil { return err }`)
			w.Line(`      u.%s = &unionCase`, item.Name.PascalCase())
		}
		w.Line(`    default:`)
		w.Line(`      return errors.New(fmt.Sprintf("unexpected union discriminator field %s value: %s", discriminator.Value))`, *model.OneOf.Discriminator, "%s")
		w.Line(`  }`)
		w.Line(`  return nil`)
		w.Line(`}`)
	}
}
