package gents

import (
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
)

func generateIoTsModels(specification *spec.Spec, generatePath string) []gen.TextFile {
	files := []gen.TextFile{}
	for _, version := range specification.Versions {
		w := NewTsWriter()
		generateIoTsVersionModels(w, &version)
		versionModelsFile := gen.TextFile{Path: filepath.Join(generatePath, versionFilename(&version, "models", "ts")), Content: w.String()}
		files = append(files, versionModelsFile)
	}
	staticCode := generateIoTsStaticCode(filepath.Join(generatePath, "io-ts.ts"))
	files = append(files, *staticCode)
	return files
}

func generateIoTsVersionModels(w *gen.Writer, version *spec.Version) {
	w.Line("/* eslint-disable @typescript-eslint/camelcase */")
	w.Line("/* eslint-disable @typescript-eslint/no-magic-numbers */")
	w.Line(importIoTsEncoding)
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
				w.Line("    %s: %s", tsIdentifier(field.Name.Source), IoTsType(&field.Type.Definition))
			}
		}
		w.Line("  }),")
		w.Line("  t.partial({")
		for _, field := range model.Object.Fields {
			if field.Type.Definition.IsNullable() {
				w.Line("    %s: %s,", tsIdentifier(field.Name.Source), IoTsType(&field.Type.Definition))
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
			w.Line("  %s: %s,", tsIdentifier(field.Name.Source), IoTsType(&field.Type.Definition))
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
		w.Line("  t.interface({%s: %s}),", tsIdentifier(item.Name.Source), IoTsType(&item.Type.Definition))
	}
	w.Line("])")
	w.EmptyLine()
	w.Line("export type %s = t.TypeOf<typeof T%s>", model.Name.PascalCase(), model.Name.PascalCase())
}
