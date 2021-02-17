package genruby

import (
	"fmt"
	"github.com/specgen-io/spec"
)

func RubyType(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return PlainRubyType(typ.Plain)
	case spec.NullableType:
		child := RubyType(typ.Child)
		result := "T.nilable(" + child + ")"
		return result
	case spec.ArrayType:
		child := RubyType(typ.Child)
		result := "T.array(" + child + ")"
		return result
	case spec.MapType:
		child := RubyType(typ.Child)
		result := "T.hash(String, " + child + ")"
		return result
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func PlainRubyType(typ string) string {
	switch typ {
	case spec.TypeByte:
		return "Integer"
	case spec.TypeInt16:
		return "Integer"
	case spec.TypeInt32:
		return "Integer"
	case spec.TypeInt64:
		return "Integer"
	case spec.TypeFloat:
		return "Float"
	case spec.TypeDouble:
		return "Float"
	case spec.TypeDecimal:
		return "Float"
	case spec.TypeBoolean:
		return "Boolean"
	case spec.TypeChar:
		return "String"
	case spec.TypeString:
		return "String"
	case spec.TypeUuid:
		return "UUID"
	case spec.TypeDate:
		return "Date"
	case spec.TypeDateTime:
		return "DateTime"
	case spec.TypeTime:
		return "TimeOfDay"
	case spec.TypeJson:
		return "T.hash(String, Untyped)"
	default:
		return typ
	}
}

