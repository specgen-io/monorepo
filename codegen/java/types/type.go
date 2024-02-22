package types

import (
	"fmt"
	"spec"
)

const TextType = `String`
const EmptyType = `void`

type Types struct {
	RawJsonType string
	BinaryType  BinaryType
}

type BinaryType struct {
	RequestType  string
	ResponseType string
}

func NewTypes(rawJsonType, requestBinaryType, responseBinaryType string) *Types {
	return &Types{RawJsonType: rawJsonType, BinaryType: BinaryType{RequestType: requestBinaryType, ResponseType: responseBinaryType}}
}

func (types *Types) RequestBodyJavaType(body *spec.RequestBody) string {
	switch body.Kind() {
	case spec.BodyBinary:
		return types.BinaryType.RequestType
	case spec.BodyText:
		return TextType
	case spec.BodyEmpty:
		return EmptyType
	case spec.BodyJson:
		return types.Java(&body.Type.Definition)
	default:
		panic(fmt.Sprintf("Unknown response body kind: %v", body.Kind()))
	}
}

func (types *Types) ResponseBodyJavaType(body *spec.ResponseBody) string {
	switch body.Kind() {
	case spec.BodyBinary:
		return types.BinaryType.ResponseType
	case spec.BodyText:
		return TextType
	case spec.BodyEmpty:
		return EmptyType
	case spec.BodyJson:
		return types.Java(&body.Type.Definition)
	default:
		panic(fmt.Sprintf("Unknown response body kind: %v", body.Kind()))
	}
}

func (t *Types) Java(typ *spec.TypeDef) string {
	javaType, _ := t.javaType(typ, false)
	return javaType
}

func (t *Types) IsReference(typ *spec.TypeDef) bool {
	_, isReference := t.javaType(typ, false)
	return isReference
}

func (t *Types) javaType(typ *spec.TypeDef, referenceTypesOnly bool) (string, bool) {
	switch typ.Node {
	case spec.PlainType:
		return t.plainJavaType(typ.Plain, referenceTypesOnly)
	case spec.NullableType:
		return t.javaType(typ.Child, true)
	case spec.ArrayType:
		child, _ := t.javaType(typ.Child, true)
		result := "List<" + child + ">"
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
func (t *Types) plainJavaType(typ string, referenceTypesOnly bool) (string, bool) {
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
		return t.RawJsonType, true
	case spec.TypeEmpty:
		return "void", false
	default:
		return typ, true
	}
}

func (t *Types) Imports() []string {
	return []string{
		`java.math.BigDecimal`,
		`java.time.*`,
		`java.util.*`,
		`java.io.*`,
	}
}
