package models

import (
	"fmt"
	"generator"
	"golang/module"
	"golang/types"
	"golang/walkers"
	"golang/writer"
	"spec"
	"strings"
)

func NewEncodingJsonGenerator(types *types.Types, modules *Modules, mode bool) *EncodingJsonGenerator {
	return &EncodingJsonGenerator{types, modules, mode}
}

type EncodingJsonGenerator struct {
	Types      *types.Types
	Modules    *Modules
	strictMode bool
}

func (g *EncodingJsonGenerator) Models(version *spec.Version) []generator.CodeFile {
	return g.models(version.ResolvedModels, g.Modules.Models(version))
}

func (g *EncodingJsonGenerator) ErrorModels(httperrors *spec.HttpErrors) []generator.CodeFile {
	return g.models(httperrors.ResolvedModels, g.Modules.HttpErrorsModels)
}

func (g *EncodingJsonGenerator) models(models []*spec.NamedModel, modelsModule module.Module) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, model := range models {
		if model.IsObject() {
			files = append(files, *g.objectModel(modelsModule, model))
		} else if model.IsOneOf() {
			files = append(files, *g.oneOfModel(modelsModule, model))
		} else if model.IsEnum() {
			files = append(files, *g.enumModel(modelsModule, model))
		}
	}
	return files
}

func (g *EncodingJsonGenerator) requiredFieldsList(object *spec.Object) string {
	requiredFields := []string{}
	for _, field := range object.Fields {
		if !field.Type.Definition.IsNullable() {
			requiredFields = append(requiredFields, fmt.Sprintf(`"%s"`, field.Name.Source))
		}
	}
	return strings.Join(requiredFields, ", ")
}

func (g *EncodingJsonGenerator) requiredFields(model *spec.NamedModel) string {
	return fmt.Sprintf(`%sRequiredFields`, model.Name.CamelCase())
}

