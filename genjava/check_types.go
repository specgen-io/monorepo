package genjava

import (
	"fmt"
	"github.com/specgen-io/spec"
)

func checkType(fieldType *spec.TypeDef, typ string) bool {
	if fieldType.Plain != typ {
		switch fieldType.Node {
		case spec.PlainType:
			return false
		case spec.NullableType:
			return checkType(fieldType.Child, typ)
		case spec.ArrayType:
			return checkType(fieldType.Child, typ)
		case spec.MapType:
			return checkType(fieldType.Child, typ)
		default:
			panic(fmt.Sprintf("Unknown type: %v", typ))
		}
	}
	return true
}

func modelsHasType(model *spec.NamedModel, typ string) bool {
	if model.IsObject() {
		for _, field := range model.Object.Fields {
			if checkType(&field.Type.Definition, typ) {
				return true
			}
		}
	}
	return false
}

func checkMapType(fieldType *spec.TypeDef) bool {
	if fieldType.Node == spec.MapType {
		return true
	}
	return false
}

func modelsHasMapType(model *spec.NamedModel) bool {
	if model.IsObject() {
		for _, field := range model.Object.Fields {
			if checkMapType(&field.Type.Definition) {
				return true
			}
		}
	}
	return false
}
