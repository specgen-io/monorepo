package genkotlin

import "github.com/specgen-io/specgen/v2/spec"

func isKotlinArrayType(model *spec.NamedModel) bool {
	for _, field := range model.Object.Fields {
		if _isKotlinArrayType(&field.Type.Definition) {
			return true
		}
	}
	return false
}

func _isKotlinArrayType(typ *spec.TypeDef) bool {
	switch typ.Node {
	case spec.ArrayType:
		return true
	case spec.NullableType:
		return _isKotlinArrayType(typ.Child)
	default:
		return false
	}
}
