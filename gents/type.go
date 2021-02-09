package gents

import (
	"fmt"
	"github.com/ModaOperandi/spec"
)

func TsType(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return PlainTsType(typ.Plain)
	case spec.NullableType:
		child := TsType(typ.Child)
		result := child + " | null"
		return result
	case spec.ArrayType:
		child := TsType(typ.Child)
		result := child + "[]"
		return result
	case spec.MapType:
		child := TsType(typ.Child)
		result := "Record<string, " + child + ">"
		return result
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func PlainTsType(typ string) string {
	switch typ {
	case spec.TypeByte:
		return "number"
	case spec.TypeInt16:
		return "number"
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
	case spec.TypeChar:
		return "string"
	case spec.TypeString:
		return "string"
	case spec.TypeUuid:
		return "string"
	case spec.TypeDate:
		return "string"
	case spec.TypeDateTime:
		return "string"
	case spec.TypeTime:
		return "string"
	case spec.TypeJson:
		return "unknown"
	default:
		return typ
	}
}

