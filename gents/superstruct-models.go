package gents

import (
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
)

func generateSuperstructModels(specification *spec.Spec, generatePath string) []gen.TextFile {
	files := []gen.TextFile{}
	for _, version := range specification.Versions {
		files = append(files, *generateSuperstructVersionModels(&version, generatePath))
	}
	staticCode := generateSuperstructStaticCode(filepath.Join(generatePath, Superstruct+".ts"))
	files = append(files, *staticCode)
	return files
}

func generateSuperstructVersionModels(version *spec.Version, generatePath string) *gen.TextFile {
	filePath := versionedPath(generatePath, version, "models.ts")
	w := NewTsWriter()
	w.Line(`import * as t from '%s'`, importPath(filepath.Join(generatePath, Superstruct), filePath))
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
	return &gen.TextFile{Path: filePath, Content: w.String()}
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
