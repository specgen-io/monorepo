package genjava

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/spec"
)

func JavaType(typ *spec.TypeDef, jsonlib string) string {
	javaType, _ := javaType(typ, false, jsonlib)
	return javaType
}

func JavaIsReferenceType(typ *spec.TypeDef, jsonlib string) bool {
	_, isReference := javaType(typ, false, jsonlib)
	return isReference
}

func javaType(typ *spec.TypeDef, referenceTypesOnly bool, jsonlib string) (string, bool) {
	switch typ.Node {
	case spec.PlainType:
		return PlainJavaType(typ.Plain, referenceTypesOnly, jsonlib)
	case spec.NullableType:
		return javaType(typ.Child, true, jsonlib)
	case spec.ArrayType:
		child, _ := javaType(typ.Child, false, jsonlib)
		result := child + "[]"
		return result, true
	case spec.MapType:
		child, _ := javaType(typ.Child, true, jsonlib)
		result := "Map<String, " + child + ">"
		return result, true
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

// PlainJavaType Returns Java type name and boolean indicating if the type is reference or value type
func PlainJavaType(typ string, referenceTypesOnly bool, jsonlib string) (string, bool) {
	switch typ {
	case spec.TypeInt32:
		if referenceTypesOnly {
			return "Integer", true
		} else {
			return "int", false
		}
	case spec.TypeInt64:
		if referenceTypesOnly {
			return "Long", true
		} else {
			return "long", false
		}
	case spec.TypeFloat:
		if referenceTypesOnly {
			return "Float", true
		} else {
			return "float", false
		}
	case spec.TypeDouble:
		if referenceTypesOnly {
			return "Double", true
		} else {
			return "double", false
		}
	case spec.TypeDecimal:
		return "BigDecimal", true
	case spec.TypeBoolean:
		if referenceTypesOnly {
			return "Boolean", true
		} else {
			return "boolean", false
		}
	case spec.TypeString:
		return "String", true
	case spec.TypeUuid:
		return "UUID", true
	case spec.TypeDate:
		return "LocalDate", true
	case spec.TypeDateTime:
		return "LocalDateTime", true
	case spec.TypeJson:
		var result string
		if jsonlib == Jackson {
			result = "JsonNode"
		}
		if jsonlib == Moshi {
			result = "Map<String, Object>"
		}
		return result, true
	case spec.TypeEmpty:
		return "void", false
	default:
		return typ, true
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
