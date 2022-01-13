package genjava

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/spec"
)

type Types struct {
	Jsonlib string
}

func NewTypes(jsonlib string) *Types {
	return &Types{jsonlib}
}

func (t *Types) JavaType(typ *spec.TypeDef) string {
	javaType, _ := t.javaType(typ, false)
	return javaType
}

func (t *Types) JavaIsReferenceType(typ *spec.TypeDef) bool {
	_, isReference := t.javaType(typ, false)
	return isReference
}

func (t *Types) javaType(typ *spec.TypeDef, referenceTypesOnly bool) (string, bool) {
	switch typ.Node {
	case spec.PlainType:
		return t.PlainJavaType(typ.Plain, referenceTypesOnly)
	case spec.NullableType:
		return t.javaType(typ.Child, true)
	case spec.ArrayType:
		child, _ := t.javaType(typ.Child, false)
		result := child + "[]"
		return result, true
	case spec.MapType:
		child, _ := t.javaType(typ.Child, true)
		result := "Map<String, " + child + ">"
		return result, true
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

// PlainJavaType Returns Java type name and boolean indicating if the type is reference or value type
func (t *Types) PlainJavaType(typ string, referenceTypesOnly bool) (string, bool) {
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
		if t.Jsonlib == Jackson {
			result = "JsonNode"
		}
		if t.Jsonlib == Moshi {
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
