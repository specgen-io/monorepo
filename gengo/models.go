package gengo

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func GenerateModels(specification *spec.Spec, moduleName string, generatePath string) *gen.Sources {
	sources := gen.NewSources()

	rootModule := Module(moduleName, generatePath)

	for _, version := range specification.Versions {
		versionModule := rootModule.Submodule(version.Version.FlatCase())
		modelsModule := versionModule.Submodule(modelsPackage)
		sources.AddGeneratedAll(generateVersionModels(&version, modelsModule))
	}
	return sources
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

func getRequiredFields(object *spec.Object) string {
	requiredFields := []string{}
	for _, field := range object.Fields {
		if !field.Type.Definition.IsNullable() {
			requiredFields = append(requiredFields, fmt.Sprintf(`"%s"`, field.Name.Source))
		}
	}
	return strings.Join(requiredFields, ", ")
}

func generateObjectModel(w *gen.Writer, model *spec.NamedModel) {
	w.Line("type %s struct {", model.Name.PascalCase())
	for _, field := range model.Object.Fields {
		w.Line("  %s %s `json:\"%s\"`", field.Name.PascalCase(), GoTypeSamePackage(&field.Type.Definition), field.Name.Source)
	}
	w.Line("}")
	w.EmptyLine()
	w.Line(`type %s %s`, model.Name.CamelCase(), model.Name.PascalCase())
	w.EmptyLine()
	w.Line(`var %sRequiredFields = []string{%s}`, model.Name.CamelCase(), getRequiredFields(model.Object))
	w.EmptyLine()
	w.Line(`func (obj *%s) UnmarshalJSON(data []byte) error {`, model.Name.PascalCase())
	w.Line(`	jsonObj := %s(*obj)`, model.Name.CamelCase())
	w.Line(`	err := json.Unmarshal(data, &jsonObj)`)
	w.Line(`	if err != nil {`)
	w.Line(`		return err`)
	w.Line(`	}`)
	w.Line(`	var rawMap map[string]json.RawMessage`)
	w.Line(`	err = json.Unmarshal(data, &rawMap)`)
	w.Line(`	if err != nil {`)
	w.Line(`		return errors.New("failed to check fields in json: " + err.Error())`)
	w.Line(`	}`)
	w.Line(`	for _, name := range %sRequiredFields {`, model.Name.CamelCase())
	w.Line(`		if _, found := rawMap[name]; !found {`)
	w.Line(`			return errors.New("required field missing: " + name)`)
	w.Line(`		}`)
	w.Line(`	}`)
	w.Line(`	*obj = %s(jsonObj)`, model.Name.PascalCase())
	w.Line(`	return nil`)
	w.Line(`}`)
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
	if model.OneOf.Discriminator != nil {
		generateOneOfModelDiscriminator(w, model)
	} else {
		generateOneOfModelWrapper(w, model)
	}
}

func getCaseChecks(oneof *spec.OneOf) string {
	caseChecks := []string{}
	for _, item := range oneof.Items {
		caseChecks = append(caseChecks, fmt.Sprintf(`u.%s == nil`, item.Name.PascalCase()))
	}
	return strings.Join(caseChecks, " && ")
}

func generateOneOfModelWrapper(w *gen.Writer, model *spec.NamedModel) {
	caseChecks := getCaseChecks(model.OneOf)
	w.Line("type %s struct {", model.Name.PascalCase())
	for _, item := range model.OneOf.Items {
		w.Line("  %s *%s `json:\"%s,omitempty\"`", item.Name.PascalCase(), GoTypeSamePackage(&item.Type.Definition), item.Name.Source)
	}
	w.Line("}")
	w.EmptyLine()
	w.Line(`type %s %s`, model.Name.CamelCase(), model.Name.PascalCase())
	w.EmptyLine()
	w.Line(`func (u %s) MarshalJSON() ([]byte, error) {`, model.Name.PascalCase())
	w.Line(`	if %s {`, caseChecks)
	w.Line(`		return nil, errors.New("union case is not set")`)
	w.Line(`	}`)
	w.Line(`	return json.Marshal(%s(u))`, model.Name.CamelCase())
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`func (u *%s) UnmarshalJSON(data []byte) error {`, model.Name.PascalCase())
	w.Line(`	jsonObj := %s(*u)`, model.Name.CamelCase())
	w.Line(`	err := json.Unmarshal(data, &jsonObj)`)
	w.Line(`	if err != nil {`)
	w.Line(`		return err`)
	w.Line(`	}`)
	w.Line(`	*u = %s(jsonObj)`, model.Name.PascalCase())
	w.Line(`	if %s {`, caseChecks)
	w.Line(`		return errors.New("union case is not set")`)
	w.Line(`	}`)
	w.Line(`	return nil`)
	w.Line(`}`)
}

func generateOneOfModelDiscriminator(w *gen.Writer, model *spec.NamedModel) {
	w.Line("type %s struct {", model.Name.PascalCase())
	for _, item := range model.OneOf.Items {
		w.Line("  %s *%s `json:\"%s,omitempty\"`", item.Name.PascalCase(), GoTypeSamePackage(&item.Type.Definition), item.Name.Source)
	}
	w.Line("}")
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
	w.Line(`  return nil, errors.New("union case is not set")`)
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`func (u *%s) UnmarshalJSON(data []byte) error {`, model.Name.PascalCase())
	w.Line(`  var discriminator struct {`)
	w.Line("    Value *string `json:\"%s\"`", *model.OneOf.Discriminator)
	w.Line(`  }`)
	w.Line(`  err := json.Unmarshal(data, &discriminator)`)
	w.Line(`  if err != nil {`)
	w.Line(`    return errors.New("failed to parse discriminator field %s: " + err.Error())`, *model.OneOf.Discriminator)
	w.Line(`  }`)
	w.Line(`  if discriminator.Value == nil {`)
	w.Line(`    return errors.New("discriminator field %s not found")`, *model.OneOf.Discriminator)
	w.Line(`  }`)
	w.EmptyLine()
	w.Line(`  switch *discriminator.Value {`)
	for _, item := range model.OneOf.Items {
		w.Line(`  case "%s":`, item.Name.Source)
		w.Line(`    unionCase := %s{}`, GoTypeSamePackage(&item.Type.Definition))
		w.Line(`    err := json.Unmarshal(data, &unionCase)`)
		w.Line(`    if err != nil {`)
		w.Line(`      return err`)
		w.Line(`    }`)
		w.Line(`    u.%s = &unionCase`, item.Name.PascalCase())
	}
	w.Line(`  default:`)
	w.Line(`    return errors.New("unexpected union discriminator field %s value: " + *discriminator.Value)`, *model.OneOf.Discriminator)
	w.Line(`  }`)
	w.Line(`  return nil`)
	w.Line(`}`)
}
