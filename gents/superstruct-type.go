package gents

import (
	"fmt"
	spec "github.com/specgen-io/spec"
)

func SsType(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return PlainSsType(typ.Plain)
	case spec.NullableType:
		child := SsType(typ.Child)
		result := "t.optional(t.nullable(" + child + "))"
		return result
	case spec.ArrayType:
		child := SsType(typ.Child)
		result := "t.array(" + child + ")"
		return result
	case spec.MapType:
		child := SsType(typ.Child)
		result := "t.record(t.string(), " + child + ")"
		return result
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func PlainSsType(typ string) string {
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
		return "t.DateTime"
	case spec.TypeJson:
		return "t.unknown()"
	default:
		return "T" + typ
	}
}


func StringSsType(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return StringPlainSsType(typ.Plain)
	case spec.NullableType:
		if typ.Child.Node != spec.PlainType {
			panic(fmt.Sprintf("Unsupported string type: %v", typ))
		}
		child := StringPlainSsType(typ.Child.Plain)
		result := "t.optional(" + child + ")"
		return result
	case spec.ArrayType:
		if typ.Child.Node != spec.PlainType {
			panic(fmt.Sprintf("Unsupported string type: %v", typ))
		}
		child := StringPlainSsType(typ.Child.Plain)
		result := "t.array(" + child + ")"
		return result
	default:
		panic(fmt.Sprintf("Unsupported string type: %v", typ))
	}
}

func StringPlainSsType(typ string) string {
	switch typ {
	case spec.TypeInt32:
		return "t.StrInteger"
	case spec.TypeInt64:
		return "t.StrInteger"
	case spec.TypeFloat:
		return "t.StrFloat"
	case spec.TypeDouble:
		return "t.StrFloat"
	case spec.TypeDecimal:
		return "t.StrFloat"
	case spec.TypeBoolean:
		return "t.StrBoolean"
	case spec.TypeString:
		return "t.string()"
	case spec.TypeUuid:
		return "t.string()"
	case spec.TypeDate:
		return "t.string()"
	case spec.TypeDateTime:
		return "t.DateTime"
	default:
		return "models.T"+typ
	}
}