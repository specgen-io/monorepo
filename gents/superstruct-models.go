package gents

import (
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
)

func generateSuperstructModels(specification *spec.Spec, generatePath string) []gen.TextFile {
	files := []gen.TextFile{}
	for _, version := range specification.Versions {
		w := NewTsWriter()
		generateSuperstructVersionModels(w, &version)
		versionModelsFile := gen.TextFile{Path: filepath.Join(generatePath, versionFilename(&version, "models", "ts")), Content: w.String()}
		files = append(files, versionModelsFile)
	}
	staticCode := generateSuperstructStaticCode(filepath.Join(generatePath, "superstruct.ts"))
	files = append(files, *staticCode)
	return files
}

func generateSuperstructVersionModels(w *gen.Writer, version *spec.Version) {
	w.Line(importSuperstructEncoding)
	for _, model := range version.ResolvedModels {
		w.EmptyLine()
		if model.IsObject() {
			generateSuperstructObjectModel(w, model)
		} else if model.IsEnum() {
			generateSuperstructEnumModel(w, model)
		} else if model.IsOneOf() {
			generateSuperstructUnionModel(w, model)
		}
	}
}

func generateSuperstructObjectModel(w *gen.Writer, model *spec.NamedModel) {
	w.Line("export const T%s = t.type({", model.Name.PascalCase())
	for _, field := range model.Object.Fields {
		w.Line("  %s: %s,", field.Name.Source, SuperstructType(&field.Type.Definition))
	}
	w.Line("})")
	w.Line("")
	w.Line("export type %s = t.Infer<typeof T%s>", model.Name.PascalCase(), model.Name.PascalCase())
}

func generateSuperstructEnumModel(w *gen.Writer, model *spec.NamedModel) {
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

func generateSuperstructUnionModel(w *gen.Writer, model *spec.NamedModel) {
	if model.OneOf.Discriminator != nil {
		w.Line("export const T%s = t.union([", model.Name.PascalCase())
		for _, item := range model.OneOf.Items {
			w.Line("  t.intersection([t.type({%s: t.literal('%s')}), %s]),", tsIdentifier(*model.OneOf.Discriminator), item.Name.Source, SuperstructType(&item.Type.Definition))
		}
		w.Line("])")
	} else {
		w.Line("export const T%s = t.union([", model.Name.PascalCase())
		for _, item := range model.OneOf.Items {
			w.Line("  t.object({%s: %s}),", item.Name.Source, SuperstructType(&item.Type.Definition))
		}
		w.Line("])")
	}
	w.EmptyLine()
	w.Line("export type %s = t.Infer<typeof T%s>", model.Name.PascalCase(), model.Name.PascalCase())
}
