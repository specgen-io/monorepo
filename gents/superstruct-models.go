package gents

import (
	spec "github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/gen"
	"path/filepath"
)

func GenerateSuperstructModels(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}
	files := []gen.TextFile{}
	for _, version := range specification.Versions {
		w := NewTsWriter()
		generateSuperstructModels(w, &version)
		filename := "index.ts"
		if version.Version.Source != "" {
			filename = version.Version.FlatCase() + ".ts"
		}
		files = append(files, gen.TextFile{Path: filepath.Join(generatePath, filename), Content: w.String()})
	}
	superstruct := generateSuperstruct(filepath.Join(generatePath, "superstruct.ts"))
	files = append(files, *superstruct)
	return gen.WriteFiles(files, true)
}

func generateSuperstructModels(w *gen.Writer, version *spec.Version) {
	w.Line(`import * as t from './superstruct'`)
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
	w.Line("export const T%s = t.object({", model.Name.PascalCase())
	for _, field := range model.Object.Fields {
		w.Line("  %s: %s,", field.Name.Source, SsType(&field.Type.Definition))
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
	w.Line("export const T%s = t.union([", model.Name.PascalCase())
	for _, item := range model.OneOf.Items {
		w.Line("  t.object({%s: %s}),", item.Name.Source, SsType(&item.Type.Definition))
	}
	w.Line("])")
	w.Line("")
	w.Line("export type %s = t.Infer<typeof T%s>", model.Name.PascalCase(), model.Name.PascalCase())
}
