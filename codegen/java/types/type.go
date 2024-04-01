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
	FileType    FileType
}

type BinaryType struct {
	RequestType  string
	ResponseType string
}

type FileType struct {
	RequestType  string
	ResponseType string
}

func NewTypes(rawJsonType, requestBinaryType, responseBinaryType, requestFileType, responseFileType string) *Types {
	return &Types{
		RawJsonType: rawJsonType,
		BinaryType:  BinaryType{RequestType: requestBinaryType, ResponseType: responseBinaryType},
		FileType:    FileType{RequestType: requestFileType, ResponseType: responseFileType},
	}
}

func (t *Types) RequestBodyJavaType(body *spec.RequestBody) string {
	switch body.Kind() {
	case spec.BodyText:
		return TextType
	case spec.BodyEmpty:
		return EmptyType
	case spec.BodyBinary:
		return t.BinaryType.RequestType
	case spec.BodyJson:
		return t.Java(&body.Type.Definition)
	default:
		panic(fmt.Sprintf("Unknown response body kind: %v", body.Kind()))
	}
}

func (t *Types) ResponseBodyJavaType(body *spec.ResponseBody) string {
	switch body.Kind() {
	case spec.BodyText:
		return TextType
	case spec.BodyEmpty:
		return EmptyType
	case spec.BodyBinary:
		return t.BinaryType.ResponseType
	case spec.BodyFile:
		return t.FileType.ResponseType
	case spec.BodyJson:
		return t.Java(&body.Type.Definition)
	default:
		panic(fmt.Sprintf("Unknown response body kind: %v", body.Kind()))
	}
}

func (t *Types) ParamJavaType(param *spec.NamedParam) string {
	if param.Type.Definition.String() == spec.TypeFile {
		return t.FileType.RequestType
	} else {
		return t.Java(&param.Type.Definition)
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
		return EmptyType, false
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
