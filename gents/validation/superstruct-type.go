package validation

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/spec"
)

func SuperstructType(typ *spec.TypeDef) string {
	return SuperstructTypeFromPackage("", typ)
}

func SuperstructTypeFromPackage(customTypesPackage string, typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return plainSuperstructType(customTypesPackage, typ.Plain)
	case spec.NullableType:
		child := SuperstructTypeFromPackage(customTypesPackage, typ.Child)
		result := "t.optional(t.nullable(" + child + "))"
		return result
	case spec.ArrayType:
		child := SuperstructTypeFromPackage(customTypesPackage, typ.Child)
		result := "t.array(" + child + ")"
		return result
	case spec.MapType:
		child := SuperstructTypeFromPackage(customTypesPackage, typ.Child)
		result := "t.record(t.string(), " + child + ")"
		return result
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func plainSuperstructType(customTypesPackage string, typ string) string {
	switch typ {
	case spec.TypeInt32:
		return "t.number()"
	case spec.TypeInt64:
		return "t.number()"
	case spec.TypeFloat:
		return "t.number()"
	case spec.TypeDouble:
		return "t.number()"
	case spec.TypeDecimal:
		return "t.number()"
	case spec.TypeBoolean:
		return "t.boolean()"
	case spec.TypeString:
		return "t.string()"
	case spec.TypeUuid:
		return "t.string()"
	case spec.TypeDate:
		return "t.string()"
	case spec.TypeDateTime:
		return "t.StrDateTime"
	case spec.TypeJson:
		return "t.unknown()"
	default:
		if customTypesPackage == "" {
			return "T" + typ
		} else {
			return customTypesPackage + ".T" + typ
		}
	}
}
