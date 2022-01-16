package validation

import (
	"github.com/specgen-io/specgen/v2/gents/common"
	"github.com/specgen-io/specgen/v2/gents/modules"
	"github.com/specgen-io/specgen/v2/gents/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func (v *ioTsValidation) GenerateVersionModels(version *spec.Version, codecModule modules.Module, module modules.Module) *sources.CodeFile {
	w := writer.NewTsWriter()
	w.Line("/* eslint-disable @typescript-eslint/camelcase */")
	w.Line("/* eslint-disable @typescript-eslint/no-magic-numbers */")
	w.Line(`import * as t from '%s'`, codecModule.GetImport(module))
	for _, model := range version.ResolvedModels {
		w.EmptyLine()
		if model.IsObject() {
			v.generateIoTsObjectModel(w, model)
		} else if model.IsEnum() {
			v.generateIoTsEnumModel(w, model)
		} else if model.IsOneOf() {
			v.generateIoTsUnionModel(w, model)
		}
	}
	return &sources.CodeFile{Path: module.GetPath(), Content: w.String()}
}

func kindOfFields(objectModel *spec.NamedModel) (bool, bool) {
	var hasRequiredFields = false
	var hasOptionalFields = false
	for _, field := range objectModel.Object.Fields {
		if !field.Type.Definition.IsNullable() {
			hasRequiredFields = true
		} else {
			hasOptionalFields = true
		}
	}
	return hasRequiredFields, hasOptionalFields
}

func (v *ioTsValidation) generateIoTsObjectModel(w *sources.Writer, model *spec.NamedModel) {
	hasRequiredFields, hasOptionalFields := kindOfFields(model)
	if hasRequiredFields && hasOptionalFields {
		w.Line("export const T%s = t.intersection([", model.Name.PascalCase())
		w.Line("  t.interface({")
		for _, field := range model.Object.Fields {
			if !field.Type.Definition.IsNullable() {
				w.Line("    %s: %s,", field.Name.Source, v.RuntimeType(&field.Type.Definition))
			}
		}
		w.Line("  }),")
		w.Line("  t.partial({")
		for _, field := range model.Object.Fields {
			if field.Type.Definition.IsNullable() {
				w.Line("    %s: %s,", field.Name.Source, v.RuntimeType(&field.Type.Definition))
			}
		}
		w.Line("  })")
		w.Line("])")
	} else {
		var modelTsType = "t.interface"
		if hasOptionalFields {
			modelTsType = "t.partial"
		}
		w.Line("export const T%s = %s({", model.Name.PascalCase(), modelTsType)
		for _, field := range model.Object.Fields {
			w.Line("  %s: %s,", field.Name.Source, v.RuntimeType(&field.Type.Definition))
		}
		w.Line("})")
	}
	w.EmptyLine()
	w.Line("export type %s = t.TypeOf<typeof T%s>", model.Name.PascalCase(), model.Name.PascalCase())
}

func (v *ioTsValidation) generateIoTsEnumModel(w *sources.Writer, model *spec.NamedModel) {
	w.Line("export enum %s {", model.Name.PascalCase())
	for _, item := range model.Enum.Items {
		w.Line(`  %s = "%s",`, item.Name.UpperCase(), item.Value)
	}
	w.Line("}")
	w.EmptyLine()
	w.Line("export const T%s = t.enum(%s)", model.Name.PascalCase(), model.Name.PascalCase())
}

func (v *ioTsValidation) generateIoTsUnionModel(w *sources.Writer, model *spec.NamedModel) {
	if model.OneOf.Discriminator != nil {
		w.Line("export const T%s = t.union([", model.Name.PascalCase())
		for _, item := range model.OneOf.Items {
			w.Line("  t.intersection([t.type({%s: t.literal('%s')}), %s]),", common.TSIdentifier(*model.OneOf.Discriminator), item.Name.Source, v.RuntimeType(&item.Type.Definition))
		}
		w.Line("])")
	} else {
		w.Line("export const T%s = t.union([", model.Name.PascalCase())
		for _, item := range model.OneOf.Items {
			w.Line("  t.interface({%s: %s}),", item.Name.Source, v.RuntimeType(&item.Type.Definition))
		}
		w.Line("])")
	}
	w.EmptyLine()
	w.Line("export type %s = t.TypeOf<typeof T%s>", model.Name.PascalCase(), model.Name.PascalCase())
}
