package gents

import (
	"fmt"
	spec "github.com/specgen-io/spec"
)

func TsType(typ *spec.TypeDef) string {
	return PackagedTsType(typ, nil)
}

func PackagedTsType(typ *spec.TypeDef, modelsPackage *string) string {
	switch typ.Node {
	case spec.PlainType:
		return PlainTsType(typ.Plain, modelsPackage)
	case spec.NullableType:
		child := PackagedTsType(typ.Child, modelsPackage)
		result := child + " | undefined"
		return result
	case spec.ArrayType:
		child := PackagedTsType(typ.Child, modelsPackage)
		result := child + "[]"
		return result
	case spec.MapType:
		child := PackagedTsType(typ.Child, modelsPackage)
		result := "Record<string, " + child + ">"
		return result
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func PlainTsType(typ string, modelsPackage *string) string {
	switch typ {
	case spec.TypeInt32:
		return "number"
	case spec.TypeInt64:
		return "number"
	case spec.TypeFloat:
		return "number"
	case spec.TypeDouble:
		return "number"
	case spec.TypeDecimal:
		return "number"
	case spec.TypeBoolean:
		return "boolean"
	case spec.TypeString:
		return "string"
	case spec.TypeUuid:
		return "string"
	case spec.TypeDate:
		return "string"
	case spec.TypeDateTime:
		return "Date"
	case spec.TypeJson:
		return "unknown"
	default:
		if modelsPackage != nil {
			return fmt.Sprintf("%s.%s", *modelsPackage, typ)
		} else {
			return typ
		}
	}
}
