package types

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/spec"
)

func DefaultValue(typ *spec.TypeDef, value string) string {
	switch typ.Node {
	case spec.ArrayType:
		if value == "[]" {
			return "[]"
		} else {
			panic(fmt.Sprintf("Array type %s default value %s is not supported", typ.Name, value))
		}
	case spec.MapType:
		if value == "{}" {
			return "{}"
		} else {
			panic(fmt.Sprintf("Map type %s default value %s is not supported", typ.Name, value))
		}
	case spec.PlainType:
		switch typ.Plain {
		case
			spec.TypeInt32,
			spec.TypeInt64,
			spec.TypeFloat,
			spec.TypeDouble,
			spec.TypeDecimal,
			spec.TypeBoolean:
			return value
		case
			spec.TypeString,
			spec.TypeUuid,
			spec.TypeDate:
			return `"` + value + `"`
		case spec.TypeDateTime:
			return `new Date("` + value + `")`
		default:
			model := typ.Info.Model
			if model != nil && model.IsEnum() {
				return `"` + value + `"`
			} else {
				panic(fmt.Sprintf("Type: %s does not support default value", typ.Name))
			}
		}
	default:
		panic(fmt.Sprintf("Type: %s does not support default value", typ.Name))
	}
}
