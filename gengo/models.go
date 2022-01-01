package gengo

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func GenerateModels(specification *spec.Spec, moduleName string, generatePath string) *sources.Sources {
	sources := sources.NewSources()

	rootModule := Module(moduleName, generatePath)

	for _, version := range specification.Versions {
		versionModule := rootModule.Submodule(version.Version.FlatCase())
		modelsModule := versionModule.Submodule(modelsPackage)
		sources.AddGeneratedAll(generateVersionModels(&version, modelsModule))
	}
	return sources
}

func generateVersionModels(version *spec.Version, module module) []sources.CodeFile {
	return []sources.CodeFile{
		*generateVersionModelsCode(version, module),
		*generateEnumsHelperFunctions(module),
	}
}

func generateVersionModelsCode(version *spec.Version, module module) *sources.CodeFile {
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
	return &sources.CodeFile{Path: module.GetPath("models.go"), Content: w.String()}
}

func requiredFieldsList(object *spec.Object) string {
	requiredFields := []string{}
	for _, field := range object.Fields {
		if !field.Type.Definition.IsNullable() {
			requiredFields = append(requiredFields, fmt.Sprintf(`"%s"`, field.Name.Source))
		}
	}
	return strings.Join(requiredFields, ", ")
}

func requiredFields(model *spec.NamedModel) string {
	return fmt.Sprintf(`%sRequiredFields`, model.Name.CamelCase())
}

func generateObjectModel(w *sources.Writer, model *spec.NamedModel) {
	w.Line("type %s struct {", model.Name.PascalCase())
	fields := [][]string{}
	for _, field := range model.Object.Fields {
		jsonAttributes := []string{field.Name.Source}
		if field.Type.Definition.IsNullable() {
			jsonAttributes = append(jsonAttributes, "omitempty")
		}
		fields = append(fields, []string{
			field.Name.PascalCase(),
			GoTypeSamePackage(&field.Type.Definition),
			fmt.Sprintf("`json:\"%s\"`", strings.Join(jsonAttributes, ",")),
		})
	}
	WriteAlignedLines(w.Indented(), fields)
	w.Line("}")
	w.EmptyLine()
	w.Line(`type %s %s`, model.Name.CamelCase(), model.Name.PascalCase())
	w.EmptyLine()
	w.Line(`var %s = []string{%s}`, requiredFields(model), requiredFieldsList(model.Object))
	w.EmptyLine()
	w.Line(`func (obj %s) MarshalJSON() ([]byte, error) {`, model.Name.PascalCase())
	w.Line(`	data, err := json.Marshal(%s(obj))`, model.Name.CamelCase())
	w.Line(`	if err != nil {`)
	w.Line(`		return nil, err`)
	w.Line(`	}`)
	w.Line(`	var rawMap map[string]json.RawMessage`)
	w.Line(`	err = json.Unmarshal(data, &rawMap)`)
	w.Line(`	for _, name := range %s {`, requiredFields(model))
	w.Line(`		value, found := rawMap[name]`)
	w.Line(`		if !found {`)
	w.Line(`			return nil, errors.New("required field missing: " + name)`)
	w.Line(`		}`)
	w.Line(`		if string(value) == "null" {`)
	w.Line(`			return nil, errors.New("required field doesn't have value: " + name)`)
	w.Line(`		}`)
	w.Line(`	}`)
	w.Line(`	return data, nil`)
	w.Line(`}`)
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
	w.Line(`	for _, name := range %s {`, requiredFields(model))
	w.Line(`		value, found := rawMap[name]`)
	w.Line(`		if !found {`)
	w.Line(`			return errors.New("required field missing: " + name)`)
	w.Line(`		}`)
	w.Line(`		if string(value) == "null" {`)
	w.Line(`			return errors.New("required field doesn't have value: " + name)`)
	w.Line(`		}`)
	w.Line(`	}`)
	w.Line(`	*obj = %s(jsonObj)`, model.Name.PascalCase())
	w.Line(`	return nil`)
	w.Line(`}`)
}

func generateEnumModel(w *sources.Writer, model *spec.NamedModel) {
	w.Line("type %s %s", model.Name.PascalCase(), "string")
	w.EmptyLine()
	w.Line("const (")
	modelName := model.Name.PascalCase()
	choiceValuesStringsParams := []string{}
	choiceValuesParams := []string{}
	items := [][]string{}
	for _, enumItem := range model.Enum.Items {
		enumConstName := modelName + enumItem.Name.PascalCase()
		items = append(items, []string{enumConstName, fmt.Sprintf(`%s = "%s"`, modelName, enumItem.Value)})
		choiceValuesStringsParams = append(choiceValuesStringsParams, fmt.Sprintf("string(%s)", enumConstName))
		choiceValuesParams = append(choiceValuesParams, fmt.Sprintf("%s", enumConstName))
	}
	WriteAlignedLines(w.Indented(), items)
	w.Line(")")
	w.EmptyLine()
	w.Line("var %s = []string{%s}", enumValuesStrings(model), JoinDelimParams(choiceValuesStringsParams))
	w.Line("var %s = []%s{%s}", enumValues(model), modelName, JoinDelimParams(choiceValuesParams))
	w.EmptyLine()
	w.Line("func (self *%s) UnmarshalJSON(b []byte) error {", modelName)
	w.Line("  str, err := readEnumStringValue(b, %sValuesStrings)", modelName)
	w.Line("  if err != nil {")
	w.Line("    return err")
	w.Line("  }")
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

func generateOneOfModel(w *sources.Writer, model *spec.NamedModel) {
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

func generateOneOfModelWrapper(w *sources.Writer, model *spec.NamedModel) {
	caseChecks := getCaseChecks(model.OneOf)
	items := [][]string{}
	w.Line("type %s struct {", model.Name.PascalCase())
	for _, item := range model.OneOf.Items {
		items = append(items, []string{
			item.Name.PascalCase(),
			GoTypeSamePackage(spec.Nullable(&item.Type.Definition)),
			fmt.Sprintf("`json:\"%s,omitempty\"`", item.Name.Source),
		})
	}
	WriteAlignedLines(w.Indented(), items)
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

func generateOneOfModelDiscriminator(w *sources.Writer, model *spec.NamedModel) {
	w.Line("type %s struct {", model.Name.PascalCase())
	items := [][]string{}
	for _, item := range model.OneOf.Items {
		items = append(items, []string{
			item.Name.PascalCase(),
			GoTypeSamePackage(spec.Nullable(&item.Type.Definition)),
		})
	}
	WriteAlignedLines(w.Indented(), items)
	w.Line("}")
	w.EmptyLine()
	w.Line(`func (u %s) MarshalJSON() ([]byte, error) {`, model.Name.PascalCase())
	for _, item := range model.OneOf.Items {
		w.Line(`  if u.%s != nil {`, item.Name.PascalCase())
		w.Line(`    data, err := json.Marshal(u.%s)`, item.Name.PascalCase())
		w.Line(`    if err != nil {`)
		w.Line(`      return nil, err`)
		w.Line(`    }`)
		w.Line(`    var rawMap map[string]json.RawMessage`)
		w.Line(`    json.Unmarshal(data, &rawMap)`)
		w.Line("    rawMap[\"%s\"] = []byte(`\"%s\"`)", *model.OneOf.Discriminator, item.Name.Source)
		w.Line(`    return json.Marshal(rawMap)`)
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
