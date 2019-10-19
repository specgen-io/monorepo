package gengql

import (
	"fmt"
	"github.com/ModaOperandi/spec"
)

func GraphqlType(typ *spec.Type) string {
	switch typ.Node {
	case spec.PlainType:
		return PlainGraphqlType(typ.PlainType)
	case spec.NullableType:
		child := GraphqlType(typ.Child)
		result := "Option[" + child + "]"
		return result
	case spec.ArrayType:
		child := GraphqlType(typ.Child)
		result := "List[" + child + "]"
		return result
	case spec.MapType:
		child := GraphqlType(typ.Child)
		result := "Map[String, " + child + "]"
		return result
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func PlainGraphqlType(typ string) string {
	return typ
	//switch typ {
	//case spec.TypeByte:
	//	return "Byte"
	//case spec.TypeShort, spec.TypeInt16:
	//	return "Short"
	//case spec.TypeInt, spec.TypeInt32:
	//	return "Int"
	//case spec.TypeLong, spec.TypeInt64:
	//	return "Long"
	//case spec.TypeFloat:
	//	return "Float"
	//case spec.TypeDouble:
	//	return "Double"
	//case spec.TypeDecimal:
	//	return "BigDecimal"
	//case spec.TypeBool, spec.TypeBoolean:
	//	return "Boolean"
	//case spec.TypeChar:
	//	return "Char"
	//case spec.TypeString, spec.TypeStr:
	//	return "String"
	//case spec.TypeDate:
	//	return "LocalDate"
	//case spec.TypeDateTime:
	//	return "LocalDateTime"
	//case spec.TypeTime:
	//	return "LocalTime"
	//case spec.TypeJson:
	//	return "JsonNode"
	//default:
	//	return typ
	//}
}
