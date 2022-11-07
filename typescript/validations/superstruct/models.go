package superstruct

import (
	"generator"
	"spec"
	"typescript/common"
	"typescript/module"
	"typescript/writer"
)

func (g *Generator) Models(version *spec.Version) *generator.CodeFile {
	modelsModule := g.Modules.Models(version)
	return g.models(version.ResolvedModels, modelsModule)
}

func (g *Generator) ErrorModels(httpErrors *spec.HttpErrors) *generator.CodeFile {
	return g.models(httpErrors.ResolvedModels, g.Modules.ErrorModules)
}

func (g *Generator) models(models []*spec.NamedModel, modelsModule module.Module) *generator.CodeFile {
	w := writer.New(modelsModule)
	imports := w.Imports()
	imports.Star(g.Modules.Validation, `t`)
	imports.Write(w)
	for _, model := range models {
		w.EmptyLine()
		if model.IsObject() {
			g.objectModel(w, model)
		} else if model.IsEnum() {
			g.enumModel(w, model)
		} else if model.IsOneOf() {
			g.unionModel(w, model)
		}
	}
	return w.ToCodeFile()
}

func (g *Generator) objectModel(w *writer.Writer, model *spec.NamedModel) {
	w.Line("export const T%s = t.type({", model.Name.PascalCase())
	for _, field := range model.Object.Fields {
		w.Line("  %s: %s,", field.Name.Source, g.RuntimeTypeSamePackage(&field.Type.Definition))
	}
	w.Line("})")
	w.Line("")
	w.Line("export type %s = t.Infer<typeof T%s>", model.Name.PascalCase(), model.Name.PascalCase())
}

func (g *Generator) enumModel(w *writer.Writer, model *spec.NamedModel) {
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

func (g *Generator) unionModel(w *writer.Writer, model *spec.NamedModel) {
	if model.OneOf.Discriminator != nil {
		w.Line("export const T%s = t.union([", model.Name.PascalCase())
		for _, item := range model.OneOf.Items {
			w.Line("  t.intersection([t.type({%s: t.literal('%s')}), %s]),", common.TSIdentifier(*model.OneOf.Discriminator), item.Name.Source, g.RuntimeTypeSamePackage(&item.Type.Definition))
		}
		w.Line("])")
	} else {
		w.Line("export const T%s = t.union([", model.Name.PascalCase())
		for _, item := range model.OneOf.Items {
			w.Line("  t.object({%s: %s}),", item.Name.Source, g.RuntimeTypeSamePackage(&item.Type.Definition))
		}
		w.Line("])")
	}
	w.EmptyLine()
	w.Line("export type %s = t.Infer<typeof T%s>", model.Name.PascalCase(), model.Name.PascalCase())
}
