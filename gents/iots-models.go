package gents

import (
	spec "github.com/specgen-io/spec.v2"
	"path/filepath"
	"specgen/gen"
)

func GenerateIoTsModels(serviceFile string, generatePath string) error {
	spec, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}

	files := []gen.TextFile{}
	for _, version := range spec.Versions {
		w := NewTsWriter()
		generateIoTsModels(w, &version)
		filename := "index.ts"
		if version.Version.Source != "" {
			filename = version.Version.FlatCase()+".ts"
		}
		files = append(files, gen.TextFile{Path: filepath.Join(generatePath, filename), Content: w.String()})
	}
	err = gen.WriteFiles(files, true)
	if err != nil {
		return err
	}

	return nil
}

func generateIoTsModels(w *gen.Writer, version *spec.Version) {
	w.Line("/* eslint-disable @typescript-eslint/camelcase */")
	w.Line("/* eslint-disable @typescript-eslint/no-magic-numbers */")
	w.Line("import * as t from './io-ts'")
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

func generateIoTsObjectModel(w *gen.Writer, model *spec.NamedModel) {
	hasRequiredFields, hasOptionalFields := kindOfFields(model)
	if hasRequiredFields && hasOptionalFields {
		w.Line("export const T%s = t.intersection([", model.Name.PascalCase())
		w.Line("  t.interface({")
		for _, field := range model.Object.Fields {
			if !field.Type.Definition.IsNullable() {
				w.Line("    %s: %s", field.Name.Source, IoTsType(&field.Type.Definition))
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
		w.Line("")
		w.Line("export const T%s = %s({", model.Name.PascalCase(), ioTsType)
		for _, field := range model.Object.Fields {
			w.Line("  %s: %s,", field.Name.Source, IoTsType(&field.Type.Definition))
		}
		w.Line("})")
	}
	w.EmptyLine()
	w.Line("export type %s = t.TypeOf<typeof T%s>", model.Name.PascalCase(), model.Name.PascalCase())
}

func generateIoTsEnumModel(w *gen.Writer, model *spec.NamedModel) {
	w.Line("export enum %s {", model.Name.PascalCase())
	for _, item := range model.Enum.Items {
		w.Line(`  %s = "%s",`, item.Name.UpperCase(), item.Value)
	}
	w.Line("}")
	w.EmptyLine()
	w.Line("export const T%s = t.enum(%s)", model.Name.PascalCase(), model.Name.PascalCase())
}

func generateIoTsUnionModel(w *gen.Writer, model *spec.NamedModel) {
	w.Line("export const T%s = t.union([", model.Name.PascalCase())
	for _, item := range model.OneOf.Items {
		w.Line("  t.interface({%s: %s}),", item.Name.Source, IoTsType(&item.Type.Definition))
	}
	w.Line("])")
	w.EmptyLine()
	w.Line("export type %s = t.TypeOf<typeof T%s>", model.Name.PascalCase(), model.Name.PascalCase())
}
