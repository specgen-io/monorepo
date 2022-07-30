package scala

import (
	"fmt"
	"github.com/specgen-io/specgen/spec/v2"
)

func ScalaType(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return PlainScalaType(typ.Plain)
	case spec.NullableType:
		child := ScalaType(typ.Child)
		result := "Option[" + child + "]"
		return result
	case spec.ArrayType:
		child := ScalaType(typ.Child)
		result := "List[" + child + "]"
		return result
	case spec.MapType:
		child := ScalaType(typ.Child)
		result := "Map[String, " + child + "]"
		return result
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func PlainScalaType(typ string) string {
	switch typ {
	case spec.TypeInt32:
		return "Int"
	case spec.TypeInt64:
		return "Long"
	case spec.TypeFloat:
		return "Float"
	case spec.TypeDouble:
		return "Double"
	case spec.TypeDecimal:
		return "BigDecimal"
	case spec.TypeBoolean:
		return "Boolean"
	case spec.TypeString:
		return "String"
	case spec.TypeUuid:
		return "java.util.UUID"
	case spec.TypeDate:
		return "java.time.LocalDate"
	case spec.TypeDateTime:
		return "java.time.LocalDateTime"
	case spec.TypeJson:
		return "io.circe.Json"
	case spec.TypeEmpty:
		return "Unit"
	default:
		return typ
	}
}
