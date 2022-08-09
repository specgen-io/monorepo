package superstruct

import (
	"generator"
	"typescript/common"
	"typescript/modules"
	"typescript/writer"
	"spec"
)

func (g *Generator) VersionModels(version *spec.Version, superstructModule modules.Module, module modules.Module) *generator.CodeFile {
	w := writer.NewTsWriter()
	w.Line(`import * as t from '%s'`, superstructModule.GetImport(module))
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

func (g *Generator) objectModel(w *generator.Writer, model *spec.NamedModel) {
	w.Line("export const T%s = t.type({", model.Name.PascalCase())
	for _, field := range model.Object.Fields {
		w.Line("  %s: %s,", field.Name.Source, g.RuntimeType(&field.Type.Definition))
	}
	w.Line("})")
	w.Line("")
	w.Line("export type %s = t.Infer<typeof T%s>", model.Name.PascalCase(), model.Name.PascalCase())
}

func (g *Generator) enumModel(w *generator.Writer, model *spec.NamedModel) {
	w.Line("export const T%s = t.enums ([", model.Name.PascalCase())
	for _, item := range model.Enum.Items {
		w.Line(`  "%s",`, item.Value)
	}
	w.Line("])")
	w.EmptyLine()
	w.Line("export type %s = t.Infer<typeof T%s>", model.Name.PascalCase(), model.Name.PascalCase())
	w.EmptyLine()
	w.Line("export const %s = {", model.Name.PascalCase())
	for _, item := range model.Enum.Items {
		w.Line(`  %s: <%s>"%s",`, item.Name.UpperCase(), model.Name.PascalCase(), item.Value)
	}
	w.Line("}")
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
			w.Line("  t.object({%s: %s}),", item.Name.Source, g.RuntimeType(&item.Type.Definition))
		}
		w.Line("])")
	}
	w.EmptyLine()
	w.Line("export type %s = t.Infer<typeof T%s>", model.Name.PascalCase(), model.Name.PascalCase())
}
