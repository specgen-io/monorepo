package validation

import (
	"github.com/specgen-io/specgen/v2/gents/common"
	"github.com/specgen-io/specgen/v2/gents/modules"
	"github.com/specgen-io/specgen/v2/gents/writer"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func generateIoTsVersionModels(version *spec.Version, iotsModule modules.Module, module modules.Module) *sources.CodeFile {
	w := writer.NewTsWriter()
	w.Line("/* eslint-disable @typescript-eslint/camelcase */")
	w.Line("/* eslint-disable @typescript-eslint/no-magic-numbers */")
	w.Line(`import * as t from '%s'`, iotsModule.GetImport(module))
	for _, model := range version.ResolvedModels {
		w.EmptyLine()
		if model.IsObject() {
			generateIoTsObjectModel(w, model)
		} else if model.IsEnum() {
			generateIoTsEnumModel(w, model)
		} else if model.IsOneOf() {
			generateIoTsUnionModel(w, model)
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

func generateIoTsObjectModel(w *sources.Writer, model *spec.NamedModel) {
	hasRequiredFields, hasOptionalFields := kindOfFields(model)
	if hasRequiredFields && hasOptionalFields {
		w.Line("export const T%s = t.intersection([", model.Name.PascalCase())
		w.Line("  t.interface({")
		for _, field := range model.Object.Fields {
			if !field.Type.Definition.IsNullable() {
				w.Line("    %s: %s,", field.Name.Source, IoTsType(&field.Type.Definition))
			}
		}
		w.Line("  }),")
		w.Line("  t.partial({")
		for _, field := range model.Object.Fields {
			if field.Type.Definition.IsNullable() {
				w.Line("    %s: %s,", field.Name.Source, IoTsType(&field.Type.Definition))
			}
		}
		w.Line("  })")
		w.Line("])")
	} else {
		var ioTsType = "t.interface"
		if hasOptionalFields {
			ioTsType = "t.partial"
		}
		w.Line("export const T%s = %s({", model.Name.PascalCase(), ioTsType)
		for _, field := range model.Object.Fields {
			w.Line("  %s: %s,", field.Name.Source, IoTsType(&field.Type.Definition))
		}
		w.Line("})")
	}
	w.EmptyLine()
	w.Line("export type %s = t.TypeOf<typeof T%s>", model.Name.PascalCase(), model.Name.PascalCase())
}

func generateIoTsEnumModel(w *sources.Writer, model *spec.NamedModel) {
	w.Line("export enum %s {", model.Name.PascalCase())
	for _, item := range model.Enum.Items {
		w.Line(`  %s = "%s",`, item.Name.UpperCase(), item.Value)
	}
	w.Line("}")
	w.EmptyLine()
	w.Line("export const T%s = t.enum(%s)", model.Name.PascalCase(), model.Name.PascalCase())
}

func generateIoTsUnionModel(w *sources.Writer, model *spec.NamedModel) {
	if model.OneOf.Discriminator != nil {
		w.Line("export const T%s = t.union([", model.Name.PascalCase())
		for _, item := range model.OneOf.Items {
			w.Line("  t.intersection([t.type({%s: t.literal('%s')}), %s]),", common.TSIdentifier(*model.OneOf.Discriminator), item.Name.Source, SuperstructType(&item.Type.Definition))
		}
		w.Line("])")
	} else {
		w.Line("export const T%s = t.union([", model.Name.PascalCase())
		for _, item := range model.OneOf.Items {
			w.Line("  t.interface({%s: %s}),", item.Name.Source, IoTsType(&item.Type.Definition))
		}
		w.Line("])")
	}
	w.EmptyLine()
	w.Line("export type %s = t.TypeOf<typeof T%s>", model.Name.PascalCase(), model.Name.PascalCase())
}