func (g *EncodingJsonGenerator) objectModel(modelsModule module.Module, model *spec.NamedModel) *generator.CodeFile {
	w := writer.New(modelsModule, model.Name.SnakeCase()+`.go`)
	if g.strictMode {
		w.Imports.Add("errors")
		if !walkers.ModelHasType(model, spec.TypeJson) {
			w.Imports.Add("encoding/json")
		}
	}
	if walkers.ModelHasType(model, spec.TypeJson) {
		w.Imports.Add("encoding/json")
	}
	if walkers.ModelHasType(model, spec.TypeDate) {
		w.Imports.Add("cloud.google.com/go/civil")
	}
	if walkers.ModelHasType(model, spec.TypeUuid) {
		w.Imports.Add("github.com/google/uuid")
	}
	if walkers.ModelHasType(model, spec.TypeDecimal) {
		w.Imports.Add("github.com/shopspring/decimal")
	}
	w.Line("type %s struct {", model.Name.PascalCase())
	w.Indent()
	for _, field := range model.Object.Fields {
		jsonAttributes := []string{field.Name.Source}
		if field.Type.Definition.IsNullable() {
			jsonAttributes = append(jsonAttributes, "omitempty")
		}
		w.LineAligned("%s %s `json:\"%s\"`",
			field.Name.PascalCase(),
			g.Types.GoTypeSamePackage(&field.Type.Definition),
			strings.Join(jsonAttributes, ","))
	}
	w.Unindent()
	w.Line("}")
	if g.strictMode {
		w.EmptyLine()
		w.Line(`type %s %s`, model.Name.CamelCase(), model.Name.PascalCase())
		w.EmptyLine()
		w.Line(`var %s = []string{%s}`, g.requiredFields(model), g.requiredFieldsList(model.Object))
		w.EmptyLine()
		w.Line(`func (obj %s) MarshalJSON() ([]byte, error) {`, model.Name.PascalCase())
		w.Line(`	data, err := json.Marshal(%s(obj))`, model.Name.CamelCase())
		w.Line(`	if err != nil {`)
		w.Line(`		return nil, err`)
		w.Line(`	}`)
		w.Line(`	var rawMap map[string]json.RawMessage`)
		w.Line(`	err = json.Unmarshal(data, &rawMap)`)
		w.Line(`	for _, name := range %s {`, g.requiredFields(model))
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
		w.Line(`	for _, name := range %s {`, g.requiredFields(model))
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
	return w.ToCodeFile()
}

func (g *EncodingJsonGenerator) enumModel(modelsModule module.Module, model *spec.NamedModel) *generator.CodeFile {
	w := writer.New(modelsModule, model.Name.SnakeCase()+`.go`)
	if g.strictMode {
		w.Imports.Module(g.Modules.Enums)
	}
	w.Line("type %s %s", model.Name.PascalCase(), "string")
	w.EmptyLine()
	w.Line("const (")
	w.Indent()
	for _, enumItem := range model.Enum.Items {
		w.LineAligned(`%s%s = "%s"`, model.Name.PascalCase(), enumItem.Name.PascalCase(), enumItem.Value)
	}
	w.Unindent()
	w.Line(")")
	w.EmptyLine()
	choiceValuesStringsParams := []string{}
	choiceValuesParams := []string{}
	for _, enumItem := range model.Enum.Items {
		enumConstName := model.Name.PascalCase() + enumItem.Name.PascalCase()
		choiceValuesStringsParams = append(choiceValuesStringsParams, enumConstName)
		choiceValuesParams = append(choiceValuesParams, fmt.Sprintf("%s", enumConstName))
	}
	w.Line("var %s = []string{%s}", g.EnumValuesStrings(model), strings.Join(choiceValuesStringsParams, ", "))
	w.Line("var %s = []%s{%s}", g.enumValues(model), model.Name.PascalCase(), strings.Join(choiceValuesParams, ", "))
	if g.strictMode {
		w.EmptyLine()
		w.Line("func (self *%s) UnmarshalJSON(b []byte) error {", model.Name.PascalCase())
		w.Line("  str, err := %s.ReadStringValue(b, %sValuesStrings)", g.Modules.Enums.Name, model.Name.PascalCase())
		w.Line("  if err != nil {")
		w.Line("    return err")
		w.Line("  }")
		w.Line("  *self = %s(str)", model.Name.PascalCase())
		w.Line("  return nil")
		w.Line("}")
	}
	return w.ToCodeFile()
}

func (g *EncodingJsonGenerator) EnumValuesStrings(model *spec.NamedModel) string {
	return fmt.Sprintf("%sValuesStrings", model.Name.PascalCase())
}

func (g *EncodingJsonGenerator) enumValues(model *spec.NamedModel) string {
	return fmt.Sprintf("%sValues", model.Name.PascalCase())
}

func (g *EncodingJsonGenerator) oneOfModel(modelsModule module.Module, model *spec.NamedModel) *generator.CodeFile {
	if model.OneOf.Discriminator != nil {
		return g.oneOfModelDiscriminator(modelsModule, model)
	} else {
		return g.oneOfModelWrapper(modelsModule, model)
	}
}

func (g *EncodingJsonGenerator) getCaseChecks(oneof *spec.OneOf) string {
	caseChecks := []string{}
	for _, item := range oneof.Items {
		caseChecks = append(caseChecks, fmt.Sprintf(`u.%s == nil`, item.Name.PascalCase()))
	}
	return strings.Join(caseChecks, " && ")
}

func (g *EncodingJsonGenerator) oneOfModelWrapper(modelsModule module.Module, model *spec.NamedModel) *generator.CodeFile {
	w := writer.New(modelsModule, model.Name.SnakeCase()+`.go`)
	if g.strictMode {
		w.Imports.Add("errors")
		if !walkers.ModelHasType(model, spec.TypeJson) {
			w.Imports.Add("encoding/json")
		}
	}
	if walkers.ModelHasType(model, spec.TypeJson) {
		w.Imports.Add("encoding/json")
	}
	if walkers.ModelHasType(model, spec.TypeDate) {
		w.Imports.Add("cloud.google.com/go/civil")
	}
	if walkers.ModelHasType(model, spec.TypeUuid) {
		w.Imports.Add("github.com/google/uuid")
	}
	if walkers.ModelHasType(model, spec.TypeDecimal) {
		w.Imports.Add("github.com/shopspring/decimal")
	}
	w.Line("type %s struct {", model.Name.PascalCase())
	w.Indent()
	for _, item := range model.OneOf.Items {
		w.LineAligned("%s %s `json:\"%s,omitempty\"`",
			item.Name.PascalCase(),
			g.Types.GoTypeSamePackage(spec.Nullable(&item.Type.Definition)),
			item.Name.Source)
	}
	w.Unindent()
	w.Line("}")
	if g.strictMode {
		w.EmptyLine()
		w.Line(`type %s %s`, model.Name.CamelCase(), model.Name.PascalCase())
		w.EmptyLine()
		w.Line(`func (u %s) MarshalJSON() ([]byte, error) {`, model.Name.PascalCase())
		w.Line(`	if %s {`, g.getCaseChecks(model.OneOf))
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
		w.Line(`	if %s {`, g.getCaseChecks(model.OneOf))
		w.Line(`		return errors.New("union case is not set")`)
		w.Line(`	}`)
		w.Line(`	return nil`)
		w.Line(`}`)
	}
	return w.ToCodeFile()
}

func (g *EncodingJsonGenerator) oneOfModelDiscriminator(modelsModule module.Module, model *spec.NamedModel) *generator.CodeFile {
	w := writer.New(modelsModule, model.Name.SnakeCase()+`.go`)
	w.Imports.Add("errors")
	w.Imports.Add("encoding/json")
	if walkers.ModelHasType(model, spec.TypeDate) {
		w.Imports.Add("cloud.google.com/go/civil")
	}
	if walkers.ModelHasType(model, spec.TypeUuid) {
		w.Imports.Add("github.com/google/uuid")
	}
	if walkers.ModelHasType(model, spec.TypeDecimal) {
		w.Imports.Add("github.com/shopspring/decimal")
	}
	w.Line("type %s struct {", model.Name.PascalCase())
	w.Indent()
	for _, item := range model.OneOf.Items {
		w.LineAligned(`%s %s`, item.Name.PascalCase(), g.Types.GoTypeSamePackage(spec.Nullable(&item.Type.Definition)))
	}
	w.Unindent()
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
		w.Line(`    unionCase := %s{}`, g.Types.GoTypeSamePackage(&item.Type.Definition))
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
	return w.ToCodeFile()
}
