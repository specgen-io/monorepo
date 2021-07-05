package gents

import (
	"fmt"
	"github.com/specgen-io/spec"
)

func StringSsTypeDefaulted(param *spec.NamedParam) string {
	theType := StringSsType(&param.Type.Definition)
	if param.Default != nil {
		theType = fmt.Sprintf("t.defaulted(%s, %s)", theType, DefaultValue(&param.Type.Definition, *param.Default))
	}
	return theType
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
