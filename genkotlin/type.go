package genkotlin

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/spec"
)

func KotlinType(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return PlainKotlinType(typ.Plain)
	case spec.NullableType:
		return KotlinType(typ.Child) + "?"
	case spec.ArrayType:
		child := KotlinType(typ.Child)
		return "List<" + child + ">"
	case spec.MapType:
		child := KotlinType(typ.Child)
		result := "Map<String, " + child + ">"
		return result
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func PlainKotlinType(typ string) string {
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
		return "UUID"
	case spec.TypeDate:
		return "LocalDate"
	case spec.TypeDateTime:
		return "LocalDateTime"
	case spec.TypeJson:
		return "JsonNode"
	case spec.TypeEmpty:
		return "Void"
	default:
		return typ
	}
}

func dateFormatAnnotation(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return dateFormatAnnotationPlain(typ.Plain)
	case spec.NullableType:
		return dateFormatAnnotation(typ.Child)
	case spec.ArrayType:
		return dateFormatAnnotation(typ.Child)
	case spec.MapType:
		return dateFormatAnnotation(typ.Child)
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func dateFormatAnnotationPlain(typ string) string {
	switch typ {
	case spec.TypeDate:
		return "@DateTimeFormat(iso = DateTimeFormat.ISO.DATE)"
	case spec.TypeDateTime:
		return "@DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME)"
	default:
		return ""
	}
}
