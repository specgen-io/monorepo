package models

import (
	"fmt"
	"strings"

	"generator"
	"golang/module"
	"golang/types"
	"golang/writer"
	"spec"
)

func NewEncodingJsonGenerator(types *types.Types, modules *Modules) *EncodingJsonGenerator {
	return &EncodingJsonGenerator{types, modules}
}

type EncodingJsonGenerator struct {
	Types   *types.Types
	Modules *Modules
}

func (g *EncodingJsonGenerator) Models(version *spec.Version) *generator.CodeFile {
	return g.models(version.ResolvedModels, g.Modules.Models(version))
}

func (g *EncodingJsonGenerator) ErrorModels(httperrors *spec.HttpErrors) *generator.CodeFile {
	return g.models(httperrors.ResolvedModels, g.Modules.HttpErrorsModels)
}

func (g *EncodingJsonGenerator) models(models []*spec.NamedModel, modelsModule module.Module) *generator.CodeFile {
	w := writer.New(modelsModule, "models.go")

	w.Imports.AddModelsTypes(models)
	if types.ModelsHasEnum(models) {
		w.Imports.Module(g.Modules.Enums)
	}

	for _, model := range models {
		w.EmptyLine()
		if model.IsObject() {
			g.objectModel(w, model)
		} else if model.IsOneOf() {
			g.oneOfModel(w, model)
		} else if model.IsEnum() {
			g.enumModel(w, model)
		}
	}

	return w.ToCodeFile()
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

func (g *EncodingJsonGenerator) objectModel(w *writer.Writer, model *spec.NamedModel) {
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
	w.Line(`			return nil, errors.NewImports("required field missing: " + name)`)
	w.Line(`		}`)
	w.Line(`		if string(value) == "null" {`)
	w.Line(`			return nil, errors.NewImports("required field doesn't have value: " + name)`)
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
	w.Line(`		return errors.NewImports("failed to check fields in json: " + err.Error())`)
	w.Line(`	}`)
	w.Line(`	for _, name := range %s {`, g.requiredFields(model))
	w.Line(`		value, found := rawMap[name]`)
	w.Line(`		if !found {`)
	w.Line(`			return errors.NewImports("required field missing: " + name)`)
	w.Line(`		}`)
	w.Line(`		if string(value) == "null" {`)
	w.Line(`			return errors.NewImports("required field doesn't have value: " + name)`)
	w.Line(`		}`)
	w.Line(`	}`)
	w.Line(`	*obj = %s(jsonObj)`, model.Name.PascalCase())
	w.Line(`	return nil`)
	w.Line(`}`)
}

func (g *EncodingJsonGenerator) enumModel(w *writer.Writer, model *spec.NamedModel) {
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
		choiceValuesStringsParams = append(choiceValuesStringsParams, fmt.Sprintf("string(%s)", enumConstName))
		choiceValuesParams = append(choiceValuesParams, fmt.Sprintf("%s", enumConstName))
	}
	w.Line("var %s = []string{%s}", g.EnumValuesStrings(model), strings.Join(choiceValuesStringsParams, ", "))
	w.Line("var %s = []%s{%s}", g.enumValues(model), model.Name.PascalCase(), strings.Join(choiceValuesParams, ", "))
	w.EmptyLine()
	w.Line("func (self *%s) UnmarshalJSON(b []byte) error {", model.Name.PascalCase())
	w.Line("  str, err := enums.ReadStringValue(b, %sValuesStrings)", model.Name.PascalCase())
	w.Line("  if err != nil {")
	w.Line("    return err")
	w.Line("  }")
	w.Line("  *self = %s(str)", model.Name.PascalCase())
	w.Line("  return nil")
	w.Line("}")
}

func (g *EncodingJsonGenerator) EnumValuesStrings(model *spec.NamedModel) string {
	return fmt.Sprintf("%sValuesStrings", model.Name.PascalCase())
}

func (g *EncodingJsonGenerator) enumValues(model *spec.NamedModel) string {
	return fmt.Sprintf("%sValues", model.Name.PascalCase())
}

func (g *EncodingJsonGenerator) oneOfModel(w *writer.Writer, model *spec.NamedModel) {
	if model.OneOf.Discriminator != nil {
		g.oneOfModelDiscriminator(w, model)
	} else {
		g.oneOfModelWrapper(w, model)
	}
}

func (g *EncodingJsonGenerator) getCaseChecks(oneof *spec.OneOf) string {
	caseChecks := []string{}
	for _, item := range oneof.Items {
		caseChecks = append(caseChecks, fmt.Sprintf(`u.%s == nil`, item.Name.PascalCase()))
	}
	return strings.Join(caseChecks, " && ")
}

func (g *EncodingJsonGenerator) oneOfModelWrapper(w *writer.Writer, model *spec.NamedModel) {
	caseChecks := g.getCaseChecks(model.OneOf)
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
	w.EmptyLine()
	w.Line(`type %s %s`, model.Name.CamelCase(), model.Name.PascalCase())
	w.EmptyLine()
	w.Line(`func (u %s) MarshalJSON() ([]byte, error) {`, model.Name.PascalCase())
	w.Line(`	if %s {`, caseChecks)
	w.Line(`		return nil, errors.NewImports("union case is not set")`)
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
	w.Line(`		return errors.NewImports("union case is not set")`)
	w.Line(`	}`)
	w.Line(`	return nil`)
	w.Line(`}`)
}

func (g *EncodingJsonGenerator) oneOfModelDiscriminator(w *writer.Writer, model *spec.NamedModel) {
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
	w.Line(`  return nil, errors.NewImports("union case is not set")`)
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`func (u *%s) UnmarshalJSON(data []byte) error {`, model.Name.PascalCase())
	w.Line(`  var discriminator struct {`)
	w.Line("    Value *string `json:\"%s\"`", *model.OneOf.Discriminator)
	w.Line(`  }`)
	w.Line(`  err := json.Unmarshal(data, &discriminator)`)
	w.Line(`  if err != nil {`)
	w.Line(`    return errors.NewImports("failed to parse discriminator field %s: " + err.Error())`, *model.OneOf.Discriminator)
	w.Line(`  }`)
	w.Line(`  if discriminator.Value == nil {`)
	w.Line(`    return errors.NewImports("discriminator field %s not found")`, *model.OneOf.Discriminator)
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
	w.Line(`    return errors.NewImports("unexpected union discriminator field %s value: " + *discriminator.Value)`, *model.OneOf.Discriminator)
	w.Line(`  }`)
	w.Line(`  return nil`)
	w.Line(`}`)
}
