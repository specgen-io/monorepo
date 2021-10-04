package genjava

import (
	"fmt"
	"github.com/specgen-io/spec"
)

func JavaType(typ *spec.TypeDef) string {
	return javaType(typ, false)
}

func javaType(typ *spec.TypeDef, referenceTypesOnly bool) string {
	switch typ.Node {
	case spec.PlainType:
		return PlainJavaType(typ.Plain, referenceTypesOnly)
	case spec.NullableType:
		result := javaType(typ.Child, true)
		return result
	case spec.ArrayType:
		child := javaType(typ.Child, false)
		result := child + "[]"
		return result
	case spec.MapType:
		child := javaType(typ.Child, true)
		result := "Map<String, " + child + ">"
		return result
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func PlainJavaType(typ string, referenceTypesOnly bool) string {
	switch typ {
	case spec.TypeInt32:
		if referenceTypesOnly {
			return "Integer"
		} else {
			return "int"
		}
	case spec.TypeInt64:
		if referenceTypesOnly {
			return "Long"
		} else {
			return "long"
		}
	case spec.TypeFloat:
		if referenceTypesOnly {
			return "Float"
		} else {
			return "float"
		}
	case spec.TypeDouble:
		if referenceTypesOnly {
			return "Double"
		} else {
			return "double"
		}
	case spec.TypeDecimal:
		return "BigDecimal"
	case spec.TypeBoolean:
		if referenceTypesOnly {
			return "Boolean"
		} else {
			return "boolean"
		}
	case spec.TypeString:
		return "String"
	case spec.TypeUuid:
		return "UUID"
	case spec.TypeDate:
		return "LocalDate"
	case spec.TypeDateTime:
		return "LocalDateTime"
	case spec.TypeJson:
		return "Object"
	case spec.TypeEmpty:
		return "void"
	default:
		return typ
	}
}

func checkDateType(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return addDateFormatAnnotation(typ.Plain)
	case spec.NullableType:
		return checkDateType(typ.Child)
	case spec.ArrayType:
		return checkDateType(typ.Child)
	case spec.MapType:
		return checkDateType(typ.Child)
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func addDateFormatAnnotation(typ string) string {
	switch typ {
	case spec.TypeDate:
		return "@DateTimeFormat(iso = DateTimeFormat.ISO.DATE)"
	case spec.TypeDateTime:
		return "@DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME)"
	default:
		return ""
	}
}
