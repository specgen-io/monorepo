package validation

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/spec"
)

func IoTsType(typ *spec.TypeDef) string {
	return IoTsTypeFromPackage("", typ)
}

func IoTsTypeFromPackage(customTypesPackage string, typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return plainIoTsType(customTypesPackage, typ.Plain)
	case spec.NullableType:
		child := IoTsTypeFromPackage(customTypesPackage, typ.Child)
		result := "t.union([" + child + ", t.null])"
		return result
	case spec.ArrayType:
		child := IoTsTypeFromPackage(customTypesPackage, typ.Child)
		result := "t.array(" + child + ")"
		return result
	case spec.MapType:
		child := IoTsTypeFromPackage(customTypesPackage, typ.Child)
		result := "t.record(t.string, " + child + ")"
		return result
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func plainIoTsType(customTypesPackage string, typ string) string {
	switch typ {
	case spec.TypeInt32:
		return "t.number"
	case spec.TypeInt64:
		return "t.number"
	case spec.TypeFloat:
		return "t.number"
	case spec.TypeDouble:
		return "t.number"
	case spec.TypeDecimal:
		return "t.number"
	case spec.TypeBoolean:
		return "t.boolean"
	case spec.TypeString:
		return "t.string"
	case spec.TypeUuid:
		return "t.string"
	case spec.TypeDate:
		return "t.string"
	case spec.TypeDateTime:
		return "t.DateISOStringNoTimezone"
	case spec.TypeJson:
		return "t.unknown"
	default:
		if customTypesPackage == "" {
			return "T" + typ
		} else {
			return customTypesPackage + ".T" + typ
		}
	}
}
