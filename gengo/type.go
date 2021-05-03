package gengo

import (
	"fmt"
	spec "github.com/specgen-io/spec.v2"
)

func GoType(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return PlainGoType(typ.Plain)
	case spec.NullableType:
		child := GoType(typ.Child)
		result := "*" + child
		return result
	case spec.ArrayType:
		child := GoType(typ.Child)
		result := "[]" + child
		return result
	case spec.MapType:
		child := GoType(typ.Child)
		result := "map[string]" + child
		return result
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func PlainGoType(typ string) string {
	switch typ {
	case spec.TypeInt32:
		return "int32"
	case spec.TypeInt64:
		return "int64"
	case spec.TypeFloat:
		return "float32"
	case spec.TypeDouble:
		return "float64"
	case spec.TypeDecimal:
		return "github.com/shopspring/decimal"
	case spec.TypeBoolean:
		return "bool"
	case spec.TypeString:
		return "string"
	case spec.TypeUuid:
		return "github.com/google/uuid"
	case spec.TypeDate:
		return "time"
	case spec.TypeDateTime:
		return "time"
	case spec.TypeJson:
		return "encoding/json"
	default:
		return typ
	}
}
