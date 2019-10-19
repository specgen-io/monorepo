package genscala

import (
	"fmt"
	"github.com/ModaOperandi/spec"
)

func ScalaType(typ *spec.Type) string {
	switch typ.Node {
	case spec.PlainType:
		return PlainScalaType(typ.PlainType)
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
	case spec.TypeByte:
		return "Byte"
	case spec.TypeShort, spec.TypeInt16:
		return "Short"
	case spec.TypeInt, spec.TypeInt32:
		return "Int"
	case spec.TypeLong, spec.TypeInt64:
		return "Long"
	case spec.TypeFloat:
		return "Float"
	case spec.TypeDouble:
		return "Double"
	case spec.TypeDecimal:
		return "BigDecimal"
	case spec.TypeBool, spec.TypeBoolean:
		return "Boolean"
	case spec.TypeChar:
		return "Char"
	case spec.TypeString, spec.TypeStr:
		return "String"
	case spec.TypeUuid:
		return "java.util.UUID"
	case spec.TypeDate:
		return "java.time.LocalDate"
	case spec.TypeDateTime:
		return "java.time.LocalDateTime"
	case spec.TypeTime:
		return "java.time.LocalTime"
	case spec.TypeJson:
		return "JsonNode"
	default:
		return typ
	}
}
