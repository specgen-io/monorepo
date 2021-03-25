package gents

import (
	"bytes"
	spec "github.com/specgen-io/spec.v1"
	"specgen/gen"
)

func GenerateIoTsModels(spec *spec.Spec, outPath string) *gen.TextFile {
	w := new(bytes.Buffer)
	generateIoTsModels(spec, w)
	return &gen.TextFile{Path: outPath, Content: w.String()}
}

func kindOfFields(objectModel spec.NamedModel) (bool, bool) {
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