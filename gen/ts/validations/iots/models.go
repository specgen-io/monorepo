package iots

import (
	"github.com/specgen-io/specgen/v2/gen/ts/common"
	"github.com/specgen-io/specgen/v2/gen/ts/modules"
	"github.com/specgen-io/specgen/v2/gen/ts/writer"
	"github.com/specgen-io/specgen/v2/generator"
	"github.com/specgen-io/specgen/v2/spec"
)

func (g *Generator) VersionModels(version *spec.Version, codecModule modules.Module, module modules.Module) *generator.CodeFile {
	w := writer.NewTsWriter()
	w.Line("/* eslint-disable @typescript-eslint/camelcase */")
	w.Line("/* eslint-disable @typescript-eslint/no-magic-numbers */")
	w.Line(`import * as t from '%s'`, codecModule.GetImport(module))
	for _, model := range version.ResolvedModels {
		w.EmptyLine()
		if model.IsObject() {
			g.objectModel(w, model)
		} else if model.IsEnum() {
			g.enumModel(w, model)
		} else if model.IsOneOf() {
			g.unionModel(w, model)
		}
	}
	return &generator.CodeFile{Path: module.GetPath(), Content: w.String()}
}

func kindOfFields(fields spec.NamedDefinitions) (bool, bool) {
	var hasRequiredFields = false
	var hasOptionalFields = false
	for _, field := range fields {
		if !field.Type.Definition.IsNullable() {
			hasRequiredFields = true
		} else {
			hasOptionalFields = true
		}
	}
	return hasRequiredFields, hasOptionalFields
}

func (g *Generator) objectModel(w *generator.Writer, model *spec.NamedModel) {
	hasRequiredFields, hasOptionalFields := kindOfFields(model.Object.Fields)
	if hasRequiredFields && hasOptionalFields {
		w.Line("export const T%s = t.intersection([", model.Name.PascalCase())
		w.Line("  t.interface({")
		for _, field := range model.Object.Fields {
			if !field.Type.Definition.IsNullable() {
				w.Line("    %s: %s,", field.Name.Source, g.RuntimeType(&field.Type.Definition))
			}
		}
		w.Line("  }),")
		w.Line("  t.partial({")
		for _, field := range model.Object.Fields {
			if field.Type.Definition.IsNullable() {
				w.Line("    %s: %s,", field.Name.Source, g.RuntimeType(&field.Type.Definition))
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
			w.Line("  %s: %s,", field.Name.Source, g.RuntimeType(&field.Type.Definition))
		}
		w.Line("})")
	}
	w.EmptyLine()
	w.Line("export type %s = t.TypeOf<typeof T%s>", model.Name.PascalCase(), model.Name.PascalCase())
}

func (g *Generator) enumModel(w *generator.Writer, model *spec.NamedModel) {
	w.Line("export enum %s {", model.Name.PascalCase())
	for _, item := range model.Enum.Items {
		w.Line(`  %s = "%s",`, item.Name.UpperCase(), item.Value)
	}
	w.Line("}")
	w.EmptyLine()
	w.Line("export const T%s = t.enum(%s)", model.Name.PascalCase(), model.Name.PascalCase())
}

func (g *Generator) unionModel(w *generator.Writer, model *spec.NamedModel) {
	if model.OneOf.Discriminator != nil {
		w.Line("export const T%s = t.union([", model.Name.PascalCase())
		for _, item := range model.OneOf.Items {
			w.Line("  t.intersection([t.type({%s: t.literal('%s')}), %s]),", common.TSIdentifier(*model.OneOf.Discriminator), item.Name.Source, g.RuntimeType(&item.Type.Definition))
		}
		w.Line("])")
	} else {
		w.Line("export const T%s = t.union([", model.Name.PascalCase())
		for _, item := range model.OneOf.Items {
			w.Line("  t.interface({%s: %s}),", item.Name.Source, g.RuntimeType(&item.Type.Definition))
		}
		w.Line("])")
	}
	w.EmptyLine()
	w.Line("export type %s = t.TypeOf<typeof T%s>", model.Name.PascalCase(), model.Name.PascalCase())
}
