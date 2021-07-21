package genjava

import (
	"fmt"
	"github.com/specgen-io/spec"
)

func checkType(fieldType *spec.TypeDef, typ string) bool {
	switch fieldType.Node {
	case spec.PlainType:
		if fieldType.Plain != typ {
			return false
		}
	case spec.NullableType:
		return checkType(fieldType.Child, typ)
	case spec.ArrayType:
		return checkType(fieldType.Child, typ)
	case spec.MapType:
		return checkType(fieldType.Child, typ)
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
	return true
}
