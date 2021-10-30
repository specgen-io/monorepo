package gengo

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/spec"
	"github.com/specgen-io/specgen/v2/gen"
)

func GenerateModels(specification *spec.Spec, moduleName string, generatePath string) error {
	files := []gen.TextFile{}

	for _, version := range specification.Versions {
		modelsModule := Module(moduleName, createPath(generatePath, version.Version.FlatCase(), modelsPackage))
		files = append(files, generateVersionModels(&version, modelsModule)...)
	}
	return gen.WriteFiles(files, true)
}

func generateVersionModels(version *spec.Version, module module) []gen.TextFile {
	return []gen.TextFile{
		*generateVersionModelsCode(version, module),
		*generateEnumsHelperFunctions(module),
	}
}

func generateVersionModelsCode(version *spec.Version, module module) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", module.Name)

	imports := Imports()
	imports.AddModelsTypes(version)
	imports.Write(w)

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
	return &gen.TextFile{Path: module.GetPath("models.go"), Content: w.String()}
}

func generateObjectModel(w *gen.Writer, model *spec.NamedModel) {
	w.Line("type %s struct {", model.Name.PascalCase())
	for _, field := range model.Object.Fields {
		w.Line("  %s %s `json:\"%s\"`", field.Name.PascalCase(), GoTypeSamePackage(&field.Type.Definition), field.Name.Source)
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
	w.Line("var %s = []string{%s}", enumValuesStrings(model), JoinDelimParams(choiceValuesStringsParams))
	w.Line("var %s = []%s{%s}", enumValues(model), modelName, JoinDelimParams(choiceValuesParams))
	w.EmptyLine()
	w.Line("func (self *%s) UnmarshalJSON(b []byte) error {", modelName)
	w.Line("  str, err := readEnumStringValue(b, %sValuesStrings)", modelName)
	w.Line("  if err != nil { return err }")
	w.Line("  *self = %s(str)", modelName)
	w.Line("  return nil")
	w.Line("}")
}

func enumValuesStrings(model *spec.NamedModel) string {
	return fmt.Sprintf("%sValuesStrings", model.Name.PascalCase())
}

func enumValues(model *spec.NamedModel) string {
	return fmt.Sprintf("%sValues", model.Name.PascalCase())
}

func generateOneOfModel(w *gen.Writer, model *spec.NamedModel) {
	w.Line("type %s struct {", model.Name.PascalCase())
	for _, item := range model.OneOf.Items {
		w.Line("  %s *%s `json:\"%s,omitempty\"`", item.Name.PascalCase(), GoTypeSamePackage(&item.Type.Definition), item.Name.Source)
	}
	w.Line("}")
	if model.OneOf.Discriminator != nil {
		w.EmptyLine()
		w.Line(`func (u %s) MarshalJSON() ([]byte, error) {`, model.Name.PascalCase())
		for _, item := range model.OneOf.Items {
			w.Line(`  if u.%s != nil {`, item.Name.PascalCase())
			w.Line(`    return json.Marshal(&struct {`)
			w.Line("      Discriminator string `json:\"%s\"`", *model.OneOf.Discriminator)
			w.Line(`      *%s`, GoTypeSamePackage(&item.Type.Definition))
			w.Line(`    }{`)
			w.Line(`      Discriminator: "%s",`, item.Name.Source)
			w.Line(`      %s: u.%s,`, GoTypeSamePackage(&item.Type.Definition), item.Name.PascalCase())
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
			w.Line(`      unionCase := %s{}`, GoTypeSamePackage(&item.Type.Definition))
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
